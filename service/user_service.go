package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"video-api/model"
	"video-api/pkg/utils"
	"video-api/repository"
	"video-api/types"

	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

type IUserService interface {
	Register(req *types.RegisterRequest) (*types.LoginResponse, error)
	Login(req *types.LoginRequest) (*types.LoginResponse, error)
	GetUserInfo(currentUserID uint, targetUserId uint) (*types.UserInfoResponse, error)
	UploadAvatar(userID uint, avatarUrl string) error
	GenerateMFA(userID uint) (*types.MfaGenerateResponse, error)
	BindMFA(userID uint, secret string, code string) error
}
type UserService struct {
	userRepo   repository.IUserRepository
	socialRepo repository.ISocialRepository //fillFollowStatus,获得关注的数据
	rdb        *redis.Client
	ctx        context.Context
	//防止缓存击穿
	sf singleflight.Group
}

func (s *UserService) GenerateMFA(userID uint) (*types.MfaGenerateResponse, error) {
	user, _ := s.userRepo.FindUserByID(userID)
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "FanOneVideo",
		AccountName: user.Username,
	})
	if err != nil {
		return nil, err
	}
	return &types.MfaGenerateResponse{
		Secret: key.Secret(),
		Qrcode: key.URL(),
	}, nil
}

func (s *UserService) BindMFA(userID uint, secret string, code string) error {
	valid := totp.Validate(code, secret)
	if !valid {
		return errors.New("invalid mfa code")
	}
	return s.userRepo.EnableMFA(userID, secret)
}

func (s *UserService) generateAndstoreTokens(user *model.User) (*types.LoginResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokens(user)
	if err != nil {
		return nil, errors.New("error generating tokens")
	}
	rtKey := fmt.Sprintf("refreshToken:%d", user.ID)
	err = s.rdb.Set(s.ctx, rtKey, accessToken, time.Hour*7*24).Err()
	if err != nil {
		return nil, errors.New("error storing tokens")
	}
	return &types.LoginResponse{
		UserID:       user.ID,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}
func (u *UserService) Register(req *types.RegisterRequest) (*types.LoginResponse, error) {
	_, err := u.userRepo.FindUserByUsername(req.Username)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("error hashing password")
	}
	user := &model.User{
		Username: req.Username,
		Password: hashPassword,
		Name:     req.Username,
	}
	if err := u.userRepo.CreateUser(user); err != nil {
		return nil, errors.New("error creating user")
	}
	return u.generateAndstoreTokens(user)
}

func (u *UserService) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	user, err := u.userRepo.FindUserByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("error getting user from db")
	}
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid password")
	}
	return u.generateAndstoreTokens(user)
}

func (u *UserService) GetUserInfo(currentUserID uint, targetUserID uint) (*types.UserInfoResponse, error) {
	/**user, err := u.userRepo.FindUserByID(uint(targetUserID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("error getting user from db")
	}
	isFollow := false
	return &types.UserInfoResponse{
		ID:            user.ID,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		Avatar:        user.Avatar,
		IsFollow:      isFollow,
	}, nil**/
	//
	cacheKey := fmt.Sprintf("user:info:%d", targetUserID)
	//查redis
	val, err := u.rdb.Get(u.ctx, cacheKey).Result()
	if err == nil {
		//缓存命中
		var userInfo types.UserInfoResponse
		if json.Unmarshal([]byte(val), &userInfo) != nil {
			u.fillFollowStatus(&userInfo, currentUserID)
			return &userInfo, nil
		}
	}
	//缓存未命中，使用singleflight防止缓存击穿
	//Do方法接受两个参数，key：标识这是同一个请求的字符串，fn：干活的函数，查库写缓存
	v, err, _ := u.sf.Do(cacheKey, func() (interface{}, error) {
		//  DB
		fmt.Println("查数据库", targetUserID)
		user, err := u.userRepo.FindUserByID(targetUserID)
		if err != nil {
			return nil, err
		}
		// return
		info := types.UserInfoResponse{
			ID:            user.ID,
			Name:          user.Username,
			FollowerCount: user.FollowerCount,
			FollowCount:   user.FollowCount,
			Avatar:        user.Avatar,
			IsFollow:      false,
		}
		//  设置缓存1小时过期，加上一点随机时间防止雪崩
		data, err := json.Marshal(info)
		ttl := time.Hour + time.Duration(time.Now().UnixNano()%60000)*time.Millisecond
		u.rdb.Set(u.ctx, cacheKey, data, ttl)
		return &info, nil
	})
	if err != nil {
		return nil, err
	}
	//类型断言，将interface{}转回UserInfoResponse
	userInfo := v.(*types.UserInfoResponse)
	//针对IsFollow重新计算
	u.fillFollowStatus(userInfo, targetUserID)
	return userInfo, nil
}

// 辅助函数,关注状态
// 基础信息是公共的，所有人查看都一样，比如名字，粉丝数，这部分可以使用缓存
// isFollow是私有的，对不同的用户来说是有’已关注‘’未关注`两个状态
// 拿到缓存后，再动态的覆盖IsFollow字段
func (u *UserService) fillFollowStatus(info *types.UserInfoResponse, currentUserID uint) {
	if currentUserID == 0 || currentUserID == info.ID {
		info.IsFollow = false
		return
	}
	//使用socialRepo查数据库
	isFollow, err := u.socialRepo.IsFollowing(currentUserID, info.ID)
	if err != nil {
		info.IsFollow = false
	} else {
		info.IsFollow = isFollow
	}
}
func (u *UserService) UploadAvatar(userID uint, avatarUrl string) error {
	return u.userRepo.UpdateAvatar(userID, avatarUrl)
}

func NewUserService(uRepo repository.IUserRepository, sRepo repository.ISocialRepository, rdb *redis.Client, ctx context.Context) IUserService {
	return &UserService{
		userRepo:   uRepo,
		socialRepo: sRepo,
		rdb:        rdb,
		ctx:        ctx,
		//sf 不需要初始化，零值可用
	}
}
