package postgresql

import (
	"context"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/jackc/pgx/v5/pgxpool"
)

type seenMessageRepository struct {
	db  *pgxpool.Pool
	obs *observability.Observability
}

var _ domain.SeenMessageRepository = (*seenMessageRepository)(nil)

func NewSeenMessageRepository(db *pgxpool.Pool, obs *observability.Observability) *seenMessageRepository {
	return &seenMessageRepository{db: db, obs: obs}
}

func (r *seenMessageRepository) CreateSeenMessage(ctx context.Context, seenMessage *domain.SeenMessage) error {
	ctx, span := r.obs.StartSpan(ctx, "seenMessageRepository.CreateSeenMessage")
	defer span()
	logger := r.obs.Logger.WithContext(ctx)
	query := `
		INSERT INTO seen_message (id, message_id, user_id, conversation_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, conversation_id) DO UPDATE SET message_id = EXCLUDED.message_id, updated_at = current_timestamp
	`
	_, err := r.db.Exec(ctx, query, seenMessage.ID, seenMessage.MessageID, seenMessage.UserID, seenMessage.ConversationID)
	if err != nil {
		logger.Error("failed to upsert seen message", err, seenMessage)
		return err
	}
	return nil
}

func (r *seenMessageRepository) GetListSeenMessageByConversationID(ctx context.Context, conversationID string) ([]*domain.SeenMessage, error) {
	ctx, span := r.obs.StartSpan(ctx, "seenMessageRepository.GetListSeenMessageByConversationID")
	defer span()
	logger := r.obs.Logger.WithContext(ctx)
	query := `
		SELECT id, message_id, user_id, conversation_id, created_at, updated_at
		FROM seen_message
		WHERE conversation_id = $1
	`
	rows, err := r.db.Query(ctx, query, conversationID)
	if err != nil {
		logger.Error("failed to get list seen message by conversation id", err, conversationID)
		return nil, err
	}
	defer rows.Close()

	var seenMessages []*domain.SeenMessage
	for rows.Next() {
		var seenMessage domain.SeenMessage
		err := rows.Scan(&seenMessage.ID, &seenMessage.MessageID, &seenMessage.UserID, &seenMessage.ConversationID, &seenMessage.CreatedAt, &seenMessage.UpdatedAt)
		if err != nil {
			logger.Error("failed to scan seen message", err, seenMessage)
			continue
		}
		seenMessages = append(seenMessages, &seenMessage)
	}
	return seenMessages, nil
}
