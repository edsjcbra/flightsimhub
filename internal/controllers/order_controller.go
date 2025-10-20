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

func (oc *OrderController) CreateOrder(c *gin.Context) {
	userID := c.GetInt("user_id")
	order, err := oc.Service.CreateOrder(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (oc *OrderController) GetOrders(c *gin.Context) {
	userID := c.GetInt("user_id")
	orders, err := oc.Service.GetOrders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}
