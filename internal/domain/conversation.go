package domain

import "time"

const (
	ConversationTypeDM    = "DM"
	ConversationTypeGroup = "GROUP"
)

type Conversation struct {
	ID            string      `json:"id,omitempty"`
	CreatedAt     *time.Time  `json:"created_at,omitempty"`
	Type          string      `json:"type,omitempty"`
	Title         string      `json:"title,omitempty"`
	Avatar        string      `json:"avatar,omitempty"`
	UpdatedAt     *time.Time  `json:"updated_at,omitempty"`
	DeletedAt     *time.Time  `json:"deleted_at,omitempty"`
	LastMessageID string      `json:"last_message_id,omitempty"`
	LastMessage   *Message    `json:"-"`
	Members       []*UserInfo `json:"members,omitempty"`
}

func (c *Conversation) TableName() string {
	return "conversation"
}

func (c *Conversation) MapFields() ([]string, []any) {
	return []string{
			"id",
			"created_at",
			"type",
			"title",
			"avatar",
			"updated_at",
			"deleted_at",
			"last_message_id",
		}, []any{
			&c.ID,
			&c.CreatedAt,
			&c.Type,
			&c.Title,
			&c.Avatar,
			&c.UpdatedAt,
			&c.DeletedAt,
			&c.LastMessageID,
		}
}
