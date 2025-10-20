package controllers

import (
	"net/http"

	"github.com/edsjcbra/flightsimhub/internal/services"
	"github.com/gin-gonic/gin"
)

type CartController struct {
	Service *services.CartService
}

func NewCartController(s *services.CartService) *CartController {
	return &CartController{Service: s}
}

// Add item body: { "product_id": 1, "quantity": 2 }
type addItemReq struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func (cc *CartController) AddToCart(c *gin.Context) {
	var req addItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := cc.Service.AddItem(userID, req.ProductID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item added"})
}

func (cc *CartController) GetCart(c *gin.Context) {
	userID := c.GetInt("user_id")
	items, err := cc.Service.GetCartItems(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (cc *CartController) RemoveFromCart(c *gin.Context) {
	var req struct {
		ItemID int `json:"item_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := cc.Service.RemoveItem(req.ItemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}
