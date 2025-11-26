package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	//"video-api/model"
	//"video-api/repository"
	"video-api/service"
	"video-api/types"

	"github.com/gin-gonic/gin"
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
		Error(c, http.StatusUnauthorized, "AUTH_ERROR", "无法获取用户信息")
		return 0, false
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		Error(c, http.StatusInternalServerError, "CONTEXT_ERROR", "用户信息格式错误")
		return 0, false
	}
	return userID, true
}
func (h *UserHandler) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		ValidationError(c, "INVALID", "错误的请求参数", err.Error())
		return
	}
	resq, err := h.userService.Register(&req)
	if err != nil {
		Error(c, http.StatusConflict, "REGISTER_FAILD", err.Error())
		return
	}
	Success(c, http.StatusOK, resq)
}
func (h *UserHandler) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		ValidationError(c, "VALIDATION_ERROR", "参数错误", err.Error())
		return
	}
	resq, err := h.userService.Login(&req)
	if err != nil {
		Error(c, http.StatusConflict, "LOGIN_ERROR", err.Error())
	}
	Success(c, http.StatusOK, resq)
}
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	fmt.Println("5.进入了Handler")
	currentUserId, _ := getUserID(c)
	val, exists := c.Get("user_id")
	fmt.Println("6.Handler中获取的user_id：", val, "是否存在：", exists)
	targetUserIDstr := c.Query("user_id")
	if targetUserIDstr == "" {
		Error(c, http.StatusBadRequest, "Auth_FAILD", "缺少user_id")
		return
	}
	targetUserID, _ := strconv.ParseUint(targetUserIDstr, 10, 64)
	resq, err := h.userService.GetUserInfo(currentUserId, uint(targetUserID))
	if err != nil {
		Error(c, http.StatusNotFound, "USER_NOT_FOUND", err.Error())
	}
	Success(c, http.StatusOK, resq)
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
		Error(c, http.StatusBadRequest, "UPLOAD_FAILD", err.Error())
		return
	}
	//构建路径
	filename := fmt.Sprintf("%d-%s", userID, file.Filename)
	saveDir := "./uploads/avatar/"
	dst := filepath.Join(saveDir, filename)
	//确保目录存在
	fmt.Println("保存头像的路径：", dst)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		Error(c, http.StatusInternalServerError, "FILE_SAVE_ERROR", "文件创建失败")
		return
	}
	if err := c.SaveUploadedFile(file, dst); err != nil {
		fmt.Println("保存文件出错：", err)
		Error(c, http.StatusInternalServerError, "UPLOAD_FAILD", "文件保存失败")
		return
	}
	avatarURL := "/static/avatars/" + filename
	// 更新数据库里面的avatar字段
	if err := h.userService.UploadAvatar(userID, avatarURL); err != nil {
		Error(c, http.StatusInternalServerError, "UPLOAD_FAILD", err.Error())
		return
	}
	Success(c, http.StatusOK, gin.H{"avatar_url": avatarURL})
}

//GET /user/mfa/qrcode
func (h *UserHandler) GenerateMFA(c *gin.Context) {
	userId, ok := getUserID(c)
	if !ok {
		return
	}
	resp, err := h.userService.GenerateMFA(userId)
	if err != nil {
		Error(c, http.StatusInternalServerError, "MFA_GENERATE_FAILED", err.Error())
		return
	}
	Success(c, http.StatusOK, resp)
}

//POST /user/mfa/bind
func (h *UserHandler) BindMFA(c *gin.Context) {
	userId, ok := getUserID(c)
	if !ok {
		return
	}
	var req types.MfaBindRequest
	if err := c.ShouldBind(&req); err != nil {
		ValidationError(c, "INVALID_PARAM", "无效的请求参数", err.Error())
		return
	}
	err := h.userService.BindMFA(userId, req.Secret, req.Code)
	if err != nil {
		Error(c, http.StatusBadRequest, "MFA_BIND_FAILED", err.Error())
		return
	}
	Success(c, http.StatusOK, "MFA绑定成功")
}
