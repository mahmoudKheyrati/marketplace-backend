package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type NotificationRepo interface {
	GetAvailableNotifications(ctx context.Context, userId int64) ([]int64, error)
	GetPendingNotifications(ctx context.Context, userId int64) ([]int64, error)
	Subscribe(ctx context.Context, productId, userId int64) error
	MarkNotificationAsSeen(ctx context.Context, userId int64, notificationId int64) error
}
type NotificationRepoImpl struct {
	db *pgxpool.Pool
}

func NewNotificationRepoImpl(db *pgxpool.Pool) *NotificationRepoImpl {
	return &NotificationRepoImpl{db: db}
}

func (n *NotificationRepoImpl) GetAvailableNotifications(ctx context.Context, userId int64) ([]int64, error) {
	query := `select id, product_id, user_id, created_at, is_notification_sent, available_status
from product_available_subscription where user_id = $1 and is_notification_sent = false and available_status = true`
	productIds := make([]int64, 0)
	rows, err := n.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var productId int64
		err := rows.Scan(&productId)
		if err != nil {
			return nil, err
		}
		if productId != 0 {
			productIds = append(productIds, productId)
		}
	}
	return productIds, nil
}

func (n *NotificationRepoImpl) GetPendingNotifications(ctx context.Context, userId int64) ([]int64, error) {
	query := `select product_id, user_id, created_at, is_notification_sent
from product_available_subscription
where user_id = $1
  and available_status = false
  and is_notification_sent = false`
	productIds := make([]int64, 0)
	rows, err := n.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var productId int64
		err := rows.Scan(&productId)
		if err != nil {
			return nil, err
		}
		if productId != 0 {
			productIds = append(productIds, productId)
		}
	}
	return productIds, nil
}

func (n *NotificationRepoImpl) Subscribe(ctx context.Context, productId, userId int64) error {
	query := `insert into product_available_subscription(product_id, user_id)
values ($1, $2)`

	_, err := n.db.Query(ctx, query, productId, userId)
	return err
}

func (n *NotificationRepoImpl) MarkNotificationAsSeen(ctx context.Context, userId int64, notificationId int64) error {
	query := `update product_available_subscription
set is_notification_sent = true
where user_id = $1
  and id = $2`
	_, err := n.db.Query(ctx, query, userId, notificationId)
	return err
}
