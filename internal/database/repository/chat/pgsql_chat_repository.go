package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kurnosovmak/test/internal/database"
	"github.com/kurnosovmak/test/internal/database/repository"
)

type PgsqlChatRepository struct {
	db *database.Database
}

func NewPgsqlChatRepository(db *database.Database) *PgsqlChatRepository {
	return &PgsqlChatRepository{db: db}
}

func (r *PgsqlChatRepository) Create(ctx context.Context, chat *repository.Chat) error {
	query := `
		INSERT INTO chats (id, type, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
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

	_, err = tx.Exec(ctx, query, chat.ID, chat.Type, chat.Title, chat.CreatedAt, chat.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create chat: %w", err)
	}

	return nil
}

func (r *PgsqlChatRepository) GetByID(ctx context.Context, id uuid.UUID) (*repository.Chat, error) {
	query := `
		SELECT id, type, title, created_at, updated_at
		FROM chats
		WHERE id = $1
	`
	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	row := tx.QueryRow(ctx, query, id)
	chat := &repository.Chat{}
	if err := row.Scan(&chat.ID, &chat.Type, &chat.Title, &chat.CreatedAt, &chat.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("chat not found")
		}
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	return chat, nil
}

func (r *PgsqlChatRepository) Update(ctx context.Context, chat *repository.Chat) error {
	query := `
		UPDATE chats
		SET type = $1, title = $2, updated_at = $3
		WHERE id = $4
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

	result, err := tx.Exec(ctx, query, chat.Type, chat.Title, chat.UpdatedAt, chat.ID)
	if err != nil {
		return fmt.Errorf("failed to update chat: %w", err)
	}

	if affected := result.RowsAffected(); affected == 0 {
		return fmt.Errorf("chat not found")
	}

	return nil
}

func (r *PgsqlChatRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM chats WHERE id = $1`

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

	result, err := tx.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	if affected := result.RowsAffected(); affected == 0 {
		return fmt.Errorf("chat not found")
	}

	return nil
}

func (r *PgsqlChatRepository) AddParticipant(ctx context.Context, participant *repository.ChatParticipant) error {
	query := `
		INSERT INTO chat_participants (chat_id, user_id, is_admin, joined_at)
		VALUES ($1, $2, $3, $4)
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

	_, err = tx.Exec(ctx, query, participant.ChatID, participant.UserID, participant.IsAdmin, participant.JoinedAt)
	if err != nil {
		return fmt.Errorf("failed to add participant: %w", err)
	}

	return nil
}

func (r *PgsqlChatRepository) RemoveParticipant(ctx context.Context, chatID, userID uuid.UUID) error {
	query := `DELETE FROM chat_participants WHERE chat_id = $1 AND user_id = $2`

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

	result, err := tx.Exec(ctx, query, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}

	if affected := result.RowsAffected(); affected == 0 {
		return fmt.Errorf("participant not found")
	}

	return nil
}

func (r *PgsqlChatRepository) GetParticipants(ctx context.Context, chatID uuid.UUID) ([]repository.ChatParticipant, error) {
	query := `
		SELECT chat_id, user_id, is_admin, joined_at
		FROM chat_participants
		WHERE chat_id = $1
	`
	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	rows, err := tx.Query(ctx, query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}
	defer rows.Close()

	var participants []repository.ChatParticipant
	for rows.Next() {
		var participant repository.ChatParticipant
		if err := rows.Scan(&participant.ChatID, &participant.UserID, &participant.IsAdmin, &participant.JoinedAt); err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}
		participants = append(participants, participant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating participants: %w", err)
	}

	return participants, nil
}

func (r *PgsqlChatRepository) GetUserChats(ctx context.Context, userID uuid.UUID) ([]repository.Chat, error) {
	query := `
		SELECT c.id, c.type, c.title, c.created_at, c.updated_at
		FROM chats c
		JOIN chat_participants cp ON c.id = cp.chat_id
		WHERE cp.user_id = $1
	`
	tx, err := r.db.GetPool().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	rows, err := tx.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user chats: %w", err)
	}
	defer rows.Close()

	var chats []repository.Chat
	for rows.Next() {
		var chat repository.Chat
		if err := rows.Scan(&chat.ID, &chat.Type, &chat.Title, &chat.CreatedAt, &chat.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan chat: %w", err)
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating chats: %w", err)
	}

	return chats, nil
}
