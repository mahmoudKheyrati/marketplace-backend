package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"github.com/spf13/cast"
)

type StoreHandler struct {
	storeRepo repository.StoreRepo
}

func NewStoreHandler(storeRepo repository.StoreRepo) *StoreHandler {
	return &StoreHandler{storeRepo: storeRepo}
}

type StoreRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AvatarUrl   string `json:"avatar_url"`
	Owner       int64  `json:"owner"`
	Creator     int64  `json:"creator"`
}

func (s *StoreHandler) CreateStore(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)

	permissionName := c.Locals(pkg.UserPermissionNameKey).(string)
	if permissionName != "store-admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "you don't have access"})
	}
	var request StoreRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	storeId, err := s.storeRepo.CreateStore(ctx, request.Name, request.Description, request.AvatarUrl, userId, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"store_id": storeId, "status": "ok"})
}

func (s *StoreHandler) UpdateStore(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	storeId := cast.ToInt64(c.Params("storeId"))

	permissionName := c.Locals(pkg.UserPermissionNameKey).(string)
	if permissionName != "store-admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "you don't have access"})
	}
	var request StoreRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	storeId, err := s.storeRepo.UpdateStore(ctx, request.Name, request.Description, request.AvatarUrl, userId, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"store_id": storeId, "status": "ok"})
}

func (s *StoreHandler) DeleteStore(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	storeId := cast.ToInt64(c.Params("storeId"))

	err := s.storeRepo.DeleteStore(ctx, userId, storeId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"status": "ok"})

}

func (s *StoreHandler) GetMyStores(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)

	stores, err := s.storeRepo.GetStoresByUserId(ctx, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"stores": stores, "status": "ok"})
}

func (s *StoreHandler) GetStoreByStoreId(c *fiber.Ctx) error {
	ctx := context.Background()

	storeId := cast.ToInt64(c.Params("storeId"))

	store, err := s.storeRepo.GetStoreByStoreId(ctx, storeId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"store": store, "status": "ok"})

}

func (s *StoreHandler) GetAllStores(c *fiber.Ctx) error {
	ctx := context.Background()

	store, err := s.storeRepo.GetAllStores(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"stores": store, "status": "ok"})
}

func (s *StoreHandler) GetAllProductsByStoreId(c *fiber.Ctx) error {
	ctx := context.Background()

	storeId := cast.ToInt64(c.Params("storeId"))

	storeProducts, err := s.storeRepo.GetAllProductsByStoreId(ctx, storeId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"store_products": storeProducts, "status": "ok"})

}

type StoreAddressRequest struct {
	Country    string `json:"country"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Street     string `json:"street"`
	PostalCode string `json:"postal_code"`
}

func (s *StoreHandler) AddStoreAddress(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	storeId := cast.ToInt64(c.Params("storeId"))

	permissionName := c.Locals(pkg.UserPermissionNameKey).(string)
	if permissionName != "store-admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "you don't have access"})
	}
	var request StoreAddressRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	storeAddressId, err := s.storeRepo.AddStoreAddress(ctx, userId, storeId, request.Country, request.Province, request.City, request.Street, request.PostalCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"store_address_id": storeAddressId, "status": "ok"})

}

func (s *StoreHandler) GetStoreAddressesByStoreId(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	storeId := cast.ToInt64(c.Params("storeId"))

	permissionName := c.Locals(pkg.UserPermissionNameKey).(string)
	if permissionName != "store-admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "you don't have access"})
	}
	storeAddresses, err := s.storeRepo.GetStoreAddressesByStoreId(ctx, userId, storeId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"store_addresses": storeAddresses, "status": "ok"})

}

func (s *StoreHandler) UpdateStoreAddresses(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	storeId := cast.ToInt64(c.Params("storeId"))

	permissionName := c.Locals(pkg.UserPermissionNameKey).(string)
	if permissionName != "store-admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "you don't have access"})
	}
	var request StoreAddressRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	storeAddressId, err := s.storeRepo.UpdateStoreAddresses(ctx, userId, storeId, request.Country, request.Province, request.City, request.Street, request.PostalCode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"store_address_id": storeAddressId, "status": "ok"})
}

func (s *StoreHandler) AddStoreCategory(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	storeId := cast.ToInt64(c.Params("storeId"))
	categoryId := cast.ToInt64(c.Params("categoryId"))

	permissionName := c.Locals(pkg.UserPermissionNameKey).(string)
	if permissionName != "store-admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "you don't have access"})
	}

	err := s.storeRepo.AddStoreCategory(ctx, userId, storeId, categoryId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"status": "ok"})

}

func (s *StoreHandler) DeleteStoreCategory(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	storeId := cast.ToInt64(c.Params("storeId"))
	categoryId := cast.ToInt64(c.Params("categoryId"))

	permissionName := c.Locals(pkg.UserPermissionNameKey).(string)
	if permissionName != "store-admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "you don't have access"})
	}

	err := s.storeRepo.DeleteStoreCategory(ctx, userId, storeId, categoryId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"status": "ok"})

}

type StoreProductRequest struct {
	ProductId      int64   `json:"product_id"`
	WarrantyId     int64   `json:"warranty_id"`
	Price          float64 `json:"price"`
	OffPercent     float64 `json:"off_percent"`
	MaxOffPrice    float64 `json:"max_off_price"`
	AvailableCount int     `json:"available_count"`
}

func (s *StoreHandler) AddProductToStore(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	storeId := cast.ToInt64(c.Params("storeId"))

	permissionName := c.Locals(pkg.UserPermissionNameKey).(string)
	if permissionName != "store-admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "you don't have access"})
	}
	var request StoreProductRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	err := s.storeRepo.AddProductToStore(ctx, userId, storeId, request.ProductId, request.WarrantyId, request.Price, request.OffPercent, request.MaxOffPrice,
		request.AvailableCount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"status": "ok"})

}

func (s *StoreHandler) UpdateStoreProduct(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	storeId := cast.ToInt64(c.Params("storeId"))

	permissionName := c.Locals(pkg.UserPermissionNameKey).(string)
	if permissionName != "store-admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "you don't have access"})
	}
	var request StoreProductRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	err := s.storeRepo.UpdateStoreProduct(ctx, userId, storeId, request.ProductId, request.WarrantyId, request.Price, request.OffPercent, request.MaxOffPrice,
		request.AvailableCount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not "})
	}

	return c.JSON(fiber.Map{"status": "ok"})

}
