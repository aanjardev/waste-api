package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/aanjardev/waste-api/internal/model"
)

type PaymentRepository struct {
	mock.Mock
}

func (m *PaymentRepository) Create(ctx context.Context, p *model.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *PaymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *PaymentRepository) FindAll(ctx context.Context, status string, householdID *uuid.UUID, from, to *string, limit, offset int) ([]model.Payment, int, error) {
	args := m.Called(ctx, status, householdID, from, to, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]model.Payment), args.Int(1), args.Error(2)
}

func (m *PaymentRepository) Update(ctx context.Context, p *model.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *PaymentRepository) HasPendingByHousehold(ctx context.Context, householdID uuid.UUID) (bool, error) {
	args := m.Called(ctx, householdID)
	return args.Bool(0), args.Error(1)
}

func (m *PaymentRepository) ExistsByWasteID(ctx context.Context, wasteID uuid.UUID) (bool, error) {
	args := m.Called(ctx, wasteID)
	return args.Bool(0), args.Error(1)
}

func (m *PaymentRepository) CreateTx(ctx context.Context, tx pgx.Tx, p *model.Payment) error {
	args := m.Called(ctx, tx, p)
	return args.Error(0)
}

func (m *PaymentRepository) ExistsByWasteIDTx(ctx context.Context, tx pgx.Tx, wasteID uuid.UUID) (bool, error) {
	args := m.Called(ctx, tx, wasteID)
	return args.Bool(0), args.Error(1)
}