package handler

import (
	"net/http"

	"github.com/aanjardev/waste-api/internal/dto"
	"github.com/aanjardev/waste-api/internal/response"
	"github.com/aanjardev/waste-api/internal/service"
)

type HouseholdHandler struct {
	svc service.HouseholdService
}

func NewHouseholdHandler(svc service.HouseholdService) *HouseholdHandler {
	return &HouseholdHandler{svc: svc}
}

// POST /api/households
func (h *HouseholdHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateHouseholdRequest
	if err := decodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if errs := validateStruct(req); errs != nil {
		response.Error(w, http.StatusUnprocessableEntity, "validation error", errs)
		return
	}

	res, err := h.svc.Create(r.Context(), req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "household created", res)
}

// GET /api/households?page=&limit=
func (h *HouseholdHandler) List(w http.ResponseWriter, r *http.Request) {
	page := parseIntQuery(r, "page", 1, 1, 1_000_000)
	limit := parseIntQuery(r, "limit", 10, 1, 100)

	items, total, err := h.svc.FindAll(r.Context(), page, limit)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	totalPages := (total + limit - 1) / limit
	meta := &response.Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	response.SuccessWithMeta(w, http.StatusOK, "households retrieved", items, meta)
}

// GET /api/households/{id}
func (h *HouseholdHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, r, "id")
	if !ok {
		return
	}

	res, err := h.svc.FindByID(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "household retrieved", res)
}

// DELETE /api/households/{id}
func (h *HouseholdHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, r, "id")
	if !ok {
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "household deleted", nil)
}