package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/aanjardev/waste-api/internal/dto"
	"github.com/aanjardev/waste-api/internal/model"
	"github.com/aanjardev/waste-api/internal/service"
	"github.com/aanjardev/waste-api/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPickupService_Create_RejectsWhenHouseholdHasPendingPayment(t *testing.T) {
	pickupRepo := new(mocks.PickupRepository)
	paymentRepo := new(mocks.PaymentRepository)
	householdRepo := new(mocks.HouseholdRepository)

	householdID := uuid.New()
	householdRepo.On("Exists", mock.Anything, householdID).Return(true, nil)
	paymentRepo.On("HasPendingByHousehold", mock.Anything, householdID).Return(true, nil)

	svc := service.NewPickupService(pickupRepo, paymentRepo, householdRepo, nil, nil)

	_, err := svc.Create(context.Background(), dto.CreatePickupRequest{
		HouseholdID: householdID.String(),
		Type:        "organic",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrHouseholdHasDebt)
	pickupRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPickupService_Create_RejectsWhenHouseholdNotFound(t *testing.T) {
	pickupRepo := new(mocks.PickupRepository)
	paymentRepo := new(mocks.PaymentRepository)
	householdRepo := new(mocks.HouseholdRepository)

	householdID := uuid.New()
	householdRepo.On("Exists", mock.Anything, householdID).Return(false, nil)

	svc := service.NewPickupService(pickupRepo, paymentRepo, householdRepo, nil, nil)

	_, err := svc.Create(context.Background(), dto.CreatePickupRequest{
		HouseholdID: householdID.String(),
		Type:        "organic",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrNotFound)
}

func TestPickupService_Create_Success(t *testing.T) {
	pickupRepo := new(mocks.PickupRepository)
	paymentRepo := new(mocks.PaymentRepository)
	householdRepo := new(mocks.HouseholdRepository)

	householdID := uuid.New()
	householdRepo.On("Exists", mock.Anything, householdID).Return(true, nil)
	paymentRepo.On("HasPendingByHousehold", mock.Anything, householdID).Return(false, nil)
	pickupRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.WastePickup")).Return(nil)

	svc := service.NewPickupService(pickupRepo, paymentRepo, householdRepo, nil, nil)

	res, err := svc.Create(context.Background(), dto.CreatePickupRequest{
		HouseholdID: householdID.String(),
		Type:        "organic",
	})

	require.NoError(t, err)
	assert.Equal(t, "organic", res.Type)
	assert.Equal(t, "pending", res.Status)
	pickupRepo.AssertExpectations(t)
}

func TestPickupService_Schedule_RejectsWhenStatusNotPending(t *testing.T) {
	pickupRepo := new(mocks.PickupRepository)
	paymentRepo := new(mocks.PaymentRepository)
	householdRepo := new(mocks.HouseholdRepository)

	id := uuid.New()
	pickupRepo.On("FindByID", mock.Anything, id).Return(&model.WastePickup{
		ID:     id,
		Type:   model.WasteOrganic,
		Status: model.PickupCompleted,
	}, nil)

	svc := service.NewPickupService(pickupRepo, paymentRepo, householdRepo, nil, nil)

	_, err := svc.Schedule(context.Background(), id, dto.SchedulePickupRequest{
		PickupDate: time.Now().Add(24 * time.Hour),
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrInvalidStatus)
}

func TestPickupService_Schedule_RejectsElectronicWithoutSafetyCheck(t *testing.T) {
	pickupRepo := new(mocks.PickupRepository)
	paymentRepo := new(mocks.PaymentRepository)
	householdRepo := new(mocks.HouseholdRepository)

	id := uuid.New()
	pickupRepo.On("FindByID", mock.Anything, id).Return(&model.WastePickup{
		ID:          id,
		Type:        model.WasteElectronic,
		Status:      model.PickupPending,
		SafetyCheck: false,
	}, nil)

	svc := service.NewPickupService(pickupRepo, paymentRepo, householdRepo, nil, nil)

	_, err := svc.Schedule(context.Background(), id, dto.SchedulePickupRequest{
		PickupDate: time.Now().Add(24 * time.Hour),
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrSafetyCheckNeeded)
}

func TestPickupService_Schedule_Success(t *testing.T) {
	pickupRepo := new(mocks.PickupRepository)
	paymentRepo := new(mocks.PaymentRepository)
	householdRepo := new(mocks.HouseholdRepository)

	id := uuid.New()
	pickupRepo.On("FindByID", mock.Anything, id).Return(&model.WastePickup{
		ID:          id,
		Type:        model.WasteElectronic,
		Status:      model.PickupPending,
		SafetyCheck: true,
	}, nil)
	pickupRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.WastePickup")).Return(nil)

	svc := service.NewPickupService(pickupRepo, paymentRepo, householdRepo, nil, nil)

	date := time.Now().Add(24 * time.Hour)
	res, err := svc.Schedule(context.Background(), id, dto.SchedulePickupRequest{
		PickupDate: date,
	})

	require.NoError(t, err)
	assert.Equal(t, "scheduled", res.Status)
	assert.NotNil(t, res.PickupDate)
	pickupRepo.AssertExpectations(t)
}

func TestPickupService_Complete_CreatesPaymentWithCorrectAmount(t *testing.T) {
	pickupRepo := new(mocks.PickupRepository)
	paymentRepo := new(mocks.PaymentRepository)
	householdRepo := new(mocks.HouseholdRepository)
	paymentSvc := new(mocks.PaymentService)
	txBeginner := new(mocks.TxBeginner)
	mockTx := new(mocks.Tx)

	id := uuid.New()
	householdID := uuid.New()

	pickupRepo.On("FindByID", mock.Anything, id).Return(&model.WastePickup{
		ID:          id,
		HouseholdID: householdID,
		Type:        model.WasteElectronic,
		Status:      model.PickupScheduled,
	}, nil)

	txBeginner.On("Begin", mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)
	mockTx.On("Commit", mock.Anything).Return(nil)

	pickupRepo.On("UpdateStatusIfTx", mock.Anything, mockTx, id, model.PickupScheduled, model.PickupCompleted, (*time.Time)(nil)).
		Return(&model.WastePickup{
			ID:          id,
			HouseholdID: householdID,
			Type:        model.WasteElectronic,
			Status:      model.PickupCompleted,
		}, nil)

	paymentSvc.On("CreateForPickupTx", mock.Anything, mockTx, householdID, id, float64(100000)).
		Return(&model.Payment{}, nil)

	svc := service.NewPickupService(pickupRepo, paymentRepo, householdRepo, paymentSvc, txBeginner)

	res, err := svc.Complete(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, "completed", res.Status)
	paymentSvc.AssertExpectations(t)
	mockTx.AssertExpectations(t)
	txBeginner.AssertExpectations(t)
}

func TestPickupService_Complete_RejectsWhenAlreadyCompletedByOther(t *testing.T) {
	pickupRepo := new(mocks.PickupRepository)
	paymentRepo := new(mocks.PaymentRepository)
	householdRepo := new(mocks.HouseholdRepository)
	paymentSvc := new(mocks.PaymentService)
	txBeginner := new(mocks.TxBeginner)
	mockTx := new(mocks.Tx)

	id := uuid.New()
	householdID := uuid.New()

	pickupRepo.On("FindByID", mock.Anything, id).Return(&model.WastePickup{
		ID:          id,
		HouseholdID: householdID,
		Type:        model.WasteOrganic,
		Status:      model.PickupScheduled,
	}, nil)

	txBeginner.On("Begin", mock.Anything).Return(mockTx, nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)

	pickupRepo.On("UpdateStatusIfTx", mock.Anything, mockTx, id, model.PickupScheduled, model.PickupCompleted, (*time.Time)(nil)).
		Return(nil, nil)

	svc := service.NewPickupService(pickupRepo, paymentRepo, householdRepo, paymentSvc, txBeginner)

	_, err := svc.Complete(context.Background(), id)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrInvalidStatus)
	paymentSvc.AssertNotCalled(t, "CreateForPickupTx", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockTx.AssertNotCalled(t, "Commit", mock.Anything)
}

