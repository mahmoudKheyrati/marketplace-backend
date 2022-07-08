package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type AddressRepo interface {
	CreateAddress(ctx context.Context, userId int64, country, province, city, street, postalCode, homePhoneNumber string) (int64, error)
	GetUserAddressesByUserId(ctx context.Context, userId int64) ([]model.Address, error)
	UpdateUserAddress(ctx context.Context, userId int64, addressId int64, country, province, city, street, postalCode, homePhoneNumber string) (int64, error)
	DeleteUserAddress(ctx context.Context, userId int64, addressId int64) error
}
type AddressRepoImpl struct {
	db *pgxpool.Pool
}

func (a *AddressRepoImpl) CreateAddress(ctx context.Context, userId int64, country, province, city, street, postalCode, homePhoneNumber string) (int64, error) {
	query := `insert into address(user_id, country, province, city, street, postal_code, home_phone_number)
values ($1, $2, $3, $4, $5, $6, $7) returning id`

	rows, err := a.db.Query(ctx, query, userId, country, province, city, street, postalCode, homePhoneNumber)
	if err != nil {
		return -1, err
	}

	var id int64 = -1
	err = rows.Scan(&id)
	return id, err
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
			address.Id,
			address.UserId,
			address.Country,
			address.Province,
			address.City,
			address.Street,
			address.PostalCode,
			address.HomePhoneNumber,
			address.CreatedAt)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func (a *AddressRepoImpl) UpdateUserAddress(ctx context.Context, userId int64, addressId int64, country, province, city, street, postalCode, homePhoneNumber string) (int64, error) {
	query := `
begin;
update address
set is_last_version = false
where id = $1;

insert into address(user_id, country, province, city, street, postal_code, home_phone_number)
values ($2, $3, $4, $5, $6, $7, $8);
commit;
`
	rows, err := a.db.Query(ctx, query, addressId, userId, country, province, city, street, postalCode, homePhoneNumber)
	if err != nil {
		return -1, err
	}

	var id int64 = -1
	err = rows.Scan(&id)
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
