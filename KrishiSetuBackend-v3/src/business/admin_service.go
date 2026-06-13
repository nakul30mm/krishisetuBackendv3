package business

import (
	"krishisetu-backend/src/constants"
	"krishisetu-backend/src/dto"
	"krishisetu-backend/models"
)

type AdminService struct {
	adminRepo     AdminRepository
	userRepo      UserRepository
	productRepo   ProductRepository
	equipmentRepo EquipmentRepository
	reviewRepo    ReviewRepository
}

func NewAdminService(
	adminRepo AdminRepository,
	userRepo UserRepository,
	productRepo ProductRepository,
	equipmentRepo EquipmentRepository,
	reviewRepo ReviewRepository,
) *AdminService {
	return &AdminService{
		adminRepo:     adminRepo,
		userRepo:      userRepo,
		productRepo:   productRepo,
		equipmentRepo: equipmentRepo,
		reviewRepo:    reviewRepo,
	}
}

func (s *AdminService) GetStats() (*dto.AdminStatsDTO, error) {
	uCount, pCount, eCount, rCount, err := s.adminRepo.GetStats()
	if err != nil {
		return nil, constants.NewAppError(500, "Failed to fetch admin stats")
	}

	return &dto.AdminStatsDTO{
		TotalUsers:      uCount,
		TotalProducts:   pCount,
		TotalEquipments: eCount,
		TotalRequests:   rCount,
	}, nil
}

func (s *AdminService) GetUsers() ([]models.User, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, constants.NewAppError(500, "Failed to fetch users")
	}
	return users, nil
}

func (s *AdminService) BlockUser(id uint) error {
	err := s.userRepo.UpdateFields(id, map[string]interface{}{"is_blocked": true})
	if err != nil {
		return constants.NewAppError(500, "Failed to block user")
	}
	return nil
}

func (s *AdminService) UnblockUser(id uint) error {
	err := s.userRepo.UpdateFields(id, map[string]interface{}{"is_blocked": false})
	if err != nil {
		return constants.NewAppError(500, "Failed to unblock user")
	}
	return nil
}

func (s *AdminService) DeleteUser(id uint) error {
	err := s.userRepo.Delete(id)
	if err != nil {
		return constants.NewAppError(500, "Failed to delete user")
	}
	return nil
}

func (s *AdminService) DeleteProduct(id uint) error {
	err := s.productRepo.Delete(id)
	if err != nil {
		return constants.NewAppError(500, "Failed to delete product")
	}
	return nil
}

func (s *AdminService) DeleteEquipment(id uint) error {
	err := s.equipmentRepo.Delete(id)
	if err != nil {
		return constants.NewAppError(500, "Failed to delete equipment")
	}
	return nil
}

func (s *AdminService) DeleteReview(id uint) error {
	err := s.reviewRepo.Delete(id)
	if err != nil {
		return constants.NewAppError(500, "Failed to delete review")
	}
	return nil
}

func (s *AdminService) PromoteToAdmin(identifier string) (*models.User, error) {
	user, err := s.userRepo.FindByIdentifier(identifier)
	if err != nil {
		return nil, constants.NewAppError(404, "User with "+identifier+" not found.")
	}

	err = s.userRepo.UpdateFields(user.ID, map[string]interface{}{"is_admin": true})
	if err != nil {
		return nil, constants.NewAppError(500, "Failed to promote user")
	}

	user.IsAdmin = true
	return user, nil
}
