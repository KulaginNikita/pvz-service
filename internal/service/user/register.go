package userservice

import (
	"context"

	"github.com/KulaginNikita/pvz-service/internal/domain/user"
	"github.com/KulaginNikita/pvz-service/internal/repository/userrepo/converter"
)

func (s *userService) Register(ctx context.Context, u *user.User) error {
	// TODO: hash password before saving (optional now)
	return s.repo.Create(ctx, converter.ToDB(u))
}

