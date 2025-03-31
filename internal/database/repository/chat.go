package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID
	Type      string
	Title     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChatParticipant struct {
	ChatID   uuid.UUID
	UserID   uuid.UUID
	IsAdmin  bool
	JoinedAt time.Time
}

type ChatRepository interface {
	Create(ctx context.Context, chat *Chat) error
	GetByID(ctx context.Context, id uuid.UUID) (*Chat, error)
	Update(ctx context.Context, chat *Chat) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddParticipant(ctx context.Context, participant *ChatParticipant) error
	RemoveParticipant(ctx context.Context, chatID, userID uuid.UUID) error
	GetParticipants(ctx context.Context, chatID uuid.UUID) ([]ChatParticipant, error)
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]Chat, error)
}
