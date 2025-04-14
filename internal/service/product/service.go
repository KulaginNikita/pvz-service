package productservice

import (
	"github.com/KulaginNikita/pvz-service/internal/repository/productrepo"
	"github.com/KulaginNikita/pvz-service/pkg/jwtutil"
)

type productService struct {
	repo       productrepo.ProductRepository
	jwtManager *jwtutil.Manager
}

func NewProductService(repo productrepo.ProductRepository, jwtManager *jwtutil.Manager) *productService {
	return &productService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}