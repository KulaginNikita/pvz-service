package converter

import (
	"github.com/KulaginNikita/pvz-service/internal/domain/pvz"
	"github.com/KulaginNikita/pvz-service/internal/repository/pvzrepo/model"
)
func ToDB(p *pvz.PVZ) *model.PVZ {
	return &model.PVZ{
		ID:   p.ID,
		City: string(p.City),
		RegisteredAt: p.RegisteredAt,
	}
}

func FromDB(m *model.PVZ) *pvz.PVZ {
	return &pvz.PVZ{
		ID:   m.ID,
		City: pvz.City(m.City),
		RegisteredAt: m.RegisteredAt,
	}
}