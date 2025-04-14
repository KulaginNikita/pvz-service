package pvzservice

import (
	"context"
	"errors"
	"time"

	"github.com/KulaginNikita/pvz-service/internal/domain/pvz"
	"github.com/KulaginNikita/pvz-service/internal/middleware"
	"github.com/KulaginNikita/pvz-service/internal/repository/pvzrepo/converter"
	"github.com/google/uuid"
)


func (s *pvzService) CreatePVZ(ctx context.Context, pvz *pvz.PVZ) (*pvz.PVZ, error) {
	role, ok := ctx.Value(middleware.RoleContextKey).(string)
	if !ok {
		return nil, errors.New("missing role in context")
	}

	if role != "moderator" {
		return nil, errors.New("user does not have 'moderator' role")
	}

	if pvz.ID == uuid.Nil {
		pvz.ID = uuid.New()
	}

	if pvz.RegisteredAt.IsZero() {
		pvz.RegisteredAt = time.Now()
	}

	if pvz.City != "Москва" && pvz.City != "Санкт-Петербург" && pvz.City != "Казань" {
		return nil, errors.New("PVZ can only be created in Москва, Санкт-Петербург, or Казань")
	}

	err := s.repo.CreatePVZ(ctx, converter.ToDB(pvz))
	if err != nil {
		return nil, err
	}

	return pvz, nil
}

