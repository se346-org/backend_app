package domain

import "time"

type Contact struct {
	ID        string     `json:"id,omitempty"`
	UserID    string     `json:"user_id,omitempty"`
	FriendID  string     `json:"friend_id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (c *Contact) TableName() string {
	return "contact"
}

func (c *Contact) MapFields() ([]string, []any) {
	return []string{
			"id",
			"user_id",
			"friend_id",
			"created_at",
			"updated_at", 
		}, []any{
			&c.ID,
			&c.UserID,
			&c.FriendID,
			&c.CreatedAt,
			&c.UpdatedAt,
		}
}
