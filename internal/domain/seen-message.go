package domain

import "time"

type SeenMessage struct {
	ID             string     `json:"id,omitempty"`
	MessageID      string     `json:"message_id,omitempty"`
	UserID         string     `json:"user_id,omitempty"`
	ConversationID string     `json:"conversation_id,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

func (s *SeenMessage) TableName() string {
	return "seen_message"
}

func (s *SeenMessage) MapFields() ([]string, []any) {
	return []string{
			"id",
			"message_id",
			"user_id",
			"conversation_id",
			"created_at",
			"updated_at",
		}, []any{
			&s.ID,
			&s.MessageID,
			&s.UserID,
			&s.ConversationID,
			&s.CreatedAt,
			&s.UpdatedAt,
		}
}
