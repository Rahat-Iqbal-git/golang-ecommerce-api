package repository

import (
	"github.com/rahat-iqbal/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type CartRepository interface {
	FindByUserAndProduct(userID, productID uint) (*models.CartItem, error)
	FindByIDAndUser(id, userID uint) (*models.CartItem, error)
	ListByUser(userID uint) ([]models.CartItem, error)
	Create(item *models.CartItem) error
	Save(item *models.CartItem) error
	Delete(item *models.CartItem) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) FindByUserAndProduct(userID, productID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error
	return &item, err
}

func (r *cartRepository) FindByIDAndUser(id, userID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Preload("Product").Where("id = ? AND user_id = ?", id, userID).First(&item).Error
	return &item, err
}

func (r *cartRepository) ListByUser(userID uint) ([]models.CartItem, error) {
	var items []models.CartItem
	err := r.db.Preload("Product").Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

func (r *cartRepository) Create(item *models.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartRepository) Save(item *models.CartItem) error {
	return r.db.Save(item).Error
}

func (r *cartRepository) Delete(item *models.CartItem) error {
	return r.db.Delete(item).Error
}
