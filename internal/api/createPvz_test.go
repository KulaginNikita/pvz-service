package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mocks "github.com/KulaginNikita/pvz-service/internal/api/mocks"
	"github.com/KulaginNikita/pvz-service/internal/domain/pvz"
	pvzservice "github.com/KulaginNikita/pvz-service/internal/service/pvz"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
)

func TestPostPvz(t *testing.T) {
	now := time.Now()
	id := uuid.New()

	type fields struct {
		result *pvz.PVZ
		err    error
		expect bool
	}
	tests := []struct {
		name           string
		requestBody    string
		fields         fields
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "успешное создание ПВЗ",
			requestBody: `{"city":"Казань"}`,
			fields: fields{
				result: &pvz.PVZ{
					ID:           id,
					City:         "Казань",
					RegisteredAt: now,
				},
				expect: true,
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "", 
		},
		{
			name:           "невалидный JSON",
			requestBody:    `{"city":`,
			fields:         fields{expect: false},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid JSON body\n",
		},
		{
			name:        "запрещённый город",
			requestBody: `{"city":"Запретный"}`,
			fields: fields{
				err:    pvzservice.ErrForbiddenCity,
				expect: true,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   pvzservice.ErrForbiddenCity.Error() + "\n",
		},
		{
			name:        "неавторизованный пользователь",
			requestBody: `{"city":"Москва"}`,
			fields: fields{
				err:    pvzservice.ErrUnauthorized,
				expect: true,
			},
			expectedStatus: http.StatusForbidden,
			expectedBody:   pvzservice.ErrUnauthorized.Error() + "\n",
		},
		{
			name:        "внутренняя ошибка",
			requestBody: `{"city":"Санкт-Петербург"}`,
			fields: fields{
				err:    errors.New("unexpected"),
				expect: true,
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "unexpected\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			pvzMock := mocks.NewPVZServiceMock(mc)
			if tt.fields.expect {
				pvzMock.CreatePVZMock.Return(tt.fields.result, tt.fields.err)
			}

			apiInstance := NewAPI(
				mocks.NewUserServiceMock(mc),
				pvzMock,
				mocks.NewReceptionServiceMock(mc),
				mocks.NewProductServiceMock(mc),
				nil,
			)

			req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()

			apiInstance.PostPvz(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" && w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}

			if tt.expectedStatus == http.StatusCreated && tt.fields.result != nil {
				var resp PVZ
				err := json.NewDecoder(w.Body).Decode(&resp)
				if err != nil {
					t.Errorf("failed to decode response JSON: %v", err)
				}
				if *resp.Id != tt.fields.result.ID {
					t.Errorf("expected ID %v, got %v", tt.fields.result.ID, resp.Id)
				}
				if resp.City != PVZCity(tt.fields.result.City) {
					t.Errorf("expected city %v, got %v", tt.fields.result.City, resp.City)
				}
			}
		})
	}
}
