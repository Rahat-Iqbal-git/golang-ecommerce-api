package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rahat-iqbal/ecommerce-api/internal/service"
)

type CartHandler struct {
	svc service.CartService
}

func NewCartHandler(svc service.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

func userIDFromContext(c *gin.Context) (uint, bool) {
	raw, exists := c.Get("claims")
	if !exists {
		return 0, false
	}
	claims, ok := raw.(jwt.MapClaims)
	if !ok {
		return 0, false
	}
	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, false
	}
	return uint(sub), true
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return uint(id), true
}

type addToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

func (h *CartHandler) AddItem(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	var req addToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.svc.AddItem(userID, req.ProductID, req.Quantity)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrProductNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrInsufficientStock):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add to cart"})
		}
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *CartHandler) ListItems(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	items, err := h.svc.ListItems(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch cart"})
		return
	}
	c.JSON(http.StatusOK, items)
}

type updateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	id, ok := parseID(c)
	if !ok {
		return
	}

	var req updateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.svc.UpdateItem(id, userID, req.Quantity)
	if err != nil {
		if errors.Is(err, service.ErrCartItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update cart item"})
		}
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.svc.RemoveItem(id, userID); err != nil {
		if errors.Is(err, service.ErrCartItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not remove cart item"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}
