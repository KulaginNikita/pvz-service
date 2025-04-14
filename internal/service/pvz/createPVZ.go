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


func (s *pvzService) CreatePVZ(ctx context.Context, pvz *pvz.PVZ) error {
	// Получаем роль из контекста
	role, ok := ctx.Value(middleware.RoleContextKey).(string)
	if !ok {
		return errors.New("missing role in context")
	}

	// Проверка роли
	if role != "moderator" {
		return errors.New("user does not have 'moderator' role")
	}

	if pvz.ID == uuid.Nil {
		pvz.ID = uuid.New()
	}

	if pvz.RegisteredAt.IsZero() {
		pvz.RegisteredAt = time.Now()
	}

	// Проверка на допустимый город
	if pvz.City != "Москва" && pvz.City != "Санкт-Петербург" && pvz.City != "Казань" {
		return errors.New("PVZ can only be created in Москва, Санкт-Петербург, or Казань")
	}

	// Вызов репозитория
	err := s.repo.CreatePVZ(ctx, converter.ToDB(pvz))
	if err != nil {
		return err
	}

	return nil
}
