package productrepo

import (
	"context"

	"errors"

	"github.com/KulaginNikita/pvz-service/internal/repository/productrepo/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName         = "products"
	columnID          = "id"
	columnReceptionID = "reception_id"
	columnCreatedAt   = "created_at"
	columnType        = "type"
)

type repository struct {
	db *pgxpool.Pool
}


func NewProductRepository(pool *pgxpool.Pool) *repository {
	return &repository{db: pool}
}

func (r *repository) CreateProduct(ctx context.Context, p *model.Product) error {	
	query, args, err := sq.Insert(tableName).
		Columns(columnID, columnReceptionID, columnCreatedAt, columnType).
		Values(p.ID, p.ReceptionID, p.DateTime, p.Type).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}
	
	_, err = r.db.Exec(ctx, query, args...)
	return err
}

func (r *repository) GetOpenReceptionID(ctx context.Context, pvzID uuid.UUID) (uuid.UUID, error) {
	query, args, err := sq.Select("id").
		From("receptions").
		Where(sq.Eq{"pvz_id": pvzID, "status": "in_progress",}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	var id uuid.UUID
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, errors.New("no open reception found")
		}
		return uuid.Nil, err
	}

	return id, nil
}


func (r *repository) GetLastProductIDByReceptionID(ctx context.Context, receptionID uuid.UUID) (uuid.UUID, error) {
	query, args, err := sq.
		Select("id").
		From("products").
		Where(sq.Eq{"reception_id": receptionID}).
		OrderBy("created_at DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	var productID uuid.UUID
	err = r.db.QueryRow(ctx, query, args...).Scan(&productID)
	if err != nil {
		return uuid.Nil, errors.New("no products found to delete")
	}

	return productID, nil
}


func (r *repository) DeleteProductByID(ctx context.Context, productID uuid.UUID) error {
	query, args, err := sq.
		Delete("products").
		Where(sq.Eq{"id": productID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("product not found or already deleted")
	}

	return nil
}
