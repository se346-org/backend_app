package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type messageRepository struct {
	db *pgxpool.Pool
}

// CreateMessage implements domain.MessageRepository.
func (m *messageRepository) CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error) {
	query := `INSERT INTO message (id, conversation_id, user_id, type, body, created_at, updated_at, deleted_at, reply_to) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := m.db.Exec(ctx, query, message.ID, message.ConversationID, message.UserID, message.Type, message.Body, message.CreatedAt, message.UpdatedAt, message.DeletedAt, message.ReplyTo)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// GetListMessageByConversationID implements domain.MessageRepository.
func (m *messageRepository) GetListMessageByConversationID(ctx context.Context, conversationID string, lastID string, limit int) ([]*domain.Message, error) {
	fields := []string{
		"m.id",
		"m.conversation_id",
		"m.user_id",
		"m.type",
		"m.body",
		"m.created_at",
		"m.updated_at",
		"m.deleted_at",
		"m.reply_to",
		"u.id",
		"u.full_name",
		"u.avatar",
		"u.type",
	}
	condition := "conversation_id = $1"
	if lastID != "" {
		condition = fmt.Sprintf("%s AND m.id < $2", condition)
	}
	query := fmt.Sprintf(`SELECT %s FROM message AS m JOIN user_info AS u ON m.user_id = u.id WHERE %s ORDER BY m.id DESC LIMIT %d`, strings.Join(fields, ","), condition, limit)
	var params []any
	params = append(params, conversationID)
	if lastID != "" {
		params = append(params, lastID)
	}
	rows, err := m.db.Query(ctx, query, params...)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	if err == pgx.ErrNoRows {
		return []*domain.Message{}, nil
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var message domain.Message
		var user domain.UserInfo
		values := []any{
			&message.ID,
			&message.ConversationID,
			&message.UserID,
			&message.Type,
			&message.Body,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.DeletedAt,
			&message.ReplyTo,
			&user.ID,
			&user.FullName,
			&user.Avatar,
			&user.Type,
		}
		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		message.User = &user
		messages = append(messages, &message)
	}
	return messages, nil
}

// GetMessageByID implements domain.MessageRepository.
func (m *messageRepository) GetMessageByID(ctx context.Context, id string) (*domain.Message, error) {
	fields := []string{
		"m.id",
		"m.conversation_id",
		"m.user_id",
		"m.type",
		"m.body",
		"m.created_at",
		"m.updated_at",
		"m.deleted_at",
		"m.reply_to",
		"u.id",
		"u.full_name",
		"u.avatar",
		"u.type",
	}
	query := fmt.Sprintf(`SELECT %s FROM message AS m JOIN user_info AS u ON m.user_id = u.id WHERE m.id = $1`, strings.Join(fields, ","))
	row := m.db.QueryRow(ctx, query, id)
	var message domain.Message
	var user domain.UserInfo
	values := []any{
		&message.ID,
		&message.ConversationID,
		&message.UserID,
		&message.Type,
		&message.Body,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.DeletedAt,
		&message.ReplyTo,
		&user.ID,
		&user.FullName,
		&user.Avatar,
		&user.Type,
	}
	err := row.Scan(values...)
	if err != nil {
		return nil, err
	}
	message.User = &user
	return &message, nil
}

var _ domain.MessageRepository = &messageRepository{}

func NewMessageRepository(db *pgxpool.Pool) domain.MessageRepository {
	return &messageRepository{db: db}
}
