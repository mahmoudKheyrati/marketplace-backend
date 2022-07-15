package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"github.com/spf13/cast"
)

type OrderHandler struct {
	orderRepo repository.OrderRepo
}

func NewOrderHandler(orderRepo repository.OrderRepo) *OrderHandler {
	return &OrderHandler{orderRepo: orderRepo}
}

func (o *OrderHandler) IsUserPaidTheOrder(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	orderId := cast.ToInt64(c.Params("orderId"))

	isUserPaidTheOrder, err := o.orderRepo.IsUserPaidTheOrder(ctx, userId, orderId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"is_user_paid_the_order": isUserPaidTheOrder, "status": "ok"})
}

func (o *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)

	orderId, err := o.orderRepo.CreateOrder(ctx, userId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"order_id": orderId, "status": "ok"})
}
func (o *OrderHandler) GetOrderByOrderId(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	orderId := cast.ToInt64(c.Params("orderId"))

	order, err := o.orderRepo.GetOrderByOrderId(ctx, userId, orderId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"order": order, "status": "ok"})
}

func (o *OrderHandler) DeleteOrder(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	orderId := cast.ToInt64(c.Params("orderId"))

	err := o.orderRepo.DeleteOrder(ctx, userId, orderId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

type ProductOrderRequest struct {
	ProductId int64 `json:"product_id"`
	StoreId   int64 `json:"store_id"`
}

func (o *OrderHandler) AddProductToOrder(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	orderId := cast.ToInt64(c.Params("orderId"))

	var request ProductOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	err := o.orderRepo.AddProductToOrder(ctx, userId, orderId, request.ProductId, request.StoreId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

func (o *OrderHandler) RemoveProductFromOrder(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	orderId := cast.ToInt64(c.Params("orderId"))
	storeId := cast.ToInt64(c.Params("storeId"))
	productId := cast.ToInt64(c.Params("productId"))

	err := o.orderRepo.RemoveProductFromOrder(ctx, userId, orderId, productId, storeId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

type UpdateProductOrderQuantityRequest struct {
	ProductId int64 `json:"product_id"`
	StoreId   int64 `json:"store_id"`
	Quantity  int64 `json:"quantity"`
}

func (o *OrderHandler) UpdateProductOrderQuantity(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	orderId := cast.ToInt64(c.Params("orderId"))

	var request UpdateProductOrderQuantityRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	err := o.orderRepo.UpdateProductOrderQuantity(ctx, userId, orderId, request.Quantity, request.ProductId, request.StoreId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

func (o *OrderHandler) GetAllProductsInTheOrder(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	orderId := cast.ToInt64(c.Params("orderId"))

	orderProducts, err := o.orderRepo.GetAllProductsInTheOrder(ctx, userId, orderId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"order_products": orderProducts, "status": "ok"})
}

func (o *OrderHandler) GetAllOrdersByUserId(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)

	orders, err := o.orderRepo.GetAllOrdersByUserId(ctx, userId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"orders": orders, "status": "ok"})
}

type PayOrderRequest struct {
	PayedPrice float64 `json:"pay_price"`
}

func (o *OrderHandler) PayOrder(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	orderId := cast.ToInt64(c.Params("orderId"))

	var request PayOrderRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	err := o.orderRepo.PayOrder(ctx, userId, orderId, request.PayedPrice)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

type PromotionCodeRequest struct {
	ProductId int64 `json:"product_id"`
	StoreId   int64 `json:"store_id"`
}

func (o *OrderHandler) ApplyPromotionCodeToOrder(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	orderId := cast.ToInt64(c.Params("orderId"))

	var request PromotionCodeRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	err := o.orderRepo.AddProductToOrder(ctx, userId, orderId, request.ProductId, request.StoreId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

func (o *OrderHandler) DeletePromotionCodeFromOrder(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	orderId := cast.ToInt64(c.Params("orderId"))

	err := o.orderRepo.DeletePromotionCodeFromOrder(ctx, userId, orderId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}

func (o *OrderHandler) GetShippingMethod(c *fiber.Ctx) error {
	ctx := context.Background()

	shippingMethods, err := o.orderRepo.GetShippingMethod(ctx)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"shipping_methods": shippingMethods, "status": "ok"})
}

type UpdateShippingMethodRequest struct {
	ShippingMethodName string `json:"shipping_method_name"`
}

func (o *OrderHandler) UpdateShippingMethod(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	orderId := cast.ToInt64(c.Params("orderId"))

	var request UpdateShippingMethodRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	err := o.orderRepo.UpdateShippingMethod(ctx, userId, orderId, request.ShippingMethodName)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get "})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}
