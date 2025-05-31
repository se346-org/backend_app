package domain

import "time"

type UserOnline struct {
	ID           string     `json:"id,omitempty"`
	UserID       string     `json:"user_id,omitempty"`
	ConnectionID string     `json:"connection_id,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	User         *UserInfo  `json:"user,omitempty"`
}
