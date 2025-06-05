package usecase

import (
	"context"
	"errors"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/pkg/observability"
)

type FriendUseCase interface {
	SendFriendRequest(ctx context.Context, userID string, friendID string) error
	AcceptFriendRequest(ctx context.Context, userID string, friendID string) error
	RejectFriendRequest(ctx context.Context, userID string, friendID string) error
	GetFriends(ctx context.Context, userID string, status domain.FriendStatus, limit int, lastID string) ([]*domain.FriendWithUser, error)
	Unfriend(ctx context.Context, userID string, friendID string) error
	GetFriendRequests(ctx context.Context, userID string, limit int, lastID string) ([]*domain.FriendWithUser, error)
	GetFriendRequestsReceived(ctx context.Context, userID string, limit int, lastID string) ([]*domain.FriendWithUser, error)
}

type friendUseCase struct {
	friendRepository domain.FriendRepository
	userRepository   domain.UserRepository
	obs             *observability.Observability
}

func NewFriendUseCase(friendRepository domain.FriendRepository, userRepository domain.UserRepository, obs *observability.Observability) FriendUseCase {
	return &friendUseCase{
		friendRepository: friendRepository,
		userRepository:   userRepository,
		obs:             obs,
	}
}

func (uc *friendUseCase) SendFriendRequest(ctx context.Context, userID string, friendID string) error {
	logger := uc.obs.Logger.WithContext(ctx)
	logger.Info("Sending friend request in use case", map[string]interface{}{
		"user_id":   userID,
		"friend_id": friendID,
	})

	// Check if user exists
	user, err := uc.userRepository.GetUserByID(ctx, friendID)
	if err != nil {
		logger.WithError(err).Error("Failed to get user by ID", map[string]interface{}{
			"friend_id": friendID,
		})
		return err
	}
	if user == nil {
		logger.Error("User not found", map[string]interface{}{
			"friend_id": friendID,
		})
		return errors.New("user not found")
	}

	// Check if friendship already exists
	exists, err := uc.friendRepository.CheckFriendshipExists(ctx, userID, friendID)
	if err != nil {
		logger.WithError(err).Error("Failed to check friendship exists", map[string]interface{}{
			"user_id":   userID,
			"friend_id": friendID,
		})
		return err
	}
	if exists {
		logger.Error("Friendship already exists", map[string]interface{}{
			"user_id":   userID,
			"friend_id": friendID,
		})
		return errors.New("friendship already exists")
	}

	// Create friend request
	friend := &domain.Friend{
		UserID:   userID,
		FriendID: friendID,
		Status:   domain.FriendStatusPending,
	}
	err = uc.friendRepository.CreateFriend(ctx, friend)
	if err != nil {
		logger.WithError(err).Error("Failed to create friend request", map[string]interface{}{
			"user_id":   userID,
			"friend_id": friendID,
		})
		return err
	}

	logger.Info("Friend request created successfully", map[string]interface{}{
		"user_id":   userID,
		"friend_id": friendID,
	})
	return nil
}

func (uc *friendUseCase) AcceptFriendRequest(ctx context.Context, userID string, friendID string) error {
	friend, err := uc.friendRepository.GetFriendByUserIDs(ctx, userID, friendID)
	if err != nil {
		return err
	}
	if friend == nil {
		return errors.New("friend request not found")
	}
	if friend.Status != domain.FriendStatusPending {
		return errors.New("invalid friend request status")
	}
	if friend.FriendID != userID {
		return errors.New("unauthorized to accept this friend request")
	}

	return uc.friendRepository.UpdateFriendStatus(ctx, friend.ID, domain.FriendStatusAccepted)
}

func (uc *friendUseCase) RejectFriendRequest(ctx context.Context, userID string, friendID string) error {
	friend, err := uc.friendRepository.GetFriendByUserIDs(ctx, userID, friendID)
	if err != nil {
		return err
	}
	if friend == nil {
		return errors.New("friend request not found")
	}
	if friend.Status != domain.FriendStatusPending {
		return errors.New("invalid friend request status")
	}
	if friend.FriendID != userID {
		return errors.New("unauthorized to reject this friend request")
	}

	return uc.friendRepository.UpdateFriendStatus(ctx, friend.ID, domain.FriendStatusRejected)
}

func (uc *friendUseCase) GetFriends(ctx context.Context, userID string, status domain.FriendStatus, limit int, lastID string) ([]*domain.FriendWithUser, error) {
	return uc.friendRepository.GetListFriendsByUserID(ctx, userID, status, limit, lastID)
}

func (uc *friendUseCase) Unfriend(ctx context.Context, userID string, friendID string) error {
	friend, err := uc.friendRepository.GetFriendByUserIDs(ctx, userID, friendID)
	if err != nil {
		return err
	}
	if friend == nil {
		return errors.New("friendship not found")
	}
	if friend.Status != domain.FriendStatusAccepted {
		return errors.New("not friends")
	}

	return uc.friendRepository.DeleteFriend(ctx, friend.ID)
}

func (uc *friendUseCase) GetFriendRequests(ctx context.Context, userID string, limit int, lastID string) ([]*domain.FriendWithUser, error) {
	logger := uc.obs.Logger.WithContext(ctx)
	logger.Info("Getting friend requests", map[string]interface{}{
		"user_id": userID,
		"limit":   limit,
		"last_id": lastID,
	})

	return uc.friendRepository.GetListFriendsByUserID(ctx, userID, domain.FriendStatusPending, limit, lastID)
}

func (uc *friendUseCase) GetFriendRequestsReceived(ctx context.Context, userID string, limit int, lastID string) ([]*domain.FriendWithUser, error) {
	logger := uc.obs.Logger.WithContext(ctx)
	logger.Info("Getting friend requests received", map[string]interface{}{
		"user_id": userID,
		"limit":   limit,
		"last_id": lastID,
	})

	return uc.friendRepository.GetListFriendRequestsReceived(ctx, userID, limit, lastID)
} 