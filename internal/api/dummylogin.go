package api

import (
	"encoding/json"
	"net/http"
	"time"

)

func (a *API) PostDummyLogin(w http.ResponseWriter, r *http.Request) {
	var req PostDummyLoginJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if req.Role != PostDummyLoginJSONBodyRoleModerator && req.Role != PostDummyLoginJSONBodyRoleEmployee {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	token, err := a.jwtManager.GenerateToken(string(req.Role), time.Hour)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Token(token))
}
