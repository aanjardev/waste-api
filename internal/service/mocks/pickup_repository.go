package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/aanjardev/waste-api/internal/model"
)

type PickupRepository struct {
	mock.Mock
}

func (m *PickupRepository) Create(ctx context.Context, p *model.WastePickup) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *PickupRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.WastePickup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.WastePickup), args.Error(1)
}

func (m *PickupRepository) FindAll(ctx context.Context, status string, householdID *uuid.UUID, limit, offset int) ([]model.WastePickup, int, error) {
	args := m.Called(ctx, status, householdID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]model.WastePickup), args.Int(1), args.Error(2)
}

func (m *PickupRepository) Update(ctx context.Context, p *model.WastePickup) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}