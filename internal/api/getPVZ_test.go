package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mocks "github.com/KulaginNikita/pvz-service/internal/api/mocks"
	"github.com/KulaginNikita/pvz-service/internal/domain/pvz"
	"github.com/gojuno/minimock/v3"
)

func TestGetPvz(t *testing.T) {
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	type fields struct {
		result []pvz.PVZ
		err    error
		expect bool
	}
	tests := []struct {
		name           string
		params         GetPvzParams
		fields         fields
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "успешное получение списка ПВЗ",
			params: GetPvzParams{
				StartDate: &startDate,
				EndDate:   &endDate,
				Limit:     intPtr(5),
				Page:      intPtr(2),
			},
			fields: fields{
				result: []pvz.PVZ{
					{City: "Казань"},
					{City: "Москва"},
				},
				expect: true,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "", // валидируем отдельно
		},
		{
			name: "отсутствует дата",
			params: GetPvzParams{
				StartDate: nil,
				EndDate:   nil,
			},
			fields:         fields{expect: false},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "startDate and endDate are required\n",
		},
		{
			name: "ошибка от сервиса",
			params: GetPvzParams{
				StartDate: &startDate,
				EndDate:   &endDate,
			},
			fields: fields{
				err:    errors.New("db error"),
				expect: true,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "db error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			pvzMock := mocks.NewPVZServiceMock(mc)
			if tt.fields.expect {
				pvzMock.GetPVZMock.Return(tt.fields.result, tt.fields.err)
			}

			apiInstance := NewAPI(
				mocks.NewUserServiceMock(mc),
				pvzMock,
				mocks.NewReceptionServiceMock(mc),
				mocks.NewProductServiceMock(mc),
				nil,
			)

			req := httptest.NewRequest(http.MethodGet, "/pvz", nil)
			w := httptest.NewRecorder()

			apiInstance.GetPvz(w, req, tt.params)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" && w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}

			if tt.expectedStatus == http.StatusOK && len(tt.fields.result) > 0 {
				var decoded []pvz.PVZ
				if err := json.NewDecoder(w.Body).Decode(&decoded); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if len(decoded) != len(tt.fields.result) {
					t.Errorf("expected %d items, got %d", len(tt.fields.result), len(decoded))
				}
			}
		})
	}
}

func intPtr(v int) *int {
	return &v
}
