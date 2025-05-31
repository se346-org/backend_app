package redis

import (
	"context"
	"fmt"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/redis/go-redis/v9"
)

type userCacheRepository struct {
	client *redis.Client
}

// GetUserIDByAccountID implements domain.UserCacheRepository.
func (u *userCacheRepository) GetUserIDByAccountID(ctx context.Context, accountID string) (string, error) {
	userID, err := u.client.Get(ctx, fmt.Sprintf("user_id:%s", accountID)).Result()
	if err != nil {
		return "", err
	}
	return userID, nil
}

// SetUserIDByAccountID implements domain.UserCacheRepository.
func (u *userCacheRepository) SetUserIDByAccountID(ctx context.Context, accountID string, userID string) error {
	_, err := u.client.Set(ctx, fmt.Sprintf("user_id:%s", accountID), userID, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

func NewUserCacheRepository(client *redis.Client) *userCacheRepository {
	return &userCacheRepository{
		client: client,
	}
}

var _ domain.UserCacheRepository = &userCacheRepository{}
