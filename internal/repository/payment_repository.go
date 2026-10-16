package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository interface {
	HasPendingByHousehold(ctx context.Context, householdID uuid.UUID) (bool, error)
	// method lain akan ditambahkan di sesi berikutnya
}

type paymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) PaymentRepository {
	return &paymentRepository{pool: pool}
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