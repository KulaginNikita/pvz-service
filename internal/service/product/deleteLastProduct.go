package productservice

import (
	"context"
	"errors"

	"github.com/KulaginNikita/pvz-service/internal/middleware"
	"github.com/google/uuid"
)

func (p *productService) DeleteProduct(ctx context.Context, pvzID uuid.UUID) error {

	role, ok := ctx.Value(middleware.RoleContextKey).(string)
	if !ok {
		return errors.New("missing role in context")
	}

	if role != "employee" {
		return errors.New("user does not have 'employee' role")
	}

	receptionID, err := p.repo.GetOpenReceptionID(ctx, pvzID)
	if err != nil {
		return err
	}

	productID, err := p.repo.GetLastProductIDByReceptionID(ctx, receptionID)
	if err != nil {
		return err
	}

	err = p.repo.DeleteProductByID(ctx, productID)
	if err != nil {
		return err
	}

	return nil
}
