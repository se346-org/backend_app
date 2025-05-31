package postgresql

import (
	"context"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userOnlineRepository struct {
	db *pgxpool.Pool
}

// GetUserOnlineByConversationID implements domain.UserOnlineRepository.
func (u *userOnlineRepository) GetUserOnlineByConversationID(ctx context.Context, conversationID string) ([]*domain.UserOnline, error) {
	query := `SELECT id, user_id, connection_id, created_at FROM user_online WHERE user_id IN (SELECT user_id FROM conversation_member WHERE conversation_id = $1)`
	rows, err := u.db.Query(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userOnlines []*domain.UserOnline
	for rows.Next() {
		var userOnline domain.UserOnline
		err := rows.Scan(&userOnline.ID, &userOnline.UserID, &userOnline.ConnectionID, &userOnline.CreatedAt)
		if err != nil {
			return nil, err
		}
		userOnlines = append(userOnlines, &userOnline)
	}
	return userOnlines, nil
}

// CreateUserOnline implements domain.UserOnlineRepository.
func (u *userOnlineRepository) CreateUserOnline(ctx context.Context, userOnline *domain.UserOnline) error {
	query := `INSERT INTO user_online (id, user_id, connection_id, created_at) VALUES ($1, $2, $3, $4)`
	_, err := u.db.Exec(ctx, query, userOnline.ID, userOnline.UserID, userOnline.ConnectionID, userOnline.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserOnline implements domain.UserOnlineRepository.
func (u *userOnlineRepository) DeleteUserOnline(ctx context.Context, id string) error {
	query := `DELETE FROM user_online WHERE id = $1`
	_, err := u.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

var _ domain.UserOnlineRepository = &userOnlineRepository{}

func NewUserOnlineRepository(db *pgxpool.Pool) *userOnlineRepository {
	return &userOnlineRepository{db: db}
}
