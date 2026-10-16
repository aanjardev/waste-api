package handler

import (
	"net/http"

	"github.com/aanjardev/waste-api/internal/response"
	"github.com/aanjardev/waste-api/internal/service"
)

type ReportHandler struct {
	svc service.ReportService
}

func NewReportHandler(svc service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// GET /api/reports/waste-summary
func (h *ReportHandler) WasteSummary(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.WasteSummary(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "waste summary retrieved", items)
}

// GET /api/reports/payment-summary
func (h *ReportHandler) PaymentSummary(w http.ResponseWriter, r *http.Request) {
	resp, err := h.svc.PaymentSummary(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "payment summary retrieved", resp)
}