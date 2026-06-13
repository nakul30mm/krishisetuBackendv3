package repository

import (
	"krishisetu-backend/models"
	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *reviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) FindByID(id uint) (*models.Review, error) {
	var review models.Review
	if err := r.db.First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *reviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}

func (r *reviewRepository) GetByEquipmentID(equipmentID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("equipment_id = ?", equipmentID).Order("created_at desc").Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetByProductID(productID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("product_id = ?", productID).Order("created_at desc").Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetByRentalID(rentalID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("rental_id = ?", rentalID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetByMarketplaceRequestID(marketplaceRequestID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("marketplace_request_id = ?", marketplaceRequestID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) UpdateEquipmentRating(equipmentID uint) error {
	var reviews []models.Review
	err := r.db.Where("equipment_id = ?", equipmentID).Find(&reviews).Error
	if err != nil {
		return err
	}

	if len(reviews) == 0 {
		return r.db.Model(&models.Equipment{}).Where("id = ?", equipmentID).Updates(map[string]interface{}{
			"average_rating": 0,
			"total_reviews":  0,
		}).Error
	}

	var total int
	for _, rev := range reviews {
		total += rev.Rating
	}
	avg := float64(total) / float64(len(reviews))

	return r.db.Model(&models.Equipment{}).Where("id = ?", equipmentID).Updates(map[string]interface{}{
		"average_rating": avg,
		"total_reviews":  len(reviews),
	}).Error
}

func (r *reviewRepository) UpdateProductRating(productID uint) error {
	var reviews []models.Review
	err := r.db.Where("product_id = ?", productID).Find(&reviews).Error
	if err != nil {
		return err
	}

	if len(reviews) == 0 {
		return r.db.Model(&models.Product{}).Where("id = ?", productID).Updates(map[string]interface{}{
			"average_rating": 0,
			"total_reviews":  0,
		}).Error
	}

	var total int
	for _, rev := range reviews {
		total += rev.Rating
	}
	avg := float64(total) / float64(len(reviews))

	return r.db.Model(&models.Product{}).Where("id = ?", productID).Updates(map[string]interface{}{
		"average_rating": avg,
		"total_reviews":  len(reviews),
	}).Error
}

func (r *reviewRepository) UpdateOwnerRating(ownerID uint) error {
	var rReviews []models.Review
	err := r.db.Joins("JOIN equipments ON equipments.id = reviews.equipment_id").
		Where("equipments.owner_id = ?", ownerID).Find(&rReviews).Error
	if err != nil {
		return err
	}

	var pReviews []models.Review
	err = r.db.Joins("JOIN products ON products.id = reviews.product_id").
		Where("products.seller_id = ?", ownerID).Find(&pReviews).Error
	if err != nil {
		return err
	}

	reviews := append(rReviews, pReviews...)

	if len(reviews) == 0 {
		return r.db.Model(&models.User{}).Where("id = ?", ownerID).Updates(map[string]interface{}{
			"average_rating": 0,
			"total_ratings":  0,
		}).Error
	}

	var total int
	for _, rev := range reviews {
		total += rev.Rating
	}
	avg := float64(total) / float64(len(reviews))

	return r.db.Model(&models.User{}).Where("id = ?", ownerID).Updates(map[string]interface{}{
		"average_rating": avg,
		"total_ratings":  len(reviews),
	}).Error
}
