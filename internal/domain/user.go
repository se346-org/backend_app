package domain

import (
	"time"
)

const (
	ExternalUserType = "EXTERNAL"
	InternalUserType = "INTERNAL"
)

type UserInfo struct {
	ID             string     `json:"id,omitempty"`
	AccountID      string     `json:"account_id,omitempty"`
	Type           string     `json:"type,omitempty"`
	Email          string     `json:"email,omitempty"`
	FullName       string     `json:"full_name,omitempty"`
	Avatar         string     `json:"avatar,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	ConversationID *string    `json:"-"` // for query conversation with another user
}

func (u *UserInfo) TableName() string {
	return "user_info"
}

func (u *UserInfo) MapFields() ([]string, []any) {
	return []string{
			"id",
			"account_id",
			"type",
			"email",
			"full_name",
			"avatar",
			"created_at",
			"updated_at",
		}, []any{
			&u.ID,
			&u.AccountID,
			&u.Type,
			&u.Email,
			&u.FullName,
			&u.Avatar,
			&u.CreatedAt,
			&u.UpdatedAt,
		}
}
