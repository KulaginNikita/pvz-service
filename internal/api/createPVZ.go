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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	pvz := domain.PVZ{City: domain.City(req.City)}

	created, err := a.pvzService.CreatePVZ(r.Context(), &pvz)
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

	resp := PVZ{
		Id:               &created.ID, 
		City:             PVZCity(created.City),
		RegistrationDate: &created.RegisteredAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
