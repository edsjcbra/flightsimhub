package routes

import (
	"github.com/edsjcbra/flightsimhub/internal/controllers"
	"github.com/edsjcbra/flightsimhub/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "FlightSimHub API is running 🚀",
			})
		})

		authService := services.NewAuthService()
		authController := controllers.NewAuthController(authService)

		api.POST("/auth/signup", authController.Signup)
		api.POST("/auth/login", authController.Login)
	}
}
