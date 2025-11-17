package types

// 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 登录响应
type LoginResponse struct {
	UserID       uint   `json:"user_id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// 用户信息响应
type UserInfoResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
	Avatar        string `json:"avatar"`
}
