package message

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kurnosovmak/test/internal/database"
	"github.com/kurnosovmak/test/internal/message"
)

type PgsqlMessageRepository struct {
	db *database.Database
}

func NewPgsqlMessageRepository(db *database.Database) *PgsqlMessageRepository {
	return &PgsqlMessageRepository{db: db}
}

func (r *PgsqlMessageRepository) Create(ctx context.Context, msg *message.Message) error {
	query := `
		INSERT INTO messages (chat_id, sender_id, content, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`
	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, query, msg.ChatID, msg.UserID, msg.Content, msg.Status)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

func (r *PgsqlMessageRepository) GetByChatID(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]message.Message, int, error) {
	query := `
		SELECT id, chat_id, sender_id, content, status, created_at, updated_at
		FROM messages
		WHERE chat_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`

	countQuery := `
		SELECT COUNT(*)
		FROM messages
		WHERE chat_id = $1 AND deleted_at IS NULL
	`

	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	var total int
	if err := tx.QueryRow(ctx, countQuery, chatID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to get total messages count: %w", err)
	}

	rows, err := tx.Query(ctx, query, chatID, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []message.Message{}, 0, nil
		}
		return nil, 0, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []message.Message
	for rows.Next() {
		var msg message.Message
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.UserID, &msg.Content, &msg.Status, &msg.CreatedAt, &msg.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, total, nil
}

func (r *PgsqlMessageRepository) UpdateStatus(ctx context.Context, messageID uuid.UUID, status string) error {
	query := `
		UPDATE messages
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`

	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	result, err := tx.Exec(ctx, query, status, messageID)
	if err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return message.ErrChatNotFound
	}

	return nil
}
