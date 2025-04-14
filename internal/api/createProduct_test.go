package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/KulaginNikita/pvz-service/internal/api/mocks"
	"github.com/KulaginNikita/pvz-service/internal/domain/product"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
)

func TestPostProducts(t *testing.T) {
	type fields struct {
		createProductErr error
		expectCall       bool
	}
	tests := []struct {
		name           string
		requestBody    PostProductsJSONBody
		fields         fields
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "успешное создание продукта",
			requestBody: PostProductsJSONBody{
				PvzId: uuid.New(),
				Type:  PostProductsJSONBodyTypeОдежда,
			},
			fields:         fields{createProductErr: nil, expectCall: true},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message": "product created"}`,
		},
		{
			name: "ошибка: пустой pvzId",
			requestBody: PostProductsJSONBody{
				PvzId: uuid.Nil,
				Type:  PostProductsJSONBodyTypeОдежда,
			},
			fields:         fields{expectCall: false},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing pvzId\n",
		},
		{
			name: "ошибка авторизации",
			requestBody: PostProductsJSONBody{
				PvzId: uuid.New(),
				Type:  PostProductsJSONBodyTypeЭлектроника,
			},
			fields:         fields{createProductErr: product.ErrUnauthorized, expectCall: true},
			expectedStatus: http.StatusForbidden,
			expectedBody:   product.ErrUnauthorized.Error() + "\n",
		},
		{
			name: "другая ошибка создания продукта",
			requestBody: PostProductsJSONBody{
				PvzId: uuid.New(),
				Type:  PostProductsJSONBodyTypeОбувь,
			},
			fields:         fields{createProductErr: errors.New("db error"), expectCall: true},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "db error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			productMock := mocks.NewProductServiceMock(mc)
			if tt.fields.expectCall {
				productMock.CreateProductMock.Return(tt.fields.createProductErr)
			}

			apiInstance := NewAPI(
				mocks.NewUserServiceMock(mc),
				mocks.NewPVZServiceMock(mc),
				mocks.NewReceptionServiceMock(mc),
				productMock,
				nil,
			)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
			w := httptest.NewRecorder()

			apiInstance.PostProducts(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}
