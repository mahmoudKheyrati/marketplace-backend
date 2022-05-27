package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v4"
	"github.com/mahmoudKheyrati/marketplace-backend/config"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"go.uber.org/zap"
	"log"
)

func main() {
	zapLogger := pkg.NewLogger()
	defer func(logger *zap.SugaredLogger) {
		err := zapLogger.Sync()
		if err != nil {
			log.Fatalln(err)
		}
	}(zapLogger) // flush logs if any

	cfg := config.NewConfig()

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())

	_, _ = pgx.Connect(context.Background(), "")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Port)))
}
