package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mocks "github.com/KulaginNikita/pvz-service/internal/api/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func TestPostPvzPvzIdDeleteLastProduct(t *testing.T) {
	type fields struct {
		deleteErr error
		expect    bool
	}
	tests := []struct {
		name           string
		fields         fields
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "успешное удаление",
			fields: fields{
				deleteErr: nil,
				expect:    true,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message": "Product deleted successfully"}`,
		},
		{
			name: "ошибка при удалении",
			fields: fields{
				deleteErr: http.ErrBodyNotAllowed,
				expect:    true,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "failed to delete product: " + http.ErrBodyNotAllowed.Error() + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			productMock := mocks.NewProductServiceMock(mc)
			if tt.fields.expect {
				productMock.DeleteProductMock.Return(tt.fields.deleteErr)
			}

			apiInstance := NewAPI(
				mocks.NewUserServiceMock(mc),
				mocks.NewPVZServiceMock(mc),
				mocks.NewReceptionServiceMock(mc),
				productMock,
				nil,
			)

			pvzID := openapi_types.UUID(uuid.New())
			req := httptest.NewRequest(http.MethodPost, "/pvz/delete_last_product", nil)
			w := httptest.NewRecorder()

			apiInstance.PostPvzPvzIdDeleteLastProduct(w, req, pvzID)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}
