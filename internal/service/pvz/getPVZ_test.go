package pvzservice

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KulaginNikita/pvz-service/internal/domain/pvz"
	"github.com/KulaginNikita/pvz-service/internal/middleware"
	"github.com/KulaginNikita/pvz-service/internal/service/pvz/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetPVZ(t *testing.T) {
	t.Parallel()

	validParams := &pvz.PVZFilter{
		StartDate: time.Now().Add(-time.Hour),
		EndDate:   time.Now(),
		Limit:     5,
	}

	tests := []struct {
		name      string
		ctx       context.Context
		params    *pvz.PVZFilter
		mockSetup func(r *mocks.PVZRepositoryMock)
		wantErr   bool
	}{
		{
			name:   "успешное получение",
			ctx:    context.WithValue(context.Background(), middleware.RoleContextKey, "moderator"),
			params: validParams,
			mockSetup: func(r *mocks.PVZRepositoryMock) {
				r.GetPVZMock.Return([]pvz.PVZ{{ID: uuid.New()}}, nil)
			},
			wantErr: false,
		},
		{
			name:      "отсутствие роли",
			ctx:       context.Background(),
			params:    validParams,
			mockSetup: func(r *mocks.PVZRepositoryMock) {},
			wantErr:   true,
		},
		{
			name:      "неподдерживаемая роль",
			ctx:       context.WithValue(context.Background(), middleware.RoleContextKey, "client"),
			params:    validParams,
			mockSetup: func(r *mocks.PVZRepositoryMock) {},
			wantErr:   true,
		},
		{
			name:      "не указаны даты",
			ctx:       context.WithValue(context.Background(), middleware.RoleContextKey, "moderator"),
			params:    &pvz.PVZFilter{},
			mockSetup: func(r *mocks.PVZRepositoryMock) {},
			wantErr:   true,
		},
		{
			name:   "ошибка репозитория",
			ctx:    context.WithValue(context.Background(), middleware.RoleContextKey, "moderator"),
			params: validParams,
			mockSetup: func(r *mocks.PVZRepositoryMock) {
				r.GetPVZMock.Return(nil, errors.New("repo fail"))
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

			_, err := svc.GetPVZ(tt.ctx, tt.params)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
