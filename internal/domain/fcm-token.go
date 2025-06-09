package domain

import "time"

type FcmToken struct {
	ID        string     `json:"id,omitempty"`
	UserID    string     `json:"user_id,omitempty"`
	Token     string     `json:"token,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
