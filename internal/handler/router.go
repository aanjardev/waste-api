package handler

import (
	"net/http"

	"github.com/aanjardev/waste-api/internal/middleware"
	"github.com/aanjardev/waste-api/internal/service"
)

func NewRouter(householdSvc service.HouseholdService) http.Handler {
	mux := http.NewServeMux()

	householdHandler := NewHouseholdHandler(householdSvc)

	// Household routes
	mux.HandleFunc("POST /api/households", householdHandler.Create)
	mux.HandleFunc("GET /api/households", householdHandler.List)
	mux.HandleFunc("GET /api/households/{id}", householdHandler.Detail)
	mux.HandleFunc("DELETE /api/households/{id}", householdHandler.Delete)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Chain middleware
	return middleware.Chain(
		mux,
		middleware.Recover,
		middleware.Logging,
		middleware.RequireJSON,
	)
}