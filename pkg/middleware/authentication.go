package middleware

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mahmoudKheyrati/marketplace-backend/api"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

// Protected protect routes
func (a *AuthMiddleware) Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SuccessHandler: func(c *fiber.Ctx) error {
			// extract jwt data
			jwtData := api.JwtData{}

			token := c.Locals("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)

			raw := claims[pkg.JwtDataTokenKey].(string)
			err := json.Unmarshal([]byte(raw), &jwtData)
			if err != nil {
				pkg.Logger().Error(err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			c.Locals(pkg.JwtDataKey, jwtData)
			c.Locals(pkg.UserIdKey, jwtData.UserId)
			return c.Next()
		},
		ErrorHandler: jwtError,
		SigningKey:   []byte(a.jwtSecret),
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
