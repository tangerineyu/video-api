package service

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"video-api/model"
	"video-api/repository"
	"video-api/types"

	"github.com/redis/go-redis/v9"
)

type VideoService struct {
	videoRepo repository.IVideoRepository
	userRepo  repository.IUserRepository
	rdb       *redis.Client
	ctx       context.Context
}

func (v *VideoService) VisitVideo(videoID uint) error {
	//redis 计数+1
	//ZINCRBY video:rank:daily 1 videoID
	err := v.rdb.ZIncrBy(v.ctx, RanKey, 1, fmt.Sprintf("%d", videoID)).Err()
	if err != nil {
		return err
	}
	go v.videoRepo.IncrVisitCount(videoID)
	return nil
}

func (v *VideoService) GetPopularRank() (*types.VideoListResponse, error) {
	videoIDsStr, err := v.rdb.ZRevRange(v.ctx, RanKey, 0, 9).Result()
	if err != nil {
		return nil, err
	}
	var videos []model.Video
	for _, idStr := range videoIDsStr {
		vid, _ := strconv.Atoi(idStr)
		v, _ := v.videoRepo.GetVideosByID(uint(vid))
		if v != nil {
			videos = append(videos, *v)
		}
	}
	return &types.VideoListResponse{
		VideoList: v.packVideoInfos(videos, 0),
	}, nil
}

func (v *VideoService) GetPublishList(targetUserID uint, currentUserID uint) (*types.VideoListResponse, error) {
	videos, err := v.videoRepo.GetVideosByUserID(targetUserID)
	if err != nil {
		return nil, err
	}
	return &types.VideoListResponse{
		VideoList: v.packVideoInfos(videos, currentUserID),
	}, nil
}

func (v *VideoService) PublishVideo(userID uint, title, playPath, coverPath string) error {
	video := &model.Video{
		UserID:      userID,
		Title:       title,
		PlayURL:     playPath,
		CoverURL:    coverPath,
		PublishTime: time.Now(),
	}
	return v.videoRepo.CreateVideo(video)
}

const RanKey = "video:rank:daily"

type IVideoService interface {
	PublishVideo(userID uint, title, plyPath, coverPath string) error
	GetPublishList(targetUserID uint, currentUserID uint) (*types.VideoListResponse, error)
	Feed(latestTime int64, userID uint) (*types.FeedResponse, error)
	Search(keyword string) (*types.VideoListResponse, error)
	VisitVideo(videoID uint) error
	GetPopularRank() (*types.VideoListResponse, error)
}

func NewVideoService(uRepo repository.IUserRepository, vRepo repository.IVideoRepository, rdb *redis.Client, ctx context.Context) IVideoService {
	return &VideoService{
		videoRepo: vRepo,
		userRepo:  uRepo,
		rdb:       rdb,
		ctx:       ctx,
	}
}
func (s *VideoService) Search(keyword string) (*types.VideoListResponse, error) {
	videos, err := s.videoRepo.SearchVideos(keyword)
	if err != nil {
		return nil, err
	}
	return &types.VideoListResponse{
		VideoList: s.packVideoInfos(videos, 0),
	}, nil
}
func (s *VideoService) Feed(latestTime int64, userID uint) (*types.FeedResponse, error) {
	queryTime := time.Now()
	if latestTime != 0 {
		queryTime = time.Unix(latestTime/1000, 0)

	}
	videos, err := s.videoRepo.GetVideosByTime(queryTime, 30)
	if err != nil {
		return nil, err
	}
	nextTime := int64(0)
	if len(videos) > 0 {
		nextTime = videos[len(videos)-1].PublishTime.UnixMilli()
	}
	return &types.FeedResponse{
		NextTime:  nextTime,
		VideoList: s.packVideoInfos(videos, userID),
	}, nil
}
func (s *VideoService) packVideoInfos(videos []model.Video, currentUserID uint) []types.VideoInfo {
	infos := make([]types.VideoInfo, len(videos))
	for _, v := range videos {
		isFavorite := false
		infos = append(infos, types.VideoInfo{
			ID: v.UserID,
			Author: types.UserInfoResponse{
				ID:            v.Author.ID,
				Name:          v.Author.Name,
				FollowCount:   v.Author.FollowCount,
				FollowerCount: v.Author.FollowerCount,
				IsFollow:      false,
				Avatar:        v.Author.Avatar,
			},
			PlayURL:       v.PlayURL,
			CoverURL:      v.CoverURL,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    isFavorite,
			Title:         v.Title,
		})
	}
	return infos
}
