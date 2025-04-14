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
	"github.com/KulaginNikita/pvz-service/internal/domain/reception"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
)

func TestPostReceptions(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	type fields struct {
		result *reception.Reception
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
			name:        "успешное создание приёмки",
			requestBody: `{"pvzId":"` + id.String() + `"}`,
			fields: fields{
				result: &reception.Reception{
					ID:       id,
					PVZID:    id,
					Status:   reception.StatusInProgress,
					StartedAt: now,
				},
				expect: true,
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "", // можно валидировать отдельно
		},
		{
			name:           "невалидный JSON",
			requestBody:    `{"pvzId":`,
			fields:         fields{expect: false},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid input\n",
		},
		{
			name:        "ошибка создания приёмки",
			requestBody: `{"pvzId":"` + id.String() + `"}`,
			fields: fields{
				err:    errors.New("creation error"),
				expect: true,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "failed to create reception: creation error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			receptionMock := mocks.NewReceptionServiceMock(mc)
			if tt.fields.expect {
				receptionMock.CreateReceptionMock.Return(tt.fields.result, tt.fields.err)
			}

			apiInstance := NewAPI(
				mocks.NewUserServiceMock(mc),
				mocks.NewPVZServiceMock(mc),
				receptionMock,
				mocks.NewProductServiceMock(mc),
				nil,
			)

			req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()

			apiInstance.PostReceptions(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" && w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}

			if tt.expectedStatus == http.StatusCreated && tt.fields.result != nil {
				var resp reception.Reception
				err := json.NewDecoder(w.Body).Decode(&resp)
				if err != nil {
					t.Errorf("failed to decode response JSON: %v", err)
				}
				if resp.ID != id {
					t.Errorf("expected ID %v, got %v", id, resp.ID)
				}
			}
		})
	}
}
