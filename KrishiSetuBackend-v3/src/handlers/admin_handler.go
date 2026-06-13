package handlers

import (
	"krishisetu-backend/src/business"
	"krishisetu-backend/src/constants"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService *business.AdminService
}

func NewAdminHandler(adminService *business.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// GetAdminStats returns high-level system metrics
func (h *AdminHandler) GetAdminStats(c *gin.Context) {
	stats, err := h.adminService.GetStats()
	if err != nil {
		constants.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetAdminUsers returns a list of all users
func (h *AdminHandler) GetAdminUsers(c *gin.Context) {
	users, err := h.adminService.GetUsers()
	if err != nil {
		constants.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// BlockUser suspends a user account
func (h *AdminHandler) BlockUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		constants.HandleError(c, constants.NewAppError(http.StatusBadRequest, "Invalid ID format"))
		return
	}

	if err := h.adminService.BlockUser(uint(id)); err != nil {
		constants.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User blocked successfully"})
}

// UnblockUser restores a user account
func (h *AdminHandler) UnblockUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		constants.HandleError(c, constants.NewAppError(http.StatusBadRequest, "Invalid ID format"))
		return
	}

	if err := h.adminService.UnblockUser(uint(id)); err != nil {
		constants.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User unblocked successfully"})
}

// DeleteUser removes a user account
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		constants.HandleError(c, constants.NewAppError(http.StatusBadRequest, "Invalid ID format"))
		return
	}

	if err := h.adminService.DeleteUser(uint(id)); err != nil {
		constants.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// AdminDeleteProduct allows an admin to remove any product
func (h *AdminHandler) AdminDeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		constants.HandleError(c, constants.NewAppError(http.StatusBadRequest, "Invalid ID format"))
		return
	}

	if err := h.adminService.DeleteProduct(uint(id)); err != nil {
		constants.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product removed by admin"})
}

// AdminDeleteEquipment allows an admin to remove any equipment
func (h *AdminHandler) AdminDeleteEquipment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		constants.HandleError(c, constants.NewAppError(http.StatusBadRequest, "Invalid ID format"))
		return
	}

	if err := h.adminService.DeleteEquipment(uint(id)); err != nil {
		constants.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Equipment removed by admin"})
}

// AdminDeleteReview allows an admin to remove any review
func (h *AdminHandler) AdminDeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		constants.HandleError(c, constants.NewAppError(http.StatusBadRequest, "Invalid ID format"))
		return
	}

	if err := h.adminService.DeleteReview(uint(id)); err != nil {
		constants.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Review removed by admin"})
}

// PromoteToAdmin allows promoting a user to administrator (for setup)
func (h *AdminHandler) PromoteToAdmin(c *gin.Context) {
	identifier := c.Param("identifier")

	user, err := h.adminService.PromoteToAdmin(identifier)
	if err != nil {
		constants.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SUCCESS! " + user.Email + " IS NOW ADMIN!"})
}
