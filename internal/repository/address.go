package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type AddressRepo interface {
	CreateAddress(ctx context.Context, userId int64, country, province, city, street, postalCode, homePhoneNumber string) (int64, error)
	GetAddressById(ctx context.Context, userId, addressId int64) (*model.Address, error)
	GetUserAddressesByUserId(ctx context.Context, userId int64) ([]model.Address, error)
	UpdateUserAddress(ctx context.Context, userId int64, addressId int64, country, province, city, street, postalCode, homePhoneNumber string) (int64, error)
	DeleteUserAddress(ctx context.Context, userId int64, addressId int64) error
}
type AddressRepoImpl struct {
	db *pgxpool.Pool
}

func NewAddressRepoImpl(db *pgxpool.Pool) *AddressRepoImpl {
	return &AddressRepoImpl{db: db}
}

func (a *AddressRepoImpl) CreateAddress(ctx context.Context, userId int64, country, province, city, street, postalCode, homePhoneNumber string) (int64, error) {
	query := `insert into address(user_id, country, province, city, street, postal_code, home_phone_number)
values ($1, $2, $3, $4, $5, $6, $7) returning id`

	rows, err := a.db.Query(ctx, query, userId, country, province, city, street, postalCode, homePhoneNumber)
	if err != nil {
		return -1, err
	}

	rows.Next()
	var id int64 = -1
	err = rows.Scan(&id)
	return id, err
}

func (a *AddressRepoImpl) GetAddressById(ctx context.Context, userId, addressId int64) (*model.Address, error) {
	query := `select id,
       user_id,
       country,
       province,
       city,
       street,
       postal_code,
       home_phone_number,
       created_at
from address
where user_id = $1 and id = $2
  and is_last_version = true limit 1;
`
	row := a.db.QueryRow(ctx, query, userId, addressId)
	var address model.Address
	err := row.Scan(
		&address.Id,
		&address.UserId,
		&address.Country,
		&address.Province,
		&address.City,
		&address.Street,
		&address.PostalCode,
		&address.HomePhoneNumber,
		&address.CreatedAt)
	return &address, err
}

func (a *AddressRepoImpl) GetUserAddressesByUserId(ctx context.Context, userId int64) ([]model.Address, error) {
	query := `select id,
       user_id,
       country,
       province,
       city,
       street,
       postal_code,
       home_phone_number,
       created_at
from address
where user_id = $1
  and is_last_version = true;
`
	rows, err := a.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	var addresses = make([]model.Address, 0)
	for rows.Next() {
		var address model.Address
		err := rows.Scan(
			&address.Id,
			&address.UserId,
			&address.Country,
			&address.Province,
			&address.City,
			&address.Street,
			&address.PostalCode,
			&address.HomePhoneNumber,
			&address.CreatedAt)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func (a *AddressRepoImpl) UpdateUserAddress(ctx context.Context, userId int64, addressId int64, country, province, city, street, postalCode, homePhoneNumber string) (int64, error) {

	var id int64 = -1
	err := a.db.BeginTxFunc(ctx, pgx.TxOptions{}, func(tx pgx.Tx) error {
		query := `
update address set is_last_version = false where id = $1; 

`
		_, err := tx.Exec(ctx, query, addressId)
		if err != nil {
			return err
		}
		query = `
insert into address(user_id, country, province, city, street, postal_code, home_phone_number)
values ($1, $2, $3, $4, $5, $6, $7) returning id
`

		_, err = tx.Exec(ctx, query, userId, country, province, city, street, postalCode, homePhoneNumber)
		if err != nil {
			return err
		}

		return err

	})
	if err != nil {
		return id, err
	}
	query := `select id from address where user_id= $1 and country = $2 and province = $3 and city = $4 and street=$5 and postal_code = $6 and home_phone_number = $7 and is_last_version = true`
	err = a.db.QueryRow(ctx, query, userId, country, province, city, street, postalCode, homePhoneNumber).Scan(&id)

	return id, err
}

func (a *AddressRepoImpl) DeleteUserAddress(ctx context.Context, userId int64, addressId int64) error {
	query := `
update address
set is_last_version = false
where user_id= $1 and id = $2 ;
`
	_, err := a.db.Query(ctx, query, userId, addressId)
	return err
}
