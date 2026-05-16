package repository

import (
	"github.com/rahat-iqbal/ecommerce-api/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindByID(id uint) (*models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return &product, err
}
