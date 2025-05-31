package redis

import (
	"context"
	"fmt"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/redis/go-redis/v9"
)

type sessionCacheRepository struct {
	client *redis.Client
}

// NewSessionCacheRepository creates a new session cache repository.
func NewSessionCacheRepository(client *redis.Client) domain.SessionCacheRepository {
	return &sessionCacheRepository{
		client: client,
	}
}

// CreateSessionWithExpireTime implements domain.SessionCacheRepository.
func (s *sessionCacheRepository) CreateSessionWithExpireTime(ctx context.Context, session *domain.Session) error {
	mapSession := session.ConvertToMapString()
	_, err := s.client.HSet(ctx, fmt.Sprintf("session:%s", session.SessionToken), mapSession).Result()
	if err != nil {
		return err
	}
	return nil
}

// DeleteSession implements domain.SessionCacheRepository.
func (s *sessionCacheRepository) DeleteSession(ctx context.Context, token string) error {
	_, err := s.client.Del(ctx, fmt.Sprintf("session:%s", token)).Result()
	if err != nil {
		return err
	}
	return nil
}

// GetSessionByToken implements domain.SessionCacheRepository.
func (s *sessionCacheRepository) GetSessionByToken(ctx context.Context, token string) (*domain.Session, error) {
	sessionMap, err := s.client.HGetAll(ctx, fmt.Sprintf("session:%s", token)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if len(sessionMap) == 0 {
		return nil, nil
	}
	session := &domain.Session{}
	session.FromMap(sessionMap)
	return session, nil
}

var _ domain.SessionCacheRepository = (*sessionCacheRepository)(nil)
