package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
	"video-api/model"
	"video-api/pkg/utils"
	"video-api/repository"
	"video-api/types"
)

type IUserService interface {
	Register(req *types.RegisterRequest) (*types.LoginResponse, error)
	Login(req *types.LoginRequest) (*types.LoginResponse, error)
	GetUserInfo(currentUserID uint, targetUserId uint) (*types.UserInfoResponse, error)
	UploadAvatar(userID uint, avatarUrl string) error
}
type UserService struct {
	userRepo repository.IUserRepository
	rdb      *redis.Client
	ctx      context.Context
}

func (u UserService) Register(req *types.RegisterRequest) (*types.LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) GetUserInfo(currentUserID uint, targetUserId uint) (*types.UserInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserService) UploadAvatar(userID uint, avatarUrl string) error {
	//TODO implement me
	panic("implement me")
}

func NewUserService(repo repository.IUserRepository, rdb *redis.Client, ctx context.Context) IUserService {
	return &UserService{
		userRepo: repo,
		rdb:      rdb,
		ctx:      ctx,
	}
}
