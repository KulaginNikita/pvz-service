package productrepo

import (
	"context"

	"github.com/KulaginNikita/pvz-service/internal/repository/productrepo/model"
	"github.com/google/uuid"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, p *model.Product) error
	GetOpenReceptionID(ctx context.Context, pvzID uuid.UUID) (uuid.UUID, error)
	GetLastProductIDByReceptionID(ctx context.Context, receptionID uuid.UUID) (uuid.UUID, error)
	DeleteProductByID(ctx context.Context, productID uuid.UUID) error
}