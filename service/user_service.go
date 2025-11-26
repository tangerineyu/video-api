package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"video-api/model"
	"video-api/pkg/utils"
	"video-api/repository"
	"video-api/types"

	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
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
	userRepo repository.IUserRepository
	rdb      *redis.Client
	ctx      context.Context
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
	user, err := u.userRepo.FindUserByID(uint(targetUserID))
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
	}, nil

}

func (u *UserService) UploadAvatar(userID uint, avatarUrl string) error {
	return u.userRepo.UpdateAvatar(userID, avatarUrl)
}

func NewUserService(repo repository.IUserRepository, rdb *redis.Client, ctx context.Context) IUserService {
	return &UserService{
		userRepo: repo,
		rdb:      rdb,
		ctx:      ctx,
	}
}
