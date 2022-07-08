package api

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthHandler struct {
	authRepo repository.AuthRepo
}

func NewAuthHandler(authRepo repository.AuthRepo) *AuthHandler {
	return &AuthHandler{authRepo: authRepo}
}

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *AuthHandler) Login(c *fiber.Ctx) error {
	var request LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}

	user, err := a.authRepo.Authenticate(context.Background(), request.Email, request.Password)
	if err != nil {
		pkg.Logger().Error(err)
		return nil
	}

	if user == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "username or password are not valid!"})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.Id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("SECRET_SHOULD_READ_FROM_CONFIG"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}

type SignUpRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

func (a *AuthHandler) Signup(c *fiber.Ctx) error {
	var request SignUpRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "error in parsing request body"})
	}
	err := a.authRepo.SignUp(context.Background(), request.Email, request.Password, request.PhoneNumber, request.FirstName, request.LastName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "can not make user!. change your input"})
	}

	return c.JSON(fiber.Map{"status": "ok", "message": "success signup"})
}
