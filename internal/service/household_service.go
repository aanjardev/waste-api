package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/aanjardev/waste-api/internal/dto"
	"github.com/aanjardev/waste-api/internal/model"
	"github.com/aanjardev/waste-api/internal/repository"
)

type HouseholdService interface {
	Create(ctx context.Context, req dto.CreateHouseholdRequest) (*dto.HouseholdResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*dto.HouseholdResponse, error)
	FindAll(ctx context.Context, page, limit int) ([]dto.HouseholdResponse, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type householdService struct {
	repo repository.HouseholdRepository
}

func NewHouseholdService(repo repository.HouseholdRepository) HouseholdService {
	return &householdService{repo: repo}
}

func (s *householdService) Create(ctx context.Context, req dto.CreateHouseholdRequest) (*dto.HouseholdResponse, error) {
	h := &model.Household{
		ID:        uuid.New(),
		OwnerName: req.OwnerName,
		Address:   req.Address,
	}

	if err := s.repo.Create(ctx, h); err != nil {
		return nil, fmt.Errorf("create household: %w", err)
	}

	resp := dto.ToHouseholdResponse(h)
	return &resp, nil
}

func (s *householdService) FindByID(ctx context.Context, id uuid.UUID) (*dto.HouseholdResponse, error) {
	h, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToHouseholdResponse(h)
	return &resp, nil
}

func (s *householdService) FindAll(ctx context.Context, page, limit int) ([]dto.HouseholdResponse, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	items, total, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list households: %w", err)
	}

	return dto.ToHouseholdResponseList(items), total, nil
}

func (s *householdService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}