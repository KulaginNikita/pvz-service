package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mocks"github.com/KulaginNikita/pvz-service/internal/api/mocks"
	"github.com/gojuno/minimock/v3"
)

func TestPostLogin(t *testing.T) {
	type fields struct {
		loginToken string
		loginError error
		expectCall bool
	}
	tests := []struct {
		name           string
		requestBody    PostLoginJSONBody
		fields         fields
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "успешная авторизация",
			requestBody: PostLoginJSONBody{
				Email:    "test@example.com",
				Password: "password123",
			},
			fields:         fields{loginToken: "mocked-jwt-token", expectCall: true},
			expectedStatus: http.StatusOK,
			expectedBody:   `"mocked-jwt-token"` + "\n",
		},
		{
			name: "невалидный JSON",
			requestBody: PostLoginJSONBody{
				Email:    "",
				Password: "",
			},
			fields:         fields{expectCall: false},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid input\n",
		},
		{
			name: "ошибка авторизации",
			requestBody: PostLoginJSONBody{
				Email:    "user@example.com",
				Password: "wrongpass",
			},
			fields:         fields{loginError: errors.New("invalid credentials"), expectCall: true},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "authentication failed: invalid credentials\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			userMock := mocks.NewUserServiceMock(mc)
			if tt.fields.expectCall {
				userMock.LoginMock.Return(tt.fields.loginToken, tt.fields.loginError)
			}

			apiInstance := NewAPI(
				userMock,
				mocks.NewPVZServiceMock(mc),
				mocks.NewReceptionServiceMock(mc),
				mocks.NewProductServiceMock(mc),
				nil,
			)

			var req *http.Request
			if tt.name == "невалидный JSON" {
				req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("{"))
			} else {
				body, _ := json.Marshal(tt.requestBody)
				req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			}

			w := httptest.NewRecorder()
			apiInstance.PostLogin(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}
