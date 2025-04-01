package message

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ChatID    uuid.UUID `json:"chat_id" db:"chat_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Status    string    `json:"status" db:"status"` // delivered, read
}

type CreateMessageRequest struct {
	ChatID  uuid.UUID `json:"chat_id" validate:"required"`
	Content string    `json:"content" validate:"required"`
}

type GetMessagesRequest struct {
	ChatID uuid.UUID `json:"chat_id" validate:"required"`
	Limit  int       `json:"limit" validate:"required,min=1,max=100"`
	Offset int       `json:"offset" validate:"min=0"`
}

type MessageResponse struct {
	Message Message `json:"message"`
}

type MessagesResponse struct {
	Messages []Message `json:"messages"`
	Total    int       `json:"total"`
}

type UpdateMessageStatusRequest struct {
	MessageID uuid.UUID `json:"message_id" validate:"required"`
	Status    string    `json:"status" validate:"required,oneof=delivered read"`
}
