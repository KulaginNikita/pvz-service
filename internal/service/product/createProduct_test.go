package productservice

import (
	"context"
	"errors"
	"testing"

	domain "github.com/KulaginNikita/pvz-service/internal/domain/product"
	"github.com/KulaginNikita/pvz-service/internal/middleware"
	"github.com/KulaginNikita/pvz-service/internal/service/product/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		p   *domain.Product
	}

	tests := []struct {
		name      string
		args      args
		mockSetup func(r *mocks.ProductRepositoryMock)
		wantErr   bool
	}{
		{
			name: "успешное создание продукта",
			args: args{
				ctx: context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
				p:   &domain.Product{ReceptionID: uuid.New()},
			},
			mockSetup: func(r *mocks.ProductRepositoryMock) {
				r.GetOpenReceptionIDMock.Return(uuid.New(), nil)
				r.CreateProductMock.Return(nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка получения роли из контекста",
			args: args{
				ctx: context.Background(),
				p:   &domain.Product{},
			},
			mockSetup: func(r *mocks.ProductRepositoryMock) {},
			wantErr:   true,
		},
		{
			name: "неправильная роль",
			args: args{
				ctx: context.WithValue(context.Background(), middleware.RoleContextKey, "client"),
				p:   &domain.Product{},
			},
			mockSetup: func(r *mocks.ProductRepositoryMock) {},
			wantErr:   true,
		},
		{
			name: "ошибка GetOpenReceptionID",
			args: args{
				ctx: context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
				p:   &domain.Product{ReceptionID: uuid.New()},
			},
			mockSetup: func(r *mocks.ProductRepositoryMock) {
				r.GetOpenReceptionIDMock.Return(uuid.Nil, errors.New("fail"))
			},
			wantErr: true,
		},
		{
			name: "ошибка CreateProduct",
			args: args{
				ctx: context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
				p:   &domain.Product{ReceptionID: uuid.New()},
			},
			mockSetup: func(r *mocks.ProductRepositoryMock) {
				r.GetOpenReceptionIDMock.Return(uuid.New(), nil)
				r.CreateProductMock.Return(errors.New("db fail"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			mockRepo := mocks.NewProductRepositoryMock(mc)
			tt.mockSetup(mockRepo)

			svc := NewProductService(mockRepo, nil)

			err := svc.CreateProduct(tt.args.ctx, tt.args.p)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
