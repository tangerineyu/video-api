package service

import (
	"context"
	"errors"
	"video-api/model"
	"video-api/repository"
	"video-api/types"

	"github.com/redis/go-redis/v9"
)

type IInteractionService interface {
	FavoriteAction(userID, videoID uint, actionType int) error
	GetFavoriteList(userID uint, currentUserID uint) (*types.VideoListResponse, error)
	CommentAction(userID, videoID uint, actionType int, commentText string, commentID uint) (*types.CommentInfo, error)
	GetCommentList(videoID uint) (*types.CommentListResponse, error)
}
type InteractionService struct {
	InteractionRepo repository.InteractionRepository
	UserRepo        repository.IUserRepository
	rdb             *redis.Client
	ctx             context.Context
}

func (i *InteractionService) GetCommentList(videoID uint) (*types.CommentListResponse, error) {
	comments, err := i.InteractionRepo.GetCommentList(videoID)
	if err != nil {
		return nil, errors.New("failed to get comments")
	}
	list := make([]types.CommentInfo, 0)
	for _, comment := range comments {
		list = append(list, types.CommentInfo{
			ID: comment.ID,
			User: types.UserInfoResponse{
				ID:     comment.User.ID,
				Name:   comment.User.Name,
				Avatar: comment.User.Avatar,
			},
			Content:         comment.Content,
			CreatedDateBase: comment.CreatedAt.Format("01-02"),
		})
	}
	return &types.CommentListResponse{
		CommentList: list,
	}, nil
}

func (i *InteractionService) FavoriteAction(userID, videoID uint, actionType int) error {
	if actionType == 1 {
		if fav, _ := i.InteractionRepo.IsFavorite(userID, videoID); fav {
			return nil
		}
		return i.InteractionRepo.AddFavorite(userID, videoID)
	}
	return i.InteractionRepo.RemoveFavorite(userID, videoID)
}

func (i *InteractionService) GetFavoriteList(targetUserID uint, currentUserID uint) (*types.VideoListResponse, error) {
	videos, err := i.InteractionRepo.GetFavoriteVideoList(targetUserID)
	if err != nil {
		return nil, errors.New("failed to get videos")
	}
	list := make([]types.VideoInfo, 0)
	for _, video := range videos {
		list = append(list, types.VideoInfo{
			ID: video.UserID,
			Author: types.UserInfoResponse{
				ID:            video.Author.ID,
				Name:          video.Author.Name,
				FollowCount:   video.Author.FollowCount,
				FollowerCount: video.Author.FollowerCount,
				IsFollow:      false,
				Avatar:        video.Author.Avatar,
			},
			PlayURL:       video.PlayURL,
			CoverURL:      video.CoverURL,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		})
	}
	return &types.VideoListResponse{
		VideoList: list,
	}, nil
}

func (i *InteractionService) CommentAction(userID, videoID uint, actionType int, commentText string, commentID uint) (*types.CommentInfo, error) {
	if actionType == 1 {
		comment := &model.Comment{UserID: userID, VideoID: videoID, Content: commentText}
		if err := i.InteractionRepo.CreateComment(comment); err != nil {
			return nil, err
		}
		user, _ := i.UserRepo.FindUserByID(userID)
		return &types.CommentInfo{
			ID: comment.ID,
			User: types.UserInfoResponse{
				ID:     user.ID,
				Name:   user.Name,
				Avatar: user.Avatar,
			},
			Content:         comment.Content,
			CreatedDateBase: comment.CreatedAt.Format("01-02"),
		}, nil
	}
	//删除评论
	if actionType == 2 {
		//评论是否存在
		comment, err := i.InteractionRepo.GetCommentsByID(commentID)
		if err != nil {
			return nil, errors.New("comment not found")
		}
		//检查评论的ID与当前请求者的ID
		if comment.UserID != userID {
			return nil, errors.New("unauthorized to delete this comment")
		}
		if err := i.InteractionRepo.DeleteComment(commentID, videoID); err != nil {
			return nil, err
		}
		return nil, nil
	}
	return nil, errors.New("无效的操作类型")

}

func NewInteractionService(interactionRepo repository.InteractionRepository, userRepo repository.IUserRepository, rdb *redis.Client, ctx context.Context) IInteractionService {
	return &InteractionService{
		InteractionRepo: interactionRepo,
		UserRepo:        userRepo,
		rdb:             rdb,
		ctx:             ctx,
	}
}
