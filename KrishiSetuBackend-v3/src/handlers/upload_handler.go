package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image uploaded"})
		return
	}

	// Create unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Return the relative URL
	url := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{"url": url})
}
