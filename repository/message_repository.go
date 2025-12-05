package repository

import (
	"video-api/model"

	"gorm.io/gorm"
)

type IMessageRepository interface {
	CreateMessage(msg *model.Message) error
	GetMessages(userID, toUserID uint) ([]model.Message, error)
}
type messageRepository struct {
	db *gorm.DB
}

func (m *messageRepository) CreateMessage(msg *model.Message) error {
	return m.db.Create(msg).Error
}

func (m *messageRepository) GetMessages(userIdA, userIdB uint) ([]model.Message, error) {
	var messages []model.Message
	err := m.db.Where("(from_user_id = ? AND to_user_id = ?)"+
		" OR (from_user_id = ? AND to_user_id = ?)",
		userIdA, userIdB, userIdB, userIdA).
		Order("created_at asc").
		Find(&messages).Error
	return messages, err
}

func NewMessageRepository(db *gorm.DB) IMessageRepository {
	return &messageRepository{db: db}
}
