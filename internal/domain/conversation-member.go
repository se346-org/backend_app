package domain

import "time"

type ConversationMember struct {
	ID             string     `json:"id,omitempty"`
	ConversationID string     `json:"conversation_id,omitempty"`
	UserID         string     `json:"user_id,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	User           *UserInfo  `json:"-"`
}

func (c *ConversationMember) TableName() string {
	return "conversation_member"
}

func (c *ConversationMember) MapFields() ([]string, []any) {
	return []string{
			"id",
			"conversation_id",
			"user_id",
			"created_at",
			"updated_at",
			"deleted_at",
		}, []any{
			&c.ID,
			&c.ConversationID,
			&c.UserID,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.DeletedAt,
		}
}

type ConversationMemberWithUser struct {
	ConversationID string
	UserID         string
	FullName       string
	Avatar         string
	UserType       string
}
