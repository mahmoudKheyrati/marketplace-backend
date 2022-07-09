package repository

import (
	"context"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type OrderRepo interface {
	IsUserPaidTheOrder(ctx context.Context, orderId int64) (bool, error)
	IsOrderDeletedByUser(ctx context.Context, userId, orderId int64) (bool, error)
	IsUserOwnsTheOrder(ctx context.Context, userId, orderId int64) (bool, error)

	CreateOrder(ctx context.Context, userId int64) (int64, error)
	DeleteOrder(ctx context.Context, userId, orderId int64) error

	AddProductToOrder(ctx context.Context, userId, orderId, productId, storeId int64) error
	RemoveProductFromOrder(ctx context.Context, userId, orderId, productId, storeId int64) error

	UpdateProductOrderQuantity(ctx context.Context, userId, orderId, productId, storeId int64) error

	GetAllProductsInTheOrder(ctx context.Context, userId, orderId int64) ([]model.OrderProduct, error)

	GetAllOrdersByUserId(ctx context.Context, userId int64) ([]model.Order, error)

	PayOrder(ctx context.Context, userId, orderId int64, payedPrice float64) error
}
