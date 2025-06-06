package postgresql

import (
	"context"
	"time"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/chat-socio/backend/pkg/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contactRepository struct {
	db  *pgxpool.Pool
	obs *observability.Observability
}

// AcceptRequestFriend implements domain.ContactRepository.
func (c *contactRepository) AcceptRequestFriend(ctx context.Context, requestFriendID string) error {
	ctx, span := c.obs.StartSpan(ctx, "ContactRepository.AcceptRequestFriend")
	defer span()
	logger := c.obs.Logger.WithContext(ctx)
	tx, err := c.db.Begin(ctx)
	if err != nil {
		logger.Error("failed to begin transaction", err)
		return err
	}
	defer tx.Rollback(ctx)
	query := `select id, from_user_id, to_user_id, status, created_at, updated_at from friend_request where id = $1`
	row := tx.QueryRow(ctx, query, requestFriendID)
	var requestFriend domain.RequestFriend
	err = row.Scan(&requestFriend.ID, &requestFriend.FromUserID, &requestFriend.ToUserID, &requestFriend.Status, &requestFriend.CreatedAt, &requestFriend.UpdatedAt)
	if err != nil {
		logger.Error("failed to scan friend request", err)
		return err
	}
	query = `update friend_request set status = 'accepted', updated_at = now() where id = $1`
	_, err = tx.Exec(ctx, query, requestFriendID)
	if err != nil {
		logger.Error("failed to update friend request", err)
		return err
	}
	now := time.Now()
	fromContactID, err := uuid.NewID()
	if err != nil {
		logger.Error("failed to generate uuid", err)
		return err
	}
	toContactID, err := uuid.NewID()
	if err != nil {
		logger.Error("failed to generate uuid", err)
		return err
	}
	query = `insert into contact (id, user_id, friend_id, created_at, updated_at) values ($1, $2, $3, $4, $5), ($6, $7, $8, $9, $10)`
	_, err = tx.Exec(ctx, query, fromContactID, requestFriend.FromUserID, requestFriend.ToUserID, now, now, toContactID, requestFriend.ToUserID, requestFriend.FromUserID, now, now)
	if err != nil {
		logger.Error("failed to insert contact", err)
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		logger.Error("failed to commit transaction", err)
		return err
	}
	return nil
}

// CreateContacts implements domain.ContactRepository.
func (c *contactRepository) CreateContacts(ctx context.Context, contacts []*domain.Contact) error {
	panic("unimplemented")
}

// CreateRequestFriend implements domain.ContactRepository.
func (c *contactRepository) CreateRequestFriend(ctx context.Context, requestFriend *domain.RequestFriend) error {
	panic("unimplemented")
}

// GetListContactByUserID implements domain.ContactRepository.
func (c *contactRepository) GetListContactByUserID(ctx context.Context, userID string, limit int, lastID string) ([]*domain.Contact, error) {
	panic("unimplemented")
}

// GetListRequestFriendReceivedByUserID implements domain.ContactRepository.
func (c *contactRepository) GetListRequestFriendReceivedByUserID(ctx context.Context, userID string, limit int, lastID string) ([]*domain.RequestFriend, error) {
	panic("unimplemented")
}

// GetListRequestFriendSentByUserID implements domain.ContactRepository.
func (c *contactRepository) GetListRequestFriendSentByUserID(ctx context.Context, userID string, limit int, lastID string) ([]*domain.RequestFriend, error) {
	panic("unimplemented")
}

// UpdateRequestFriendStatus implements domain.ContactRepository.
func (c *contactRepository) UpdateRequestFriendStatus(ctx context.Context, id string, status string) error {
	panic("unimplemented")
}

var _ domain.ContactRepository = &contactRepository{}

func NewContactRepository(db *pgxpool.Pool, obs *observability.Observability) *contactRepository {
	return &contactRepository{
		db:  db,
		obs: obs,
	}
}
