package api

import (
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (a *API) PostPvzPvzIdDeleteLastProduct(w http.ResponseWriter, r *http.Request, pvzID openapi_types.UUID) {

	err := a.productService.DeleteProduct(r.Context(), uuid.UUID(pvzID))
	if err != nil {
		http.Error(w, "failed to delete product: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Product deleted successfully"}`))
}
