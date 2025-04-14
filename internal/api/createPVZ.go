package api

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/KulaginNikita/pvz-service/internal/domain/pvz"
	pvzservice "github.com/KulaginNikita/pvz-service/internal/service/pvz"
)

func (a *API) PostPvz(w http.ResponseWriter, r *http.Request) {
	var req PVZ

	// Декодируем JSON тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	// Вызываем сервисный слой, используя ctx с ролью из middleware
	pvz := domain.PVZ{City: domain.City(req.City)}

	err := a.pvzService.CreatePVZ(r.Context(), &pvz)

	if err != nil {
		switch {
		case errors.Is(err, pvzservice.ErrForbiddenCity):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, pvzservice.ErrUnauthorized):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}
