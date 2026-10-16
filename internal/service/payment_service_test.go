package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/aanjardev/waste-api/internal/model"
	"github.com/aanjardev/waste-api/internal/service"
	"github.com/aanjardev/waste-api/internal/service/mocks"
)

func TestPaymentService_Confirm_RejectsEmptyProofURL(t *testing.T) {
	paymentRepo := new(mocks.PaymentRepository)
	householdRepo := new(mocks.HouseholdRepository)
	pickupRepo := new(mocks.PickupRepository)

	svc := service.NewPaymentService(paymentRepo, householdRepo, pickupRepo)

	_, err := svc.Confirm(context.Background(), uuid.New(), "")

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrProofFileRequired)
	paymentRepo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestPaymentService_Confirm_RejectsNonPendingStatus(t *testing.T) {
	paymentRepo := new(mocks.PaymentRepository)
	householdRepo := new(mocks.HouseholdRepository)
	pickupRepo := new(mocks.PickupRepository)

	id := uuid.New()
	paymentRepo.On("FindByID", mock.Anything, id).Return(&model.Payment{
		ID:     id,
		Status: model.PaymentPaid,
	}, nil)

	svc := service.NewPaymentService(paymentRepo, householdRepo, pickupRepo)

	_, err := svc.Confirm(context.Background(), id, "/uploads/test.png")

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrInvalidStatus)
}

func TestPaymentService_Confirm_Success(t *testing.T) {
	paymentRepo := new(mocks.PaymentRepository)
	householdRepo := new(mocks.HouseholdRepository)
	pickupRepo := new(mocks.PickupRepository)

	id := uuid.New()
	paymentRepo.On("FindByID", mock.Anything, id).Return(&model.Payment{
		ID:     id,
		Status: model.PaymentPending,
	}, nil)
	paymentRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Payment")).Return(nil)

	svc := service.NewPaymentService(paymentRepo, householdRepo, pickupRepo)

	res, err := svc.Confirm(context.Background(), id, "/uploads/test.png")

	require.NoError(t, err)
	assert.Equal(t, "paid", res.Status)
	assert.NotNil(t, res.PaymentDate)
	assert.NotNil(t, res.ProofFileURL)
	assert.Equal(t, "/uploads/test.png", *res.ProofFileURL)
	paymentRepo.AssertExpectations(t)
}

func TestPaymentService_CreateForPickup_RejectsDuplicate(t *testing.T) {
	paymentRepo := new(mocks.PaymentRepository)
	householdRepo := new(mocks.HouseholdRepository)
	pickupRepo := new(mocks.PickupRepository)

	householdID := uuid.New()
	wasteID := uuid.New()

	paymentRepo.On("ExistsByWasteID", mock.Anything, wasteID).Return(true, nil)

	svc := service.NewPaymentService(paymentRepo, householdRepo, pickupRepo)

	_, err := svc.CreateForPickup(context.Background(), householdID, wasteID, 50000)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrDuplicatePayment)
}