package dto

import "time"

type CreatePaymentRequest struct {
	HouseholdID string  `json:"household_id" validate:"required,uuid"`
	WasteID     string  `json:"waste_id" validate:"required,uuid"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
}

type PaymentResponse struct {
	ID           string     `json:"id"`
	HouseholdID  string     `json:"household_id"`
	WasteID      string     `json:"waste_id"`
	Amount       float64    `json:"amount"`
	PaymentDate  *time.Time `json:"payment_date"`
	Status       string     `json:"status"`
	ProofFileURL *string    `json:"proof_file_url"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type PaymentListQuery struct {
	Status      string `json:"status" validate:"omitempty,oneof=pending paid failed"`
	HouseholdID string `json:"household_id" validate:"omitempty,uuid"`
	From        string `json:"from" validate:"omitempty,datetime=2006-01-02"`
	To          string `json:"to" validate:"omitempty,datetime=2006-01-02"`
}