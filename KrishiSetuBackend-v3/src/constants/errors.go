package constants

import (
	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		return
	}
	c.JSON(500, gin.H{"error": err.Error()})
}
