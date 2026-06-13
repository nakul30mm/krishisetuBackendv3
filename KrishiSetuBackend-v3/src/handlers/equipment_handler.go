package handlers

import (
	"net/http"
	"strconv"

	"krishisetu-backend/src/business"
	"krishisetu-backend/src/dto"
	"github.com/gin-gonic/gin"
)

type EquipmentHandler struct {
	Service            *business.EquipmentService
	AutoCompleteRentals func()
}

func NewEquipmentHandler(service *business.EquipmentService, autoCompleteRentals func()) *EquipmentHandler {
	return &EquipmentHandler{
		Service:            service,
		AutoCompleteRentals: autoCompleteRentals,
	}
}

func (h *EquipmentHandler) CreateEquipment(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateEquipmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	equipment, err := h.Service.CreateEquipment(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, equipment)
}

func (h *EquipmentHandler) GetAllEquipments(c *gin.Context) {
	search := c.Query("search")
	sort := c.Query("sort")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")

	// Map common UI sort aliases to backend values
	switch sort {
	case "priceLowHigh", "price_low_high":
		sort = "price_low_high"
	case "priceHighLow", "price_high_low":
		sort = "price_high_low"
	case "pincodeClosest", "pincode_closest":
		sort = "pincode_closest"
	case "pincodeFarthest", "pincode_farthest":
		sort = "pincode_farthest"
	}

	var userID uint
	if id, exists := c.Get("user_id"); exists {
		if uID, ok := id.(uint); ok {
			userID = uID
		}
	}

	response, err := h.Service.GetAllEquipments(userID, search, sort, minPrice, maxPrice, h.AutoCompleteRentals)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *EquipmentHandler) GetMyEquipments(c *gin.Context) {
	userID := c.GetUint("user_id")

	response, err := h.Service.GetMyEquipments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *EquipmentHandler) GetUnavailableDates(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
		return
	}

	ranges, err := h.Service.GetUnavailableDates(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ranges)
}

func (h *EquipmentHandler) UpdateEquipment(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
		return
	}

	var req dto.UpdateEquipmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	equipment, err := h.Service.UpdateEquipment(userID, uint(id), req)
	if err != nil {
		// Differentiate 404, 403, and 400
		if err.Error() == "Equipment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "Not authorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, equipment)
}

func (h *EquipmentHandler) DeleteEquipment(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid equipment ID"})
		return
	}

	err = h.Service.DeleteEquipment(userID, uint(id))
	if err != nil {
		if err.Error() == "Equipment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "Not authorized" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Equipment deleted"})
}
