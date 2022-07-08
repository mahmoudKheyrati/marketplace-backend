package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/mahmoudKheyrati/marketplace-backend/api"
	"github.com/mahmoudKheyrati/marketplace-backend/config"
	"github.com/mahmoudKheyrati/marketplace-backend/internal/repository"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"go.uber.org/zap"
	"log"
)

func main() {
	//ctx := context.Background()
	zapLogger := pkg.Logger()
	defer func(logger *zap.SugaredLogger) {
		err := zapLogger.Sync()
		if err != nil {
			log.Fatalln(err)
		}
	}(zapLogger) // flush logs if any

	cfg := config.NewConfig()
	db := pkg.CreateNewPostgresConnection(pkg.PostgresConfig{
		Host:          cfg.Postgres.Host,
		Port:          cfg.Postgres.Port,
		Username:      cfg.Postgres.Username,
		Password:      cfg.Postgres.Password,
		Database:      cfg.Postgres.Database,
		MaxConnection: cfg.Postgres.MaxConnection,
	})
	defer db.Close()
	// create repositories
	authRepo := repository.NewAuthRepoImpl(db)
	authHandler := api.NewAuthHandler(authRepo)
	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	apiRoute := app.Group("/api")
	v2 := apiRoute.Group("/v2")
	auth := v2.Group("/auth")
	{
		auth.Post("/token", authHandler.Login)
		auth.Post("/signup", authHandler.Signup)
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Port)))
}
