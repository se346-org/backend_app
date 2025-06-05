package presenter

import (
	"time"
	"github.com/chat-socio/backend/internal/domain"
)

// FriendResponse represents a friend response
type FriendResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	FriendID  string `json:"friend_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	User      struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		FullName string `json:"full_name"`
		Avatar   string `json:"avatar"`
	} `json:"user"`
}

// ToUserResponse converts a domain UserInfo to a UserResponse
func ToUserResponse(user *domain.UserInfo) UserResponse {
	return UserResponse{
		UserID:    user.ID,
		FullName:  user.FullName,
		Avatar:    user.Avatar,
		UserType:  "user", // Default type if not specified
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToFriendResponse converts a domain FriendWithUser to a FriendResponse.
func ToFriendResponse(friend domain.FriendWithUser) FriendResponse {
	return FriendResponse{
		ID:        friend.ID,
		UserID:    friend.UserID,
		FriendID:  friend.FriendID,
		Status:    string(friend.Status),
		CreatedAt: friend.CreatedAt.Format(time.RFC3339),
		UpdatedAt: friend.UpdatedAt.Format(time.RFC3339),
		User: struct {
			ID       string `json:"id"`
			Email    string `json:"email"`
			FullName string `json:"full_name"`
			Avatar   string `json:"avatar"`
		}{
			ID:       friend.User.ID,
			Email:    friend.User.Email,
			FullName: friend.User.FullName,
			Avatar:   friend.User.Avatar,
		},
	}
}

// ToFriendResponses converts a slice of domain FriendWithUser to a slice of FriendResponse.
func ToFriendResponses(friends []domain.FriendWithUser) []FriendResponse {
	var responses []FriendResponse
	for _, friend := range friends {
		responses = append(responses, ToFriendResponse(friend))
	}
	return responses
} 