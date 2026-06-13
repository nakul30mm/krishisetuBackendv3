package handlers

import (
	"net/http"
	"strconv"

	"krishisetu-backend/src/business"
	"krishisetu-backend/src/dto"
	"github.com/gin-gonic/gin"
)

type RentalHandler struct {
	Service *business.RentalService
}

func NewRentalHandler(service *business.RentalService) *RentalHandler {
	return &RentalHandler{Service: service}
}

func (h *RentalHandler) CreateRental(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateRentalDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rental, err := h.Service.CreateRental(userID, req)
	if err != nil {
		if err.Error() == "Equipment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "You cannot rent your own equipment" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, rental)
}

func (h *RentalHandler) GetMyRentalRequests(c *gin.Context) {
	userID := c.GetUint("user_id")

	response, err := h.Service.GetMyRentalRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *RentalHandler) GetRentalRequestsForOwner(c *gin.Context) {
	userID := c.GetUint("user_id")

	response, err := h.Service.GetRentalRequestsForOwner(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *RentalHandler) ApproveRental(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rental ID"})
		return
	}

	err = h.Service.ApproveRental(userID, uint(id))
	if err != nil {
		if err.Error() == "Rental not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "You are not authorized to approve this request" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rental approved successfully"})
}

func (h *RentalHandler) RejectRental(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rental ID"})
		return
	}

	rental, err := h.Service.RejectRental(userID, uint(id))
	if err != nil {
		if err.Error() == "Rental not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "You are not authorized to reject this request" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, rental)
}

func (h *RentalHandler) DeleteRental(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rental ID"})
		return
	}

	err = h.Service.DeleteRental(userID, uint(id))
	if err != nil {
		if err.Error() == "Rental request not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "Only the renter can delete this request" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rental request deleted successfully"})
}

func (h *RentalHandler) UpdateRental(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rental ID"})
		return
	}

	var req dto.UpdateRentalDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err = h.Service.UpdateRental(userID, uint(id), req)
	if err != nil {
		if err.Error() == "Rental request not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "Only the renter can edit this request" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rental request updated successfully"})
}
