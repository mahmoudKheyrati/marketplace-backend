package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo interface {
	Authenticate(ctx context.Context, email string, password string) (*model.User, error)
	SignUp(ctx context.Context, email, password, phoneNumber, firstName, lastName string) error
}

type AuthRepoImpl struct {
	db *pgxpool.Pool
}

func NewAuthRepoImpl(db *pgxpool.Pool) *AuthRepoImpl {
	return &AuthRepoImpl{db: db}
}

func (a *AuthRepoImpl) Authenticate(ctx context.Context, email string, password string) (*model.User, error) {
	var user *model.User
	rows, err := a.db.Query(ctx, `select id,
       email,
       password,
       phone_number,
       first_name,
       last_name,
       avatar_url,
       national_id,
       permission_name,
       created_at
from "user" where email = $1`, email)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		user = &model.User{}
		err = rows.Scan(
			&user.Id,
			&user.Email,
			&user.Password,
			&user.PhoneNumber,
			&user.FirstName,
			&user.LastName,
			&user.AvatarUrl,
			&user.NationalId,
			&user.PermissionName,
			&user.CreatedAt)
		if err != nil {
			return nil, err
		}
	}
	if user == nil {
		return nil, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *AuthRepoImpl) SignUp(ctx context.Context, email, password, phoneNumber, firstName, lastName string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return err
	}
	_, err = a.db.Query(ctx, `insert into "user" (email, password, phone_number, first_name, last_name, permission_name)
values ($1, $2, $3, $4, $5, 'normal-user')`,
		email, string(passwordHash), phoneNumber, firstName, lastName)
	return err
}
