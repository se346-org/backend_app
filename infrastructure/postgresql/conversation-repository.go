package postgresql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type conversationRepository struct {
	db *pgxpool.Pool
}

// CheckDMConversationExist implements domain.ConversationRepository.
func (c *conversationRepository) CheckDMConversationExist(ctx context.Context, userID1 string, userID2 string) (*domain.Conversation, error) {
	query := `
		SELECT c.id, c.created_at, c.type, c.title, c.avatar, c.updated_at FROM conversation c
		INNER JOIN conversation_member cm1 ON c.id = cm1.conversation_id
		INNER JOIN conversation_member cm2 ON c.id = cm2.conversation_id
		WHERE cm1.user_id = $1 AND cm2.user_id = $2 AND c.type = 'DM'
	`
	var conversation domain.Conversation
	err := c.db.QueryRow(ctx, query, userID1, userID2).Scan(&conversation.ID, &conversation.CreatedAt, &conversation.Type, &conversation.Title, &conversation.Avatar, &conversation.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

// CreateConversation implements domain.ConversationRepository.
func (c *conversationRepository) CreateConversation(ctx context.Context, conversation *domain.Conversation, conversationMembers []*domain.ConversationMember) (*domain.Conversation, error) {
	tx, err := c.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO conversation (id, created_at, type, title, avatar, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(ctx, query, conversation.ID, conversation.CreatedAt, conversation.Type, conversation.Title, conversation.Avatar, conversation.UpdatedAt)
	if err != nil {
		return nil, err
	}

	query = `
		INSERT INTO conversation_member (id, conversation_id, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, conversationMember := range conversationMembers {
		_, err = tx.Exec(ctx, query, conversationMember.ID, conversationMember.ConversationID, conversationMember.UserID, conversationMember.CreatedAt, conversationMember.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return conversation, nil
}

// GetConversationByID implements domain.ConversationRepository.
func (c *conversationRepository) GetConversationByID(ctx context.Context, id string) (*domain.Conversation, []*domain.ConversationMemberWithUser, error) {
	var conversation domain.Conversation
	var conversationMembers []*domain.ConversationMemberWithUser

	fields, values := conversation.MapFields()
	// query get conversation and conversation members
	query := fmt.Sprintf(`SELECT %s FROM conversation WHERE id = $1`, strings.Join(fields, ", "))
	row := c.db.QueryRow(ctx, query, id)

	if err := row.Scan(values...); err != nil {
		return nil, nil, err
	}

	query = `SELECT cm.conversation_id, cm.user_id, ui.full_name, ui.avatar, ui.type FROM conversation_member cm
		INNER JOIN user_info ui ON cm.user_id = ui.id
		WHERE cm.conversation_id = $1`
	rows, err := c.db.Query(ctx, query, id)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var conversationMember domain.ConversationMemberWithUser
		if err := rows.Scan(&conversationMember.ConversationID, &conversationMember.UserID, &conversationMember.FullName, &conversationMember.Avatar, &conversationMember.UserType); err != nil {
			return nil, nil, err
		}
		conversationMembers = append(conversationMembers, &conversationMember)
	}

	return &conversation, conversationMembers, nil
}

// GetListConversationByUserID implements domain.ConversationRepository.
func (c *conversationRepository) GetListConversationByUserID(ctx context.Context, userID string, lastMessageID string, limit int) ([]*domain.Conversation, error) {
	var conversations []*domain.Conversation
	var conditionLastMessageID string
	var params []any
	params = append(params, userID)
	if lastMessageID != "" {
		conditionLastMessageID = `AND last_message_id < $2`
		params = append(params, lastMessageID)
	}
	// Add NULL handling for last_message_id
	// fieldsWithCoalesce := []string{
	// 	"c.id",
	// 	"c.created_at",
	// 	"c.type",
	// 	"c.title",
	// 	"c.avatar",
	// 	"c.updated_at",
	// 	"c.deleted_at",
	// 	"COALESCE(c.last_message_id::text, '')", // Handle NULL last_message_id
	// 	"COALESCE(m.id::text, '')",              // Handle NULL message id
	// 	"COALESCE(m.conversation_id::text, '')",
	// 	"COALESCE(m.user_id::text, '')",
	// 	"COALESCE(m.type::text, '')",
	// 	"COALESCE(m.body::text, '')",
	// 	"COALESCE(m.created_at, NULL)",
	// 	"COALESCE(m.updated_at, NULL)",
	// 	"COALESCE(m.reply_to::text, '')",
	// 	"COALESCE(ui.full_name::text, '')",
	// 	"COALESCE(ui.avatar::text, '')",
	// 	"COALESCE(ui.type::text, '')",
	// }

	// query := fmt.Sprintf(`SELECT %s FROM conversation c
	// 	LEFT JOIN message m ON c.last_message_id = m.id
	// 	LEFT JOIN user_info ui ON m.user_id = ui.id
	// 	WHERE c.id IN (SELECT DISTINCT conversation_id FROM conversation_member WHERE user_id = $1) %s
	// 	ORDER BY last_message_id DESC LIMIT %d`, strings.Join(fieldsWithCoalesce, ", "), conditionLastMessageID, limit)
	// Add seen_message check to determine if the last message has been read
	query := fmt.Sprintf(`
		WITH conversation_data AS (
			SELECT c.id, c.created_at, c.type, c.title, c.avatar, c.updated_at, c.deleted_at, 
				COALESCE(c.last_message_id::text, '') as last_message_id,
				COALESCE(m.id::text, '') as message_id,
				COALESCE(m.conversation_id::text, '') as message_conversation_id,
				COALESCE(m.user_id::text, '') as message_user_id,
				COALESCE(m.type::text, '') as message_type,
				COALESCE(m.body::text, '') as message_body,
				COALESCE(m.created_at, NULL) as message_created_at,
				COALESCE(m.updated_at, NULL) as message_updated_at,
				COALESCE(m.reply_to::text, '') as message_reply_to,
				COALESCE(ui.id::text, '') as user_id,
				COALESCE(ui.full_name::text, '') as user_full_name,
				COALESCE(ui.avatar::text, '') as user_avatar,
				COALESCE(ui.type::text, '') as user_type,
				CASE 
					WHEN EXISTS (
						SELECT 1 FROM seen_message sm 
						WHERE sm.conversation_id = c.id 
						AND sm.user_id = $1
						AND sm.message_id = c.last_message_id
					) THEN true
					ELSE false
				END as is_read
			FROM conversation c
			LEFT JOIN message m ON c.last_message_id = m.id
			LEFT JOIN user_info ui ON m.user_id = ui.id
			WHERE c.id IN (
				SELECT DISTINCT conversation_id 
				FROM conversation_member 
				WHERE user_id = $1
			) %s
			ORDER BY c.last_message_id DESC
			LIMIT %d
		),
		conversation_members AS (
			SELECT 
				cm.conversation_id,
				json_agg(
					json_build_object(
						'id', ui.id,
						'full_name', ui.full_name,
						'avatar', ui.avatar,
						'type', ui.type
					)
				) as members
			FROM conversation_member cm
			INNER JOIN user_info ui ON cm.user_id = ui.id
			WHERE cm.conversation_id IN (SELECT id FROM conversation_data)
			GROUP BY cm.conversation_id
		)
		SELECT 
			cd.*,
			COALESCE(cm.members::text, '[]') as members
		FROM conversation_data cd
		LEFT JOIN conversation_members cm ON cd.id = cm.conversation_id
		ORDER BY cd.last_message_id DESC`, conditionLastMessageID, limit)

	rows, err := c.db.Query(ctx, query, params...)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	if err == pgx.ErrNoRows {
		return []*domain.Conversation{}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var conversation domain.Conversation
		var message domain.Message
		var userInfo domain.UserInfo
		var members []*domain.UserInfo
		var membersString string
		values := []any{
			&conversation.ID,
			&conversation.CreatedAt,
			&conversation.Type,
			&conversation.Title,
			&conversation.Avatar,
			&conversation.UpdatedAt,
			&conversation.DeletedAt,
			&conversation.LastMessageID,
			&message.ID,
			&message.ConversationID,
			&message.UserID,
			&message.Type,
			&message.Body,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.ReplyTo,
			&userInfo.ID,
			&userInfo.FullName,
			&userInfo.Avatar,
			&userInfo.Type,
			&message.IsRead,
			&membersString,
		}
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		if message.ID == "" {
			conversation.LastMessage = nil
		} else {
			if userInfo.FullName == "" && userInfo.Avatar == "" && userInfo.Type == "" {
				message.User = nil
			} else {
				message.User = &userInfo
			}
			conversation.LastMessage = &message
		}
		fmt.Println(membersString)
		err = json.Unmarshal([]byte(membersString), &members)
		if err != nil {
			return nil, err
		}
		conversation.Members = members
		fmt.Println(conversation.Members)
		conversations = append(conversations, &conversation)
	}

	return conversations, nil
}

// UpdateLastMessageID implements domain.ConversationRepository.
func (c *conversationRepository) UpdateLastMessageID(ctx context.Context, conversationID string, lastMessageID string) error {
	var isCommited bool
	tx, err := c.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if !isCommited {
			tx.Rollback(ctx)
		}
	}()
	var temp int
	query := `SELECT 1 FROM conversation WHERE id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, query, conversationID).Scan(&temp)
	if err != nil {
		return err
	}
	query = `
		UPDATE conversation SET last_message_id = $1, updated_at = $2 WHERE id = $3
	`

	_, err = tx.Exec(ctx, query, lastMessageID, time.Now(), conversationID)
	if err != nil {
		return err
	}

	isCommited = true
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

var _ domain.ConversationRepository = &conversationRepository{}

func NewConversationRepository(db *pgxpool.Pool) domain.ConversationRepository {
	return &conversationRepository{db: db}
}

func (c *conversationRepository) CheckIsMemberOfConversation(ctx context.Context, userID string, conversationID string) (bool, error) {
	var isMember int
	query := `SELECT 1 FROM conversation_member WHERE user_id = $1 AND conversation_id = $2`
	err := c.db.QueryRow(ctx, query, userID, conversationID).Scan(&isMember)
	if err != nil {
		return false, err
	}
	return isMember == 1, nil
}
