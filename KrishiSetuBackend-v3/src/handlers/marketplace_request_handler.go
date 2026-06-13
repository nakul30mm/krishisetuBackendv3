package handlers

import (
	"net/http"
	"strconv"

	"krishisetu-backend/src/business"
	"krishisetu-backend/src/dto"
	"github.com/gin-gonic/gin"
)

type MarketplaceRequestHandler struct {
	Service *business.MarketplaceRequestService
}

func NewMarketplaceRequestHandler(service *business.MarketplaceRequestService) *MarketplaceRequestHandler {
	return &MarketplaceRequestHandler{Service: service}
}

func (h *MarketplaceRequestHandler) CreateMarketplaceRequest(c *gin.Context) {
	buyerID := c.GetUint("user_id")

	var req dto.CreateMarketplaceRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	marketplaceRequest, err := h.Service.CreateMarketplaceRequest(buyerID, req)
	if err != nil {
		if err.Error() == "Product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, marketplaceRequest)
}

func (h *MarketplaceRequestHandler) GetSentMarketplaceRequests(c *gin.Context) {
	buyerID := c.GetUint("user_id")
	search := c.Query("search")

	requests, err := h.Service.GetSentMarketplaceRequests(buyerID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func (h *MarketplaceRequestHandler) GetReceivedMarketplaceRequests(c *gin.Context) {
	sellerID := c.GetUint("user_id")
	search := c.Query("search")

	requests, err := h.Service.GetReceivedMarketplaceRequests(sellerID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, requests)
}

func (h *MarketplaceRequestHandler) UpdateMarketplaceRequestStatus(c *gin.Context) {
	reqIDStr := c.Param("id")
	action := c.Param("action") // "approve" or "reject"
	sellerID := c.GetUint("user_id")

	reqID, err := strconv.ParseUint(reqIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	req, err := h.Service.UpdateMarketplaceRequestStatus(sellerID, uint(reqID), action)
	if err != nil {
		if err.Error() == "Request not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "Only the product seller can map this action" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	var actionText string
	if action == "approve" {
		actionText = "APPROVED"
	} else {
		actionText = "REJECTED"
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request " + actionText, "request": req})
}

func (h *MarketplaceRequestHandler) DeleteMarketplaceRequest(c *gin.Context) {
	reqIDStr := c.Param("id")
	buyerID := c.GetUint("user_id")

	reqID, err := strconv.ParseUint(reqIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	err = h.Service.DeleteMarketplaceRequest(buyerID, uint(reqID))
	if err != nil {
		if err.Error() == "Request not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "Only the buyer can delete this request" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request deleted successfully"})
}

func (h *MarketplaceRequestHandler) ConfirmMarketplaceTransaction(c *gin.Context) {
	reqIDStr := c.Param("id")
	reqID, err := strconv.ParseUint(reqIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var input dto.ConfirmTransactionDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	buyerID := c.GetUint("user_id")

	req, err := h.Service.ConfirmMarketplaceTransaction(buyerID, uint(reqID), input)
	if err != nil {
		if err.Error() == "Request not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "Only the buyer can confirm this transaction" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction status updated", "request": req})
}

func (h *MarketplaceRequestHandler) UpdateMarketplaceRequest(c *gin.Context) {
	reqIDStr := c.Param("id")
	reqID, err := strconv.ParseUint(reqIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	buyerID := c.GetUint("user_id")

	var input dto.UpdateMarketplaceRequestDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	req, err := h.Service.UpdateMarketplaceRequest(buyerID, uint(reqID), input)
	if err != nil {
		if err.Error() == "Request not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "Only the buyer can edit this request" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request updated successfully", "request": req})
}
