package business

import (
	"errors"
	"time"

	"krishisetu-backend/models"
	"krishisetu-backend/src/dto"
	"krishisetu-backend/src/utils"
)

type AuthService struct {
	userRepo UserRepository
}

func NewAuthService(userRepo UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(req dto.RegisterDTO) error {
	if req.FullName == "" || req.Email == "" || req.Password == "" {
		return errors.New("All fields are required")
	}

	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return errors.New("Email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("Password hashing failed")
	}

	user := models.User{
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  hashedPassword,
		Phone:     req.Phone,
		Age:       req.Age,
		Gender:    req.Gender,
		City:      req.City,
		District:  req.District,
		State:     req.State,
		Pincode:   req.Pincode,
		Location:  req.Location,
		IsAdmin:   false,
		IsBlocked: false,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return errors.New("User creation failed")
	}

	return nil
}

func (s *AuthService) Login(req dto.LoginDTO) (string, map[string]interface{}, error) {
	if req.Email == "" || req.Password == "" {
		return "", nil, errors.New("Email and password are required")
	}

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", nil, errors.New("Invalid email or password")
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return "", nil, errors.New("Invalid email or password")
	}

	if user.IsBlocked {
		return "", nil, errors.New("Account suspended. Please contact support.")
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", nil, errors.New("Token generation failed")
	}

	userData := map[string]interface{}{
		"id":    user.ID,
		"name":  user.FullName,
		"email": user.Email,
	}

	return token, userData, nil
}

func (s *AuthService) ForgotPassword(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("Email not registered")
	}

	otp, expiry := utils.GenerateOTP()

	user.ResetOTP = &otp
	user.ResetOTPExpiry = &expiry
	user.ResetOTPVerified = false

	if err := s.userRepo.Update(user); err != nil {
		return errors.New("Failed to save OTP")
	}

	if err := utils.SendOTPEmail(user.Email, otp); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyOTP(email, otp string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("User not found")
	}

	if user.ResetOTP == nil || *user.ResetOTP != otp {
		return errors.New("Invalid OTP")
	}

	if user.ResetOTPExpiry == nil || time.Now().After(*user.ResetOTPExpiry) {
		return errors.New("OTP expired")
	}

	user.ResetOTPVerified = true

	if err := s.userRepo.Update(user); err != nil {
		return errors.New("Failed to verify OTP")
	}

	return nil
}

func (s *AuthService) ResetPassword(email, newPassword string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("User not found")
	}

	if !user.ResetOTPVerified {
		return errors.New("OTP not verified")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("Password hashing failed")
	}

	updates := map[string]interface{}{
		"password":           hashedPassword,
		"reset_otp":          nil,
		"reset_otp_expiry":   nil,
		"reset_otp_verified": false,
	}

	if err := s.userRepo.UpdateFields(user.ID, updates); err != nil {
		return errors.New("Failed to update password")
	}

	return nil
}

func (s *AuthService) Profile(userID uint) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("User not found")
	}

	profile := map[string]interface{}{
		"full_name":       user.FullName,
		"email":           user.Email,
		"phone":           user.Phone,
		"age":             user.Age,
		"gender":          user.Gender,
		"city":            user.City,
		"district":        user.District,
		"state":           user.State,
		"pincode":         user.Pincode,
		"location":        user.Location,
		"profile_picture": user.ProfilePicture,
	}

	return profile, nil
}

func (s *AuthService) UpdateProfile(userID uint, req dto.UpdateProfileDTO) error {
	updates := map[string]interface{}{
		"full_name":       req.FullName,
		"phone":           req.Phone,
		"age":             req.Age,
		"gender":          req.Gender,
		"city":            req.City,
		"district":        req.District,
		"state":           req.State,
		"pincode":         req.Pincode,
		"location":        req.Location,
		"profile_picture": req.ProfilePicture,
	}

	if err := s.userRepo.UpdateFields(userID, updates); err != nil {
		return errors.New("Failed to update profile")
	}

	return nil
}
