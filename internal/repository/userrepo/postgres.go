package userrepo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/KulaginNikita/pvz-service/internal/repository/userrepo/model"
)

const (
	tableName     = "users"
	columnID      = "id"
	columnEmail   = "email"
	columnPassword = "password_hash"
	columnRole     = "role"
)

type repository struct {
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *repository {
	return &repository{db: pool}
}

func (r *repository) Create(ctx context.Context, u *model.User) error {
	query, args, err := sq.Insert(tableName).
		Columns(columnEmail, columnPassword, columnRole).
		Values(u.Email, u.Password, u.Role).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query, args, err := sq.Select(columnID, columnEmail, columnPassword, columnRole).
		From(tableName).
		Where(sq.Eq{columnEmail: email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	
	u := &model.User{}
	err = r.db.QueryRow(ctx, query, args...).Scan(&u.ID, &u.Email, &u.Password, &u.Role)
	if err != nil {
		return nil, err
	}

	return u, nil
}
