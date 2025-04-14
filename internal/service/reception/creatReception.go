package receptionservice

import (
	"context"
	"errors"
	"time"
	"github.com/KulaginNikita/pvz-service/internal/domain/reception"
	"github.com/KulaginNikita/pvz-service/internal/middleware"
	"github.com/KulaginNikita/pvz-service/internal/repository/receptionrepo/converter"
	"github.com/google/uuid"
)

// Создание приемки товаров
func (s *receptionService) CreateReception(ctx context.Context, rec *reception.Reception) (*reception.Reception, error) {
	role, ok := ctx.Value(middleware.RoleContextKey).(string)
	if !ok {
		return nil, errors.New("missing role in context")
	}
	if role != "employee" {
		return nil, errors.New("user does not have 'employee' role")
	}

	hasOpenReception, err := s.repo.HasOpenReception(ctx, converter.ToDBID(rec.PVZID))
	if err != nil {
		return nil, err
	}
	if hasOpenReception {
		return nil, errors.New("there is already an open reception for this PVZ")
	}

	rec.Status = "in_progress"
	rec.StartedAt = time.Now()
	if rec.ID == uuid.Nil {
		rec.ID = uuid.New()
	}

	err = s.repo.Create(ctx, converter.ToDB(rec))
	if err != nil {
		return nil, err
	}

	return rec, nil
}
