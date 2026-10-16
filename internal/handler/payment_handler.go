package handler

import (
	"net/http"

	"github.com/aanjardev/waste-api/internal/dto"
	"github.com/aanjardev/waste-api/internal/response"
	"github.com/aanjardev/waste-api/internal/service"
)

type PaymentHandler struct {
	svc       service.PaymentService
	uploadDir string
}

func NewPaymentHandler(svc service.PaymentService, uploadDir string) *PaymentHandler {
	return &PaymentHandler{svc: svc, uploadDir: uploadDir}
}

// POST /api/payments
func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePaymentRequest
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

	response.Success(w, http.StatusCreated, "payment created", res)
}

// GET /api/payments?status=&household_id=&from=&to=&page=&limit=
func (h *PaymentHandler) List(w http.ResponseWriter, r *http.Request) {
	q := dto.PaymentListQuery{
		Status:      r.URL.Query().Get("status"),
		HouseholdID: r.URL.Query().Get("household_id"),
		From:        r.URL.Query().Get("from"),
		To:          r.URL.Query().Get("to"),
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

	response.SuccessWithMeta(w, http.StatusOK, "payments retrieved", items, meta)
}

// PUT /api/payments/{id}/confirm
// Content-Type: multipart/form-data
// Form field: proof_file
func (h *PaymentHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUID(w, r, "id")
	if !ok {
		return
	}

	proofURL, err := saveUploadedFile(r, h.uploadDir)
	if err != nil {
		response.Error(w, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	res, err := h.svc.Confirm(r.Context(), id, proofURL)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "payment confirmed", res)
}