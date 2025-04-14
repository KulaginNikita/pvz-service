package pvzservice

import (
	"github.com/KulaginNikita/pvz-service/internal/repository/pvzrepo"
	"github.com/KulaginNikita/pvz-service/pkg/jwtutil"
	"errors"
)

var (
	ErrForbiddenCity = errors.New("PVZ can only be created in Москва, Санкт-Петербург, or Казань")
	ErrUnauthorized  = errors.New("user is not authorized")
)

type pvzService struct {
	repo       pvzrepo.PVZRepository
	jwtManager *jwtutil.Manager
}

func NewPVZService(repo pvzrepo.PVZRepository, jwtManager *jwtutil.Manager) *pvzService {
	return &pvzService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}