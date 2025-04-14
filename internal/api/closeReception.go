package api

import (
	"net/http"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (a *API) PostPvzPvzIdCloseLastReception(w http.ResponseWriter, r *http.Request, pvzId openapi_types.UUID) {

	err := a.receptionService.CloseReception(r.Context(), uuid.UUID(pvzId))
	if err != nil {
		http.Error(w, "failed to close reception: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Reception closed successfully"}`))
}