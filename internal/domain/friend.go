package domain

import "time"

type FriendStatus string

const (
	FriendStatusPending  FriendStatus = "pending"
	FriendStatusAccepted FriendStatus = "accepted"
	FriendStatusRejected FriendStatus = "rejected"
)

type Friend struct {
	ID        string       `json:"id"`
	UserID    string       `json:"user_id"`
	FriendID  string       `json:"friend_id"`
	Status    FriendStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type FriendWithUser struct {
	Friend
	User *UserInfo `json:"user"`
} 