package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/aanjardev/waste-api/internal/dto"
	"github.com/aanjardev/waste-api/internal/model"
	"github.com/aanjardev/waste-api/internal/repository"
)

type PickupService interface {
	Create(ctx context.Context, req dto.CreatePickupRequest) (*dto.PickupResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*dto.PickupResponse, error)
	FindAll(ctx context.Context, q dto.PickupListQuery, page, limit int) ([]dto.PickupResponse, int, error)
	Schedule(ctx context.Context, id uuid.UUID, req dto.SchedulePickupRequest) (*dto.PickupResponse, error)
	Complete(ctx context.Context, id uuid.UUID) (*dto.PickupResponse, error)
	Cancel(ctx context.Context, id uuid.UUID) (*dto.PickupResponse, error)
}

type pickupService struct {
	repo          repository.PickupRepository
	paymentRepo   repository.PaymentRepository
	householdRepo repository.HouseholdRepository
	paymentSvc    PaymentService
}

func NewPickupService(
	repo repository.PickupRepository,
	paymentRepo repository.PaymentRepository,
	householdRepo repository.HouseholdRepository,
	paymentSvc PaymentService,
) PickupService {
	return &pickupService{
		repo:          repo,
		paymentRepo:   paymentRepo,
		householdRepo: householdRepo,
		paymentSvc:    paymentSvc,
	}
}

func (s *pickupService) Create(ctx context.Context, req dto.CreatePickupRequest) (*dto.PickupResponse, error) {
	householdID, err := uuid.Parse(req.HouseholdID)
	if err != nil {
		return nil, fmt.Errorf("parse household id: %w", err)
	}

	hasPending, err := s.paymentRepo.HasPendingByHousehold(ctx, householdID)
	if err != nil {
		return nil, fmt.Errorf("check pending payment: %w", err)
	}
	if hasPending {
		return nil, model.ErrHouseholdHasDebt
	}

	p := &model.WastePickup{
		ID:          uuid.New(),
		HouseholdID: householdID,
		Type:        model.WasteType(req.Type),
		Status:      model.PickupPending,
		SafetyCheck: req.SafetyCheck,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create pickup: %w", err)
	}

	resp := dto.ToPickupResponse(p)
	return &resp, nil
}

func (s *pickupService) FindByID(ctx context.Context, id uuid.UUID) (*dto.PickupResponse, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToPickupResponse(p)
	return &resp, nil
}

func (s *pickupService) FindAll(ctx context.Context, q dto.PickupListQuery, page, limit int) ([]dto.PickupResponse, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	var householdIDPtr *uuid.UUID
	if q.HouseholdID != "" {
		id, err := uuid.Parse(q.HouseholdID)
		if err != nil {
			return nil, 0, fmt.Errorf("parse household id: %w", err)
		}
		householdIDPtr = &id
	}

	items, total, err := s.repo.FindAll(ctx, q.Status, householdIDPtr, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list pickups: %w", err)
	}

	return dto.ToPickupResponseList(items), total, nil
}

func (s *pickupService) Schedule(ctx context.Context, id uuid.UUID, req dto.SchedulePickupRequest) (*dto.PickupResponse, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p.Status != model.PickupPending {
		return nil, model.ErrInvalidStatus
	}

	if p.Type == model.WasteElectronic && !p.SafetyCheck {
		return nil, model.ErrSafetyCheckNeeded
	}

	p.Status = model.PickupScheduled
	p.PickupDate = &req.PickupDate

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update pickup: %w", err)
	}

	resp := dto.ToPickupResponse(p)
	return &resp, nil
}

func (s *pickupService) Complete(ctx context.Context, id uuid.UUID) (*dto.PickupResponse, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p.Status != model.PickupScheduled {
		return nil, model.ErrInvalidStatus
	}

	p.Status = model.PickupCompleted
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update pickup: %w", err)
	}

	amount := pickupAmount(p.Type)
	if _, err := s.paymentSvc.CreateForPickup(ctx, p.HouseholdID, p.ID, amount); err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	resp := dto.ToPickupResponse(p)
	return &resp, nil
}

func pickupAmount(t model.WasteType) float64 {
	switch t {
	case model.WasteElectronic:
		return 100000
	default:
		return 50000
	}
}

func (s *pickupService) Cancel(ctx context.Context, id uuid.UUID) (*dto.PickupResponse, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p.Status != model.PickupPending && p.Status != model.PickupScheduled {
		return nil, model.ErrInvalidStatus
	}

	p.Status = model.PickupCanceled
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update pickup: %w", err)
	}

	resp := dto.ToPickupResponse(p)
	return &resp, nil
}

var _ = time.Now