package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/KulaginNikita/pvz-service/internal/api/mocks"

	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func TestPostPvzPvzIdCloseLastReception(t *testing.T) {
	type fields struct {
		receptionServiceErr error
	}
	tests := []struct {
		name           string
		fields         fields
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "успешное закрытие приёмки",
			fields: fields{
				receptionServiceErr: nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "Reception closed successfully"}`,
		},
		{
			name: "ошибка при закрытии приёмки",
			fields: fields{
				receptionServiceErr: errors.New("failed to close"),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "failed to close reception: failed to close\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			receptionMock := mocks.NewReceptionServiceMock(mc)
			receptionMock.CloseReceptionMock.Return(tt.fields.receptionServiceErr)

			userMock := mocks.NewUserServiceMock(mc)
			pvzMock := mocks.NewPVZServiceMock(mc)
			productMock := mocks.NewProductServiceMock(mc)

			apiInstance := NewAPI(userMock, pvzMock, receptionMock, productMock, nil)

			req := httptest.NewRequest(http.MethodPost, "/pvz/123/close_last_reception", nil)
			w := httptest.NewRecorder()
			pvzID := openapi_types.UUID(uuid.New())

			apiInstance.PostPvzPvzIdCloseLastReception(w, req, pvzID)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}
