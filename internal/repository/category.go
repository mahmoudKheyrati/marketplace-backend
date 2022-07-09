package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type CategoryRepo interface {
	GetMainCategories(ctx context.Context) ([]model.Category, error)
	GetAllCategories(ctx context.Context) ([]model.Category, error)
	GetSubCategoriesByCategoryId(ctx context.Context, categoryId int64) ([]model.Category, error)
	GetParentsByCategoryId(ctx context.Context, categoryId int64) ([]model.Category, error)
}

type CategoryRepoImpl struct {
	db *pgxpool.Pool
}

func (c *CategoryRepoImpl) GetMainCategories(ctx context.Context) ([]model.Category, error) {
	query := `select id, name, parent
from category
where parent is null`
	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var categories = make([]model.Category, 0)
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Parent)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *CategoryRepoImpl) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	query := `select id, name, parent from category order by id, parent`
	rows, err := c.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var categories = make([]model.Category, 0)
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Parent)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *CategoryRepoImpl) GetSubCategoriesByCategoryId(ctx context.Context, categoryId int64) ([]model.Category, error) {
	query := `
with recursive cte as (
    select id,name,parent from category where id = $1
    union
    select c2.id,c2.name,c2.parent from category c2 join cte as c1 on c2.parent=c1.id
) select id, name, parent from cte
`
	rows, err := c.db.Query(ctx, query, categoryId)
	if err != nil {
		return nil, err
	}
	var categories = make([]model.Category, 0)
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Parent)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *CategoryRepoImpl) GetParentsByCategoryId(ctx context.Context, categoryId int64) ([]model.Category, error) {
	query := `
with recursive cte as (
    select id,name,parent from category where id = $1
    union
    select c2.id,c2.name,c2.parent from category c2 join cte as c1 on c2.id=c1.parent
) select id, name, parent from cte
`
	rows, err := c.db.Query(ctx, query, categoryId)
	if err != nil {
		return nil, err
	}
	var categories = make([]model.Category, 0)
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.Id, &category.Name, &category.Parent)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
