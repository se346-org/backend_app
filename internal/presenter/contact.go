package presenter

import "time"

type RequestFriendRequest struct {
	FromUserID   string `json:"-"`
	TargetUserID string `json:"target_user_id,omitempty"`
}

type RequestFriendResponse struct {
	IsSuccess bool `json:"is_success"`
}

type RejectRequestFriendRequest struct {
	RequestFriendID string `json:"request_friend_id,omitempty"`
}

type RejectRequestFriendResponse struct {
	IsSuccess bool `json:"is_success"`
}

type AcceptRequestFriendRequest struct {
	RequestFriendID string `json:"request_friend_id,omitempty"`
}

type AcceptRequestFriendResponse struct {
	IsSuccess bool `json:"is_success"`
}

type GetListRequestFriendSentRequest struct {
	UserID string `json:"user_id,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	LastID string `json:"last_id,omitempty"`
}

type GetListRequestFriendSentResponse struct {
	FriendRequestID string        `json:"friend_request_id,omitempty"`
	TargetUser      *UserResponse `json:"target_user,omitempty"`
	Status          string        `json:"status,omitempty"`
	CreatedAt       *time.Time    `json:"created_at,omitempty"`
	UpdatedAt       *time.Time    `json:"updated_at,omitempty"`
}

type GetListRequestFriendReceivedRequest struct {
	UserID string `json:"user_id,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	LastID string `json:"last_id,omitempty"`
}

type GetListRequestFriendReceivedResponse struct {
	FriendRequestID string        `json:"friend_request_id,omitempty"`
	FromUser        *UserResponse `json:"from_user,omitempty"`
	Status          string        `json:"status,omitempty"`
	CreatedAt       *time.Time    `json:"created_at,omitempty"`
	UpdatedAt       *time.Time    `json:"updated_at,omitempty"`
}

type GetListContactRequest struct {
	UserID string `json:"user_id,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	LastID string `json:"last_id,omitempty"`
}

type GetListContactResponse struct {
	ContactID string        `json:"contact_id,omitempty"`
	User      *UserResponse `json:"user,omitempty"`
	CreatedAt *time.Time    `json:"created_at,omitempty"`
	UpdatedAt *time.Time    `json:"updated_at,omitempty"`
}
