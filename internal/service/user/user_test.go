package userservice

import (
	"context"
	"errors"
	"testing"

	"github.com/KulaginNikita/pvz-service/internal/domain/user"
	"github.com/KulaginNikita/pvz-service/internal/repository/userrepo/converter"
	"github.com/KulaginNikita/pvz-service/internal/service/user/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		email     string
		password  string
		mockSetup func(repo *mocks.UserRepositoryMock, jwt *mocks.JwtManagerMock)
		wantErr   bool
		wantToken string
	}{
		{
			name:     "успешный логин",
			email:    "test@mail.com",
			password: "1234",
			mockSetup: func(repo *mocks.UserRepositoryMock, jwt *mocks.JwtManagerMock) {
				repo.GetByEmailMock.Return(converter.ToDB(&user.User{
					Email:    "test@mail.com",
					Password: "1234",
					Role:     "employee",
				}), nil)
				jwt.GenerateTokenMock.Return("mock-token", nil)
			},
			wantErr:   false,
			wantToken: "mock-token",
		},
		{
			name:     "неверный пароль",
			email:    "fail@mail.com",
			password: "wrong",
			mockSetup: func(repo *mocks.UserRepositoryMock, jwt *mocks.JwtManagerMock) {
				repo.GetByEmailMock.Return(converter.ToDB(&user.User{
					Email:    "fail@mail.com",
					Password: "1234",
					Role:     "employee",
				}), nil)
			},
			wantErr: true,
		},
		{
			name:     "ошибка базы",
			email:    "err@mail.com",
			password: "123",
			mockSetup: func(repo *mocks.UserRepositoryMock, jwt *mocks.JwtManagerMock) {
				repo.GetByEmailMock.Return(nil, errors.New("db fail"))
			},
			wantErr: true,
		},
		{
			name:     "ошибка генерации токена",
			email:    "test@mail.com",
			password: "1234",
			mockSetup: func(repo *mocks.UserRepositoryMock, jwt *mocks.JwtManagerMock) {
				repo.GetByEmailMock.Return(converter.ToDB(&user.User{
					Email:    "test@mail.com",
					Password: "1234",
					Role:     "employee",
				}), nil)
				jwt.GenerateTokenMock.Return("", errors.New("jwt fail"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // захват переменной
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			defer mc.Finish()

			repo := mocks.NewUserRepositoryMock(mc)
			jwt := mocks.NewJwtManagerMock(mc)
			tt.mockSetup(repo, jwt)

			svc := NewUserService(repo, jwt)

			token, err := svc.Login(context.Background(), tt.email, tt.password)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantToken, token)
			}
		})
	}
}
