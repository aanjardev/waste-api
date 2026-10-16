package model

import (
	"time"

	"github.com/google/uuid"
)

type WasteType string

const (
	WasteOrganic    WasteType = "organic"
	WastePlastic    WasteType = "plastic"
	WastePaper      WasteType = "paper"
	WasteElectronic WasteType = "electronic"
)

type PickupStatus string

const (
	PickupPending   PickupStatus = "pending"
	PickupScheduled PickupStatus = "scheduled"
	PickupCompleted PickupStatus = "completed"
	PickupCanceled  PickupStatus = "canceled"
)

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "pending"
	PaymentPaid    PaymentStatus = "paid"
	PaymentFailed  PaymentStatus = "failed"
)

type Household struct {
	ID        uuid.UUID `json:"id"`
	OwnerName string    `json:"owner_name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WastePickup struct {
	ID          uuid.UUID    `json:"id"`
	HouseholdID uuid.UUID    `json:"household_id"`
	Type        WasteType    `json:"type"`
	Status      PickupStatus `json:"status"`
	PickupDate  *time.Time   `json:"pickup_date"`
	SafetyCheck bool         `json:"safety_check"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Payment struct {
	ID           uuid.UUID     `json:"id"`
	HouseholdID  uuid.UUID     `json:"household_id"`
	WasteID      uuid.UUID     `json:"waste_id"`
	Amount       float64       `json:"amount"`
	PaymentDate  *time.Time    `json:"payment_date"`
	Status       PaymentStatus `json:"status"`
	ProofFileURL *string       `json:"proof_file_url"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}