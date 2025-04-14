package api

import (
	"encoding/json"
	"net/http"

	"github.com/KulaginNikita/pvz-service/internal/domain/pvz"
)

func (a *API) GetPvz(w http.ResponseWriter, r *http.Request, params GetPvzParams) {
	// Валидация даты
	if params.StartDate == nil || params.EndDate == nil {
		http.Error(w, "startDate and endDate are required", http.StatusBadRequest)
		return
	}

	limit := int64(10)
	if params.Limit != nil && *params.Limit > 0 {
		limit = int64(*params.Limit)
	}

	page := 1
	if params.Page != nil && *params.Page > 0 {
		page = *params.Page
	}

	offset := int64((page - 1)) * limit

	// Сформировать фильтр
	filter := &pvz.PVZFilter{
		StartDate: *params.StartDate,
		EndDate:   *params.EndDate,
		Limit:     limit,
		Offset:    offset,
	}


	result, err := a.pvzService.GetPVZ(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
