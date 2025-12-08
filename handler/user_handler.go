package handler

import (
	"strconv"
	"video-api/pkg/errno"
	"video-api/pkg/log"
	"video-api/pkg/upload"

	//"video-api/model"
	//"video-api/repository"
	"video-api/service"
	"video-api/types"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(svc service.IUserService) *UserHandler {
	return &UserHandler{
		userService: svc,
	}
}
func getUserID(c *gin.Context) (uint, bool) {
	userIDVal, exist := c.Get("userID")
	if !exist {
		log.Log.Error("Context中缺少userID",
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
		)
		Error(c, errno.AuthorizationFailedErr)
		return 0, false
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		log.Log.Error("userID类型断言失败",
			zap.String("path", c.Request.URL.Path),
			zap.Any("userID", userIDVal))
		Error(c, errno.ServiceErr)
		return 0, false
	}
	return userID, true
}
func (h *UserHandler) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Log.Warn("用户注册参数绑定失败",
			zap.String("ip", c.ClientIP()),
			zap.Error(err), // 记录具体是哪个字段错了
		)

		// 2. Response: 明确告诉前端参数错了 (40001)
		// 这里的 err.Error() 会包含 "Field validation for 'Username' failed..." 等详细信息
		SendResponse(c, errno.ParamErr, err.Error())
		return
	}
	resq, err := h.userService.Register(&req)
	if err != nil {
		//		Error(c, http.StatusConflict, "REGISTER_FAILD", err.Error())
		log.Log.Error("注册失败",
			zap.String("username", req.Username),
			zap.Error(err))
		SendResponse(c, errno.RegisterErr, err.Error())
		return
	}
	Success(c, resq)
}
func (h *UserHandler) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Log.Warn("参数错误")
		Error(c, errno.ParamErr)
		return
	}
	resq, err := h.userService.Login(&req)
	if err != nil {
		log.Log.Warn("用户登录失败",
			zap.String("username", req.Username),
			zap.String("ip", c.ClientIP()),
			zap.Error(err),
		)
		Error(c, errno.LoginErr)
		return
	}
	log.Log.Info("用户登录成功",
		zap.Uint("user_id", resq.UserID),
		zap.String("username", req.Username))
	Success(c, resq)
}
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	//fmt.Println("5.进入了Handler")
	currentUserId, _ := getUserID(c)
	//val, exists := c.Get("user_id")
	//fmt.Println("6.Handler中获取的user_id：", val, "是否存在：", exists)
	targetUserIDstr := c.Query("user_id")
	if targetUserIDstr == "" {
		log.Log.Warn("目标用户ID不存在",
			zap.String("targetUserID", targetUserIDstr))
		Error(c, errno.ParamErr)
		return
	}
	targetUserID, _ := strconv.ParseUint(targetUserIDstr, 10, 64)
	resq, err := h.userService.GetUserInfo(currentUserId, uint(targetUserID))
	if err != nil {
		log.Log.Warn("用户不存在",
			zap.String("targetUserID", targetUserIDstr),
			zap.Error(err))
		Error(c, errno.UserNotFoundErr)
		return
	}
	Success(c, resq)
}
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	//获取userid
	userID, ok := getUserID(c)
	if !ok {
		//直接返回，因为在getuserid这个函数已经处理了
		return
	}
	file, err := c.FormFile("data")
	if err != nil {
		log.Log.Warn("解析上传文件失败",
			zap.Error(err))
		Error(c, errno.FileUploadErr)
		return
	}
	avatarURL, err := upload.UploadToOSS(file, userID)
	if err != nil {
		//Error(c, http.StatusInternalServerError, "UPLOAD_FAILD", err.Error())
		log.Log.Error("OSS上传失败",
			zap.Uint("user_id", userID),
			zap.Error(err))
		Error(c, errno.FileUploadErr)
		return
	}
	if err := h.userService.UploadAvatar(userID, avatarURL); err != nil {
		log.Log.Error("头像数据入库失败",
			zap.String("url", avatarURL),
			zap.Error(err))
		Error(c, errno.ServiceErr)
		return
	}
	Success(c, gin.H{"avatar_url": avatarURL})
}

// GET /user/mfa/qrcode
func (h *UserHandler) GenerateMFA(c *gin.Context) {
	userId, ok := getUserID(c)
	if !ok {
		return
	}
	resp, err := h.userService.GenerateMFA(userId)
	if err != nil {
		log.Log.Warn("MFA_GENERATE_FAILED")
		Error(c, errno.MFAGenerateErr)
		return
	}
	Success(c, resp)
}

// POST /user/mfa/bind
func (h *UserHandler) BindMFA(c *gin.Context) {
	userId, ok := getUserID(c)
	if !ok {
		return
	}
	var req types.MfaBindRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Log.Warn("无效的请求参数")
		//ValidationError(c, "INVALID_PARAM", "无效的请求参数", err.Error())
		Error(c, errno.ParamErr)
		return
	}
	err := h.userService.BindMFA(userId, req.Secret, req.Code)
	if err != nil {
		log.Log.Warn("MFA绑定失败")
		Error(c, errno.MFABindErr)
		return
	}
	Success(c, gin.H{"msg": "MFA绑定成功"})
}

// POST /tool/search_image
