package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type WarrantyRepo interface {
	CreateWarranty(ctx context.Context, name, warrantyType string, months int) (int64, error)
	GetWarrantyByWarrantyId(ctx context.Context, warrantyId int64) (*model.Warranty, error)
	GetStoreProductWarranty(ctx context.Context, storeId, productId int64) (*model.Warranty, error)
}

type WarrantyRepoImpl struct {
	db *pgxpool.Pool
}

func NewWarrantyRepoImpl(db *pgxpool.Pool) *WarrantyRepoImpl {
	return &WarrantyRepoImpl{db: db}
}

func (w *WarrantyRepoImpl) CreateWarranty(ctx context.Context, name, warrantyType string, months int) (int64, error) {
	query := `
insert into warranty(name, type, month)
values ($1, $2, $3) returning id
`
	row := w.db.QueryRow(ctx, query, name, warrantyType, months)
	var id int64 = -1
	err := row.Scan(&id)
	return id, err
}

func (w *WarrantyRepoImpl) GetWarrantyByWarrantyId(ctx context.Context, warrantyId int64) (*model.Warranty, error) {
	query := `
select id, name, type, month, created_at
from warranty where id = $1
`
	row := w.db.QueryRow(ctx, query, warrantyId)
	var warranty model.Warranty
	err := row.Scan(&warranty.Id, &warranty.Name, &warranty.WarrantyType, &warranty.Month, &warranty.CreatedAt)
	return &warranty, err
}

func (w *WarrantyRepoImpl) GetStoreProductWarranty(ctx context.Context, storeId, productId int64) (*model.Warranty, error) {
	query := `
select id, name, type, month, created_at
from warranty
where id in (select warranty_id from store_product where store_id = $1 and product_id = $2)
limit 1
`
	row := w.db.QueryRow(ctx, query, storeId, productId)
	var warranty model.Warranty
	err := row.Scan(&warranty.Id, &warranty.Name, &warranty.WarrantyType, &warranty.Month, &warranty.CreatedAt)
	return &warranty, err
}
