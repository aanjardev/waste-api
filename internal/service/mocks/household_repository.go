package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/aanjardev/waste-api/internal/model"
)

type HouseholdRepository struct {
	mock.Mock
}

func (m *HouseholdRepository) Create(ctx context.Context, h *model.Household) error {
	args := m.Called(ctx, h)
	return args.Error(0)
}

func (m *HouseholdRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Household, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Household), args.Error(1)
}

func (m *HouseholdRepository) FindAll(ctx context.Context, limit, offset int) ([]model.Household, int, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]model.Household), args.Int(1), args.Error(2)
}

func (m *HouseholdRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *HouseholdRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}