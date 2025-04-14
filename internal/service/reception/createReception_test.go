package receptionservice

import (
	"context"
	"errors"
	"testing"

	"github.com/KulaginNikita/pvz-service/internal/domain/reception"
	"github.com/KulaginNikita/pvz-service/internal/middleware"
	"github.com/KulaginNikita/pvz-service/internal/service/reception/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateReception(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		ctx       context.Context
		input     *reception.Reception
		mockSetup func(r *mocks.ReceptionRepositoryMock)
		wantErr   bool
	}{
		{
			name:  "успешное создание при отсутствии открытого приёма",
			ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
			input: &reception.Reception{PVZID: uuid.New()},
			mockSetup: func(r *mocks.ReceptionRepositoryMock) {
				r.HasOpenReceptionMock.Return(false, nil)
				r.CreateMock.Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "отсутствует роль в контексте",
			ctx:       context.Background(),
			input:     &reception.Reception{},
			mockSetup: func(r *mocks.ReceptionRepositoryMock) {},
			wantErr:   true,
		},
		{
			name:  "роль не employee",
			ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "moderator"),
			input: &reception.Reception{},
			mockSetup: func(r *mocks.ReceptionRepositoryMock) {},
			wantErr: true,
		},
		{
			name:  "ошибка HasOpenReception",
			ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
			input: &reception.Reception{PVZID: uuid.New()},
			mockSetup: func(r *mocks.ReceptionRepositoryMock) {
				r.HasOpenReceptionMock.Return(false, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:  "уже есть открытая приёмка",
			ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
			input: &reception.Reception{PVZID: uuid.New()},
			mockSetup: func(r *mocks.ReceptionRepositoryMock) {
				r.HasOpenReceptionMock.Return(true, nil)
			},
			wantErr: true,
		},
		{
			name:  "ошибка Create",
			ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
			input: &reception.Reception{PVZID: uuid.New()},
			mockSetup: func(r *mocks.ReceptionRepositoryMock) {
				r.HasOpenReceptionMock.Return(false, nil)
				r.CreateMock.Return(errors.New("insert fail"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			repo := mocks.NewReceptionRepositoryMock(mc)
			tt.mockSetup(repo)

			svc := NewReceptionService(repo, nil)

			_, err := svc.CreateReception(tt.ctx, tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
