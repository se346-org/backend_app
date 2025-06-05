package postgresql

import (
	"context"
	"time"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type friendRepository struct {
	db *pgxpool.Pool
}

func NewFriendRepository(db *pgxpool.Pool) domain.FriendRepository {
	return &friendRepository{db: db}
}

func (r *friendRepository) CreateFriend(ctx context.Context, friend *domain.Friend) error {
	query := `
		INSERT INTO friends (id, user_id, friend_id, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	friend.ID = uuid.New().String()
	friend.CreatedAt = time.Now()
	friend.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		friend.ID,
		friend.UserID,
		friend.FriendID,
		friend.Status,
		friend.CreatedAt,
		friend.UpdatedAt,
	)
	return err
}

func (r *friendRepository) GetFriendByID(ctx context.Context, id string) (*domain.Friend, error) {
	query := `
		SELECT id, user_id, friend_id, status, created_at, updated_at
		FROM friends
		WHERE id = $1
	`
	friend := &domain.Friend{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&friend.ID,
		&friend.UserID,
		&friend.FriendID,
		&friend.Status,
		&friend.CreatedAt,
		&friend.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return friend, nil
}

func (r *friendRepository) GetFriendByUserIDs(ctx context.Context, userID string, friendID string) (*domain.Friend, error) {
	query := `
		SELECT id, user_id, friend_id, status, created_at, updated_at
		FROM friends
		WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
	`
	friend := &domain.Friend{}
	err := r.db.QueryRow(ctx, query, userID, friendID).Scan(
		&friend.ID,
		&friend.UserID,
		&friend.FriendID,
		&friend.Status,
		&friend.CreatedAt,
		&friend.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return friend, nil
}

func (r *friendRepository) GetListFriendsByUserID(ctx context.Context, userID string, status domain.FriendStatus, limit int, lastID string) ([]*domain.FriendWithUser, error) {
	query := `
		SELECT f.id, f.user_id, f.friend_id, f.status, f.created_at, f.updated_at,
			   u.id, u.email, u.full_name, u.avatar, u.created_at, u.updated_at
		FROM friends f
		JOIN user_info u ON (f.friend_id = u.id AND f.user_id = $1) OR (f.user_id = u.id AND f.friend_id = $1)
		WHERE f.status = $2
		AND ($3 = '' OR f.id < $3)
		ORDER BY f.id DESC
		LIMIT $4
	`
	rows, err := r.db.Query(ctx, query, userID, status, lastID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*domain.FriendWithUser
	for rows.Next() {
		friend := &domain.FriendWithUser{
			User: &domain.UserInfo{},
		}
		err := rows.Scan(
			&friend.ID,
			&friend.UserID,
			&friend.FriendID,
			&friend.Status,
			&friend.CreatedAt,
			&friend.UpdatedAt,
			&friend.User.ID,
			&friend.User.Email,
			&friend.User.FullName,
			&friend.User.Avatar,
			&friend.User.CreatedAt,
			&friend.User.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

func (r *friendRepository) UpdateFriendStatus(ctx context.Context, id string, status domain.FriendStatus) error {
	query := `
		UPDATE friends
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(ctx, query, status, time.Now(), id)
	return err
}

func (r *friendRepository) DeleteFriend(ctx context.Context, id string) error {
	query := `
		DELETE FROM friends
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *friendRepository) CheckFriendshipExists(ctx context.Context, userID string, friendID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM friends
			WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
		)
	`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, friendID).Scan(&exists)
	return exists, err
}

func (r *friendRepository) GetListFriendRequestsReceived(ctx context.Context, userID string, limit int, lastID string) ([]*domain.FriendWithUser, error) {
	query := `
		SELECT f.id, f.user_id, f.friend_id, f.status, f.created_at, f.updated_at,
			   u.id, u.email, u.full_name, u.avatar, u.created_at, u.updated_at
		FROM friends f
		JOIN user_info u ON f.user_id = u.id
		WHERE f.friend_id = $1 AND f.status = $2
		AND ($3 = '' OR f.id < $3)
		ORDER BY f.id DESC
		LIMIT $4
	`
	rows, err := r.db.Query(ctx, query, userID, domain.FriendStatusPending, lastID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*domain.FriendWithUser
	for rows.Next() {
		friend := &domain.FriendWithUser{
			User: &domain.UserInfo{},
		}
		err := rows.Scan(
			&friend.ID,
			&friend.UserID,
			&friend.FriendID,
			&friend.Status,
			&friend.CreatedAt,
			&friend.UpdatedAt,
			&friend.User.ID,
			&friend.User.Email,
			&friend.User.FullName,
			&friend.User.Avatar,
			&friend.User.CreatedAt,
			&friend.User.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		friends = append(friends, friend)
	}
	return friends, nil
}

const getFriendsQuery = `
	SELECT f.id, f.user_id, f.friend_id, f.status, f.created_at, f.updated_at,
		   u.user_id, u.full_name, u.avatar, u.user_type, u.created_at, u.updated_at
	FROM friends f
	JOIN users u ON f.friend_id = u.user_id
	WHERE f.user_id = $1 AND f.status = $2
	ORDER BY f.created_at DESC
	LIMIT $3
` 