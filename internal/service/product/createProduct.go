package productservice

import (
	"context"
	"errors"
	"time"

	"github.com/KulaginNikita/pvz-service/internal/domain/product"
	"github.com/KulaginNikita/pvz-service/internal/middleware"
	"github.com/KulaginNikita/pvz-service/internal/repository/productrepo/converter"
	"github.com/google/uuid"
)

func (s *productService) CreateProduct(ctx context.Context, p *product.Product) error {

	role, ok := ctx.Value(middleware.RoleContextKey).(string)
	if !ok {
		return errors.New("missing role in context")
	}
	if role != "employee" {
		return errors.New("user does not have 'employee' role")
	}

	receptionID, err := s.repo.GetOpenReceptionID(ctx, uuid.UUID(p.ReceptionID))
	if err != nil {
		return err
	}

	p.DateTime = time.Now()
	p.ID = uuid.New()
	p.ReceptionID = receptionID

	err = s.repo.CreateProduct(ctx, converter.ToDB(p))
	if err != nil {
		return err
	}

	return nil
}
