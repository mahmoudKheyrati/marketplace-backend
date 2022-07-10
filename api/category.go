package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"github.com/spf13/cast"
)

type CategoryHandler struct {
	categoryRepo repository.CategoryRepo
}

func NewCategoryHandler(categoryRepo repository.CategoryRepo) *CategoryHandler {
	return &CategoryHandler{categoryRepo: categoryRepo}
}

func (ca *CategoryHandler) GetMainCategories(c *fiber.Ctx) error {
	ctx := context.Background()

	categories, err := ca.categoryRepo.GetMainCategories(ctx)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get main categories"})
	}
	return c.JSON(fiber.Map{"categories": categories})
}

func (ca *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	ctx := context.Background()

	categories, err := ca.categoryRepo.GetAllCategories(ctx)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get all categories"})
	}
	return c.JSON(fiber.Map{"categories": categories})
}

func (ca *CategoryHandler) GetSubCategoriesByCategoryId(c *fiber.Ctx) error {
	ctx := context.Background()

	categoryId := cast.ToInt64(c.Params("categoryId"))

	categories, err := ca.categoryRepo.GetSubCategoriesByCategoryId(ctx, categoryId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get categories by categoryId"})
	}
	return c.JSON(fiber.Map{"categories": categories})
}

func (ca *CategoryHandler) GetParentsByCategoryId(c *fiber.Ctx) error {
	ctx := context.Background()

	categoryId := cast.ToInt64(c.Params("categoryId"))

	categories, err := ca.categoryRepo.GetParentsByCategoryId(ctx, categoryId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get categories by categoryId"})
	}
	return c.JSON(fiber.Map{"categories": categories})
}
