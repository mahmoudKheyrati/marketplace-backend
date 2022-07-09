package repository

import (
	"context"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type ProductRepo interface {
	GetProductByProductId(ctx context.Context, productId int64) (*model.Product, error)
	GetAllStoreProductsByProductId(ctx context.Context, productId int64) ([]model.StoreProduct, error) // product available in the stores

	SearchProductByName(ctx context.Context, name string) ([]model.Product, error)
	SearchProductByBrand(ctx context.Context, brand string) ([]model.Product, error)
	SearchProductByCategory(ctx context.Context, parentCategoryId int64) ([]model.Product, error)

	GetBrandsByCategoryId(ctx context.Context, categoryId int64) ([]string, error)
	FilterByBrand(ctx context.Context, brand string) ([]model.Product, error)

	GetPriceRangeByCategoryId(ctx context.Context, categoryId int64) (float64, float64, error)
	FilterByPrice(ctx context.Context, categoryId int64, min, max float64) ([]model.Product, error)

	GetSpecificationsByCategoryId(ctx context.Context, categoryId int64) ([]string, error)
	FilterBySpecifications(ctx context.Context, categoryId int64, specificationKey string, specificationValue interface{}) ([]model.Product, error)

	SortByPriceCheapestToMostExpensive(ctx context.Context, categoryId int64) ([]model.Product, error)
	SortByPriceMostExpensiveToCheapest(ctx context.Context, categoryId int64) ([]model.Product, error)

	SortByRatingHighestToLowest(ctx context.Context, categoryId int64) ([]model.Product, error)

	SortByRecentlyAdded(ctx context.Context, categoryId int64) ([]model.Product, error)
}
