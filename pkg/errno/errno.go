package errno

type ErrNo struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func (e ErrNo) Code() int64 {
	return e.StatusCode
}
func (e ErrNo) Msg() string {
	return e.StatusMsg
}
func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{code, msg}
}

var (
	Success = NewErrNo(0, "SUCCESS")
	//400
	ParamErr            = NewErrNo(40001, "请求参数错误")
	LoginErr            = NewErrNo(40002, "登录发生错误")
	RegisterErr         = NewErrNo(40003, "注册发生错误")
	UserAlreadyExistErr = NewErrNo(40004, "用户名已经存在")
	TokenInvalidErr     = NewErrNo(40005, "token无效或过期")
	UserNotFoundErr     = NewErrNo(40006, "无法获取用户信息")
	FileUploadErr       = NewErrNo(40007, "文件上传失败")
	PublishErr          = NewErrNo(40008, "视频发布时出现错误")
	MFAGenerateErr      = NewErrNo(40010, "MFA生成失败")
	MFABindErr          = NewErrNo(40011, "MFA绑定失败")
	FavoriteActionErr   = NewErrNo(40012, "关注操作失败")
	GetFeedErr          = NewErrNo(40013, "获取视频流失败")
	//401 unauthorized
	AuthorizationFailedErr = NewErrNo(40101, "用户未登录")
	//403 Forbidden
	PermissionDeniedErr = NewErrNo(40301, "操作非法")
	//404
	CommentNotFoundErr = NewErrNo(40401, "评论不存在")
	VideoNotFoundErr   = NewErrNo(40402, "视频不存在")
	//500
	ServiceErr = NewErrNo(50001, "服务器内部错误")
)

func ConvertErr(err error) ErrNo {
	if err == nil {
		return Success
	}
	return NewErrNo(ServiceErr.StatusCode, err.Error())
}
