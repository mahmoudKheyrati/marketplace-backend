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
	"github.com/mahmoudKheyrati/marketplace-backend/pkg/metric"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg/middleware"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg/prometheus"
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

	metrics := metric.GetMetrics()
	prom := prometheus.NewPrometheus(3000)
	go prom.RunHTTPServer()

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
	categoryRepo := repository.NewCategoryRepoImpl(db)
	productRepo := repository.NewProductRepoImpl(db)
	warrantyRepo := repository.NewWarrantyRepoImpl(db)
	storeRepo := repository.NewStoreRepoImpl(db)
	orderRepo := repository.NewOrderRepoImpl(db)

	// create api handlers
	authHandler := api.NewAuthHandler(authRepo, cfg)
	notificationHandler := api.NewNotificationHandler(notificationRepo)
	ticketHandler := api.NewTicketHandler(ticketRepo)
	userHandler := api.NewUserHandler(userRepo)
	addressHandler := api.NewAddressHandler(addressRepo)
	voteHandler := api.NewVoteHandler(voteRepo)
	reviewHandler := api.NewReviewHandler(reviewRepo)
	categoryHandler := api.NewCategoryHandler(categoryRepo)
	productHandler := api.NewProductHandler(productRepo)
	warrantyHandler := api.NewWarrantyHandler(warrantyRepo)
	storeHandler := api.NewStoreHandler(storeRepo)
	orderHandler := api.NewOrderHandler(orderRepo)

	// create middlewares
	authMiddleware := middleware.NewAuthMiddleware(cfg.JwtSecret)
	metricMiddleware := middleware.NewMetricsMiddleware(metrics)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(metricMiddleware.MetricsMiddleware)
	//app.Use(metricsMiddleware.MetricsMiddleware)
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
		notification.Post("/subscribe/:productId", notificationHandler.SubscribeToProductNotification)
		notification.Delete("/unsubscribe/:productId", notificationHandler.UnSubscribeToProductNotification)
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
		users.Get("/:userId", userHandler.GetUserByUserId)
	}
	addresses := v2.Group("addresses", authMiddleware.Protected())
	{
		addresses.Get("/", addressHandler.GetAllAddresses)
		addresses.Get("/:addressId", addressHandler.GetAddressByAddressId)
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
	categories := v2.Group("/categories")
	{
		categories.Get("/", categoryHandler.GetAllCategories)
		categories.Get("/main", categoryHandler.GetMainCategories)
		categories.Get("/subs/:categoryId", categoryHandler.GetSubCategoriesByCategoryId)
		categories.Get("/parents/:categoryId", categoryHandler.GetParentsByCategoryId)
	}
	products := v2.Group("/products")
	{
		products.Get("/:productId", productHandler.GetProductByProductId)
		products.Get("/similar/:productId", productHandler.GetSimilarProducts)
		products.Get("/frequently_bought_together/:productId", productHandler.GetFrequentlyBoughtTogetherProducts)
		products.Get("/category/:categoryId", productHandler.GetProductsByCategoryId)
		products.Get("/stores/:productId", productHandler.GetAllStoreProductsByProductId)
		products.Get("/brands/:categoryId", productHandler.GetBrandsByCategoryId)
		products.Get("/price_range/:categoryId", productHandler.GetPriceRangeByCategoryId)
		products.Get("/specifications/:categoryId", productHandler.GetSpecificationsByCategoryId)
	}
	warranty := v2.Group("/warranty")
	{
		warranty.Post("/create", authMiddleware.Protected(), warrantyHandler.CreateWarranty)
		warranty.Get("/:warrantyId", warrantyHandler.GetWarrantyByWarrantyId)
		warranty.Get("/product/:productId/:storeId", warrantyHandler.GetStoreProductWarranty)
	}
	stores := v2.Group("/stores", authMiddleware.Protected())
	{
		stores.Post("/create", storeHandler.CreateStore)
		stores.Post("/update/:storeId", storeHandler.UpdateStore)
		stores.Delete("/delete/:storeId", storeHandler.DeleteStore)
		stores.Get("/me", storeHandler.GetMyStores)
		stores.Get("/:storeId", storeHandler.GetStoreByStoreId)
		stores.Get("/", storeHandler.GetAllStores)
		stores.Get("/products/:storeId", storeHandler.GetAllProductsByStoreId)
		stores.Post("/products/:storeId/create", storeHandler.AddProductToStore)
		stores.Post("/products/:storeId/:addressId/update", storeHandler.UpdateStoreProduct)
		stores.Post("/addresses/:storeId/create", storeHandler.AddStoreAddress)
		stores.Get("/addresses/:storeId", storeHandler.GetStoreAddressesByStoreId)
		stores.Post("/addresses/:storeId/update", storeHandler.UpdateStoreAddresses)
		stores.Post("/categories/:storeId/:categoryId/create", storeHandler.AddStoreCategory)
		stores.Delete("/categories/:storeId/:categoryId/delete", storeHandler.DeleteStoreCategory)
	}
	orders := v2.Group("/orders", authMiddleware.Protected())
	{
		orders.Get("/check_payment_status/:orderId", orderHandler.IsUserPaidTheOrder)
		orders.Post("/create", orderHandler.CreateOrder)
		orders.Delete("/delete/:orderId", orderHandler.DeleteOrder)
		orders.Post("/add_product/:orderId", orderHandler.AddProductToOrder)
		orders.Delete("/remove_product/:orderId/:storeId/:productId", orderHandler.RemoveProductFromOrder)
		orders.Post("/update_quantity/:orderId", orderHandler.UpdateProductOrderQuantity)
		orders.Get("/:orderId", orderHandler.GetAllProductsInTheOrder)
		orders.Get("/me", orderHandler.GetAllOrdersByUserId)
		orders.Post("/pay/:orderId", orderHandler.PayOrder)
		orders.Post("/promotion_code/apply/:orderId", orderHandler.ApplyPromotionCodeToOrder)
		orders.Delete("/promotion_code/delete/:orderId", orderHandler.DeletePromotionCodeFromOrder)
		orders.Get("/shipping_methods", orderHandler.GetShippingMethod)
		orders.Post("/shipping_methods/update/:orderId", orderHandler.UpdateShippingMethod)
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Port)))
}
