package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/KulaginNikita/pvz-service/internal/domain/product"
)

func (a *API) PostProducts(w http.ResponseWriter, r *http.Request) {
	var req PostProductsJSONBody

	// 1. Парсим тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Проверяем, что pvzId не нулевой
	if req.PvzId == uuid.Nil {
		http.Error(w, "missing pvzId", http.StatusBadRequest)
		return
	}

	// 3. Формируем доменную модель
	p := &product.Product{
		ReceptionID: req.PvzId,
		Type:        product.ProductType(req.Type),
	}

	// 4. Вызываем сервис
	if err := a.productService.CreateProduct(r.Context(), p); err != nil {
		if errors.Is(err, product.ErrUnauthorized) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 5. Ответ
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "product created"}`))
}
