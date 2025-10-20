package routes

import (
	"github.com/edsjcbra/flightsimhub/internal/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra todas as rotas da aplicação
func RegisterRoutes(
	router *gin.Engine,
	authController *controllers.AuthController,
	productController *controllers.ProductController,
	cartController *controllers.CartController,
	orderController *controllers.OrderController,
) {
	api := router.Group("/api/v1")
	{
		// ------------------ AUTH ------------------
		auth := api.Group("/auth")
		auth.POST("/signup", authController.Signup)
		auth.POST("/login", authController.Login)

		// ------------------ PRODUCTS ------------------
		products := api.Group("/products")
		products.POST("/", productController.CreateProduct)
		products.GET("/", productController.GetAllProducts)
		products.GET("/:id", productController.GetProduct) // <-- corrigido para GetProduct
		products.PUT("/:id", productController.UpdateProduct)
		products.DELETE("/:id", productController.DeleteProduct)

		// ------------------ CART ------------------
		cart := api.Group("/cart")
		cart.GET("/", cartController.GetCart)
		cart.POST("/add", cartController.AddToCart)
		cart.POST("/remove", cartController.RemoveFromCart)

		// ------------------ ORDERS ------------------
		orders := api.Group("/orders")
		orders.GET("/", orderController.GetOrders)
		orders.POST("/", orderController.CreateOrder)

		// ------------------ HEALTH ------------------
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}
}
