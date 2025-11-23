package service

import (
	"context"
	"errors"
	"video-api/model"
	"video-api/repository"
	"video-api/types"

	"github.com/redis/go-redis/v9"
)

type ISocialService interface {
	RelationAction(userID, toUserID uint, actionType int) error
	GetFollowList(userID, currentUserID uint) (*types.UserListResponse, error)
	GetFollowerList(userID, currentUserID uint) (*types.UserListResponse, error)
	GetFriendList(userID, currentUserID uint) (*types.UserListResponse, error)
}
type SocialService struct {
	SocialRepo repository.ISocialRepository
	UsrRepo    repository.IUserRepository
	rdb        *redis.Client
	ctx        context.Context
}

func (s SocialService) RelationAction(userID, toUserID uint, actionType int) error {
	if userID == toUserID {
		return errors.New("不能关注自己")
	}
	_, err := s.UsrRepo.FindUserByID(toUserID)
	if err != nil {
		return errors.New("用户不存在")
	}
	if actionType == 1 {
		isFollow, _ := s.SocialRepo.IsFollowing(userID, toUserID)
		if isFollow {
			return nil
		}
		return s.SocialRepo.ActionFollow(userID, toUserID)
	}
	if actionType == 2 {
		isFollow, _ := s.SocialRepo.IsFollowing(userID, toUserID)
		if !isFollow {
			return errors.New("未关注该用户，无法取消关注")
		}
		return s.SocialRepo.ActionUnfollow(userID, toUserID)
	}
	return errors.New("无效的操作类型")

}

func (s *SocialService) GetFollowList(userID, currentUserID uint) (*types.UserListResponse, error) {
	users, err := s.SocialRepo.GetFollowList(userID)
	if err != nil {
		return nil, err
	}
	return s.packUserList(users, currentUserID)
}

func (s SocialService) GetFollowerList(userID, currentUserID uint) (*types.UserListResponse, error) {
	users, err := s.SocialRepo.GetFollowerList(userID)
	if err != nil {
		return nil, err
	}
	return s.packUserList(users, currentUserID)
}

func (s SocialService) GetFriendList(userID, currentUserID uint) (*types.UserListResponse, error) {
	users, err := s.SocialRepo.GetFriendList(userID)
	if err != nil {
		return nil, err
	}
	return s.packUserList(users, currentUserID)
}

func NewSocialService(socialRepo repository.ISocialRepository, userRepo repository.IUserRepository, rdb *redis.Client, ctx context.Context) ISocialService {
	return &SocialService{
		SocialRepo: socialRepo,
		UsrRepo:    userRepo,
		rdb:        rdb,
		ctx:        ctx,
	}
}
func (s *SocialService) packUserList(users []model.User, currentUserID uint) (*types.UserListResponse, error) {
	userList := make([]types.UserInfoResponse, 0, len(users))
	for _, u := range users {
		isFollow := false
		if currentUserID != 0 {
			if currentUserID == u.ID {
				isFollow = false
			} else {
				isFollow, _ = s.SocialRepo.IsFollowing(currentUserID, u.ID)
			}
		}
		userList = append(userList, types.UserInfoResponse{
			ID:            u.ID,
			Name:          u.Name,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
			IsFollow:      isFollow,
			Avatar:        u.Avatar,
		})
	}
	return &types.UserListResponse{
		UserList: userList,
	}, nil
}
