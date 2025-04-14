package userservice

import (
	"github.com/KulaginNikita/pvz-service/internal/repository/userrepo"
	"github.com/KulaginNikita/pvz-service/pkg/jwtutil"
)

type userService struct {
	repo       userrepo.UserRepository
	jwtManager jwtutil.TokenManager
}

func NewUserService(repo userrepo.UserRepository, jwt jwtutil.TokenManager) *userService {
	return &userService{
		repo:       repo,
		jwtManager: jwt,
	}
}
