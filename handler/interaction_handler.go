package handler

import (
	"strconv"
	"video-api/pkg/errno"
	"video-api/pkg/log"
	"video-api/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InteractionHandler struct {
	interactionService service.IInteractionService
}

func NewInteractionHandler(svc service.IInteractionService) *InteractionHandler {
	return &InteractionHandler{
		interactionService: svc,
	}
}
func (h *InteractionHandler) FavoriteAction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	videoIDStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")
	videoID, _ := strconv.ParseUint(videoIDStr, 10, 64)
	actionType, _ := strconv.Atoi(actionTypeStr)
	if videoID == 0 || (actionType != 1 && actionType != 2) {
		log.Log.Warn("无效的请求参数, video_id或action_type参数错误",
			zap.String("video_id", videoIDStr),
			zap.String("action_type", actionTypeStr))
		Error(c, errno.ParamErr)
		//ValidationError(c, "INVALID_PARAM", "无效的请求参数", "video_id或action_type参数错误")
		return
	}
	err := h.interactionService.FavoriteAction(userID, uint(videoID), actionType)
	if err != nil {
		log.Log.Warn("点赞/取消操作失败",
			zap.Uint("userID", userID),
			zap.String("videoID", videoIDStr),
			zap.Error(err))
		Error(c, errno.FavoriteActionErr)
		return
	}
	Success(c, gin.H{"msg": "操作成功"})

}
func (h *InteractionHandler) FavoriteList(c *gin.Context) {
	targetUserIDStr := c.Query("user_id")
	targetUserID, _ := strconv.ParseUint(targetUserIDStr, 10, 64)
	if targetUserID == 0 {
		log.Log.Warn("无效的请求参数，user_id参数错误",
			zap.String("user_id", targetUserIDStr))
		Error(c, errno.ParamErr)
		return
	}
	var currentUserID uint = 0
	if id, exists := c.Get("use_id"); exists {
		currentUserID = id.(uint)
	}
	resp, err := h.interactionService.GetFavoriteList(uint(targetUserID), currentUserID)
	if err != nil {
		log.Log.Error("获取点赞列表失败",
			zap.String("target_user_id", targetUserIDStr),
			zap.Error(err))
		Error(c, errno.ServiceErr)
		return
	}
	SuccessList(c, resp.VideoList, nil)
}
func (h *InteractionHandler) CommentAction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	videoIDStr := c.Query("video_id")
	actionTypeStr := c.Query("action_type")

	videoID, _ := strconv.ParseUint(videoIDStr, 10, 64)
	actionType, _ := strconv.Atoi(actionTypeStr)
	if videoID == 0 {
		log.Log.Warn("无效的请求参数, video_id参数错误")
		Error(c, errno.ParamErr)
		return
	}
	var content string
	var commentID uint
	if actionType == 1 {
		content = c.Query("content")
		if content == "" {
			log.Log.Warn("参数错误, 评论内容不能为空",
				zap.Uint("user_id", userID))
			Error(c, errno.ParamErr)
			return
		}
	} else if actionType == 2 {
		commentIDStr := c.Query("comment_id")
		cid, _ := strconv.ParseUint(commentIDStr, 10, 64)
		commentID = uint(cid)
		if commentID == 0 {
			log.Log.Warn("评论ID不能为空",
				zap.Uint("user_id", userID))
			Error(c, errno.ParamErr)
			return
		}

	} else {
		log.Log.Warn("action_type参数错误")
		Error(c, errno.ParamErr)
		return
	}
	commentInfo, err := h.interactionService.CommentAction(userID, uint(videoID), actionType, content, commentID)
	if err != nil {
		log.Log.Error("发布/删除评论操作失败",
			zap.Uint("userID", userID),
			zap.String("videoID", videoIDStr),
			zap.Int("action_type", actionType),
			zap.Error(err))
		Error(c, errno.ServiceErr)
		return
	}
	if actionType == 1 {
		Success(c, commentInfo)
	} else {
		Success(c, gin.H{"msg": "删除评论成功"})
	}

}
func (h *InteractionHandler) CommentList(c *gin.Context) {
	videoIDStr := c.Query("video_id")
	videoID, _ := strconv.ParseUint(videoIDStr, 10, 64)
	if videoID == 0 {
		log.Log.Warn("video_id参数错误")
		Error(c, errno.ParamErr)
		return
	}
	resp, err := h.interactionService.GetCommentList(uint(videoID))
	if err != nil {
		log.Log.Error("获取评论列表失败",
			zap.String("video_id", videoIDStr),
			zap.Error(err))
		Error(c, errno.ServiceErr)
		return
	}
	SuccessList(c, resp.CommentList, nil)
}
