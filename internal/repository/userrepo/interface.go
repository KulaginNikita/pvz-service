package userrepo

import (
	"context"
	"github.com/KulaginNikita/pvz-service/internal/repository/userrepo/model"
)


type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

