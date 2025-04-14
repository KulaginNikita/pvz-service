package api

import (
	"encoding/json"
	"net/http"
)

func (a *API) PostLogin(w http.ResponseWriter, r *http.Request) {
	var req PostLoginJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	token, err := a.userService.Login(r.Context(), string(req.Email), req.Password)
	if err != nil {
		http.Error(w, "authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Token(token))
}
