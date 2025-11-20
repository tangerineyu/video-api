package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"video-api/service"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoService service.IVideoService
}

func NewVideoHandler(svc service.IVideoService) *VideoHandler {
	return &VideoHandler{
		videoService: svc,
	}
}
func (h *VideoHandler) Publish(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		Error(c, http.StatusBadRequest, "INVALID_PARAM", "无效的请求参数")
		return
	}
	filename := fmt.Sprintf("%d_%d_%s", userID, time.Now().UnixNano(), data.Filename)
	videoSaveDir := "./uploads/videos/"
	videoPath := filepath.Join(videoSaveDir, filename)
	if err := os.MkdirAll(videoSaveDir, 0755); err != nil {
		Error(c, http.StatusInternalServerError, "SERVER_ERROR", "服务器错误")
		return
	}
	if err := c.SaveUploadedFile(data, videoPath); err != nil {
		Error(c, http.StatusInternalServerError, "SAVE_FILE_ERROR", "保存文件失败")
		return
	}
	coverPath := "./uploads/covers/" + filename + ".jpg"
	playURL := "/static/videos/" + filename
	err = h.videoService.PublishVideo(userID, title, playURL, coverPath)
	if err != nil {
		Error(c, http.StatusInternalServerError, "PUBLISH_ERROR", "发布视频失败")
		return
	}
	Success(c, http.StatusOK, "视频发布成功")
}
func (h *VideoHandler) Feed(c *gin.Context) {
	var userID uint = 0
	if id, exists := c.Get("userID"); exists {
		userID = id.(uint)
	}
	latestTimeStr := c.Query("latest_time")
	var latestTime int64 = 0
	if latestTimeStr != "" {
		latestTime, _ = strconv.ParseInt(latestTimeStr, 10, 64)
	}
	resp, err := h.videoService.Feed(latestTime, userID)
	if err != nil {
		Error(c, http.StatusInternalServerError, "DB_ERROR", "获取视频流失败")
		return
	}
	Success(c, http.StatusOK, resp)

}
func (h *VideoHandler) List(c *gin.Context) {
	targetUserIDstr := c.Query("user_id")
	targerUserID, _ := strconv.ParseUint(targetUserIDstr, 10, 64)
	var currentUserID uint = 0
	if id, exists := c.Get("userID"); exists {
		currentUserID = id.(uint)
	}
	resp, err := h.videoService.GetPublishList(uint(targerUserID), currentUserID)
	if err != nil {
		Error(c, http.StatusInternalServerError, "DB_ERROR", "获取发布列表失败")
		return
	}
	Success(c, http.StatusOK, resp)
}
