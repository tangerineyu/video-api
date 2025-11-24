package repository

import (
	"video-api/model"

	"gorm.io/gorm"
)

type ISocialRepository interface {
	IsFollowing(userID, followID uint) (bool, error)
	ActionFollow(userID, followID uint) error
	ActionUnfollow(userID, followID uint) error
	GetFollowList(userID uint) ([]model.User, error)
	GetFollowerList(userID uint) ([]model.User, error)
	GetFriendList(userID uint) ([]model.User, error)
}
type socialRepository struct {
	db *gorm.DB
}

func (s *socialRepository) IsFollowing(userID, followID uint) (bool, error) {
	var count int64
	err := s.db.Model(&model.UserRelation{}).Where("user_id = ? AND follow_id = ?", userID, followID).
		Count(&count).Error
	return count > 0, err
}

func (s socialRepository) ActionFollow(userID, followID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		//创建关系
		if err := tx.Create(&model.UserRelation{UserID: userID, FollowID: followID}).Error; err != nil {
			return err
		}
		//我的关注数加1
		if err := tx.Model(&model.User{}).Where("id = ?", userID).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1)).Error; err != nil {
			return err
		}

		//对方粉丝数加1
		return tx.Model(model.User{}).Where("id = ?", followID).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1)).Error
	})
}

func (s socialRepository) ActionUnfollow(userID, followID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		//删除关系
		if err := tx.Delete(&model.UserRelation{UserID: userID, FollowID: followID}).Error; err != nil {
			return err
		}
		//我的关注数减1
		if err := tx.Model(&model.User{}).Where("id = ?", followID).UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1)).Error; err != nil {
			return err
		}
		//对方粉丝数减1
		return tx.Model(model.User{}).Where("id = ?", followID).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1)).Error

	})
}

func (s socialRepository) GetFollowList(userID uint) ([]model.User, error) {
	//TODO implement me
	var users []model.User
	err := s.db.Joins("JOIN user_relations ON user_relations.follow_id = users.id").
		Where("user_relations.user_id = ?", userID).
		Find(&users).Error
	return users, err
}

func (s socialRepository) GetFollowerList(userID uint) ([]model.User, error) {
	//TODO implement me
	var users []model.User
	err := s.db.Joins("JOIN user_relations ON user_relations.user_id = users.id").
		Where("user_relations.follow_id = ?", userID).
		Find(&users).Error
	return users, err
}

func (s socialRepository) GetFriendList(userID uint) ([]model.User, error) {

	var friends []model.User
	err := s.db.Raw(`
		SELECT u.* FROM users u
		JOIN user_relations r1 ON u.id = r1.follow_id
		JOIN user_relations r2 ON u.id = r2.user_id
		WHERE r1.user_id = ? AND r2.follow_id = ?`, userID, userID).Scan(&friends).Error
	return friends, err

}

func NewSocialRepository(db *gorm.DB) ISocialRepository {
	return &socialRepository{db: db}
}
