package userservice

import (
	"github.com/KulaginNikita/pvz-service/internal/repository/userrepo"
	"github.com/KulaginNikita/pvz-service/internal/service"
	"github.com/KulaginNikita/pvz-service/pkg/jwtutil"
)
type userService struct {
	repo       userrepo.UserRepository
	jwtManager *jwtutil.Manager
}

func NewUserService(repo userrepo.UserRepository, jwtManager *jwtutil.Manager) service.UserService {
	return &userService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}