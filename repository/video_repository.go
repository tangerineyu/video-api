package repository

import (
	"gorm.io/gorm"
	"time"
	"video-api/model"
)

type videoRepository struct {
	db *gorm.DB
}

func (v videoRepository) CreateVideo(video *model.Video) error {
	return v.db.Create(video).Error
}

func (v videoRepository) GetVideosByTime(lastTime time.Time, limit int) ([]model.Video, error) {
	var videos []model.Video
	result := v.db.Preload("Author").
		Where("created_at < ?", lastTime).
		Order("created_at desc").
		Limit(limit).Find(&videos)
	return videos, result.Error
}

func (v videoRepository) GetVideosByUserID(userID uint) ([]model.Video, error) {
	var videos []model.Video
	result := v.db.Preload("Author").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&videos)
	return videos, result.Error
}

func (v videoRepository) SearchVideos(keywords string) ([]model.Video, error) {
	var videos []model.Video
	kw := "%" + keywords + "%"
	result := v.db.Preload("Author").
		Where("title LIKE ?", kw).
		Order("created_at desc").
		Find(&videos)
	return videos, result.Error
}

func (v videoRepository) GetVideosByID(videoID uint) (*model.Video, error) {
	var video model.Video
	result := v.db.First(&video, videoID)
	return &video, result.Error
}

type IVideoRepository interface {
	CreateVideo(video *model.Video) error
	GetVideosByTime(lastTime time.Time, limit int) ([]model.Video, error)
	GetVideosByUserID(userID uint) ([]model.Video, error)
	SearchVideos(keywords string) ([]model.Video, error)
	GetVideosByID(videoID uint) (*model.Video, error)
}

func NewVideoRepository(db *gorm.DB) IVideoRepository {
	return &videoRepository{db: db}

}
