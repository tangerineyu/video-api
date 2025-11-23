package handler

import (
	"net/http"
	"strconv"
	"video-api/service"

	"github.com/gin-gonic/gin"
)

type SocialHandler struct {
	socialService service.ISocialService
}

func NewSocialHandler(svc service.ISocialService) *SocialHandler {
	return &SocialHandler{
		socialService: svc,
	}
}

// 关注/取消关注
func (h *SocialHandler) RelationAction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	toUserID, ok := getUserID(c)
	if !ok {
		return
	}
	actionTypeStr := c.Query("action_type")
	actionType, _ := strconv.Atoi(actionTypeStr)
	if (actionType != 1 && actionType != 2) || toUserID == 0 {
		ValidationError(c, "INVALID_PARAM", "无效的请求参数", "action_type参数错误")
		return
	}
	if uint(toUserID) == userID {
		Error(c, http.StatusBadRequest, "ACTION_INVALID", "不能关注自己")
		return
	}
	err := h.socialService.RelationAction(userID, toUserID, actionType)
	if err != nil {
		Error(c, http.StatusInternalServerError, "ACTION_FAILED", err.Error())
		return
	}
	Success(c, http.StatusOK, "操作成功")
}

// 关注列表
func (h *SocialHandler) FollowList(c *gin.Context) {
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
	resp, err := h.socialService.GetFollowList(uint(targetUserID), currentUserID)
	if err != nil {
		Error(c, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}
	SuccessList(c, resp.UserList, nil)
}

// 粉丝列表
func (h *SocialHandler) FollowerList(c *gin.Context) {
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
	resp, err := h.socialService.GetFollowerList(uint(targetUserID), currentUserID)
	if err != nil {
		Error(c, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}
	SuccessList(c, resp.UserList, nil)
}

// 好友列表
//GET /relation/friend/list/?user_id=1
func (h *SocialHandler) FriendList(c *gin.Context) {
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
	resp, err := h.socialService.GetFriendList(uint(targetUserID), currentUserID)
	if err != nil {
		Error(c, http.StatusInternalServerError, "DB_ERROR", err.Error())
		return
	}
	SuccessList(c, resp.UserList, nil)
}
