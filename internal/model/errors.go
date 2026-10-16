package model

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrHouseholdHasDebt  = errors.New("household has pending payment")
	ErrInvalidStatus     = errors.New("invalid status transition")
	ErrSafetyCheckNeeded = errors.New("electronic pickup requires safety check")
	ErrProofFileRequired = errors.New("proof file is required")
	ErrDuplicatePayment  = errors.New("payment already exists for this pickup")
)