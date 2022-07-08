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
	"github.com/mahmoudKheyrati/marketplace-backend/pkg/middleware"
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
	notificationRepo := repository.NewNotificationRepoImpl(db)

	authHandler := api.NewAuthHandler(authRepo, cfg)
	notificationHandler := api.NewNotificationHandler(notificationRepo)

	authMiddleware := middleware.NewAuthMiddleware(cfg.JwtSecret)

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
	notification := v2.Group("/notifications", authMiddleware.Protected())
	{
		notification.Get("/", notificationHandler.GetAvailableNotifications)
		notification.Get("/pending", notificationHandler.GetPendingNotifications)
		notification.Post("/subscribe/:productId", notificationHandler.SubscribeToProduct)
		notification.Post("/seen/:productId", notificationHandler.SeenNotification)
	}
	v2.Get("/test", authMiddleware.Protected(), func(c *fiber.Ctx) error {
		data := c.Locals(pkg.JwtDataKey).(api.JwtData)
		fmt.Println("api:::: ", data)
		return c.SendStatus(fiber.StatusAccepted)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Port)))
}
