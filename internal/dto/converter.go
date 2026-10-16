package dto

import "github.com/aanjardev/waste-api/internal/model"

func ToHouseholdResponse(h *model.Household) HouseholdResponse {
	return HouseholdResponse{
		ID:        h.ID.String(),
		OwnerName: h.OwnerName,
		Address:   h.Address,
		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
	}
}

func ToHouseholdResponseList(items []model.Household) []HouseholdResponse {
	result := make([]HouseholdResponse, 0, len(items))
	for i := range items {
		result = append(result, ToHouseholdResponse(&items[i]))
	}
	return result
}
func ToPickupResponse(p *model.WastePickup) PickupResponse {
	return PickupResponse{
		ID:          p.ID.String(),
		HouseholdID: p.HouseholdID.String(),
		Type:        string(p.Type),
		Status:      string(p.Status),
		PickupDate:  p.PickupDate,
		SafetyCheck: p.SafetyCheck,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ToPickupResponseList(items []model.WastePickup) []PickupResponse {
	result := make([]PickupResponse, 0, len(items))
	for i := range items {
		result = append(result, ToPickupResponse(&items[i]))
	}
	return result
}