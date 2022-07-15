package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type OrderRepo interface {
	IsUserPaidTheOrder(ctx context.Context, userId, orderId int64) (bool, error)
	IsOrderDeletedByUser(ctx context.Context, userId, orderId int64) (bool, error)
	IsUserOwnsTheOrder(ctx context.Context, userId, orderId int64) (bool, error)

	CreateOrder(ctx context.Context, userId int64) (int64, error)
	DeleteOrder(ctx context.Context, userId, orderId int64) error

	AddProductToOrder(ctx context.Context, userId, orderId, productId, storeId int64) error
	RemoveProductFromOrder(ctx context.Context, userId, orderId, productId, storeId int64) error

	UpdateProductOrderQuantity(ctx context.Context, userId, orderId, quantity, productId, storeId int64) error

	GetAllProductsInTheOrder(ctx context.Context, userId, orderId int64) ([]model.OrderProduct, error)

	GetAllOrdersByUserId(ctx context.Context, userId int64) ([]model.Order, error)

	PayOrder(ctx context.Context, userId, orderId int64, payedPrice float64) error

	ApplyPromotionCodeToOrder(ctx context.Context, userId, orderId int64, promotionCode string) error
	DeletePromotionCodeFromOrder(ctx context.Context, userId, orderId int64) error

	GetShippingMethod(ctx context.Context) ([]model.ShippingMethod, error)
	UpdateShippingMethod(ctx context.Context, userId, orderId int64, ShippingMethodName string) error
}

type OrderRepoImpl struct {
	db *pgxpool.Pool
}

func NewOrderRepoImpl(db *pgxpool.Pool) *OrderRepoImpl {
	return &OrderRepoImpl{db: db}
}

func (o *OrderRepoImpl) IsUserPaidTheOrder(ctx context.Context, userId, orderId int64) (bool, error) {
	query := `
select is_paid from "order" where id = $1 and user_id= $2 and deleted_at is null
`
	row := o.db.QueryRow(ctx, query, userId, orderId)
	var isPaid bool
	err := row.Scan(&isPaid)
	return isPaid, err
}

func (o *OrderRepoImpl) IsOrderDeletedByUser(ctx context.Context, userId, orderId int64) (bool, error) {
	query := `
select deleted_at is not null from "order" where id= $1 and user_id = $2
`
	row := o.db.QueryRow(ctx, query, orderId, userId)
	var isDeleted bool
	err := row.Scan(&isDeleted)
	return isDeleted, err
}

func (o *OrderRepoImpl) IsUserOwnsTheOrder(ctx context.Context, userId, orderId int64) (bool, error) {
	query := `
select exists(select 1 from "order" where id = $1 and user_id= $2 )
`
	row := o.db.QueryRow(ctx, query, orderId, userId)
	var isUserOwnsTheOrder bool
	err := row.Scan(&isUserOwnsTheOrder)
	return isUserOwnsTheOrder, err
}

func (o *OrderRepoImpl) mustUserOwnsTheOrder(ctx context.Context, userId, orderId int64) error {
	isOwnsTheOrder, err := o.IsUserOwnsTheOrder(ctx, userId, orderId)
	if err != nil {
		return err
	}
	if !isOwnsTheOrder {
		return errors.New("invalid orderId")
	}
	return nil
}

func (o *OrderRepoImpl) mustOrderInProgress(ctx context.Context, userId, orderId int64) error {
	err := o.mustUserOwnsTheOrder(ctx, userId, orderId)
	if err != nil {
		return err
	}
	isOrderDeleted, err := o.IsOrderDeletedByUser(ctx, userId, orderId)
	if err != nil {
		return err
	}
	if isOrderDeleted {
		return errors.New("the order deleted")
	}
	return nil

}

func (o *OrderRepoImpl) CreateOrder(ctx context.Context, userId int64) (int64, error) {
	query := `
insert into "order"(user_id )
values ($1) returning id;
`
	row := o.db.QueryRow(ctx, query, userId)
	var orderId int64 = -1
	err := row.Scan(&orderId)
	return orderId, err
}

func (o *OrderRepoImpl) DeleteOrder(ctx context.Context, userId, orderId int64) error {
	err := o.mustOrderInProgress(ctx, userId, orderId)
	if err != nil {
		return err
	}
	query := `
delete from "order" where id = $1 and user_id= $2 and deleted_at is null
`
	_, err = o.db.Exec(ctx, query, orderId, userId)
	return err
}

func (o *OrderRepoImpl) AddProductToOrder(ctx context.Context, userId, orderId, productId, storeId int64) error {
	err := o.mustOrderInProgress(ctx, userId, orderId)
	if err != nil {
		return err
	}

	query := `
insert into product_order(product_id, store_id, order_id)
values ($1, $2, $3);
`
	_, err = o.db.Exec(ctx, query, productId, storeId, orderId)
	return err
}

func (o *OrderRepoImpl) RemoveProductFromOrder(ctx context.Context, userId, orderId, productId, storeId int64) error {
	err := o.mustOrderInProgress(ctx, userId, orderId)
	if err != nil {
		return err
	}

	query := `
delete
from product_order
where order_id = $1
 and product_id = $2
 and store_id = $3
`
	_, err = o.db.Exec(ctx, query, orderId, productId, storeId)
	return err
}

func (o *OrderRepoImpl) UpdateProductOrderQuantity(ctx context.Context, userId, orderId, quantity, productId, storeId int64) error {
	err := o.mustOrderInProgress(ctx, userId, orderId)
	if err != nil {
		return err
	}

	query := `
update product_order
set quantity = $1
where order_id = $2
 and product_id = $3
 and store_id = $4
`
	_, err = o.db.Exec(ctx, query, quantity, orderId, productId, storeId)
	return err
}

func (o *OrderRepoImpl) GetAllProductsInTheOrder(ctx context.Context, userId, orderId int64) ([]model.OrderProduct, error) {
	err := o.mustOrderInProgress(ctx, userId, orderId)
	if err != nil {
		return nil, err
	}

	query := `
select product_id, store_id, order_id, quantity, created_at
from product_order
where order_id = $1
`
	rows, err := o.db.Query(ctx, query, orderId)
	var orderProducts = make([]model.OrderProduct, 0)
	for rows.Next() {
		var orderProduct model.OrderProduct
		err := rows.Scan(&orderProduct.ProductId, &orderProduct.StoreId, &orderProduct.OrderId, &orderProduct.Quantity, &orderProduct.CreatedAt)
		if err != nil {
			return nil, err
		}
		orderProducts = append(orderProducts, orderProduct)
	}
	return orderProducts, nil
}

func (o *OrderRepoImpl) GetAllOrdersByUserId(ctx context.Context, userId int64) ([]model.Order, error) {
	query := `
select id,
       status,
       tracking_code,
       user_id,
       address_id,
       shipping_method_id,
       applied_promotion_code,
       payed_price,
       is_paid,
       pay_date,
       created_at
from "order"
where user_id = $1 and deleted_at is null;
`
	rows, err := o.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	var orders = make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.Id, &order.Status, &order.TrackingCode,
			&order.UserId, &order.AddressId, &order.ShippingMethodId,
			&order.AppliedPromotionCode, &order.PayedPrice, &order.IsPaid, &order.PayDate, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (o *OrderRepoImpl) PayOrder(ctx context.Context, userId, orderId int64, payedPrice float64) error {
	err := o.mustOrderInProgress(ctx, userId, orderId)
	if err != nil {
		return err
	}

	query := `
update "order"
set is_paid = true,
    status='confirmed',
    payed_price = $1,
    pay_date=now()
where id = $2
  and user_id = $3;
`
	_, err = o.db.Exec(ctx, query, payedPrice, orderId, userId)
	return err
}

func (o *OrderRepoImpl) ApplyPromotionCodeToOrder(ctx context.Context, userId, orderId int64, promotionCode string) error {
	err := o.mustOrderInProgress(ctx, userId, orderId)
	if err != nil {
		return err
	}

	query := `
update "order" set applied_promotion_code = $1 where id = $2 and user_id = $3 ;

`
	_, err = o.db.Exec(ctx, query, promotionCode, orderId, userId)
	return err
}

func (o *OrderRepoImpl) DeletePromotionCodeFromOrder(ctx context.Context, userId, orderId int64) error {
	err := o.mustOrderInProgress(ctx, userId, orderId)
	if err != nil {
		return err
	}

	query := `
update "order" set applied_promotion_code = null where id= ? and user_id = ? ;
`
	_, err = o.db.Exec(ctx, query, orderId, userId)
	return err
}

func (o *OrderRepoImpl) GetShippingMethod(ctx context.Context) ([]model.ShippingMethod, error) {
	query := `
select name, expected_arrival_working_days, base_cost, created_at
from shipping_method;
`
	rows, err := o.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var shippingMethods = make([]model.ShippingMethod, 0)
	for rows.Next() {
		var method model.ShippingMethod
		err := rows.Scan(&method.Name, &method.ExpectedArrivalWorkingDays, &method.BaseCost, &method.CreatedAt)
		if err != nil {
			return nil, err
		}
		shippingMethods = append(shippingMethods, method)
	}
	return shippingMethods, nil
}

func (o *OrderRepoImpl) UpdateShippingMethod(ctx context.Context, userId, orderId int64, ShippingMethodName string) error {
	err := o.mustOrderInProgress(ctx, userId, orderId)
	if err != nil {
		return err
	}

	query := `
update "order" set shipping_method_id = $1 where id= $2 and user_id = $3 ;
`
	_, err = o.db.Exec(ctx, query, ShippingMethodName, orderId, userId)
	return err
}
