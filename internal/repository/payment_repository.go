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

type PaymentRepository interface {
	Create(ctx context.Context, p *model.Payment) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Payment, error)
	FindAll(ctx context.Context, status string, householdID *uuid.UUID, from, to *string, limit, offset int) ([]model.Payment, int, error)
	Update(ctx context.Context, p *model.Payment) error
	HasPendingByHousehold(ctx context.Context, householdID uuid.UUID) (bool, error)
	ExistsByWasteID(ctx context.Context, wasteID uuid.UUID) (bool, error)
	CreateTx(ctx context.Context, tx pgx.Tx, p *model.Payment) error
	ExistsByWasteIDTx(ctx context.Context, tx pgx.Tx, wasteID uuid.UUID) (bool, error)
}

type paymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) PaymentRepository {
	return &paymentRepository{pool: pool}
}

func (r *paymentRepository) Create(ctx context.Context, p *model.Payment) error {
	query := `
		INSERT INTO payments (id, household_id, waste_id, amount, payment_date, status, proof_file_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		p.ID, p.HouseholdID, p.WasteID, p.Amount,
		p.PaymentDate, p.Status, p.ProofFileURL,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}
	return nil
}

func (r *paymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Payment, error) {
	query := `
		SELECT id, household_id, waste_id, amount, payment_date, status, proof_file_url, created_at, updated_at
		FROM payments
		WHERE id = $1
	`
	var p model.Payment
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.HouseholdID, &p.WasteID, &p.Amount,
		&p.PaymentDate, &p.Status, &p.ProofFileURL,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("find payment: %w", err)
	}
	return &p, nil
}

func (r *paymentRepository) FindAll(
	ctx context.Context,
	status string,
	householdID *uuid.UUID,
	from, to *string,
	limit, offset int,
) ([]model.Payment, int, error) {
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
	if from != nil {
		where += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, *from)
		argIdx++
	}
	if to != nil {
		where += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, *to)
		argIdx++
	}

	countQuery := "SELECT COUNT(*) FROM payments " + where
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count payments: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT id, household_id, waste_id, amount, payment_date, status, proof_file_url, created_at, updated_at
		FROM payments
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list payments: %w", err)
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(
			&p.ID, &p.HouseholdID, &p.WasteID, &p.Amount,
			&p.PaymentDate, &p.Status, &p.ProofFileURL,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan payment: %w", err)
		}
		payments = append(payments, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows err: %w", err)
	}

	return payments, total, nil
}

func (r *paymentRepository) Update(ctx context.Context, p *model.Payment) error {
	query := `
		UPDATE payments
		SET status = $2, payment_date = $3, proof_file_url = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		p.ID, p.Status, p.PaymentDate, p.ProofFileURL,
	).Scan(&p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrNotFound
		}
		return fmt.Errorf("update payment: %w", err)
	}
	return nil
}

func (r *paymentRepository) HasPendingByHousehold(ctx context.Context, householdID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM payments
			WHERE household_id = $1 AND status = 'pending'
		)
	`
	var exists bool
	if err := r.pool.QueryRow(ctx, query, householdID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check pending payment: %w", err)
	}
	return exists, nil
}

func (r *paymentRepository) ExistsByWasteID(ctx context.Context, wasteID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM payments WHERE waste_id = $1)`
	var exists bool
	if err := r.pool.QueryRow(ctx, query, wasteID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check payment by waste: %w", err)
	}
	return exists, nil
}

func (r *paymentRepository) CreateTx(ctx context.Context, tx pgx.Tx, p *model.Payment) error {
	query := `
		INSERT INTO payments (id, household_id, waste_id, amount, payment_date, status, proof_file_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	err := tx.QueryRow(ctx, query,
		p.ID, p.HouseholdID, p.WasteID, p.Amount,
		p.PaymentDate, p.Status, p.ProofFileURL,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert payment in tx: %w", err)
	}
	return nil
}

func (r *paymentRepository) ExistsByWasteIDTx(ctx context.Context, tx pgx.Tx, wasteID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM payments WHERE waste_id = $1)`
	var exists bool
	if err := tx.QueryRow(ctx, query, wasteID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check payment by waste in tx: %w", err)
	}
	return exists, nil
}