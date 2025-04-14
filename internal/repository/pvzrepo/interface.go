package pvzrepo

import (
	"context"

	"github.com/KulaginNikita/pvz-service/internal/domain/pvz"
	"github.com/KulaginNikita/pvz-service/internal/repository/pvzrepo/model"
)

type PVZRepository interface {
	CreatePVZ(ctx context.Context, p *model.PVZ) error
	GetPVZ(ctx context.Context, params *pvz.PVZFilter) ([]pvz.PVZ, error)
}
