package service

import (
	"context"

	"github.com/KulaginNikita/pvz-service/internal/domain/product"
	"github.com/KulaginNikita/pvz-service/internal/domain/pvz"
	"github.com/KulaginNikita/pvz-service/internal/domain/reception"
	"github.com/KulaginNikita/pvz-service/internal/domain/user"
	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, u *user.User) error
	Login(ctx context.Context, email, password string) (string, error)
}


type ReceptionService interface {
	CreateReception(ctx context.Context, r *reception.Reception) (*reception.Reception, error)
	CloseReception(ctx context.Context, pvzid uuid.UUID) error
}

type PVZService interface {
	CreatePVZ(ctx context.Context, p *pvz.PVZ) (*pvz.PVZ, error)
	GetPVZ(ctx context.Context, params *pvz.PVZFilter) ([]pvz.PVZ, error)
}


type ProductService interface {
	CreateProduct(ctx context.Context, p *product.Product) error
	DeleteProduct(ctx context.Context, pvzID uuid.UUID) error
}