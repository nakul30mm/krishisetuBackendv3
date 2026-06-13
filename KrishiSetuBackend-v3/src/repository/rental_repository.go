package repository

import (
	"krishisetu-backend/models"
	"gorm.io/gorm"
)

type rentalRepository struct {
	db *gorm.DB
}

func NewRentalRepository(db *gorm.DB) *rentalRepository {
	return &rentalRepository{db: db}
}

func (r *rentalRepository) Create(rental *models.Rental) error {
	return r.db.Create(rental).Error
}

func (r *rentalRepository) FindByID(id uint) (*models.Rental, error) {
	var rental models.Rental
	if err := r.db.First(&rental, id).Error; err != nil {
		return nil, err
	}
	return &rental, nil
}

func (r *rentalRepository) FindByIDWithEquipmentAndOwner(id uint) (*models.Rental, error) {
	var rental models.Rental
	if err := r.db.Preload("Equipment").Preload("Equipment.Owner").First(&rental, id).Error; err != nil {
		return nil, err
	}
	return &rental, nil
}

func (r *rentalRepository) Update(rental *models.Rental) error {
	return r.db.Save(rental).Error
}

func (r *rentalRepository) Delete(id uint) error {
	return r.db.Delete(&models.Rental{}, id).Error
}

func (r *rentalRepository) GetByRenterID(renterID uint) ([]models.Rental, error) {
	var rentals []models.Rental
	err := r.db.Preload("Equipment").
		Preload("Equipment.Owner").
		Where("renter_id = ?", renterID).
		Order("created_at desc").
		Find(&rentals).Error
	return rentals, err
}

func (r *rentalRepository) GetByOwnerID(ownerID uint) ([]models.Rental, error) {
	var rentals []models.Rental
	err := r.db.Preload("Equipment").
		Preload("Equipment.Owner").
		Preload("Renter").
		Where("owner_id = ?", ownerID).
		Order("created_at desc").
		Find(&rentals).Error
	return rentals, err
}

func (r *rentalRepository) GetApprovedRentalsForEquipment(equipmentID uint) ([]models.Rental, error) {
	var rentals []models.Rental
	err := r.db.Where("equipment_id = ? AND status = ?", equipmentID, "APPROVED").Find(&rentals).Error
	return rentals, err
}

func (r *rentalRepository) GetApprovedRentalsForEquipmentExcluding(equipmentID uint, excludingID uint) ([]models.Rental, error) {
	var rentals []models.Rental
	err := r.db.Where("equipment_id = ? AND status = ? AND id != ?", equipmentID, "APPROVED", excludingID).Find(&rentals).Error
	return rentals, err
}

func (r *rentalRepository) CountApprovedRentalsForEquipment(equipmentID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Rental{}).
		Where("equipment_id = ? AND status = ?", equipmentID, "APPROVED").
		Count(&count).Error
	return count, err
}

func (r *rentalRepository) RejectOverlappingPendingRentals(equipmentID uint, ignoreRentalID uint, startDate, endDate string) error {
	// Rejects any pending rentals that overlap with the approved rental period
	// Status becomes "REJECTED"
	return r.db.Model(&models.Rental{}).
		Where(`
			equipment_id = ?
			AND status = ?
			AND id != ?
			AND NOT (
				end_date < ? OR start_date > ?
			)
		`,
			equipmentID,
			"PENDING",
			ignoreRentalID,
			startDate,
			endDate,
		).
		Update("status", "REJECTED").Error
}

func (r *rentalRepository) GetAllApprovedRentals() ([]models.Rental, error) {
	var rentals []models.Rental
	err := r.db.Where("status = ?", "APPROVED").Find(&rentals).Error
	return rentals, err
}

func (r *rentalRepository) UpdateEquipmentStatus(equipmentID uint, status string) error {
	return r.db.Model(&models.Equipment{}).Where("id = ?", equipmentID).Update("status", status).Error
}

func (r *rentalRepository) GetReviewForRental(rentalID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("rental_id = ?", rentalID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}
