package business

import (
	"errors"
	"time"

	"krishisetu-backend/models"
	"krishisetu-backend/src/dto"
)

type RentalService struct {
	rentalRepo    RentalRepository
	equipmentRepo EquipmentRepository
}

func NewRentalService(rentalRepo RentalRepository, equipmentRepo EquipmentRepository) *RentalService {
	return &RentalService{
		rentalRepo:    rentalRepo,
		equipmentRepo: equipmentRepo,
	}
}

func (s *RentalService) CreateRental(userID uint, req dto.CreateRentalDTO) (*models.Rental, error) {
	if req.StartDate > req.EndDate {
		return nil, errors.New("Invalid date range")
	}

	equipment, err := s.equipmentRepo.FindByID(req.EquipmentID)
	if err != nil {
		return nil, errors.New("Equipment not found")
	}

	if equipment.OwnerID == userID {
		return nil, errors.New("You cannot rent your own equipment")
	}

	approvedRentals, err := s.rentalRepo.GetApprovedRentalsForEquipment(req.EquipmentID)
	if err == nil {
		for _, r := range approvedRentals {
			if DatesOverlap(req.StartDate, req.EndDate, r.StartDate, r.EndDate) {
				return nil, errors.New("Equipment already booked for selected dates")
			}
		}
	}

	rental := models.Rental{
		EquipmentID: req.EquipmentID,
		RenterID:    userID,
		OwnerID:     equipment.OwnerID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      "PENDING",
	}

	if err := s.rentalRepo.Create(&rental); err != nil {
		return nil, errors.New("Failed to create rental")
	}

	return &rental, nil
}

func (s *RentalService) GetMyRentalRequests(userID uint) ([]map[string]interface{}, error) {
	s.AutoCompleteRentals()

	rentals, err := s.rentalRepo.GetByRenterID(userID)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	for _, r := range rentals {
		canReview := false
		if r.Status == "COMPLETED" {
			_, reviewErr := s.rentalRepo.GetReviewForRental(r.ID)
			if reviewErr != nil {
				canReview = true
			}
		}

		response = append(response, map[string]interface{}{
			"id":         r.ID,
			"start_date": r.StartDate,
			"end_date":   r.EndDate,
			"status":     r.Status,
			"can_review": canReview,
			"equipment":  s.FormatEquipmentResponse(r.Equipment),
		})
	}

	return response, nil
}

func (s *RentalService) GetRentalRequestsForOwner(userID uint) ([]map[string]interface{}, error) {
	s.AutoCompleteRentals()

	rentals, err := s.rentalRepo.GetByOwnerID(userID)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	for _, r := range rentals {
		response = append(response, map[string]interface{}{
			"id":         r.ID,
			"start_date": r.StartDate,
			"end_date":   r.EndDate,
			"status":     r.Status,
			"equipment":  s.FormatEquipmentResponse(r.Equipment),
			"renter": map[string]interface{}{
				"id":   r.Renter.ID,
				"name": r.Renter.FullName,
			},
		})
	}

	return response, nil
}

func (s *RentalService) ApproveRental(userID uint, rentalID uint) error {
	rental, err := s.rentalRepo.FindByIDWithEquipmentAndOwner(rentalID)
	if err != nil {
		return errors.New("Rental not found")
	}

	if rental.OwnerID != userID {
		return errors.New("You are not authorized to approve this request")
	}

	approvedRentals, err := s.rentalRepo.GetApprovedRentalsForEquipment(rental.EquipmentID)
	if err == nil {
		for _, r := range approvedRentals {
			if DatesOverlap(rental.StartDate, rental.EndDate, r.StartDate, r.EndDate) {
				return errors.New("This request overlaps with an already approved rental")
			}
		}
	}

	rental.Status = "APPROVED"
	if err := s.rentalRepo.Update(rental); err != nil {
		return err
	}

	s.rentalRepo.UpdateEquipmentStatus(rental.EquipmentID, "ENGAGED")

	// Reject other overlapping pending rentals
	s.rentalRepo.RejectOverlappingPendingRentals(rental.EquipmentID, rental.ID, rental.StartDate, rental.EndDate)

	return nil
}

func (s *RentalService) RejectRental(userID uint, rentalID uint) (*models.Rental, error) {
	rental, err := s.rentalRepo.FindByID(rentalID)
	if err != nil {
		return nil, errors.New("Rental not found")
	}

	if rental.OwnerID != userID {
		return nil, errors.New("You are not authorized to reject this request")
	}

	rental.Status = "REJECTED"
	if err := s.rentalRepo.Update(rental); err != nil {
		return nil, err
	}

	return rental, nil
}

func (s *RentalService) DeleteRental(userID uint, rentalID uint) error {
	rental, err := s.rentalRepo.FindByID(rentalID)
	if err != nil {
		return errors.New("Rental request not found")
	}

	if rental.RenterID != userID {
		return errors.New("Only the renter can delete this request")
	}

	if rental.Status != "PENDING" {
		return errors.New("Only pending requests can be deleted")
	}

	return s.rentalRepo.Delete(rentalID)
}

func (s *RentalService) UpdateRental(userID uint, rentalID uint, req dto.UpdateRentalDTO) error {
	if req.StartDate > req.EndDate {
		return errors.New("Invalid date range")
	}

	rental, err := s.rentalRepo.FindByID(rentalID)
	if err != nil {
		return errors.New("Rental request not found")
	}

	if rental.RenterID != userID {
		return errors.New("Only the renter can edit this request")
	}

	if rental.Status != "PENDING" {
		return errors.New("Only pending requests can be edited")
	}

	approvedRentals, err := s.rentalRepo.GetApprovedRentalsForEquipmentExcluding(rental.EquipmentID, rental.ID)
	if err == nil {
		for _, r := range approvedRentals {
			if DatesOverlap(req.StartDate, req.EndDate, r.StartDate, r.EndDate) {
				return errors.New("Selected dates overlap with an approved rental")
			}
		}
	}

	rental.StartDate = req.StartDate
	rental.EndDate = req.EndDate

	return s.rentalRepo.Update(rental)
}

func (s *RentalService) AutoCompleteRentals() {
	rentals, err := s.rentalRepo.GetAllApprovedRentals()
	if err != nil {
		return
	}

	now := time.Now()

	for _, r := range rentals {
		end, err := time.Parse("2006-01-02", r.EndDate)
		if err != nil {
			continue
		}

		bufferEnd := end.AddDate(0, 0, 1)

		if now.After(bufferEnd) {
			r.Status = "COMPLETED"
			s.rentalRepo.Update(&r)

			count, err := s.rentalRepo.CountApprovedRentalsForEquipment(r.EquipmentID)
			if err == nil && count == 0 {
				s.rentalRepo.UpdateEquipmentStatus(r.EquipmentID, "VACANT")
			}
		}
	}
}

func DatesOverlap(start1, end1, start2, end2 string) bool {
	return start1 <= end2 && start2 <= end1
}

func (s *RentalService) FormatEquipmentResponse(e models.Equipment) map[string]interface{} {
	return map[string]interface{}{
		"id":             e.ID,
		"name":           e.Name,
		"category":       e.Category,
		"description":    e.Description,
		"price_per_day":  e.PricePerDay,
		"price_unit":     e.PriceUnit,
		"status":         e.Status,
		"booked_from":    e.BookedFrom,
		"booked_to":      e.BookedTo,
		"owner_id":       e.OwnerID,
		"average_rating": e.AverageRating,
		"total_reviews":  e.TotalReviews,
		"image1":         e.Image1,
		"image2":         e.Image2,
		"image3":         e.Image3,
		"location": map[string]interface{}{
			"pincode":  e.LocationPincode,
			"district": e.LocationDistrict,
			"state":    e.LocationState,
			"city":     e.LocationCity,
			"location": e.LocationArea,
		},
		"owner": map[string]interface{}{
			"id":             e.Owner.ID,
			"name":           e.Owner.FullName,
			"average_rating": e.Owner.AverageRating,
			"total_ratings":  e.Owner.TotalRatings,
		},
	}
}
