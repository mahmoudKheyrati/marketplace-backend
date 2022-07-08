package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
)

type UserHandler struct {
	userRepo repository.UserRepo
}

func NewUserHandler(userRepo repository.UserRepo) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (u *UserHandler) GetMyProfile(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)

	user, err := u.userRepo.GetUserByUserId(ctx, userId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "can not get user profile."})
	}
	return c.JSON(fiber.Map{"user": user})
}
