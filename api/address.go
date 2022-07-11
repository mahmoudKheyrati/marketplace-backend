package api

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"github.com/spf13/cast"
)

type AddressHandler struct {
	addressRepo repository.AddressRepo
}

func NewAddressHandler(addressRepo repository.AddressRepo) *AddressHandler {
	return &AddressHandler{addressRepo: addressRepo}
}

func (a *AddressHandler) GetAllAddresses(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)

	id, err := a.addressRepo.GetUserAddressesByUserId(ctx, userId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get user addresses."})
	}
	return c.JSON(fiber.Map{"status": "ok", "addresses": id})
}

func (a *AddressHandler) GetAddressByAddressId(c *fiber.Ctx) error {
	ctx := context.Background()

	userId := c.Locals(pkg.UserIdKey).(int64)
	addressId := cast.ToInt64(c.Params("addressId"))
	fmt.Println(userId, addressId)

	id, err := a.addressRepo.GetAddressById(ctx, userId, addressId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not get user address by addressId."})
	}
	return c.JSON(fiber.Map{"status": "ok", "address": id})
}

type AddressRequest struct {
	Country         string `json:"country"`
	Province        string `json:"province"`
	City            string `json:"city"`
	Street          string `json:"street"`
	PostalCode      string `json:"postal_code"`
	HomePhoneNumber string `json:"home_phone_number"`
}

func (a *AddressHandler) CreateAddress(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)

	var request AddressRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	addressId, err := a.addressRepo.CreateAddress(ctx, userId, request.Country, request.Province, request.City, request.Street, request.PostalCode, request.HomePhoneNumber)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not create address"})
	}
	return c.JSON(fiber.Map{"status": "ok", "address_id": addressId})

}

func (a *AddressHandler) UpdateAddress(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)
	addressId := cast.ToInt64(c.Params("addressId"))

	var request AddressRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	newAddressId, err := a.addressRepo.UpdateUserAddress(ctx, userId, addressId, request.Country, request.Province, request.City, request.Street, request.PostalCode, request.HomePhoneNumber)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not update address"})
	}
	return c.JSON(fiber.Map{"status": "ok", "new_address_id": newAddressId})

}

func (a *AddressHandler) DeleteAddress(c *fiber.Ctx) error {
	ctx := context.Background()
	userId := c.Locals(pkg.UserIdKey).(int64)
	addressId := cast.ToInt64(c.Params("addressId"))

	err := a.addressRepo.DeleteUserAddress(ctx, userId, addressId)
	if err != nil {
		pkg.Logger().Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "can not delete address"})
	}
	return c.JSON(fiber.Map{"status": "ok"})

}
