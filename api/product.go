package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"github.com/spf13/cast"
)

type ProductHandler struct {
	productRepo repository.ProductRepo
}

func NewProductHandler(productRepo repository.ProductRepo) *ProductHandler {
	return &ProductHandler{productRepo: productRepo}
}

func (p *ProductHandler) GetProductByProductId(c *fiber.Ctx) error {
	ctx := context.Background()
	productId := cast.ToInt64(c.Params("productId"))

	product, err := p.productRepo.GetProductByProductId(ctx, productId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get product by productId"})
	}
	return c.JSON(fiber.Map{"product": product, "status": "ok"})
}

func (p *ProductHandler) GetSimilarProducts(c *fiber.Ctx) error {
	ctx := context.Background()
	productId := cast.ToInt64(c.Params("productId"))

	product, err := p.productRepo.GetSimilarProducts(ctx, productId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get product by productId"})
	}
	return c.JSON(fiber.Map{"product": product, "status": "ok"})
}

func (p *ProductHandler) GetProductsByCategoryId(c *fiber.Ctx) error {
	ctx := context.Background()
	categoryId := cast.ToInt64(c.Params("categoryId"))

	products, err := p.productRepo.GetProductsByCategoryId(ctx, categoryId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get products by categoryId"})
	}
	return c.JSON(fiber.Map{"products": products, "status": "ok"})
}

func (p *ProductHandler) GetAllStoreProductsByProductId(c *fiber.Ctx) error {
	ctx := context.Background()
	productId := cast.ToInt64(c.Params("productId"))

	storeProducts, err := p.productRepo.GetAllStoreProductsByProductId(ctx, productId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get storeProduct by productId"})
	}
	return c.JSON(fiber.Map{"store_products": storeProducts, "status": "ok"})
}

func (p *ProductHandler) GetBrandsByCategoryId(c *fiber.Ctx) error {
	ctx := context.Background()
	categoryId := cast.ToInt64(c.Params("categoryId"))

	brands, err := p.productRepo.GetBrandsByCategoryId(ctx, categoryId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get brands by categoryId"})
	}
	return c.JSON(fiber.Map{"brands": brands, "status": "ok"})
}

func (p *ProductHandler) GetPriceRangeByCategoryId(c *fiber.Ctx) error {
	ctx := context.Background()
	categoryId := cast.ToInt64(c.Params("categoryId"))

	min, max, err := p.productRepo.GetPriceRangeByCategoryId(ctx, categoryId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get price range by categoryId"})
	}
	return c.JSON(fiber.Map{"min": min, "max": max, "status": "ok"})
}

func (p *ProductHandler) GetSpecificationsByCategoryId(c *fiber.Ctx) error {
	ctx := context.Background()
	categoryId := cast.ToInt64(c.Params("categoryId"))

	specificationKeys, err := p.productRepo.GetSpecificationsByCategoryId(ctx, categoryId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get specifications by categoryId"})
	}
	return c.JSON(fiber.Map{"specification_keys": specificationKeys, "status": "ok"})
}
