package pvzrepo

import (
	"context"
	"github.com/KulaginNikita/pvz-service/internal/repository/pvzrepo/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName     = "pvz"
	columnID      = "id"
	columnCity   = "city"
	columnRegistrationDate = "registration_date"
)

type repository struct {
	db *pgxpool.Pool
}

func NewPVZRepository(pool *pgxpool.Pool) *repository {
	return &repository{db: pool}
}

func (r *repository) CreatePVZ(ctx context.Context, p *model.PVZ) error {
	query, args, err := sq.Insert(tableName).
		Columns(columnID, columnCity, columnRegistrationDate).
		Values(p.ID, p.City, p.RegisteredAt).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}
	