package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"video-api/model"
)

var JwtSecretKey = []byte("my_super_secret_key")

// Claims access token
// access token每次请求使用
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 要知道哪个用户可以刷新token
type RefreshClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateTokens(user *model.User) (accessToken, refreshToken string, err error) {
	//access token 2hours
	accessExpirationTime := time.Now().Add(time.Hour * 2)
	accessClaims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
		},
	}
	accessTokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenJwt.SignedString(JwtSecretKey)
	if err != nil {
		return "", "", err
	}
	//refresh token
	refreshExpirationTime := time.Now().Add(time.Hour * 24 * 7)
	refreshClaims := &RefreshClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
		},
	}
	refreshTokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenJwt.SignedString(JwtSecretKey)
	return accessToken, refreshToken, err
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	//ParseWithClaims 内部会：
	//取 JWT 的 header（算法）
	//取出 payload
	//用你的 JwtSecretKey 重新计算签名
	//与 tokenString 内的签名对比
	//如果不一致 → 表示 token 被改过 → 会返回 error
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
