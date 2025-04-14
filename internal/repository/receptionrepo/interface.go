package receptionrepo

import (
	"context"

	"github.com/KulaginNikita/pvz-service/internal/repository/receptionrepo/model"
)

type ReceptionRepository interface {
	Create(ctx context.Context, rec *model.Reception) error
	HasOpenReception(ctx context.Context, pvzID *model.PVZID) (bool, error)
	Close(ctx context.Context, pvzID *model.PVZID) error
}