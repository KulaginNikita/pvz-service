package receptionservice

import (
	"context"
	"errors"
	"github.com/KulaginNikita/pvz-service/internal/middleware"
	"github.com/KulaginNikita/pvz-service/internal/repository/receptionrepo/converter"
	"github.com/google/uuid"
)

// Создание приемки товаров
func (s *receptionService) CloseReception(ctx context.Context, pvzid uuid.UUID) error {
	role, ok := ctx.Value(middleware.RoleContextKey).(string)
	if !ok {
		return errors.New("missing role in context")
	}

	if role != "employee" {
		return errors.New("user does not have 'employee' role")
	}

	hasOpenReception, err := s.repo.HasOpenReception(ctx, converter.ToDBID(pvzid))
	if err != nil {
		return err
	}
	if !hasOpenReception {
		return errors.New("there is no reception in progress for this PVZ yet")
	}

	err = s.repo.Close(ctx, converter.ToDBID(pvzid))
	if err != nil {
		return err
	}

	return nil
}
