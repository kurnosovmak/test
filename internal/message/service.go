package message

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrChatNotFound = errors.New("chat not found")
	ErrUserNotFound = errors.New("user not found")
)

type Repository interface {
	Create(ctx context.Context, msg *Message) error
	GetByChatID(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]Message, int, error)
	UpdateStatus(ctx context.Context, messageID uuid.UUID, status string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateMessage(ctx context.Context, userID uuid.UUID, req *CreateMessageRequest) (*MessageResponse, error) {
	msg := &Message{
		ChatID:    req.ChatID,
		UserID:    userID,
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    "delivered",
	}

	if err := s.repo.Create(ctx, msg); err != nil {
		return nil, err
	}

	return &MessageResponse{Message: *msg}, nil
}

func (s *Service) GetMessages(ctx context.Context, req *GetMessagesRequest) (*MessagesResponse, error) {
	messages, total, err := s.repo.GetByChatID(ctx, req.ChatID, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	return &MessagesResponse{
		Messages: messages,
		Total:    total,
	}, nil
}

func (s *Service) UpdateMessageStatus(ctx context.Context, messageID uuid.UUID, status string) error {
	return s.repo.UpdateStatus(ctx, messageID, status)
}
