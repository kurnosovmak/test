package chat

import (
	"context"
	"errors"

	"github.com/kurnosovmak/test/internal/database/repository"

	"github.com/google/uuid"
)

var (
	ErrChatNotFound = errors.New("chat not found")
	ErrUnauthorized = errors.New("unauthorized access to chat")
)

type Service struct {
	chatRepo repository.ChatRepository
}

func NewService(chatRepo repository.ChatRepository) *Service {
	return &Service{chatRepo: chatRepo}
}

func (s *Service) CreateChat(ctx context.Context, req *CreateChatRequest, creatorID uuid.UUID) (*ChatResponse, error) {
	chat := &repository.Chat{
		ID:    uuid.New(),
		Type:  req.Type,
		Title: req.Title,
	}

	if err := s.chatRepo.Create(ctx, chat); err != nil {
		return nil, err
	}

	// Add creator as admin
	participant := &repository.ChatParticipant{
		ChatID:  chat.ID,
		UserID:  creatorID,
		IsAdmin: true,
	}

	if err := s.chatRepo.AddParticipant(ctx, participant); err != nil {
		return nil, err
	}

	// Add other participants
	for _, userID := range req.Users {
		if userID != creatorID {
			participant := &repository.ChatParticipant{
				ChatID: chat.ID,
				UserID: userID,
			}
			if err := s.chatRepo.AddParticipant(ctx, participant); err != nil {
				return nil, err
			}
		}
	}

	return &ChatResponse{
		ID:        chat.ID,
		Type:      chat.Type,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}, nil
}

func (s *Service) GetChat(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (*ChatResponse, error) {
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, ErrChatNotFound
	}

	// Verify user is a participant
	participants, err := s.chatRepo.GetParticipants(ctx, chatID)
	if err != nil {
		return nil, err
	}

	isParticipant := false
	for _, p := range participants {
		if p.UserID == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		return nil, ErrUnauthorized
	}

	return &ChatResponse{
		ID:        chat.ID,
		Type:      chat.Type,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt,
		UpdatedAt: chat.UpdatedAt,
	}, nil
}

func (s *Service) GetUserChats(ctx context.Context, userID uuid.UUID) (*ChatListResponse, error) {
	chats, err := s.chatRepo.GetUserChats(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := make([]ChatResponse, len(chats))
	for i, chat := range chats {
		response[i] = ChatResponse{
			ID:        chat.ID,
			Type:      chat.Type,
			Title:     chat.Title,
			CreatedAt: chat.CreatedAt,
			UpdatedAt: chat.UpdatedAt,
		}
	}

	return &ChatListResponse{Chats: response}, nil
}

func (s *Service) AddParticipant(ctx context.Context, req *AddParticipantRequest, requesterID uuid.UUID) error {
	// Verify requester is admin
	participants, err := s.chatRepo.GetParticipants(ctx, req.ChatID)
	if err != nil {
		return err
	}

	isAdmin := false
	for _, p := range participants {
		if p.UserID == requesterID && p.IsAdmin {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return ErrUnauthorized
	}

	participant := &repository.ChatParticipant{
		ChatID:  req.ChatID,
		UserID:  req.UserID,
		IsAdmin: req.IsAdmin,
	}

	return s.chatRepo.AddParticipant(ctx, participant)
}

func (s *Service) RemoveParticipant(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, requesterID uuid.UUID) error {
	// Verify requester is admin
	participants, err := s.chatRepo.GetParticipants(ctx, chatID)
	if err != nil {
		return err
	}

	isAdmin := false
	for _, p := range participants {
		if p.UserID == requesterID && p.IsAdmin {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return ErrUnauthorized
	}

	return s.chatRepo.RemoveParticipant(ctx, chatID, userID)
}
