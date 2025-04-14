package pvzrepo

import (
	"context"
	"fmt"

	"github.com/KulaginNikita/pvz-service/internal/domain/pvz"
	"github.com/KulaginNikita/pvz-service/internal/repository/pvzrepo/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName              = "pvz"
	columnID               = "id"
	columnCity             = "city"
	columnRegistrationDate = "registration_date"
)

const ()

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

func (r *repository) GetPVZ(ctx context.Context, params *pvz.PVZFilter) ([]pvz.PVZ, error) {
	query, args, err := sq.Select("DISTINCT pvz."+columnID, "pvz."+columnCity, "pvz."+columnRegistrationDate).
		From(tableName).
		Join("receptions ON pvz.id = receptions.pvz_id").
		Where(sq.And{
			sq.GtOrEq{"receptions.started_at": params.StartDate},
			sq.LtOrEq{"receptions.started_at": params.EndDate},
		}).
		PlaceholderFormat(sq.Dollar).
		Limit(uint64(params.Limit)).
		Offset(uint64(params.Offset)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("sql generation failed: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var pvzList []pvz.PVZ

	for rows.Next() {
		var p pvz.PVZ
		if err := rows.Scan(&p.ID, &p.City, &p.RegisteredAt); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		pvzList = append(pvzList, p)
	}

	return pvzList, nil
}