package service

import (
	"context"
	"video-api/repository"
	"video-api/types"

	"github.com/redis/go-redis/v9"
)

type IInteractionService interface {
	Favorite(userID, videoID int, actionType int) error
	GetFavoriteList(userID int) (*types.VideoListResponse, error)
	CommentAction(userID, videoID int, actionType int, commentText string, commentID int) (*types.CommentResponse, error)
	GetCommentList(videoID int) (*types.CommentListResponse, error)
}
type InteractionService struct {
	InteractionRepo repository.InteractionRepository
	rdb             *redis.Client
	ctx             context.Context
}

func (i InteractionService) Favorite(userID, videoID int, actionType int) error {
	//TODO implement me
	panic("implement me")
}

func (i InteractionService) GetFavoriteList(userID int) (*types.VideoListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i InteractionService) CommentAction(userID, videoID int, actionType int, commentText string, commentID int) (*types.CommentResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i InteractionService) GetCommentList(videoID int) (*types.CommentListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewInteractionService(IRepo repository.InteractionRepository, rdb *redis.Client, ctx context.Context) IInteractionService {
	return &InteractionService{
		InteractionRepo: IRepo,
		rdb:             rdb,
		ctx:             ctx,
	}

}
