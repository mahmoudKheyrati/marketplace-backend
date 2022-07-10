package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type StoreRepo interface {
	GetStoreByStoreId(ctx context.Context, storeId int64) (*model.Store, error)
	CreateStore(ctx context.Context, name, description, avatarUrl string, owner, creator int64) (int64, error)
	UpdateStore(ctx context.Context, name, description, avatarUrl string, owner, storeId int64) (int64, error)
	DeleteStore(ctx context.Context, owner, storeId int64) error
	GetStoresByUserId(ctx context.Context, userId int64) ([]model.Store, error)
	GetAllStores(ctx context.Context) ([]model.Store, error)
	GetAllProductsByStoreId(ctx context.Context, storeId int64) ([]model.Product, error)

	IsUserOwnerOfStore(ctx context.Context, userId, storeId int64) (bool, error)
	// first check if the user is creator of store or not?
	AddStoreAddress(ctx context.Context, userId, storeId int64, country, province, city, street, postalCode string) (int64, error)
	// first check if the user is creator of store or not?
	GetStoreAddressesByUserId(ctx context.Context, userId, storeId int64) ([]model.StoreAddress, error)
	// first check if the user is creator of store or not?
	UpdateStoreAddresses(ctx context.Context, userId, storeId int64, country, province, city, street, postalCode string) (int64, error)

	// first check if the user is creator of store or not?
	AddStoreCategory(ctx context.Context, userId, storeId, categoryId int64) error
	// first check if the user is creator of store or not?
	DeleteStoreCategory(ctx context.Context, userId, storeId, categoryId int64) error

	// first check if the user is creator of store or not?
	AddProductToStore(ctx context.Context, userId, storeId, productId, warrantyId int64, price, offPercent, maxOffPrice float64, availableCount int) error
	// first check if the user is creator of store or not?
	UpdateStoreProduct(ctx context.Context, userId, storeId, productId, warrantyId int64, price, offPercent, maxOffPrice float64, availableCount int) error
}

type StoreRepoImpl struct {
	db *pgxpool.Pool
}

func NewStoreRepoImpl(db *pgxpool.Pool) *StoreRepoImpl {
	return &StoreRepoImpl{db: db}
}

func (s *StoreRepoImpl) GetStoreByStoreId(ctx context.Context, storeId int64) (*model.Store, error) {
	query := `
select id,
       name,
       description,
       avatar_url,
       owner,
       creator,
       created_at
from store where id = $1
`

	row := s.db.QueryRow(ctx, query, storeId)
	var store model.Store
	err := row.Scan(&store.Id, &store.Name, &store.Description, &store.AvatarUrl, &store.Owner, &store.Creator, &store.CreatedAt)
	return &store, err
}

func (s *StoreRepoImpl) CreateStore(ctx context.Context, name, description, avatarUrl string, owner, creator int64) (int64, error) {
	query := `
insert into store(name, description, avatar_url, owner, creator)
values ($1, $2, $3, $4, $5) returning id
`
	row := s.db.QueryRow(ctx, query, name, description, owner, creator)
	var id int64
	err := row.Scan(&id)
	return id, err
}

func (s *StoreRepoImpl) UpdateStore(ctx context.Context, name, description, avatarUrl string, owner, storeId int64) (int64, error) {
	query := `
update store
set name= $1,
    description = $2,
    avatar_url = $3
where id =$4?
  and owner = $5 
returning id 
`
	row := s.db.QueryRow(ctx, query, name, description, avatarUrl, storeId, owner)
	var id int64
	err := row.Scan(&id)
	return id, err
}

func (s *StoreRepoImpl) DeleteStore(ctx context.Context, owner, storeId int64) error {
	query := `
update store
set deleted_at = now()
where id = $1
  and owner = $2
`
	_, err := s.db.Exec(ctx, query, storeId, owner)
	return err
}

func (s *StoreRepoImpl) GetStoresByUserId(ctx context.Context, userId int64) ([]model.Store, error) {
	query := `
select id, name, description, avatar_url, owner, creator, created_at
from store
where creator = $1 or owner = $1
  and deleted_at is null
`
	rows, err := s.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	var stores = make([]model.Store, 0)
	for rows.Next() {
		var store model.Store
		err := rows.Scan(&store.Id, &store.Name, &store.Description, &store.AvatarUrl, &store.Owner, &store.Creator, &store.CreatedAt)
		if err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}
	return stores, nil

}

func (s *StoreRepoImpl) GetAllStores(ctx context.Context) ([]model.Store, error) {
	query := `
select id, name, description, avatar_url, owner, creator, created_at
from store
where deleted_at is null;
`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var stores = make([]model.Store, 0)
	for rows.Next() {
		var store model.Store
		err := rows.Scan(&store.Id, &store.Name, &store.Description, &store.AvatarUrl, &store.Owner, &store.Creator, &store.CreatedAt)
		if err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}
	return stores, nil
}

func (s *StoreRepoImpl) IsUserOwnerOfStore(ctx context.Context, userId, storeId int64) (bool, error) {
	query := `
select exists(select 1 from store where owner = $1 and id = $2)
`
	row := s.db.QueryRow(ctx, query, userId, storeId)
	var isOwner bool
	err := row.Scan(&isOwner)
	return isOwner, err
}
func (s *StoreRepoImpl) mustUserOwnerOfStore(ctx context.Context, userId, storeId int64) error {
	isUserOwnerOfStore, err := s.IsUserOwnerOfStore(ctx, userId, storeId)
	if err != nil {
		return err
	}
	if !isUserOwnerOfStore {
		return errors.New("you don't access to modify store")
	}
	return nil
}

func (s *StoreRepoImpl) GetAllProductsByStoreId(ctx context.Context, storeId int64) ([]model.Product, error) {
	query := `
select id,
       category_id,
       name,
       brand,
       description,
       picture_url,
       specification,
       p.created_at as created_at
from product p
         join store_product sp on p.id = sp.product_id
where id = $1
`
	rows, err := s.db.Query(ctx, query, storeId)
	if err != nil {
		return nil, err
	}
	var products = make([]model.Product, 0)
	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.Id, &product.CategoryId, &product.Name, &product.Brand, &product.Description, &product.PictureUrl, &product.Specification, &product.Specification)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (s *StoreRepoImpl) AddStoreAddress(ctx context.Context, userId, storeId int64, country, province, city, street, postalCode string) (int64, error) {
	err := s.mustUserOwnerOfStore(ctx, userId, storeId)
	if err != nil {
		return -1, err
	}
	query := `
insert into store_address (country, province, city, street, postal_code)
values ($1, $2, $3, $4, $5) returning id
`
	row := s.db.QueryRow(ctx, query, country, province, city, street, postalCode)
	var id int64 = -1
	err = row.Scan(&id)
	return id, err

}

func (s *StoreRepoImpl) GetStoreAddressesByUserId(ctx context.Context, userId, storeId int64) ([]model.StoreAddress, error) {
	err := s.mustUserOwnerOfStore(ctx, userId, storeId)
	if err != nil {
		return nil, err
	}

	query := `
select store_id, country, province, city, street, postal_code, created_at
from store_address
where store_id = $1;
`
	rows, err := s.db.Query(ctx, query, storeId)
	if err != nil {
		return nil, err
	}
	var storeAddresses = make([]model.StoreAddress, 0)
	for rows.Next() {
		var storeAddress model.StoreAddress
		err := rows.Scan(&storeAddress.StoreId, &storeAddress.Country, &storeAddress.Province, &storeAddress.City, &storeAddress.Street, &storeAddress.PostalCode, &storeAddress.CreatedAt)
		if err != nil {
			return nil, err
		}
		storeAddresses = append(storeAddresses, storeAddress)
	}
	return storeAddresses, nil
}

func (s *StoreRepoImpl) UpdateStoreAddresses(ctx context.Context, userId, storeId int64, country, province, city, street, postalCode string) (int64, error) {
	err := s.mustUserOwnerOfStore(ctx, userId, storeId)
	if err != nil {
		return -1, err
	}
	query := `
update store_address
set country = $1,
    province= $2,
    city= $3,
    street= $4,
    postal_code = $5
where store_id = $6
returning id
`
	row := s.db.QueryRow(ctx, query, country, province, city, street, postalCode, storeId)
	var id int64 = -1
	err = row.Scan(&id)
	return id, err

}

func (s *StoreRepoImpl) AddStoreCategory(ctx context.Context, userId, storeId, categoryId int64) error {
	err := s.mustUserOwnerOfStore(ctx, userId, storeId)
	if err != nil {
		return err
	}
	query := `
insert into store_category(category_id, store_id)
values ($1, $2)
`
	_, err = s.db.Exec(ctx, query, categoryId, storeId)
	return err
}

func (s *StoreRepoImpl) DeleteStoreCategory(ctx context.Context, userId, storeId, categoryId int64) error {
	err := s.mustUserOwnerOfStore(ctx, userId, storeId)
	if err != nil {
		return err
	}
	query := `
delete
from store_category
where category_id = $1
  and store_id = $2
`
	_, err = s.db.Exec(ctx, query, categoryId, storeId)
	return err

}

func (s *StoreRepoImpl) AddProductToStore(ctx context.Context, userId, storeId, productId, warrantyId int64, price, offPercent, maxOffPrice float64, availableCount int) error {
	err := s.mustUserOwnerOfStore(ctx, userId, storeId)
	if err != nil {
		return err
	}

	query := `
insert into store_product(product_id, store_id, off_percent, max_off_price, price, available_count, warranty_id)
values ($1, $2, $3, $4, $5, $6, $7) 
`
	_, err = s.db.Exec(ctx, query, productId, storeId, offPercent, maxOffPrice, price, availableCount, warrantyId)
	return err
}

func (s *StoreRepoImpl) UpdateStoreProduct(ctx context.Context, userId, storeId, productId, warrantyId int64, price, offPercent, maxOffPrice float64, availableCount int) error {
	err := s.mustUserOwnerOfStore(ctx, userId, storeId)
	if err != nil {
		return err
	}
	query := `
update store_product
set off_percent= $1,
    max_off_price= $2,
    price = $3,
    available_count = $4,
    warranty_id = $5
where store_id = $6
  and product_id = $7
`
	_, err = s.db.Exec(ctx, query, offPercent, maxOffPrice, price, availableCount, warrantyId, storeId, productId)
	return err
}
