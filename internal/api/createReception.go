package api

import (
	"encoding/json"
	"net/http"
	"github.com/KulaginNikita/pvz-service/internal/domain/reception"
)



func (a *API) PostReceptions(w http.ResponseWriter, r *http.Request) {
	// Структура для парсинга JSON тела запроса
	var req PostReceptionsJSONBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	rec := &reception.Reception{PVZID: req.PvzId}

	// Вызов функции из сервисного слоя
	err := a.receptionService.CreateReception(r.Context(), rec)
	if err != nil {
		http.Error(w, "failed to create reception: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Устанавливаем заголовки и код ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Возвращаем успешный ответ
	w.Write([]byte(`{"message": "Reception created successfully"}`))
}
