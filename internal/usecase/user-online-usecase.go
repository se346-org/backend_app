package usecase

import (
	"context"
	"time"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/pkg/pointer"
	"github.com/chat-socio/backend/pkg/uuid"
)

type UserOnlineUsecase interface {
	CreateUserOnline(ctx context.Context, userOnline *domain.UserOnline) error
	DeleteUserOnline(ctx context.Context, id string) error
	GetUserOnlineByConversationID(ctx context.Context, conversationID string) ([]*domain.UserOnline, error)
}

type userOnlineUsecase struct {
	userOnlineRepo domain.UserOnlineRepository
}

var _ UserOnlineUsecase = &userOnlineUsecase{}

func NewUserOnlineUsecase(userOnlineRepo domain.UserOnlineRepository) *userOnlineUsecase {
	return &userOnlineUsecase{userOnlineRepo: userOnlineRepo}
}

// CreateUserOnline implements UserOnlineUsecase.
func (u *userOnlineUsecase) CreateUserOnline(ctx context.Context, userOnline *domain.UserOnline) error {
	id, err := uuid.NewID()
	if err != nil {
		return err
	}
	userOnline.ID = id
	userOnline.CreatedAt = pointer.ToPtr(time.Now())
	return u.userOnlineRepo.CreateUserOnline(ctx, userOnline)
}

// DeleteUserOnline implements UserOnlineUsecase.
func (u *userOnlineUsecase) DeleteUserOnline(ctx context.Context, id string) error {
	return u.userOnlineRepo.DeleteUserOnline(ctx, id)
}

func (u *userOnlineUsecase) GetUserOnlineByConversationID(ctx context.Context, conversationID string) ([]*domain.UserOnline, error) {
	return u.userOnlineRepo.GetUserOnlineByConversationID(ctx, conversationID)
}
