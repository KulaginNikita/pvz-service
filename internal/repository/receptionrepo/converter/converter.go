package converter

import (
	domain "github.com/KulaginNikita/pvz-service/internal/domain/reception"
	"github.com/KulaginNikita/pvz-service/internal/repository/receptionrepo/model"
	"github.com/google/uuid"
)

func ToDB(r *domain.Reception) *model.Reception {
	return &model.Reception{
		ID:       r.ID,
		PVZID:    r.PVZID,
		Status:   string(r.Status),
		StartedAt: r.StartedAt,
	}
}

func ToDBID(id uuid.UUID) *model.PVZID{
	return &model.PVZID{
		PVZID: id,
	}
}


func FromDB(m *model.Reception) *domain.Reception {
	return &domain.Reception{
		ID:       m.ID,
		PVZID:    m.PVZID,
		Status:   domain.Status(m.Status),
		StartedAt: m.StartedAt,
	}
}
