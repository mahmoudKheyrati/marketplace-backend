package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/model"
)

type UserRepo interface {
	GetUserByUserId(ctx context.Context, userId int64) (*model.User, error)
}

type UserRepoImpl struct {
	db *pgxpool.Pool
}

func NewUserRepoImpl(db *pgxpool.Pool) *UserRepoImpl {
	return &UserRepoImpl{db: db}
}

func (u *UserRepoImpl) GetUserByUserId(ctx context.Context, userId int64) (*model.User, error) {
	query := `select id,
       email,
       phone_number,
       first_name,
       last_name,
       avatar_url,
       national_id,
       permission_name,
       created_at
from "user" where id = $1`

	rows, err := u.db.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.Id,
			&user.Email,
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
		return &user, nil
	}
	return nil, errors.New("user not found")
}
