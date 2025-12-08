package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"video-api/pkg/errno"
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

// 负责把errno和data拼在一起
// 如果data是一个map，平铺
// 不是就塞进data字段，Error只返回错误吗和消息，,Success返回一个对象，符合SendResponse
func SendResponse(c *gin.Context, err errno.ErrNo, data interface{}) {
	resp := gin.H{
		"status_code": err.StatusCode,
		"status_msg":  err.StatusMsg,
	}
	if data != nil {
		if m, ok := data.(map[string]interface{}); ok {
			for k, v := range m {
				resp[k] = v
			}
		} else if m, ok := data.(gin.H); ok {
			for k, v := range m {
				resp[k] = v
			}
		} else {
			resp["data"] = data
		}
	}
	c.JSON(http.StatusOK, resp)
}

func Success(c *gin.Context, data interface{}) {
	SendResponse(c, errno.Success, data)
}
func SuccessList(c *gin.Context, list interface{}, meta interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status_code": errno.Success.StatusCode,
		"status_msg":  errno.Success.StatusMsg,
		"list":        list,
		"meta":        meta,
	})
}
func Error(c *gin.Context, err errno.ErrNo) {
	SendResponse(c, err, nil)
}
func ValidationError(c *gin.Context, detail string) {
	err := errno.ParamErr //40001
	err.StatusMsg += ": " + detail
	SendResponse(c, err, nil)
}
