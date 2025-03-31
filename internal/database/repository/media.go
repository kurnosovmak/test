package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MediaFile struct {
	ID          uuid.UUID
	MessageID   uuid.UUID
	Type        string
	FileName    string
	FileSize    int64
	MimeType    string
	StoragePath string
	CreatedAt   time.Time
}

type MediaFileRepository interface {
	Create(ctx context.Context, file *MediaFile) error
	GetByID(ctx context.Context, id uuid.UUID) (*MediaFile, error)
	GetByMessageID(ctx context.Context, messageID uuid.UUID) ([]MediaFile, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
