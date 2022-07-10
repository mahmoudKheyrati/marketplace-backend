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
	ticketRepo := repository.NewTicketRepoImpl(db)
	userRepo := repository.NewUserRepoImpl(db)
	addressRepo := repository.NewAddressRepoImpl(db)
	voteRepo := repository.NewVoteRepoImpl(db)
	reviewRepo := repository.NewReviewRepoImpl(db)

	authHandler := api.NewAuthHandler(authRepo, cfg)
	notificationHandler := api.NewNotificationHandler(notificationRepo)
	ticketHandler := api.NewTicketHandler(ticketRepo)
	userHandler := api.NewUserHandler(userRepo)
	addressHandler := api.NewAddressHandler(addressRepo)
	voteHandler := api.NewVoteHandler(voteRepo)
	reviewHandler := api.NewReviewHandler(reviewRepo)

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

	tickets := v2.Group("/tickets", authMiddleware.Protected())
	{
		tickets.Get("/", ticketHandler.GetAllTickets)
		tickets.Get("/types", ticketHandler.GetAllTicketTypes)
		tickets.Post("/create/:ticketTypeId", ticketHandler.CreateTicket)
		//tickets.Post("/done/:ticketId")
		tickets.Post("/send_message/:ticketId", ticketHandler.SendMessageToTicket)
		//tickets.Post("/received/:ticketId/:messageId")
		//tickets.Post("/seen/:ticketId/:messageId")
		tickets.Get("/:ticketId", ticketHandler.LoadTicketMessages)
		tickets.Get("/unfinished", ticketHandler.GetAllUnfinishedTickets)
	}
	users := v2.Group("/users", authMiddleware.Protected())
	{
		users.Get("/me", userHandler.GetMyProfile)
	}
	addresses := v2.Group("addresses", authMiddleware.Protected())
	{
		addresses.Get("/", addressHandler.GetAllAddresses)
		addresses.Post("/create", addressHandler.CreateAddress)
		addresses.Post("/update/:addressId", addressHandler.UpdateAddress)
		addresses.Delete("/delete/:addressId", addressHandler.DeleteAddress)
	}
	votes := v2.Group("/votes", authMiddleware.Protected())
	{
		votes.Post("/create", voteHandler.CreateVote)
		votes.Delete("/delete/:reviewId", voteHandler.DeleteVote)
	}
	reviews := v2.Group("/reviews", authMiddleware.Protected())
	{
		reviews.Post("/create", reviewHandler.CreateReview)
		reviews.Post("/update/:reviewId", reviewHandler.UpdateReview)
		reviews.Delete("/delete/:reviewId", reviewHandler.DeleteReview)
		reviews.Get("/me", reviewHandler.GetUserAllReviews)
		reviews.Get("/product/:productId", reviewHandler.GetProductReviews)
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Port)))
}
