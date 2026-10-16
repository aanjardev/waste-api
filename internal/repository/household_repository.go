package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/aanjardev/waste-api/internal/model"
)

type HouseholdRepository interface {
	Create(ctx context.Context, h *model.Household) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Household, error)
	FindAll(ctx context.Context, limit, offset int) ([]model.Household, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type householdRepository struct {
	pool *pgxpool.Pool
}

func NewHouseholdRepository(pool *pgxpool.Pool) HouseholdRepository {
	return &householdRepository{pool: pool}
}

func (r *householdRepository) Create(ctx context.Context, h *model.Household) error {
	query := `
		INSERT INTO households (id, owner_name, address, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query, h.ID, h.OwnerName, h.Address).
		Scan(&h.CreatedAt, &h.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert household: %w", err)
	}
	return nil
}

func (r *householdRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Household, error) {
	query := `
		SELECT id, owner_name, address, created_at, updated_at
		FROM households
		WHERE id = $1
	`
	var h model.Household
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&h.ID, &h.OwnerName, &h.Address, &h.CreatedAt, &h.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("find household: %w", err)
	}
	return &h, nil
}

func (r *householdRepository) FindAll(ctx context.Context, limit, offset int) ([]model.Household, int, error) {
	countQuery := `SELECT COUNT(*) FROM households`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count households: %w", err)
	}

	query := `
		SELECT id, owner_name, address, created_at, updated_at
		FROM households
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list households: %w", err)
	}
	defer rows.Close()

	var households []model.Household
	for rows.Next() {
		var h model.Household
		if err := rows.Scan(&h.ID, &h.OwnerName, &h.Address, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan household: %w", err)
		}
		households = append(households, h)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows err: %w", err)
	}

	return households, total, nil
}

func (r *householdRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM households WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete household: %w", err)
	}
	if result.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *householdRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM households WHERE id = $1)`
	var exists bool
	if err := r.pool.QueryRow(ctx, query, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("check household exists: %w", err)
	}
	return exists, nil
}