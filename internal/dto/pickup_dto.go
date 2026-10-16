package dto

import "time"

type CreatePickupRequest struct {
	HouseholdID string `json:"household_id" validate:"required,uuid"`
	Type        string `json:"type" validate:"required,oneof=organic plastic paper electronic"`
	SafetyCheck bool   `json:"safety_check"`
}

type SchedulePickupRequest struct {
	PickupDate time.Time `json:"pickup_date" validate:"required"`
}

type PickupResponse struct {
	ID          string     `json:"id"`
	HouseholdID string     `json:"household_id"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	PickupDate  *time.Time `json:"pickup_date"`
	SafetyCheck bool       `json:"safety_check"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PickupListQuery struct {
	Status      string `json:"status" validate:"omitempty,oneof=pending scheduled completed canceled"`
	HouseholdID string `json:"household_id" validate:"omitempty,uuid"`
}