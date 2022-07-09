package repository

import (
	"context"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type StoreRepository interface {
	CreateStore(ctx context.Context, name, description string, owner, creator int64) (int64, error)
	UpdateStore(ctx context.Context, name, description string, owner, storeId int64) (int64, error)
	DeleteStore(ctx context.Context, owner, storeId int64) error
	GetStoresByUserId(ctx context.Context, userId int64) ([]model.Store, error)
	GetAllStores(ctx context.Context) ([]model.Store, error)

	IsUserOwnerOfStore(ctx context.Context, userId, storeId int64) (bool, error)
	AddStoreAddress(ctx context.Context, country, province, city, street, postalCode string) (int64, error)
	// first check if the user is creator of store or not?
	GetStoreAddressesByUserId(ctx context.Context, userId, storeId int64) ([]model.StoreAddress, error)
	// first check if the user is creator of store or not?
	UpdateStoreAddresses(ctx context.Context, country, province, city, street, postalCode string) (int64, error)

	// first check if the user is creator of store or not?
	AddStoreCategory(ctx context.Context, userId, storeId, categoryId int64) error
	// first check if the user is creator of store or not?
	DeleteStoreCategory(ctx context.Context, userId, storeId, categoryId int64) error

	// first check if the user is creator of store or not?
	AddProductToStore(ctx context.Context, userId, storeId, productId, warrantyId int64, price, offPercent, maxOffPrice float64, availableCount int) (int64, error)
	// first check if the user is creator of store or not?
	UpdateStoreProduct(ctx context.Context, userId, storeId, productId, warrantyId int64, price, offPercent, maxOffPrice float64, availableCount int) (int64, error)
	GetAllProductsByStoreId(ctx context.Context) ([]model.Product, error)
}
