package routes

import (
	"github.com/edsjcbra/flightsimhub/internal/controllers"
	"github.com/edsjcbra/flightsimhub/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	authController *controllers.AuthController,
	productController *controllers.ProductController,
	cartController *controllers.CartController,
	orderController *controllers.OrderController,
) {
	api := router.Group("/api/v1")

	// health
	api.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	// auth
	auth := api.Group("/auth")
	auth.POST("/signup", authController.Signup)
	auth.POST("/login", authController.Login)

	// products - public read, protected write
	products := api.Group("/products")
	products.GET("/", productController.GetAllProducts)
	products.GET("/:id", productController.GetProduct)

	protectedProducts := api.Group("/products")
	protectedProducts.Use(middlewares.JWTAuthMiddleware())
	{
		protectedProducts.POST("/", productController.CreateProduct)
		protectedProducts.PUT("/:id", productController.UpdateProduct)
		protectedProducts.DELETE("/:id", productController.DeleteProduct)
	}

	// cart (protected)
	cart := api.Group("/cart")
	cart.Use(middlewares.JWTAuthMiddleware())
	{
		cart.GET("/", cartController.GetCart)
		cart.POST("/add", cartController.AddToCart)
		cart.POST("/remove", cartController.RemoveFromCart)
	}

	// orders (protected)
	orders := api.Group("/orders")
	orders.Use(middlewares.JWTAuthMiddleware())
	{
		orders.GET("/", orderController.GetOrders)
		orders.POST("/", orderController.CreateOrder)
	}
}
