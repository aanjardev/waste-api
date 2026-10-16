package middleware

import (
	"net/http"
	"strings"

	"github.com/aanjardev/waste-api/internal/response"
)

func RequireJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			ct := r.Header.Get("Content-Type")
			// Biarkan multipart untuk upload
			if !strings.HasPrefix(ct, "application/json") && !strings.HasPrefix(ct, "multipart/form-data") {
				response.Error(w, http.StatusUnsupportedMediaType, "content-type must be application/json or multipart/form-data", nil)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}