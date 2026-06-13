package business

import (
	"errors"
	"time"

	"krishisetu-backend/models"
	"krishisetu-backend/src/dto"
)

type ReviewService struct {
	reviewRepo             ReviewRepository
	rentalRepo             RentalRepository
	equipmentRepo          EquipmentRepository
	productRepo            ProductRepository
	marketplaceRequestRepo MarketplaceRequestRepository
	userRepo               UserRepository
}

func NewReviewService(
	reviewRepo ReviewRepository,
	rentalRepo RentalRepository,
	equipmentRepo EquipmentRepository,
	productRepo ProductRepository,
	marketplaceRequestRepo MarketplaceRequestRepository,
	userRepo UserRepository,
) *ReviewService {
	return &ReviewService{
		reviewRepo:             reviewRepo,
		rentalRepo:             rentalRepo,
		equipmentRepo:          equipmentRepo,
		productRepo:            productRepo,
		marketplaceRequestRepo: marketplaceRequestRepo,
		userRepo:               userRepo,
	}
}

func (s *ReviewService) CreateReview(userID uint, req dto.CreateReviewDTO) (*models.Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	var review models.Review

	// Handle equipment-only review (no rental or marketplace request)
	if req.EquipmentID != nil && req.RentalID == nil && req.MarketRequestID == nil {
		// Ensure equipment exists
		if _, err := s.equipmentRepo.FindByID(*req.EquipmentID); err != nil {
			return nil, errors.New("equipment not found")
		}
		// Create review for equipment
		review = models.Review{
			EquipmentID: req.EquipmentID,
			ReviewerID:  userID,
			Rating:      req.Rating,
			Comment:     req.Comment,
		}
		if err := s.reviewRepo.Create(&review); err != nil {
			return nil, errors.New("failed to create review")
		}
		s.reviewRepo.UpdateEquipmentRating(*req.EquipmentID)
		// Update owner rating based on equipment owner
		equipment, _ := s.equipmentRepo.FindByID(*req.EquipmentID)
		if equipment != nil {
			s.reviewRepo.UpdateOwnerRating(equipment.OwnerID)
		}
		return &review, nil
	}

	if req.RentalID != nil {
		rental, err := s.rentalRepo.FindByID(*req.RentalID)
		if err != nil {
			return nil, errors.New("rental not found")
		}

		if rental.RenterID != userID {
			return nil, errors.New("you are not allowed to review this rental")
		}

		if rental.Status != "COMPLETED" {
			return nil, errors.New("only completed rentals can be reviewed")
		}

		_, err = s.reviewRepo.GetByRentalID(*req.RentalID)
		if err == nil {
			return nil, errors.New("review already exists for this rental")
		}

		review = models.Review{
			EquipmentID: req.EquipmentID,
			RentalID:    req.RentalID,
			ReviewerID:  userID,
			Rating:      req.Rating,
			Comment:     req.Comment,
		}

		if err := s.reviewRepo.Create(&review); err != nil {
			return nil, errors.New("failed to create review")
		}

		s.reviewRepo.UpdateEquipmentRating(*req.EquipmentID)
		s.reviewRepo.UpdateOwnerRating(rental.OwnerID)

	} else if req.MarketRequestID != nil {
		// Existing marketplace request review logic (unchanged)
		mreq, err := s.marketplaceRequestRepo.FindByID(*req.MarketRequestID)
		if err != nil {
			return nil, errors.New("marketplace request not found")
		}

		if mreq.BuyerID != userID {
			return nil, errors.New("you are not allowed to review this request")
		}

		if mreq.TransactionStatus != "YES" {
			return nil, errors.New("only successful transactions can be reviewed")
		}

		_, err = s.reviewRepo.GetByMarketplaceRequestID(*req.MarketRequestID)
		if err == nil {
			return nil, errors.New("review already exists for this transaction")
		}

		review = models.Review{
			ProductID:            req.ProductID,
			MarketplaceRequestID: req.MarketRequestID,
			ReviewerID:           userID,
			Rating:               req.Rating,
			Comment:              req.Comment,
		}

		if err := s.reviewRepo.Create(&review); err != nil {
			return nil, errors.New("failed to create review")
		}

		s.reviewRepo.UpdateProductRating(*req.ProductID)
		s.reviewRepo.UpdateOwnerRating(mreq.SellerID)

	} else {
		return nil, errors.New("rental_id or market_request_id is required")
	}

	return &review, nil
}

func (s *ReviewService) UpdateReview(userID uint, reviewID uint, req dto.UpdateReviewDTO) (*models.Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return nil, errors.New("review not found")
	}

	if review.ReviewerID != userID {
		return nil, errors.New("you are not allowed to edit this review")
	}

	review.Rating = req.Rating
	review.Comment = req.Comment

	if err := s.reviewRepo.Update(review); err != nil {
		return nil, errors.New("failed to update review")
	}

	if review.EquipmentID != nil {
		equipment, err := s.equipmentRepo.FindByID(*review.EquipmentID)
		if err == nil {
			s.reviewRepo.UpdateEquipmentRating(*review.EquipmentID)
			s.reviewRepo.UpdateOwnerRating(equipment.OwnerID)
		}
	} else if review.ProductID != nil {
		product, err := s.productRepo.FindByID(*review.ProductID)
		if err == nil {
			s.reviewRepo.UpdateProductRating(*review.ProductID)
			s.reviewRepo.UpdateOwnerRating(product.SellerID)
		}
	}

	return review, nil
}

func (s *ReviewService) DeleteReview(userID uint, reviewID uint) error {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	if review.ReviewerID != userID {
		return errors.New("you are not allowed to delete this review")
	}

	if err := s.reviewRepo.Delete(reviewID); err != nil {
		return errors.New("failed to delete review")
	}

	if review.EquipmentID != nil {
		equipment, err := s.equipmentRepo.FindByID(*review.EquipmentID)
		if err == nil {
			s.reviewRepo.UpdateEquipmentRating(*review.EquipmentID)
			s.reviewRepo.UpdateOwnerRating(equipment.OwnerID)
		}
	} else if review.ProductID != nil {
		product, err := s.productRepo.FindByID(*review.ProductID)
		if err == nil {
			s.reviewRepo.UpdateProductRating(*review.ProductID)
			s.reviewRepo.UpdateOwnerRating(product.SellerID)
		}
	}

	return nil
}

func (s *ReviewService) GetEquipmentReviews(equipmentID uint) (map[string]interface{}, error) {
	equipment, err := s.equipmentRepo.FindByID(equipmentID)
	if err != nil {
		return nil, errors.New("equipment not found")
	}

	reviews, err := s.reviewRepo.GetByEquipmentID(equipmentID)
	if err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, r := range reviews {
		var reviewerName string
		user, err := s.userRepo.FindByID(r.ReviewerID)
		if err == nil && user != nil {
			reviewerName = user.FullName
		}

		output = append(output, map[string]interface{}{
			"id":          r.ID,
			"rating":      r.Rating,
			"comment":     r.Comment,
			"reviewer":    reviewerName,
			"reviewer_id": r.ReviewerID,
			"owner_reply": r.OwnerReply,
			"created_at":  r.CreatedAt.Format("2006-01-02"),
		})
	}

	return map[string]interface{}{
		"average_rating": equipment.AverageRating,
		"total_reviews":  equipment.TotalReviews,
		"reviews":       output,
	}, nil
}

func (s *ReviewService) GetProductReviews(productID uint) (map[string]interface{}, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	reviews, err := s.reviewRepo.GetByProductID(productID)
	if err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, r := range reviews {
		var reviewerName string
		user, err := s.userRepo.FindByID(r.ReviewerID)
		if err == nil && user != nil {
			reviewerName = user.FullName
		}

		output = append(output, map[string]interface{}{
			"id":          r.ID,
			"rating":      r.Rating,
			"comment":     r.Comment,
			"reviewer":    reviewerName,
			"reviewer_id": r.ReviewerID,
			"owner_reply": r.OwnerReply,
			"created_at":  r.CreatedAt.Format("2006-01-02"),
		})
	}

	return map[string]interface{}{
		"average_rating": product.AverageRating,
		"total_reviews":  product.TotalReviews,
		"reviews":       output,
	}, nil
}

func (s *ReviewService) ReplyToReview(userID uint, reviewID uint, req dto.OwnerReplyDTO) error {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	if review.EquipmentID != nil {
		equipment, err := s.equipmentRepo.FindByID(*review.EquipmentID)
		if err != nil || equipment.OwnerID != userID {
			return errors.New("not allowed")
		}
	} else if review.ProductID != nil {
		product, err := s.productRepo.FindByID(*review.ProductID)
		if err != nil || product.SellerID != userID {
			return errors.New("not allowed")
		}
	} else {
		return errors.New("invalid review type")
	}

	now := time.Now()
	review.OwnerReply = &req.Reply
	review.RepliedAt = &now

	return s.reviewRepo.Update(review)
}

func (s *ReviewService) DeleteReply(userID uint, reviewID uint) error {
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	if review.EquipmentID != nil {
		equipment, err := s.equipmentRepo.FindByID(*review.EquipmentID)
		if err != nil || equipment.OwnerID != userID {
			return errors.New("not allowed")
		}
	} else if review.ProductID != nil {
		product, err := s.productRepo.FindByID(*review.ProductID)
		if err != nil || product.SellerID != userID {
			return errors.New("not allowed")
		}
	} else {
		return errors.New("invalid review type")
	}

	review.OwnerReply = nil
	review.RepliedAt = nil

	return s.reviewRepo.Update(review)
}
