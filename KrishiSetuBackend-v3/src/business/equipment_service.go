package business

import (
	"errors"
	"time"

	"krishisetu-backend/models"
	"krishisetu-backend/src/dto"
)

type EquipmentService struct {
	equipmentRepo EquipmentRepository
	userRepo      UserRepository
}

func NewEquipmentService(equipmentRepo EquipmentRepository, userRepo UserRepository) *EquipmentService {
	return &EquipmentService{
		equipmentRepo: equipmentRepo,
		userRepo:      userRepo,
	}
}

func (s *EquipmentService) CreateEquipment(userID uint, req dto.CreateEquipmentDTO) (*models.Equipment, error) {
	equipment := models.Equipment{
		Name:             req.Name,
		Category:         req.Category,
		Description:      req.Description,
		PricePerDay:      req.PricePerDay,
		PriceUnit:        req.PriceUnit,
		Status:           "VACANT",
		LocationPincode:  req.Location.Pincode,
		LocationDistrict: req.Location.District,
		LocationState:    req.Location.State,
		LocationCity:     req.Location.City,
		LocationArea:     req.Location.Location,
		OwnerID:          userID,
		Image1:           req.Image1,
		Image2:           req.Image2,
		Image3:           req.Image3,
	}

	if err := s.equipmentRepo.Create(&equipment); err != nil {
		return nil, err
	}

	return &equipment, nil
}

func (s *EquipmentService) GetAllEquipments(userID uint, search, sort, minPrice, maxPrice string, autoCompleteFn func()) ([]map[string]interface{}, error) {
	if autoCompleteFn != nil {
		autoCompleteFn()
	}

	var userPincode string
	if userID != 0 {
		user, err := s.userRepo.FindByID(userID)
		if err == nil && user != nil {
			userPincode = user.Pincode
		}
	}

	equipments, err := s.equipmentRepo.GetFiltered(search, sort, minPrice, maxPrice, userPincode)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	for _, e := range equipments {
		e.Status = s.DeriveEquipmentStatus(e.ID)
		response = append(response, s.FormatEquipmentResponse(e))
	}

	return response, nil
}

func (s *EquipmentService) GetMyEquipments(userID uint) ([]map[string]interface{}, error) {
	equipments, err := s.equipmentRepo.GetByOwnerID(userID)
	if err != nil {
		return nil, err
	}

	var response []map[string]interface{}
	for _, e := range equipments {
		e.Status = s.DeriveEquipmentStatus(e.ID)
		response = append(response, s.FormatEquipmentResponse(e))
	}

	return response, nil
}

func (s *EquipmentService) GetUnavailableDates(equipmentID uint) ([]map[string]interface{}, error) {
	rentals, err := s.equipmentRepo.GetApprovedRentals(equipmentID)
	if err != nil {
		return nil, err
	}

	var ranges []map[string]interface{}
	for _, r := range rentals {
		end, _ := time.Parse("2006-01-02", r.EndDate)
		ranges = append(ranges, map[string]interface{}{
			"start_date":  r.StartDate,
			"end_date":    r.EndDate,
			"buffer_date": end.AddDate(0, 0, 1).Format("2006-01-02"),
		})
	}

	return ranges, nil
}

func (s *EquipmentService) UpdateEquipment(userID uint, equipmentID uint, req dto.UpdateEquipmentDTO) (*models.Equipment, error) {
	equipment, err := s.equipmentRepo.FindByID(equipmentID)
	if err != nil {
		return nil, errors.New("Equipment not found")
	}

	if equipment.OwnerID != userID {
		return nil, errors.New("Not authorized")
	}

	count, err := s.equipmentRepo.CountActiveRentals(equipmentID)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("Cannot edit equipment with active rentals")
	}

	equipment.Name = req.Name
	equipment.Description = req.Description
	equipment.PricePerDay = req.PricePerDay
	equipment.PriceUnit = req.PriceUnit
	equipment.LocationArea = req.Location.Location
	equipment.LocationDistrict = req.Location.District
	equipment.LocationPincode = req.Location.Pincode
	equipment.LocationState = req.Location.State
	equipment.LocationCity = req.Location.City
	equipment.Image1 = req.Image1
	equipment.Image2 = req.Image2
	equipment.Image3 = req.Image3

	if err := s.equipmentRepo.Update(equipment); err != nil {
		return nil, err
	}

	return equipment, nil
}

func (s *EquipmentService) DeleteEquipment(userID uint, equipmentID uint) error {
	equipment, err := s.equipmentRepo.FindByID(equipmentID)
	if err != nil {
		return errors.New("Equipment not found")
	}

	if equipment.OwnerID != userID {
		return errors.New("Not authorized")
	}

	// We derive equipment status dynamically
	currentStatus := s.DeriveEquipmentStatus(equipmentID)
	if currentStatus == "ENGAGED" {
		return errors.New("Cannot delete engaged equipment")
	}

	count, err := s.equipmentRepo.CountActiveRentals(equipmentID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Cannot delete equipment with active rentals")
	}

	return s.equipmentRepo.Delete(equipmentID)
}

func (s *EquipmentService) DeriveEquipmentStatus(equipmentID uint) string {
	rental, err := s.equipmentRepo.GetLatestApprovedRental(equipmentID)
	if err != nil {
		return "VACANT"
	}

	end, err := time.Parse("2006-01-02", rental.EndDate)
	if err != nil {
		return "VACANT"
	}

	bufferEnd := end.AddDate(0, 0, 1)
	today := time.Now().Truncate(24 * time.Hour)

	if today.After(bufferEnd) {
		return "VACANT"
	}

	return "ENGAGED"
}

func (s *EquipmentService) FormatEquipmentResponse(e models.Equipment) map[string]interface{} {
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
