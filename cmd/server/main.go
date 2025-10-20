package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/edsjcbra/flightsimhub/config"
	"github.com/edsjcbra/flightsimhub/internal/controllers"
	"github.com/edsjcbra/flightsimhub/internal/database"
	"github.com/edsjcbra/flightsimhub/internal/routes"
	"github.com/edsjcbra/flightsimhub/internal/services"
)

func main() {
	config.LoadConfig()
	database.Connect()
	defer database.Close()

	db := database.Pool

	// services
	authService := services.NewAuthService(db)
	productService := services.NewProductService(db)
	cartService := services.NewCartService(db)
	orderService := services.NewOrderService(db)

	// controllers
	authController := controllers.NewAuthController(authService)
	productController := controllers.NewProductController(productService)
	cartController := controllers.NewCartController(cartService)
	orderController := controllers.NewOrderController(orderService)

	router := gin.Default()
	routes.RegisterRoutes(router, authController, productController, cartController, orderController)

	port := config.AppConfig.Port
	log.Printf("🌍 Server running on http://localhost:%s\n", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
