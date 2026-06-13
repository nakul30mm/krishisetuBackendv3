package repository

import (
	"krishisetu-backend/models"
	"gorm.io/gorm"
)

type marketplaceRequestRepository struct {
	db *gorm.DB
}

func NewMarketplaceRequestRepository(db *gorm.DB) *marketplaceRequestRepository {
	return &marketplaceRequestRepository{db: db}
}

func (r *marketplaceRequestRepository) Create(req *models.MarketplaceRequest) error {
	return r.db.Create(req).Error
}

func (r *marketplaceRequestRepository) FindByID(id uint) (*models.MarketplaceRequest, error) {
	var req models.MarketplaceRequest
	if err := r.db.First(&req, id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *marketplaceRequestRepository) FindByIDWithProductAndSeller(id uint) (*models.MarketplaceRequest, error) {
	var req models.MarketplaceRequest
	if err := r.db.Preload("Product").Preload("Seller").First(&req, id).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *marketplaceRequestRepository) Update(req *models.MarketplaceRequest) error {
	return r.db.Save(req).Error
}

func (r *marketplaceRequestRepository) Delete(id uint) error {
	return r.db.Delete(&models.MarketplaceRequest{}, id).Error
}

func (r *marketplaceRequestRepository) GetByBuyerID(buyerID uint, search string) ([]models.MarketplaceRequest, error) {
	var requests []models.MarketplaceRequest

	query := r.db.Model(&models.MarketplaceRequest{}).
		Where("buyer_id = ?", buyerID).
		Preload("Product").
		Preload("Seller")

	if search != "" {
		query = query.Joins("JOIN products ON products.id = marketplace_requests.product_id").
			Where("products.name LIKE ?", "%"+search+"%")
	}

	err := query.Order("marketplace_requests.created_at desc").Find(&requests).Error
	return requests, err
}

func (r *marketplaceRequestRepository) GetBySellerID(sellerID uint, search string) ([]models.MarketplaceRequest, error) {
	var requests []models.MarketplaceRequest

	query := r.db.Model(&models.MarketplaceRequest{}).
		Where("seller_id = ?", sellerID).
		Preload("Product").
		Preload("Buyer")

	if search != "" {
		query = query.Joins("JOIN products ON products.id = marketplace_requests.product_id").
			Where("products.name LIKE ?", "%"+search+"%")
	}

	err := query.Order("marketplace_requests.created_at desc").Find(&requests).Error
	return requests, err
}
