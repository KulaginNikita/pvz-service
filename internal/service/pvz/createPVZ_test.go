package pvzservice

import (
	"context"
	"errors"
	"testing"

	"github.com/KulaginNikita/pvz-service/internal/domain/pvz"
	"github.com/KulaginNikita/pvz-service/internal/middleware"
	"github.com/KulaginNikita/pvz-service/internal/service/pvz/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreatePVZ(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		ctx       context.Context
		input     *pvz.PVZ
		mockSetup func(r *mocks.PVZRepositoryMock)
		wantErr   bool
	}{
		{
			name: "успешное создание",
			ctx:  context.WithValue(context.Background(), middleware.RoleContextKey, "moderator"),
			input: &pvz.PVZ{
				City: "Москва",
			},
			mockSetup: func(r *mocks.PVZRepositoryMock) {
				r.CreatePVZMock.Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "отсутствует роль",
			ctx:       context.Background(),
			input:     &pvz.PVZ{City: "Москва"},
			mockSetup: func(r *mocks.PVZRepositoryMock) {},
			wantErr:   true,
		},
		{
			name:      "не модератор",
			ctx:       context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
			input:     &pvz.PVZ{City: "Москва"},
			mockSetup: func(r *mocks.PVZRepositoryMock) {},
			wantErr:   true,
		},
		{
			name:      "неразрешённый город",
			ctx:       context.WithValue(context.Background(), middleware.RoleContextKey, "moderator"),
			input:     &pvz.PVZ{City: "Екатеринбург"},
			mockSetup: func(r *mocks.PVZRepositoryMock) {},
			wantErr:   true,
		},
		{
			name:  "ошибка в CreatePVZ",
			ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "moderator"),
			input: &pvz.PVZ{City: "Казань"},
			mockSetup: func(r *mocks.PVZRepositoryMock) {
				r.CreatePVZMock.Return(errors.New("db fail"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			repo := mocks.NewPVZRepositoryMock(mc)
			tt.mockSetup(repo)

			svc := NewPVZService(repo, nil)

			_, err := svc.CreatePVZ(tt.ctx, tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
