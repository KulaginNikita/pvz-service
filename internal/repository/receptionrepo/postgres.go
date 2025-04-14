package receptionrepo

import (
	"context"
	"time"

	"github.com/KulaginNikita/pvz-service/internal/repository/receptionrepo/model"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName = "receptions"
	columnPVZID   = "pvz_id"
	columnID = "id"
	columnStartedAt = "started_at"
	columnCloseAt = "closed_at"
	columnStatus = "status"
)

const(
	InProgress = "in_progress"
	Closed = "close"
)

type repo struct {
	pool *pgxpool.Pool
}

func NewReceptionRepository(pool *pgxpool.Pool) ReceptionRepository {
	return &repo{pool: pool}
}

func (r *repo) HasOpenReception(ctx context.Context, pvzID *model.PVZID) (bool, error) {
	query, args, err := squirrel.
		Select("count(*)").
		From(tableName).
		Where(squirrel.Eq{columnPVZID: pvzID.PVZID, columnStatus: InProgress}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return false, err
	}

	var count int
	err = r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}


func (r *repo) Create(ctx context.Context, rec *model.Reception) error {
	query, args, err := squirrel.
		Insert(tableName).
		Columns(columnID, columnPVZID, columnStatus, columnStartedAt).
		Values(rec.ID, rec.PVZID, rec.Status, rec.StartedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	return err
}


func (r *repo) Close(ctx context.Context, pvzID *model.PVZID) error {
	query, args, err := squirrel.
		Update(tableName).
		Set(columnStatus, Closed).
		Set(columnCloseAt, time.Now()).
		Where(squirrel.Eq{columnPVZID: pvzID.PVZID, columnStatus: InProgress}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	return err
}
