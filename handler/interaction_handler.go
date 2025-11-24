package handler

import (
	"net/http"
	"strconv"
	"video-api/service"

	"github.com/gin-gonic/gin"
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
		ValidationError(c, "INVALID_PARAM", "无效的请求参数", "video_id或action_type参数错误")
		return
	}
	err := h.interactionService.FavoriteAction(userID, uint(videoID), actionType)
	if err != nil {
		Error(c, http.StatusInternalServerError, "ACTION_FAILED", err.Error())
		return
	}
	Success(c, http.StatusOK, "操作成功")

}
func (h *InteractionHandler) FavoriteList(c *gin.Context) {
	targetUserIDStr := c.Query("user_id")
	targetUserID, _ := strconv.ParseUint(targetUserIDStr, 10, 64)
	if targetUserID == 0 {
		ValidationError(c, "INVALID_PARAM", "无效的请求参数", "user_id参数错误")
		return
	}
	var currentUserID uint = 0
	if id, exists := c.Get("user_id"); exists {
		currentUserID = id.(uint)
	}
	resp, err := h.interactionService.GetFavoriteList(uint(targetUserID), currentUserID)
	if err != nil {
		Error(c, http.StatusInternalServerError, "DB_ERROR", err.Error())
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
		ValidationError(c, "INVALID_PARAM", "无效的请求参数", "video_id参数错误")
		return
	}
	var content string
	var commentID uint
	if actionType == 1 {
		content = c.Query("content")
		if content == "" {
			ValidationError(c, "VALIDATION_FAILED", "参数错误", "评论内容不能为空")
			return
		}
	} else if actionType == 2 {
		commentIDStr := c.Query("comment_id")
		cid, _ := strconv.ParseUint(commentIDStr, 10, 64)
		commentID = uint(cid)
		if commentID == 0 {
			ValidationError(c, "VALIDATION_FAILED", "参数错误", "评论ID不能为空")
			return
		}

	} else {
		ValidationError(c, "VALIDATION_FAILED", "无效的请求参数", "action_type参数错误")
		return
	}
	commentInfo, err := h.interactionService.CommentAction(userID, uint(videoID), actionType, content, commentID)
	if err != nil {
		Error(c, http.StatusInternalServerError, "ACTION_FAILED", err.Error())
		return
	}
	if actionType == 1 {
		Success(c, http.StatusOK, commentInfo)
	} else {
		Success(c, http.StatusOK, gin.H{"msg": "删除评论成功"})
	}

}
func (h *InteractionHandler) CommentList(c *gin.Context) {
	videoIDStr := c.Query("video_id")
	videoID, _ := strconv.ParseUint(videoIDStr, 10, 64)
	if videoID == 0 {
		ValidationError(c, "INVALID_PARAM", "无效的请求参数", "video_id参数错误")
		return
	}
	resp, err := h.interactionService.GetCommentList(uint(videoID))
	if err != nil {
		Error(c, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}
	SuccessList(c, resp.CommentList, nil)
}
