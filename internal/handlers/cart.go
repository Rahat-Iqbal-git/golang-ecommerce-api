package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rahat-iqbal/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type CartHandler struct {
	db *gorm.DB
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{db: db}
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

	var product models.Product
	if err := h.db.First(&product, req.ProductID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch product"})
		return
	}

	if product.Stock < req.Quantity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "insufficient stock"})
		return
	}

	var item models.CartItem
	err := h.db.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&item).Error
	if err == nil {
		item.Quantity += req.Quantity
		if err := h.db.Save(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update cart"})
			return
		}
	} else if err == gorm.ErrRecordNotFound {
		item = models.CartItem{
			UserID:    userID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := h.db.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add to cart"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not check cart"})
		return
	}

	h.db.Preload("Product").First(&item, item.ID)
	c.JSON(http.StatusOK, item)
}

func (h *CartHandler) ListItems(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	var items []models.CartItem
	if err := h.db.Preload("Product").Where("user_id = ?", userID).Find(&items).Error; err != nil {
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

	var item models.CartItem
	if err := h.db.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch cart item"})
		return
	}

	var req updateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.Quantity = req.Quantity
	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update cart item"})
		return
	}

	h.db.Preload("Product").First(&item, item.ID)
	c.JSON(http.StatusOK, item)
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	var item models.CartItem
	if err := h.db.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "cart item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch cart item"})
		return
	}

	if err := h.db.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not remove cart item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}
