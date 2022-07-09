package repository

import (
	"context"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type WarrantyRepo interface {
	CreateWarranty(ctx context.Context, name, warrantyType string, months int) (int64, error)
	GetWarrantyByWarrantyId(ctx context.Context, warrantyId int64) (*model.Warranty, error)
	GetProductWarranty(ctx context.Context, storeId, productId int64) (*model.Warranty, error)
}
