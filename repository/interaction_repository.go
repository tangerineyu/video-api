package repository

import (
	"video-api/model"

	"gorm.io/gorm"
)

type interactionRepository struct {
	db *gorm.DB
}

func (r *interactionRepository) IsFavorite(userID uint, videoID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.UserFavorite{}).
		Where("user_id = ? AND video_id = ?", userID, videoID).
		Count(&count).Error
	return count > 0, err
}

func (r interactionRepository) AddFavorite(userID uint, videoID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&model.UserFavorite{UserID: userID, VideoID: videoID}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Video{}).Where("id = ?", videoID).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r interactionRepository) RemoveFavorite(userID uint, videoID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&model.UserFavorite{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Video{}).Where("id = ?", videoID).Delete(&model.Video{}).Error; err != nil {
			return nil
		}
		return nil
	})
}

func (r interactionRepository) GetFavoriteVideoList(userID uint) ([]model.Video, error) {
	var videos []model.Video
	err := r.db.Joins("JOIN user_favorites ON user_favorites.video_id = videos.id").
		Where("user_favorites.user_id = ?", userID).
		Preload("Author").
		Find(&videos).Error
	return videos, err
}

func (r interactionRepository) CreateComment(comment *model.Comment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		return tx.Model(&model.Video{}).Where("id = ?", comment.VideoID).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	})
}

func (r interactionRepository) DeleteComment(commentID uint, videoID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Comment{}, commentID).Error; err != nil {
			return err
		}
		return tx.Model(&model.Video{}).Where("id = ?", videoID).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error
	})
}

func (r interactionRepository) GetCommentsByVideoID(commentID uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.First(&comment, commentID).Error
	return &comment, err
}

func (r interactionRepository) GetCommentList(videoID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Preload("User").
		Where("video_id = ?", videoID).
		Order("created_at desc").
		Find(&comments).Error
	return comments, err
}

type InteractionRepository interface {
	//
	IsFavorite(userID uint, videoID uint) (bool, error)
	AddFavorite(userID uint, videoID uint) error
	RemoveFavorite(userID uint, videoID uint) error
	GetFavoriteVideoList(userID uint) ([]model.Video, error)
	//
	CreateComment(comment *model.Comment) error
	DeleteComment(commentID uint, videoID uint) error
	GetCommentsByVideoID(videoID uint) (*model.Comment, error)
	GetCommentList(videoID uint) ([]model.Comment, error)
	//
}

func NewInteractionRepository(db *gorm.DB) InteractionRepository {
	return &interactionRepository{db: db}
}
