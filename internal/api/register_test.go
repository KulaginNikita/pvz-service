package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/KulaginNikita/pvz-service/internal/api/mocks"
	"github.com/gojuno/minimock/v3"
)

func TestPostRegister(t *testing.T) {
	type fields struct {
		registerErr error
		expectCall  bool
	}
	tests := []struct {
		name           string
		requestBody    string
		fields         fields
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "успешная регистрация",
			requestBody: `{
				"email": "user@example.com",
				"password": "secret",
				"role": "employee"
			}`,
			fields: fields{
				registerErr: nil,
				expectCall:  true,
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "",
		},
		{
			name:           "невалидный JSON",
			requestBody:    `{"email": "user@example.com"`, 
			fields:         fields{expectCall: false},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid input\n",
		},
		{
			name: "ошибка при регистрации",
			requestBody: `{
				"email": "user@example.com",
				"password": "secret",
				"role": "moderator"
			}`,
			fields: fields{
				registerErr: errors.New("DB error"),
				expectCall:  true,
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "registration failed: DB error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			userMock := mocks.NewUserServiceMock(mc)
			if tt.fields.expectCall {
				userMock.RegisterMock.Return(tt.fields.registerErr)
			}

			apiInstance := NewAPI(
				userMock,
				mocks.NewPVZServiceMock(mc),
				mocks.NewReceptionServiceMock(mc),
				mocks.NewProductServiceMock(mc),
				nil,
			)

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(tt.requestBody))
			w := httptest.NewRecorder()

			apiInstance.PostRegister(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}
