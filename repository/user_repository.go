package repository

import (
	"gorm.io/gorm"
	"video-api/model"
)

type IUserRepository interface {
	CreateUser(user *model.User) error
	FindUserByUsername(username string) (*model.User, error)
	FindUserByID(userID uint) (*model.User, error)
	UpdateAvatar(userID uint, avatarUrl string) error
}
type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	result := r.db.Where("username = ?", username).First(&user)
	return &user, result.Error
}

func (r *userRepository) FindUserByID(userID uint) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, userID)
	return &user, result.Error
}

func (r *userRepository) UpdateAvatar(userID uint, avatarUrl string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("avatar", avatarUrl).Error
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}
func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}
