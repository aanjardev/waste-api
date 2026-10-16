package handler

import (
	"net/http"

	"github.com/aanjardev/waste-api/internal/middleware"
	"github.com/aanjardev/waste-api/internal/service"
)

func NewRouter(
	householdSvc service.HouseholdService,
	pickupSvc service.PickupService,
	paymentSvc service.PaymentService,
	reportSvc service.ReportService,
	uploadDir string,
) http.Handler {
	mux := http.NewServeMux()

	householdHandler := NewHouseholdHandler(householdSvc)
	pickupHandler := NewPickupHandler(pickupSvc)
	paymentHandler := NewPaymentHandler(paymentSvc, uploadDir)
	reportHandler := NewReportHandler(reportSvc)

	// Household routes
	mux.HandleFunc("POST /api/households", householdHandler.Create)
	mux.HandleFunc("GET /api/households", householdHandler.List)
	mux.HandleFunc("GET /api/households/{id}", householdHandler.Detail)
	mux.HandleFunc("DELETE /api/households/{id}", householdHandler.Delete)

	// Pickup routes
	mux.HandleFunc("POST /api/pickups", pickupHandler.Create)
	mux.HandleFunc("GET /api/pickups", pickupHandler.List)
	mux.HandleFunc("PUT /api/pickups/{id}/schedule", pickupHandler.Schedule)
	mux.HandleFunc("PUT /api/pickups/{id}/complete", pickupHandler.Complete)
	mux.HandleFunc("PUT /api/pickups/{id}/cancel", pickupHandler.Cancel)

	// Payment routes
	mux.HandleFunc("POST /api/payments", paymentHandler.Create)
	mux.HandleFunc("GET /api/payments", paymentHandler.List)
	mux.HandleFunc("PUT /api/payments/{id}/confirm", paymentHandler.Confirm)

	// Report routes
	mux.HandleFunc("GET /api/reports/waste-summary", reportHandler.WasteSummary)
	mux.HandleFunc("GET /api/reports/payment-summary", reportHandler.PaymentSummary)

	// Serve uploaded files
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return middleware.Chain(
		mux,
		middleware.Recover,
		middleware.Logging,
		middleware.RequireJSON,
	)
}