package handler

import (
	"net/http"
	// "time"

	"github.com/aanjardev/waste-api/internal/dto"
	"github.com/aanjardev/waste-api/internal/response"
	"github.com/aanjardev/waste-api/internal/service"
)

type PickupHandler struct {
	svc service.PickupService
}

func NewPickupHandler(svc service.PickupService) *PickupHandler {
	return &PickupHandler{svc: svc}
}

// POST /api/pickups
func (h *PickupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePickupRequest
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

	response.Success(w, http.StatusCreated, "pickup created", res)
}

// GET /api/pickups?status=&household_id=&page=&limit=
func (h *PickupHandler) List(w http.ResponseWriter, r *http.Request) {
	q := dto.PickupListQuery{
		Status:      r.URL.Query().Get("status"),
		HouseholdID: r.URL.Query().Get("household_id"),
	}

	if errs := validateStruct(q); errs != nil {
		response.Error(w, http.StatusUnprocessableEntity, "validation error", errs)
		return
	}

	page := parseIntQuery(r, "page", 1, 1, 1_000_000)
	limit := parseIntQuery(r, "limit", 10, 1, 100)

	items, total, err := h.svc.FindAll(r.Context(), q, page, limit)
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

	response.SuccessWithMeta(w, http.StatusOK, "pickups retrieved", items, meta)
}

// PUT /api/pickups/{id}/schedule
func (h *PickupHandler) Schedule(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, r, "id")
	if !ok {
		return
	}

	var req dto.SchedulePickupRequest
	if err := decodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if errs := validateStruct(req); errs != nil {
		response.Error(w, http.StatusUnprocessableEntity, "validation error", errs)
		return
	}

	res, err := h.svc.Schedule(r.Context(), id, req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "pickup scheduled", res)
}

// PUT /api/pickups/{id}/complete
func (h *PickupHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, r, "id")
	if !ok {
		return
	}

	res, err := h.svc.Complete(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "pickup completed", res)
}

// PUT /api/pickups/{id}/cancel
func (h *PickupHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, r, "id")
	if !ok {
		return
	}

	res, err := h.svc.Cancel(r.Context(), id)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "pickup canceled", res)
}

