package domain

import "time"

type RequestFriend struct {
	ID         string     `json:"id,omitempty"`
	FromUserID string     `json:"from_user_id,omitempty"`
	ToUserID   string     `json:"to_user_id,omitempty"`
	Status     string     `json:"status,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	FromUser   *UserInfo  `json:"-"`
	ToUser     *UserInfo  `json:"-"`
}

type RequestFriendStatus int

const (
	RequestFriendStatusPending RequestFriendStatus = iota
	RequestFriendStatusAccepted
	RequestFriendStatusRejected
)

func (rs RequestFriendStatus) String() string {
	return []string{"pending", "accepted", "rejected"}[rs]
}

func ToRequestFriendStatus(status string) RequestFriendStatus {
	switch status {
	case "pending":
		return RequestFriendStatusPending
	case "accepted":
		return RequestFriendStatusAccepted
	case "rejected":
		return RequestFriendStatusRejected
	}
	return RequestFriendStatusPending
}
