package pvzrepo

import (
	"context"
	"github.com/KulaginNikita/pvz-service/internal/repository/pvzrepo/model"
)

type PVZRepository interface {
	CreatePVZ(ctx context.Context, p *model.PVZ) error
	// GetALLPVZ(ctx context.Context) ([]pvz.PVZ, error)
}

