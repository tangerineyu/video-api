package handler

import (
	"fmt"
	"strconv"
	"video-api/pkg/errno"
	"video-api/pkg/log"
	"video-api/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	fmt.Println("社交操作，当前用户ID：", userID)
	if !ok {
		return
	}
	toUserIDStr := c.Query("to_user_id")
	toUserIDUint, err := strconv.ParseUint(toUserIDStr, 10, 64)
	if err != nil || toUserIDUint == 0 {
		log.Log.Warn("to_user_id参数错误",
			zap.String("to_user_id", toUserIDStr))
		Error(c, errno.ParamErr)
		return
	}
	toUserID := uint(toUserIDUint)
	fmt.Println("目标用户ID：", toUserID)
	actionTypeStr := c.Query("action_type")
	actionType, _ := strconv.Atoi(actionTypeStr)
	if (actionType != 1 && actionType != 2) || toUserID == 0 {
		log.Log.Warn("action_type参数错误")
		Error(c, errno.ParamErr)
		return
	}
	if toUserID == userID {
		log.Log.Warn("不能关注自己",
			zap.String("to_user_id", toUserIDStr))
		Error(c, errno.PermissionDeniedErr)
		return
	}
	err = h.socialService.RelationAction(userID, toUserID, actionType)
	if err != nil {
		log.Log.Error("关注/取消关注异常",
			zap.String("to_user_id", toUserIDStr),
			zap.Int("action_type", actionType),
			zap.Error(err))
		Error(c, errno.ServiceErr)
		return
	}
	Success(c, gin.H{"msg": "操作成功"})
}

// 关注列表
func (h *SocialHandler) FollowList(c *gin.Context) {
	targetUserIDStr := c.Query("user_id")
	targetUserID, _ := strconv.ParseUint(targetUserIDStr, 10, 64)
	if targetUserID == 0 {
		log.Log.Warn("user_id参数错误",
			zap.String("target_user_id", targetUserIDStr))
		Error(c, errno.ParamErr)
		return
	}
	var currentUserID uint = 0
	if id, exists := c.Get("userID"); exists {
		currentUserID = id.(uint)
	}
	resp, err := h.socialService.GetFollowList(uint(targetUserID), currentUserID)
	if err != nil {
		log.Log.Error("获取关注列表失败",
			zap.Uint("current_user_id", currentUserID),
			zap.String("target_user_id", targetUserIDStr),
			zap.Error(err))
		Error(c, errno.ServiceErr)
		return
	}
	SuccessList(c, resp.UserList, nil)
}

// 粉丝列表
func (h *SocialHandler) FollowerList(c *gin.Context) {
	targetUserIDStr := c.Query("user_id")
	targetUserID, _ := strconv.ParseUint(targetUserIDStr, 10, 64)
	if targetUserID == 0 {
		log.Log.Warn("user_id参数错误",
			zap.String("target_user_id", targetUserIDStr))
		Error(c, errno.ParamErr)
		return
	}
	var currentUserID uint = 0
	if id, exists := c.Get("userID"); exists {
		currentUserID = id.(uint)
	}
	resp, err := h.socialService.GetFollowerList(uint(targetUserID), currentUserID)
	if err != nil {
		log.Log.Error("获取粉丝列表失败",
			zap.Uint("current_user_id", currentUserID),
			zap.String("target_user_id", targetUserIDStr),
			zap.Error(err))
		Error(c, errno.ServiceErr)
		return
	}
	SuccessList(c, resp.UserList, nil)
}

// 好友列表
// GET /relation/friend/list/?user_id=1
func (h *SocialHandler) FriendList(c *gin.Context) {
	targetUserIDStr := c.Query("user_id")
	targetUserID, _ := strconv.ParseUint(targetUserIDStr, 10, 64)
	if targetUserID == 0 {
		log.Log.Warn("user_id参数错误",
			zap.String("target_user_id", targetUserIDStr))
		Error(c, errno.ParamErr)
		return
	}
	var currentUserID uint = 0
	if id, exists := c.Get("userID"); exists {
		currentUserID = id.(uint)
	}
	resp, err := h.socialService.GetFriendList(uint(targetUserID), currentUserID)
	if err != nil {
		log.Log.Error("获取好友列表失败",
			zap.Uint("current_user_id", currentUserID),
			zap.String("target_user_id", targetUserIDStr),
			zap.Error(err))
		Error(c, errno.ServiceErr)
		return
	}
	SuccessList(c, resp.UserList, nil)
}
