package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type ProductRepo interface {
	GetProductByProductId(ctx context.Context, productId int64) (*model.Product, error)
	GetProductsByCategoryId(ctx context.Context, categoryId int64) ([]model.Product, error)
	GetAllStoreProductsByProductId(ctx context.Context, productId int64) ([]model.StoreProduct, error) // product available in the stores

	GetSimilarProducts(ctx context.Context, productId int64) ([]model.Product, error)
	GetFrequentlyBoughtTogetherProducts(ctx context.Context, productId int64) ([]model.Product, error)

	GetBrandsByCategoryId(ctx context.Context, categoryId int64) ([]string, error)
	GetPriceRangeByCategoryId(ctx context.Context, categoryId int64) (float64, float64, error)
	GetSpecificationsByCategoryId(ctx context.Context, categoryId int64) ([]string, error)

	//SearchProductByName(ctx context.Context, name string) ([]model.Product, error)
	//SearchProductByBrand(ctx context.Context, brand string) ([]model.Product, error)
	//SearchProductByCategory(ctx context.Context, parentCategoryId int64) ([]model.Product, error)

	//FilterByBrand(ctx context.Context, brand string) ([]model.Product, error)
	//FilterByPrice(ctx context.Context, categoryId int64, min, max float64) ([]model.Product, error)
	//FilterBySpecifications(ctx context.Context, categoryId int64, specificationKey string, specificationValue interface{}) ([]model.Product, error)

	//SortByPriceCheapestToMostExpensive(ctx context.Context, categoryId int64) ([]model.Product, error)
	//SortByPriceMostExpensiveToCheapest(ctx context.Context, categoryId int64) ([]model.Product, error)
	//SortByRatingHighestToLowest(ctx context.Context, categoryId int64) ([]model.Product, error)
	//SortByRecentlyAdded(ctx context.Context, categoryId int64) ([]model.Product, error)
}
type ProductRepoImpl struct {
	db *pgxpool.Pool
}

func NewProductRepoImpl(db *pgxpool.Pool) *ProductRepoImpl {
	return &ProductRepoImpl{db: db}
}

func (p *ProductRepoImpl) GetProductByProductId(ctx context.Context, productId int64) (*model.Product, error) {
	query := `
select id,
       category_id,
       name,
       brand,
       description,
       picture_url,
       specification
from product
where id = $1
`
	row := p.db.QueryRow(ctx, query, productId)
	var product model.Product
	err := row.Scan(&product.Id, &product.CategoryId, &product.Name, &product.Brand, &product.Description, &product.PictureUrl, &product.Specification)
	return &product, err
}

func (p *ProductRepoImpl) GetProductsByCategoryId(ctx context.Context, categoryId int64) ([]model.Product, error) {
	query := `
select id,
       category_id,
       name,
       brand,
       description,
       picture_url,
       specification,
       created_at
from product
where category_id in (
    with recursive cte as (
        select id, name, parent
        from category
        where id = $1
        union
        select c2.id, c2.name, c2.parent
        from category c2
                 join cte as c1 on c2.parent = c1.id
    )
    select id
    from cte
) order by category_id;
`
	rows, err := p.db.Query(ctx, query, categoryId)
	if err != nil {
		return nil, err
	}
	var products = make([]model.Product, 0)
	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.Id, &product.CategoryId, &product.Name, &product.Brand, &product.Description, &product.PictureUrl, &product.Specification, &product.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)

	}
	return products, nil
}

func (p *ProductRepoImpl) GetAllStoreProductsByProductId(ctx context.Context, productId int64) ([]model.StoreProduct, error) {
	query := `
select product_id,
       store_id,
       off_percent,
       max_off_price,
       price,
       available_count,
       store_product.created_at, 
       w.name as warranty_name, 
       w.type as warranty_type, 
       w.month as warranty_month
from store_product
left join warranty w on store_product.warranty_id = w.id
where product_id = $1
  and available_count > 0 and is_last_version = true;
`
	rows, err := p.db.Query(ctx, query, productId)
	if err != nil {
		return nil, err
	}
	var storeProducts = make([]model.StoreProduct, 0)
	for rows.Next() {
		var storeProduct model.StoreProduct
		var warrantyName, warrantyType *string
		var warrantyMonth *int
		err := rows.Scan(&storeProduct.ProductId, &storeProduct.StoreId, &storeProduct.OffPercent, &storeProduct.MaxOffPrice, &storeProduct.Price,
			&storeProduct.AvailableCount, &storeProduct.CreatedAt, &warrantyName, &warrantyType, &warrantyMonth)
		if err != nil {
			return nil, err
		}
		if warrantyName != nil && warrantyType != nil && warrantyMonth != nil {
			storeProduct.Warranty = &model.Warranty{
				Name:         *warrantyName,
				WarrantyType: *warrantyType,
				Month:        *warrantyMonth,
			}
		}
		storeProducts = append(storeProducts, storeProduct)
	}
	return storeProducts, nil
}

func (p *ProductRepoImpl) GetSimilarProducts(ctx context.Context, productId int64) ([]model.Product, error) {
	product, err := p.GetProductByProductId(ctx, productId)
	if err != nil {
		return nil, err
	}

	query := `
with cte as (
    select product.id,
           category_id,
           name,
           brand,
           description,
           picture_url,
           specification,
           product.created_at,
           sum(2*(case when v.up_vote = true then 1 else 0 end) + (case when v.down_vote = true then 1 else 0 end)) as sum
    from product
             left join review r on product.id = r.product_id left join votes v on r.id = v.review_id
    where category_id = $1
    group by product.id, category_id, name, brand, description, picture_url, specification, product.created_at
)select id,
        category_id,
        name,
        brand,
        description,
        picture_url,
        specification,
        created_at
from cte order by sum desc limit 10;
`
	rows, err := p.db.Query(ctx, query, product.CategoryId)
	if err != nil {
		return nil, err
	}
	var products = make([]model.Product, 0)
	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.Id, &product.CategoryId, &product.Name, &product.Brand, &product.Description, &product.PictureUrl,
			&product.Specification, &product.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *ProductRepoImpl) GetFrequentlyBoughtTogetherProducts(ctx context.Context, productId int64) ([]model.Product, error) {
	product, err := p.GetProductByProductId(ctx, productId)
	if err != nil {
		return nil, err
	}

	query := `
with cte as (
    select p.id,
           category_id,
           name,
           brand,
           description,
           picture_url,
           specification,
           p.created_at,
           count(po.quantity) as co

    from product p
             left join product_order po on p.id = po.product_id
             left join "order" o on o.id = po.order_id

    where p.id =$1 and o.is_paid = true and payed_price != 0 and pay_date is not null
    group by p.id, category_id, name, brand, description, picture_url, specification, p.created_at

)select id,
        category_id,
        name,
        brand,
        description,
        picture_url,
        specification,
        created_at
from cte order by co desc limit 10;
`
	rows, err := p.db.Query(ctx, query, product.CategoryId)
	if err != nil {
		return nil, err
	}
	var products = make([]model.Product, 0)
	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.Id, &product.CategoryId, &product.Name, &product.Brand, &product.Description, &product.PictureUrl,
			&product.Specification, &product.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (p *ProductRepoImpl) GetBrandsByCategoryId(ctx context.Context, categoryId int64) ([]string, error) {
	query := `select distinct brand from product where category_id = $1`
	rows, err := p.db.Query(ctx, query, categoryId)
	if err != nil {
		return nil, err
	}
	var brands = make([]string, 0)
	for rows.Next() {
		var brand string
		err := rows.Scan(&brand)
		if err != nil {
			return nil, err
		}
		brands = append(brands, brand)
	}
	return brands, nil
}

func (p *ProductRepoImpl) GetPriceRangeByCategoryId(ctx context.Context, categoryId int64) (float64, float64, error) {
	query := `
select min(sp.price) as min, max(sp.price) as max from store_product sp join product p on p.id = sp.product_id
where p.category_id = $1
`
	row := p.db.QueryRow(ctx, query, categoryId)
	var min, max float64
	err := row.Scan(&min, &max)
	return min, max, err
}

func (p *ProductRepoImpl) GetSpecificationsByCategoryId(ctx context.Context, categoryId int64) ([]string, error) {
	query := `
select distinct (jsonb_object_keys(p.specification))
from product p
         join category c on p.category_id = c.id
where p.category_id = $1
`
	rows, err := p.db.Query(ctx, query, categoryId)
	if err != nil {
		return nil, err
	}
	var specificationKeys = make([]string, 0)
	for rows.Next() {
		var specificationKey string
		err := rows.Scan(&specificationKey)
		if err != nil {
			return nil, err
		}
		specificationKeys = append(specificationKeys, specificationKey)
	}
	return specificationKeys, nil
}
