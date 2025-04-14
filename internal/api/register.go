package api

import (
	"encoding/json"
	"net/http"

	"github.com/KulaginNikita/pvz-service/internal/domain/user"
)

func (a *API) PostRegister(w http.ResponseWriter, r *http.Request) {
	var req PostRegisterJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	newUser := &user.User{
		Email:    string(req.Email),
		Password: req.Password,
		Role:     user.Role(req.Role),
	}

	if err := a.userService.Register(r.Context(), newUser); err != nil {
		http.Error(w, "registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
