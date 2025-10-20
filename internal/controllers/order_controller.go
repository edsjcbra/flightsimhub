package controllers

import (
	"net/http"

	"github.com/edsjcbra/flightsimhub/internal/services"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	Service *services.OrderService
}

func NewOrderController(s *services.OrderService) *OrderController {
	return &OrderController{Service: s}
}

func (o *OrderController) GetOrders(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"orders": "stub"})
}

func (o *OrderController) CreateOrder(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "order created (stub)"})
}
