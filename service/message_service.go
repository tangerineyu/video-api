package service

import (
	"video-api/model"
	"video-api/repository"
)

type IMessageService interface {
	SaveMessage(fromID, toID uint, content string) error
	GetChatHistory(uidA, uidB uint) ([]model.Message, error)
}
type messageService struct {
	msgRepo repository.IMessageRepository
}

func (m *messageService) SaveMessage(fromID, toID uint, content string) error {
	return m.msgRepo.CreateMessage(&model.Message{
		FromUserID: fromID,
		ToUserID:   toID,
		Content:    content,
	})
}

func (m *messageService) GetChatHistory(uidA, uidB uint) ([]model.Message, error) {
	return m.msgRepo.GetMessages(uidA, uidB)
}

func NewMessageService(repo repository.IMessageRepository) IMessageService {
	return &messageService{msgRepo: repo}
}
