package dto

import "time"

type CreateHouseholdRequest struct {
	OwnerName string `json:"owner_name" validate:"required,min=2,max=255"`
	Address   string `json:"address" validate:"required,min=5"`
}

type HouseholdResponse struct {
	ID        string    `json:"id"`
	OwnerName string    `json:"owner_name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginationQuery struct {
	Page  int `json:"page" validate:"omitempty,min=1"`
	Limit int `json:"limit" validate:"omitempty,min=1,max=100"`
}