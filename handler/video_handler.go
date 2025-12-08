package handler

import (
	"strconv"
	"video-api/pkg/errno"
	"video-api/pkg/log"
	"video-api/pkg/upload"
	"video-api/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
		log.Log.Warn("无效的请求参数",
			zap.String("video_title", title),
			zap.Uint("user_id", userID),
			zap.Error(err))
		Error(c, errno.ParamErr)
		return
	}
	playURL, err := upload.UploadToOSS(data, userID)
	if err != nil {
		log.Log.Error("视频上传到OSS失败",
			zap.String("filename", data.Filename),
			zap.Uint("user_id", userID),
			zap.Error(err))
		Error(c, errno.ServiceErr)
		return
	}
	coverURL := playURL + "?x-oss-process=video/snapshot,t_0,f_jpg"

	err = h.videoService.PublishVideo(userID, title, playURL, coverURL)
	if err != nil {
		log.Log.Error("发布视频失败",
			zap.String("playURL", playURL),
			zap.Uint("user_id", userID),
			zap.Error(err))
		Error(c, errno.PublishErr)
		return
	}
	Success(c, gin.H{"msg": "视频发布成功"})
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
		log.Log.Error("获取视频流失败",
			zap.Uint("user_id", userID),
			zap.Error(err))
		Error(c, errno.ServiceErr)
		return
	}
	Success(c, resp)

}
func (h *VideoHandler) List(c *gin.Context) {
	targetUserIDstr := c.Query("user_id")
	targetUserID, _ := strconv.ParseUint(targetUserIDstr, 10, 64)
	var currentUserID uint = 0
	if id, exists := c.Get("userID"); exists {
		currentUserID = id.(uint)
	}
	resp, err := h.videoService.GetPublishList(uint(targetUserID), currentUserID)
	if err != nil {
		log.Log.Error("获取发布列表失败",
			zap.String("target_user_id", targetUserIDstr),
			zap.Uint("user_id", currentUserID),
			zap.Error(err))
		Error(c, errno.ServiceErr)
		return
	}
	Success(c, resp)
}

// POST /video/visit/:id
func (h *VideoHandler) VisitVideo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		log.Log.Warn("访问视频ID无效",
			zap.String("id_str", idStr))
		Error(c, errno.ParamErr)
		return
	}
	if err := h.videoService.VisitVideo(uint(id)); err != nil {
		log.Log.Error("增加访问量失败",
			zap.Int("video_id", id),
			zap.Error(err))
		Error(c, errno.ServiceErr)
	}
	Success(c, gin.H{"msg": "访问量增加成功"})
}
func (h *VideoHandler) PopularRank(c *gin.Context) {
	resp, err := h.videoService.GetPopularRank()
	if err != nil {
		log.Log.Error("获取热门视频排行失败",
			zap.Error(err))
		Error(c, errno.ServiceErr)
		return
	}
	SuccessList(c, resp.VideoList, nil)

}
