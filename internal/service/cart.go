package service

import (
	"errors"

	"github.com/rahat-iqbal/ecommerce-api/internal/models"
	"github.com/rahat-iqbal/ecommerce-api/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrCartItemNotFound  = errors.New("cart item not found")
)

type CartService interface {
	AddItem(userID, productID uint, quantity int) (*models.CartItem, error)
	ListItems(userID uint) ([]models.CartItem, error)
	UpdateItem(id, userID uint, quantity int) (*models.CartItem, error)
	RemoveItem(id, userID uint) error
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{cartRepo: cartRepo, productRepo: productRepo}
}

func (s *cartService) AddItem(userID, productID uint, quantity int) (*models.CartItem, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	if product.Stock < quantity {
		return nil, ErrInsufficientStock
	}

	item, err := s.cartRepo.FindByUserAndProduct(userID, productID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		item = &models.CartItem{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		}
		if err := s.cartRepo.Create(item); err != nil {
			return nil, err
		}
	} else {
		item.Quantity += quantity
		if err := s.cartRepo.Save(item); err != nil {
			return nil, err
		}
	}

	item.Product = *product
	return item, nil
}

func (s *cartService) ListItems(userID uint) ([]models.CartItem, error) {
	return s.cartRepo.ListByUser(userID)
}

func (s *cartService) UpdateItem(id, userID uint, quantity int) (*models.CartItem, error) {
	item, err := s.cartRepo.FindByIDAndUser(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartItemNotFound
		}
		return nil, err
	}

	item.Quantity = quantity
	if err := s.cartRepo.Save(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *cartService) RemoveItem(id, userID uint) error {
	item, err := s.cartRepo.FindByIDAndUser(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartItemNotFound
		}
		return err
	}
	return s.cartRepo.Delete(item)
}
