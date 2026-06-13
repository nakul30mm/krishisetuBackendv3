package handlers

import (
	"net/http"
	"strconv"

	"krishisetu-backend/src/business"
	"krishisetu-backend/src/dto"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	Service *business.ReviewService
}

func NewReviewHandler(service *business.ReviewService) *ReviewHandler {
	return &ReviewHandler{Service: service}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateReviewDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	review, err := h.Service.CreateReview(userID, req)
	if err != nil {
		if err.Error() == "rental not found" || err.Error() == "marketplace request not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "you are not allowed to review this rental" || err.Error() == "you are not allowed to review this request" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "review added successfully",
		"review":  review,
	})
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID := c.GetUint("user_id")
	reviewIDStr := c.Param("id")
	reviewID, err := strconv.ParseUint(reviewIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	var req dto.UpdateReviewDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	review, err := h.Service.UpdateReview(userID, uint(reviewID), req)
	if err != nil {
		if err.Error() == "review not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "you are not allowed to edit this review" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "review updated successfully",
		"review":  review,
	})
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	userID := c.GetUint("user_id")
	reviewIDStr := c.Param("id")
	reviewID, parseErr := strconv.ParseUint(reviewIDStr, 10, 32)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	err := h.Service.DeleteReview(userID, uint(reviewID))
	if err != nil {
		if err.Error() == "review not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "you are not allowed to delete this review" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "review deleted successfully"})
}

func (h *ReviewHandler) GetEquipmentReviews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid equipment ID"})
		return
	}

	response, err := h.Service.GetEquipmentReviews(uint(id))
	if err != nil {
		if err.Error() == "equipment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	response, err := h.Service.GetProductReviews(uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ReviewHandler) ReplyToReview(c *gin.Context) {
	userID := c.GetUint("user_id")
	reviewIDStr := c.Param("id")
	reviewID, err := strconv.ParseUint(reviewIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	var req dto.OwnerReplyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reply is required"})
		return
	}

	err = h.Service.ReplyToReview(userID, uint(reviewID), req)
	if err != nil {
		if err.Error() == "review not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "not allowed" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "reply saved successfully"})
}

func (h *ReviewHandler) DeleteReply(c *gin.Context) {
	userID := c.GetUint("user_id")
	reviewIDStr := c.Param("id")
	reviewID, err := strconv.ParseUint(reviewIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	err = h.Service.DeleteReply(userID, uint(reviewID))
	if err != nil {
		if err.Error() == "review not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "not allowed" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "reply deleted successfully"})
}
