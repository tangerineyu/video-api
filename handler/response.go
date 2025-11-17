package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// PaginationMeta,分页API
type Pagination struct {
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
}

func Success(c *gin.Context, httpStatus int, data interface{}) {
	c.JSON(httpStatus, gin.H{
		"success": true,
		"data":    data,
	})
}
func SuccessList(c *gin.Context, data interface{}, meta interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"meta":    meta,
	})
}
func Error(c *gin.Context, httpStatus int, code string, message string) {
	c.JSON(httpStatus, gin.H{
		"success": false,
		"error": ErrorResponse{
			Code:    code,
			Message: message,
		},
	})
}
func ValidationError(c *gin.Context, code string, message string, details interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error": ErrorResponse{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
