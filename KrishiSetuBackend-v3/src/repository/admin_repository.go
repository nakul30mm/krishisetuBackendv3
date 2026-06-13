package repository

import (
	"krishisetu-backend/models"

	"gorm.io/gorm"
)

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *adminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) GetStats() (int64, int64, int64, int64, error) {
	var userCount, productCount, equipmentCount, requestCount int64

	if err := r.db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	if err := r.db.Model(&models.Product{}).Count(&productCount).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	if err := r.db.Model(&models.Equipment{}).Count(&equipmentCount).Error; err != nil {
		return 0, 0, 0, 0, err
	}
	if err := r.db.Model(&models.MarketplaceRequest{}).Count(&requestCount).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	return userCount, productCount, equipmentCount, requestCount, nil
}
