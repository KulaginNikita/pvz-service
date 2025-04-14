package receptionservice

import (
	"context"
	"errors"
	"testing"

	"github.com/KulaginNikita/pvz-service/internal/middleware"
	"github.com/KulaginNikita/pvz-service/internal/service/reception/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCloseReception(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		ctx       context.Context
		pvzID     uuid.UUID
		mockSetup func(r *mocks.ReceptionRepositoryMock)
		wantErr   bool
	}{
		{
			name:  "успешное закрытие",
			ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
			pvzID: uuid.New(),
			mockSetup: func(r *mocks.ReceptionRepositoryMock) {
				r.HasOpenReceptionMock.Return(true, nil)
				r.CloseMock.Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "нет роли в контексте",
			ctx:       context.Background(),
			pvzID:     uuid.New(),
			mockSetup: func(r *mocks.ReceptionRepositoryMock) {},
			wantErr:   true,
		},
		{
			name:      "не employee",
			ctx:       context.WithValue(context.Background(), middleware.RoleContextKey, "moderator"),
			pvzID:     uuid.New(),
			mockSetup: func(r *mocks.ReceptionRepositoryMock) {},
			wantErr:   true,
		},
		{
			name:  "приемка не открыта",
			ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
			pvzID: uuid.New(),
			mockSetup: func(r *mocks.ReceptionRepositoryMock) {
				r.HasOpenReceptionMock.Return(false, nil)
			},
			wantErr: true,
		},
		{
			name:  "ошибка в HasOpenReception",
			ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
			pvzID: uuid.New(),
			mockSetup: func(r *mocks.ReceptionRepositoryMock) {
				r.HasOpenReceptionMock.Return(false, errors.New("db fail"))
			},
			wantErr: true,
		},
		{
			name:  "ошибка в Close",
			ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
			pvzID: uuid.New(),
			mockSetup: func(r *mocks.ReceptionRepositoryMock) {
				r.HasOpenReceptionMock.Return(true, nil)
				r.CloseMock.Return(errors.New("fail close"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			defer mc.Finish()

			repo := mocks.NewReceptionRepositoryMock(mc)
			tt.mockSetup(repo)

			svc := NewReceptionService(repo, nil)

			err := svc.CloseReception(tt.ctx, tt.pvzID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
