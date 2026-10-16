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

type PaymentService interface {
	Create(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error)
	FindByID(ctx context.Context, id uuid.UUID) (*dto.PaymentResponse, error)
	FindAll(ctx context.Context, q dto.PaymentListQuery, page, limit int) ([]dto.PaymentResponse, int, error)
	Confirm(ctx context.Context, id uuid.UUID, proofURL string) (*dto.PaymentResponse, error)
	CreateForPickup(ctx context.Context, householdID, wasteID uuid.UUID, amount float64) (*model.Payment, error)
}

type paymentService struct {
	repo          repository.PaymentRepository
	householdRepo repository.HouseholdRepository
	pickupRepo    repository.PickupRepository
}

func NewPaymentService(
	repo repository.PaymentRepository,
	householdRepo repository.HouseholdRepository,
	pickupRepo repository.PickupRepository,
) PaymentService {
	return &paymentService{
		repo:          repo,
		householdRepo: householdRepo,
		pickupRepo:    pickupRepo,
	}
}

func (s *paymentService) Create(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	householdID, err := uuid.Parse(req.HouseholdID)
	if err != nil {
		return nil, fmt.Errorf("parse household id: %w", err)
	}
	wasteID, err := uuid.Parse(req.WasteID)
	if err != nil {
		return nil, fmt.Errorf("parse waste id: %w", err)
	}

	exists, err := s.householdRepo.Exists(ctx, householdID)
	if err != nil {
		return nil, fmt.Errorf("check household: %w", err)
	}
	if !exists {
		return nil, model.ErrNotFound
	}

	pickup, err := s.pickupRepo.FindByID(ctx, wasteID)
	if err != nil {
		return nil, err
	}
	if pickup.HouseholdID != householdID {
		return nil, fmt.Errorf("pickup does not belong to household")
	}

	exists, err = s.repo.ExistsByWasteID(ctx, wasteID)
	if err != nil {
		return nil, fmt.Errorf("check payment exists: %w", err)
	}
	if exists {
		return nil, model.ErrDuplicatePayment
	}

	p := &model.Payment{
		ID:          uuid.New(),
		HouseholdID: householdID,
		WasteID:     wasteID,
		Amount:      req.Amount,
		Status:      model.PaymentPending,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	resp := dto.ToPaymentResponse(p)
	return &resp, nil
}

func (s *paymentService) FindByID(ctx context.Context, id uuid.UUID) (*dto.PaymentResponse, error) {
	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := dto.ToPaymentResponse(p)
	return &resp, nil
}

func (s *paymentService) FindAll(ctx context.Context, q dto.PaymentListQuery, page, limit int) ([]dto.PaymentResponse, int, error) {
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

	var fromPtr, toPtr *string
	if q.From != "" {
		fromPtr = &q.From
	}
	if q.To != "" {
		toPtr = &q.To
	}

	items, total, err := s.repo.FindAll(ctx, q.Status, householdIDPtr, fromPtr, toPtr, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list payments: %w", err)
	}

	return dto.ToPaymentResponseList(items), total, nil
}

func (s *paymentService) Confirm(ctx context.Context, id uuid.UUID, proofURL string) (*dto.PaymentResponse, error) {
	if proofURL == "" {
		return nil, model.ErrProofFileRequired
	}

	p, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if p.Status != model.PaymentPending {
		return nil, model.ErrInvalidStatus
	}

	now := time.Now()
	p.Status = model.PaymentPaid
	p.PaymentDate = &now
	p.ProofFileURL = &proofURL

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, fmt.Errorf("update payment: %w", err)
	}

	resp := dto.ToPaymentResponse(p)
	return &resp, nil
}

func (s *paymentService) CreateForPickup(ctx context.Context, householdID, wasteID uuid.UUID, amount float64) (*model.Payment, error) {
	exists, err := s.repo.ExistsByWasteID(ctx, wasteID)
	if err != nil {
		return nil, fmt.Errorf("check payment exists: %w", err)
	}
	if exists {
		return nil, model.ErrDuplicatePayment
	}

	p := &model.Payment{
		ID:          uuid.New(),
		HouseholdID: householdID,
		WasteID:     wasteID,
		Amount:      amount,
		Status:      model.PaymentPending,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	return p, nil
}