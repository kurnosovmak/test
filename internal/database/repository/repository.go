package repository

import (
	"context"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Chat struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IsGroup   bool      `json:"is_group"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsEdited  bool      `json:"is_edited"`
	IsDeleted bool      `json:"is_deleted"`
}

type MediaFile struct {
	ID        string    `json:"id"`
	MessageID string    `json:"message_id"`
	FileType  string    `json:"file_type"`
	FilePath  string    `json:"file_path"`
	FileSize  int64     `json:"file_size"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

type ChatRepository interface {
	Create(ctx context.Context, chat *Chat) error
	GetByID(ctx context.Context, id string) (*Chat, error)
	Update(ctx context.Context, chat *Chat) error
	Delete(ctx context.Context, id string) error
	AddParticipant(ctx context.Context, chatID, userID string, isAdmin bool) error
	RemoveParticipant(ctx context.Context, chatID, userID string) error
	GetParticipants(ctx context.Context, chatID string) ([]string, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message *Message) error
	GetByID(ctx context.Context, id string) (*Message, error)
	GetByChatID(ctx context.Context, chatID string, limit, offset int) ([]*Message, error)
	Update(ctx context.Context, message *Message) error
	Delete(ctx context.Context, id string) error
	SetMessageStatus(ctx context.Context, messageID, userID string, isRead bool) error
}

type MediaRepository interface {
	Create(ctx context.Context, media *MediaFile) error
	GetByID(ctx context.Context, id string) (*MediaFile, error)
	GetByMessageID(ctx context.Context, messageID string) ([]*MediaFile, error)
	Delete(ctx context.Context, id string) error
}
