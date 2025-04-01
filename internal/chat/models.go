package chat

import (
	"time"

	"github.com/google/uuid"
)

type CreateChatRequest struct {
	Type  string      `json:"type" validate:"required,oneof=private group"`
	Title *string     `json:"title,omitempty"`
	Users []uuid.UUID `json:"users" validate:"required,min=0"`
}

type GetChatRequest struct {
	ChatID uuid.UUID `json:"chat_id" validate:"required"`
}

type UpdateChatRequest struct {
	Title string `json:"title" validate:"required,min=1"`
}

type ChatResponse struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Title     *string   `json:"title,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ParticipantResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	IsAdmin  bool      `json:"is_admin"`
	JoinedAt time.Time `json:"joined_at"`
}

type ChatListResponse struct {
	Chats []ChatResponse `json:"chats"`
}

type AddParticipantRequest struct {
	ChatID  uuid.UUID `json:"chat_id" validate:"required"`
	UserID  uuid.UUID `json:"user_id" validate:"required"`
	IsAdmin bool      `json:"is_admin"`
}

type RemoveParticipantRequest struct {
	ChatID uuid.UUID `json:"chat_id" validate:"required"`
	UserID uuid.UUID `json:"user_id" validate:"required"`
}
