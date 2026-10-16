package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aanjardev/waste-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PickupRepository interface {
	Create(ctx context.Context, p *model.WastePickup) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.WastePickup, error)
	FindAll(ctx context.Context, status string, householdID *uuid.UUID, limit, offset int) ([]model.WastePickup, int, error)
	Update(ctx context.Context, p *model.WastePickup) error
	UpdateStatusIf(ctx context.Context, id uuid.UUID, from, to model.PickupStatus, pickupDate *time.Time) (*model.WastePickup, error)
	UpdateStatusIfTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, from, to model.PickupStatus, pickupDate *time.Time) (*model.WastePickup, error)
}

type pickupRepository struct {
	pool *pgxpool.Pool
}

func NewPickupRepository(pool *pgxpool.Pool) PickupRepository {
	return &pickupRepository{pool: pool}
}

func (r *pickupRepository) Create(ctx context.Context, p *model.WastePickup) error {
	query := `
		INSERT INTO waste_pickups (id, household_id, type, status, pickup_date, safety_check, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		p.ID, p.HouseholdID, p.Type, p.Status, p.PickupDate, p.SafetyCheck,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert pickup: %w", err)
	}
	return nil
}

func (r *pickupRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.WastePickup, error) {
	query := `
		SELECT id, household_id, type, status, pickup_date, safety_check, created_at, updated_at
		FROM waste_pickups
		WHERE id = $1
	`
	var p model.WastePickup
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.HouseholdID, &p.Type, &p.Status,
		&p.PickupDate, &p.SafetyCheck, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("find pickup: %w", err)
	}
	return &p, nil
}

func (r *pickupRepository) FindAll(ctx context.Context, status string, householdID *uuid.UUID, limit, offset int) ([]model.WastePickup, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if status != "" {
		where += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}
	if householdID != nil {
		where += fmt.Sprintf(" AND household_id = $%d", argIdx)
		args = append(args, *householdID)
		argIdx++
	}

	// Count
	countQuery := "SELECT COUNT(*) FROM waste_pickups " + where
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count pickups: %w", err)
	}

	// Data
	query := fmt.Sprintf(`
		SELECT id, household_id, type, status, pickup_date, safety_check, created_at, updated_at
		FROM waste_pickups
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list pickups: %w", err)
	}
	defer rows.Close()

	var pickups []model.WastePickup
	for rows.Next() {
		var p model.WastePickup
		if err := rows.Scan(
			&p.ID, &p.HouseholdID, &p.Type, &p.Status,
			&p.PickupDate, &p.SafetyCheck, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan pickup: %w", err)
		}
		pickups = append(pickups, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows err: %w", err)
	}

	return pickups, total, nil
}

func (r *pickupRepository) Update(ctx context.Context, p *model.WastePickup) error {
	query := `
		UPDATE waste_pickups
		SET type = $2, status = $3, pickup_date = $4, safety_check = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		p.ID, p.Type, p.Status, p.PickupDate, p.SafetyCheck,
	).Scan(&p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrNotFound
		}
		return fmt.Errorf("update pickup: %w", err)
	}
	return nil
}

func (r *pickupRepository) UpdateStatusIf(
	ctx context.Context,
	id uuid.UUID,
	from, to model.PickupStatus,
	pickupDate *time.Time,
) (*model.WastePickup, error) {
	query := `
		UPDATE waste_pickups
		SET status = $3, pickup_date = COALESCE($4, pickup_date), updated_at = NOW()
		WHERE id = $1 AND status = $2
		RETURNING id, household_id, type, status, pickup_date, safety_check, created_at, updated_at
	`
	var p model.WastePickup
	err := r.pool.QueryRow(ctx, query, id, from, to, pickupDate).Scan(
		&p.ID, &p.HouseholdID, &p.Type, &p.Status,
		&p.PickupDate, &p.SafetyCheck, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Tidak ada baris ter-update: berarti status sudah berubah duluan
			return nil, nil
		}
		return nil, fmt.Errorf("conditional update pickup: %w", err)
	}
	return &p, nil
}

func (r *pickupRepository) UpdateStatusIfTx(
	ctx context.Context,
	tx pgx.Tx,
	id uuid.UUID,
	from, to model.PickupStatus,
	pickupDate *time.Time,
) (*model.WastePickup, error) {
	query := `
		UPDATE waste_pickups
		SET status = $3, pickup_date = COALESCE($4, pickup_date), updated_at = NOW()
		WHERE id = $1 AND status = $2
		RETURNING id, household_id, type, status, pickup_date, safety_check, created_at, updated_at
	`
	var p model.WastePickup
	err := tx.QueryRow(ctx, query, id, from, to, pickupDate).Scan(
		&p.ID, &p.HouseholdID, &p.Type, &p.Status,
		&p.PickupDate, &p.SafetyCheck, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("conditional update pickup in tx: %w", err)
	}
	return &p, nil
}
