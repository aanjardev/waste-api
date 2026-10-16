package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
	"github.com/aanjardev/waste-api/internal/dto"
	"github.com/aanjardev/waste-api/internal/model"
)

type PaymentService struct {
	mock.Mock
}

func (m *PaymentService) Create(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PaymentResponse), args.Error(1)
}

func (m *PaymentService) FindByID(ctx context.Context, id uuid.UUID) (*dto.PaymentResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PaymentResponse), args.Error(1)
}

func (m *PaymentService) FindAll(ctx context.Context, q dto.PaymentListQuery, page, limit int) ([]dto.PaymentResponse, int, error) {
	args := m.Called(ctx, q, page, limit)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]dto.PaymentResponse), args.Int(1), args.Error(2)
}

func (m *PaymentService) Confirm(ctx context.Context, id uuid.UUID, proofURL string) (*dto.PaymentResponse, error) {
	args := m.Called(ctx, id, proofURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PaymentResponse), args.Error(1)
}

func (m *PaymentService) CreateForPickup(ctx context.Context, householdID, wasteID uuid.UUID, amount float64) (*model.Payment, error) {
	args := m.Called(ctx, householdID, wasteID, amount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *PaymentService) CreateForPickupTx(ctx context.Context, tx pgx.Tx, householdID, wasteID uuid.UUID, amount float64) (*model.Payment, error) {
	args := m.Called(ctx, tx, householdID, wasteID, amount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}