package response

import (
	"errors"
	"net/http"

	"github.com/aanjardev/waste-api/internal/model"
)

func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		Error(w, http.StatusNotFound, "resource not found", nil)
	case errors.Is(err, model.ErrHouseholdHasDebt):
		Error(w, http.StatusConflict, "household has pending payment", nil)
	case errors.Is(err, model.ErrInvalidStatus):
		Error(w, http.StatusConflict, "invalid status transition", nil)
	case errors.Is(err, model.ErrSafetyCheckNeeded):
		Error(w, http.StatusUnprocessableEntity, "electronic pickup requires safety check", nil)
	case errors.Is(err, model.ErrProofFileRequired):
		Error(w, http.StatusUnprocessableEntity, "proof file is required", nil)
	case errors.Is(err, model.ErrDuplicatePayment):
		Error(w, http.StatusConflict, "payment already exists for this pickup", nil)
	case errors.Is(err, model.ErrPickupMismatch):
		Error(w, http.StatusUnprocessableEntity, "pickup does not belong to household", nil)
	default:
		Error(w, http.StatusInternalServerError, "internal server error", nil)
	}
}