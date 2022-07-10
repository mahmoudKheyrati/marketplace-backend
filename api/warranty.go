package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"github.com/spf13/cast"
)

type WarrantyHandler struct {
	warrantyRepo repository.WarrantyRepo
}

func NewWarrantyHandler(warrantyRepo repository.WarrantyRepo) *WarrantyHandler {
	return &WarrantyHandler{warrantyRepo: warrantyRepo}
}

type WarrantyRequest struct {
	Name         string `json:"name"`
	WarrantyType string `json:"warranty_type"`
	Month        int    `json:"month"`
}

func (w *WarrantyHandler) CreateWarranty(c *fiber.Ctx) error {
	ctx := context.Background()

	permissionName := c.Locals(pkg.UserPermissionNameKey).(string)
	if permissionName != "store-admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "you don't have access"})
	}
	var request WarrantyRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	warrantyId, err := w.warrantyRepo.CreateWarranty(ctx, request.Name, request.WarrantyType, request.Month)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "acn not create warranty"})
	}
	return c.JSON(fiber.Map{"warranty_id": warrantyId, "status": "ok"})
}

func (w *WarrantyHandler) GetWarrantyByWarrantyId(c *fiber.Ctx) error {
	ctx := context.Background()
	warrantyId := cast.ToInt64(c.Params("warrantyId"))

	warranty, err := w.warrantyRepo.GetWarrantyByWarrantyId(ctx, warrantyId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get warranty by warrantyId"})
	}
	return c.JSON(fiber.Map{"warranty": warranty, "status": "ok"})
}

func (w *WarrantyHandler) GetStoreProductWarranty(c *fiber.Ctx) error {
	ctx := context.Background()
	productId := cast.ToInt64(c.Params("productId"))
	storeId := cast.ToInt64(c.Params("storeId"))

	warranty, err := w.warrantyRepo.GetStoreProductWarranty(ctx, storeId, productId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get warranty by productId and storeId"})
	}
	return c.JSON(fiber.Map{"warranty": warranty, "status": "ok"})
}
