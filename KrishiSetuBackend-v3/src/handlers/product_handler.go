package handlers

import (
	"net/http"
	"strconv"

	"krishisetu-backend/src/business"
	"krishisetu-backend/src/dto"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Service *business.ProductService
}

func NewProductHandler(service *business.ProductService) *ProductHandler {
	return &ProductHandler{Service: service}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.Service.CreateProduct(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	search := c.Query("search")
	sort := c.Query("sort")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")

	var userID uint
	if id, exists := c.Get("user_id"); exists {
		if uID, ok := id.(uint); ok {
			userID = uID
		}
	}

	response, err := h.Service.GetAllProducts(userID, search, sort, minPrice, maxPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetMyProducts(c *gin.Context) {
	userID := c.GetUint("user_id")
	search := c.Query("search")

	response, err := h.Service.GetMyProducts(userID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	userID := c.GetUint("user_id")
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var req dto.UpdateProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.Service.UpdateProduct(userID, uint(productID), req)
	if err != nil {
		if err.Error() == "Product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "You can only edit your own products" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "product": product})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	userID := c.GetUint("user_id")
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	err = h.Service.DeleteProduct(userID, uint(productID))
	if err != nil {
		if err.Error() == "Product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err.Error() == "You can only delete your own products" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
