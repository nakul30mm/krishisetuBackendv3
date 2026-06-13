package repository

import (
	"krishisetu-backend/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateFields(id uint, updates map[string]interface{}) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) FindByIdentifier(identifier string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ? OR email = ?", identifier, identifier).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
