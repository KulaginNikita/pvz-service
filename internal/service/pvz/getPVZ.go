package pvzservice

import (
	"context"
	"errors"

	"github.com/KulaginNikita/pvz-service/internal/domain/pvz"
	"github.com/KulaginNikita/pvz-service/internal/middleware"
)

func (s *pvzService) GetPVZ(ctx context.Context, params *pvz.PVZFilter) ([]pvz.PVZ, error) {
	// 1. Авторизация
	role, ok := ctx.Value(middleware.RoleContextKey).(string)
	if !ok {
		return nil, errors.New("missing role in context")
	}
	if role != "moderator" && role != "employee" {
		return nil, errors.New("access denied: only moderator or employee allowed")
	}

	// 2. Валидация дат
	if params.StartDate.IsZero() || params.EndDate.IsZero() {
		return nil, errors.New("startDate and endDate must be set")
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	// 3. Вызов репозитория
	return s.repo.GetPVZ(ctx, params)
}
