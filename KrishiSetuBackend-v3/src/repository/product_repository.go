package repository

import (
	"fmt"
	"krishisetu-backend/models"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetFiltered(search, sort, minPrice, maxPrice, userPincode string) ([]models.Product, error) {
	query := r.db.Model(&models.Product{}).Preload("Seller").Where("is_available = ?", true)

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR category LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}

	switch sort {
	case "pincode_closest":
		if userPincode != "" {
			query = query.Order(fmt.Sprintf("ABS(CAST(location_pincode AS SIGNED) - %s) ASC", userPincode))
		}
	case "pincode_farthest":
		if userPincode != "" {
			query = query.Order(fmt.Sprintf("ABS(CAST(location_pincode AS SIGNED) - %s) DESC", userPincode))
		}
	case "price_low_high":
		query = query.Order("price ASC")
	case "price_high_low":
		query = query.Order("price DESC")
	default:
		query = query.Order("created_at DESC")
	}

	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetBySellerID(sellerID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Seller").Where("seller_id = ?", sellerID).Order("created_at desc").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindByIDWithSeller(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Seller").First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	// Cascade delete associated marketplace requests
	r.db.Where("product_id = ?", id).Delete(&models.MarketplaceRequest{})
	return r.db.Delete(&models.Product{}, id).Error
}
