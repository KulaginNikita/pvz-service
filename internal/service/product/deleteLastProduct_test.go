package productservice

import (
	"context"
	"errors"
	"testing"

	"github.com/KulaginNikita/pvz-service/internal/middleware"
	"github.com/KulaginNikita/pvz-service/internal/service/product/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDeleteProduct(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx   context.Context
		pvzID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func(repo *mocks.ProductRepositoryMock)
		wantErr   bool
	}{
		{
			name: "успешное удаление последнего продукта",
			args: args{
				ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
				pvzID: uuid.New(),
			},
			mockSetup: func(repo *mocks.ProductRepositoryMock) {
				receptionID := uuid.New()
				productID := uuid.New()
				repo.GetOpenReceptionIDMock.Return(receptionID, nil)
				repo.GetLastProductIDByReceptionIDMock.Return(productID, nil)
				repo.DeleteProductByIDMock.Return(nil)
			},
			wantErr: false,
		},
		{
			name: "отсутствует роль в контексте",
			args: args{
				ctx:   context.Background(),
				pvzID: uuid.New(),
			},
			mockSetup: func(repo *mocks.ProductRepositoryMock) {},
			wantErr:   true,
		},
		{
			name: "роль не employee",
			args: args{
				ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "client"),
				pvzID: uuid.New(),
			},
			mockSetup: func(repo *mocks.ProductRepositoryMock) {},
			wantErr:   true,
		},
		{
			name: "ошибка при GetOpenReceptionID",
			args: args{
				ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
				pvzID: uuid.New(),
			},
			mockSetup: func(repo *mocks.ProductRepositoryMock) {
				repo.GetOpenReceptionIDMock.Return(uuid.Nil, errors.New("fail get reception"))
			},
			wantErr: true,
		},
		{
			name: "ошибка при GetLastProductIDByReceptionID",
			args: args{
				ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
				pvzID: uuid.New(),
			},
			mockSetup: func(repo *mocks.ProductRepositoryMock) {
				receptionID := uuid.New()
				repo.GetOpenReceptionIDMock.Return(receptionID, nil)
				repo.GetLastProductIDByReceptionIDMock.Return(uuid.Nil, errors.New("fail get product"))
			},
			wantErr: true,
		},
		{
			name: "ошибка при DeleteProductByID",
			args: args{
				ctx:   context.WithValue(context.Background(), middleware.RoleContextKey, "employee"),
				pvzID: uuid.New(),
			},
			mockSetup: func(repo *mocks.ProductRepositoryMock) {
				receptionID := uuid.New()
				productID := uuid.New()
				repo.GetOpenReceptionIDMock.Return(receptionID, nil)
				repo.GetLastProductIDByReceptionIDMock.Return(productID, nil)
				repo.DeleteProductByIDMock.Return(errors.New("fail delete"))
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

			err := svc.DeleteProduct(tt.args.ctx, tt.args.pvzID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
