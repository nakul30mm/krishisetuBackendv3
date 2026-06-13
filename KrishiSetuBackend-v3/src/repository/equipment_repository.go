package repository

import (
	"fmt"
	"krishisetu-backend/models"
	"gorm.io/gorm"
)

type equipmentRepository struct {
	db *gorm.DB
}

func NewEquipmentRepository(db *gorm.DB) *equipmentRepository {
	return &equipmentRepository{db: db}
}

func (r *equipmentRepository) Create(equipment *models.Equipment) error {
	return r.db.Create(equipment).Error
}

func (r *equipmentRepository) GetFiltered(search, sort, minPrice, maxPrice, userPincode string) ([]models.Equipment, error) {
	query := r.db.Model(&models.Equipment{}).Preload("Owner")

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR category LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if minPrice != "" {
		query = query.Where("price_per_day >= ?", minPrice)
	}
	if maxPrice != "" {
		query = query.Where("price_per_day <= ?", maxPrice)
	}

	// Support multiple sort aliases for better UI compatibility
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
		query = query.Order("price_per_day ASC")
	case "price_high_low":
		query = query.Order("price_per_day DESC")
	default:
		query = query.Order("created_at DESC")
	}

	var equipments []models.Equipment
	if err := query.Find(&equipments).Error; err != nil {
		return nil, err
	}
	return equipments, nil
}

func (r *equipmentRepository) GetByOwnerID(ownerID uint) ([]models.Equipment, error) {
	var equipments []models.Equipment
	err := r.db.Preload("Owner").
		Where("owner_id = ?", ownerID).
		Find(&equipments).Error
	if err != nil {
		return nil, err
	}
	return equipments, nil
}

func (r *equipmentRepository) FindByID(id uint) (*models.Equipment, error) {
	var equipment models.Equipment
	if err := r.db.First(&equipment, id).Error; err != nil {
		return nil, err
	}
	return &equipment, nil
}

func (r *equipmentRepository) FindByIDWithOwner(id uint) (*models.Equipment, error) {
	var equipment models.Equipment
	if err := r.db.Preload("Owner").First(&equipment, id).Error; err != nil {
		return nil, err
	}
	return &equipment, nil
}

func (r *equipmentRepository) Update(equipment *models.Equipment) error {
	return r.db.Save(equipment).Error
}

func (r *equipmentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Equipment{}, id).Error
}

func (r *equipmentRepository) CountActiveRentals(equipmentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Rental{}).
		Where("equipment_id = ? AND status = ?", equipmentID, "APPROVED").
		Count(&count).Error
	return count, err
}

func (r *equipmentRepository) GetApprovedRentals(equipmentID uint) ([]models.Rental, error) {
	var rentals []models.Rental
	err := r.db.Where("equipment_id = ? AND status = ?", equipmentID, "APPROVED").Find(&rentals).Error
	return rentals, err
}

func (r *equipmentRepository) GetLatestApprovedRental(equipmentID uint) (*models.Rental, error) {
	var rental models.Rental
	err := r.db.Where("equipment_id = ? AND status = ?", equipmentID, "APPROVED").
		Order("end_date desc").
		First(&rental).Error
	if err != nil {
		return nil, err
	}
	return &rental, nil
}
