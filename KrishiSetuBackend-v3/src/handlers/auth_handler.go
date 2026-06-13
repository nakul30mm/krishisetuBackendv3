package handlers

import (
    "net/http"

    "krishisetu-backend/src/business"
    "krishisetu-backend/src/dto"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    Service *business.AuthService
}

func NewAuthHandler(service *business.AuthService) *AuthHandler {
    return &AuthHandler{Service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req dto.RegisterDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    if err := h.Service.Register(req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }
    token, user, err := h.Service.Login(req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token, "user": user})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
    var body struct { Email string `json:"email"` }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    if err := h.Service.ForgotPassword(body.Email); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "OTP sent to registered email"})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
    var body struct { Email string `json:"email"`; OTP string `json:"otp"` }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    if err := h.Service.VerifyOTP(body.Email, body.OTP); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "OTP verified"})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
    var body struct { Email string `json:"email"`; NewPassword string `json:"newPassword"` }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    if err := h.Service.ResetPassword(body.Email, body.NewPassword); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func (h *AuthHandler) Profile(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    profile, err := h.Service.Profile(userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, profile)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    var req dto.UpdateProfileDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    if err := h.Service.UpdateProfile(userID.(uint), req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
