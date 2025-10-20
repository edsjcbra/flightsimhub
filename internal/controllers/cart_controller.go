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

func (c *CartController) GetCart(ctx *gin.Context) {
	// Exemplo: retornar carrinho do usuário (stub temporário)
	ctx.JSON(http.StatusOK, gin.H{"cart": "stub"})
}

func (c *CartController) AddToCart(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "item added (stub)"})
}

func (c *CartController) RemoveFromCart(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "item removed (stub)"})
}
