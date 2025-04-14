package receptionservice

import (
	"github.com/KulaginNikita/pvz-service/internal/repository/receptionrepo"
	"github.com/KulaginNikita/pvz-service/pkg/jwtutil"
)
	

type receptionService struct {
	repo        receptionrepo.ReceptionRepository
	jwtManager  *jwtutil.Manager
}

func NewReceptionService(repo receptionrepo.ReceptionRepository, jwtManager *jwtutil.Manager) *receptionService {
	return &receptionService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}