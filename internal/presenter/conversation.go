package presenter

import (
	"errors"
	"time"

	"github.com/chat-socio/backend/internal/domain"
)

type ConversationResponse struct {
	ConversationID string                        `json:"conversation_id,omitempty"`
	Title          string                        `json:"title,omitempty"`
	Avatar         string                        `json:"avatar,omitempty"`
	LastMessageID  string                        `json:"last_message_id,omitempty"`
	CreatedAt      *time.Time                    `json:"created_at,omitempty"`
	UpdatedAt      *time.Time                    `json:"updated_at,omitempty"`
	Type           string                        `json:"type,omitempty"`
	Members        []*ConversationMemberResponse `json:"members,omitempty"`
}

type ConversationMemberResponse struct {
	UserID   string `json:"user_id,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	UserType string `json:"user_type,omitempty"`
}

type CreateConversationRequest struct {
	Title   string   `json:"title,omitempty"`
	Type    string   `json:"type,omitempty"`
	Avatar  string   `json:"avatar,omitempty"`
	Members []string `json:"members,omitempty"`
}

func (c *CreateConversationRequest) Validate() error {
	if len(c.Members) < 2 {
		return errors.New("at least 2 members are required")
	}
	if c.Type != domain.ConversationTypeGroup && c.Type != domain.ConversationTypeDM {
		return errors.New("invalid conversation type")
	}
	if c.Type == domain.ConversationTypeDM && len(c.Members) != 2 {
		return errors.New("invalid number of members for DM")
	}
	if c.Type == domain.ConversationTypeGroup && len(c.Members) < 2 {
		return errors.New("invalid number of members for group")
	}
	return nil
}

type SendMessageRequest struct {
	ConversationID string `json:"conversation_id,omitempty"`
	UserID         string `json:"user_id,omitempty"`
	Type           string `json:"type,omitempty"`
	Body           string `json:"body,omitempty"`
	ReplyTo        string `json:"reply_to,omitempty"`
	UserOnlineID   string `json:"user_online_id,omitempty"` // for ignore user online id
}

func (s *SendMessageRequest) Validate() error {
	if s.ConversationID == "" {
		return errors.New("conversation_id is required")
	}
	if s.UserID == "" {
		return errors.New("user_id is required")
	}
	if s.Type == "" {
		return errors.New("type is required")
	}
	if s.Body == "" {
		return errors.New("body is required")
	}
	if s.UserOnlineID == "" {
		return errors.New("user_online_id is required")
	}
	return nil
}

type MessageResponse struct {
	MessageID      string        `json:"message_id,omitempty"`
	Body           string        `json:"body,omitempty"`
	CreatedAt      *time.Time    `json:"created_at,omitempty"`
	UpdatedAt      *time.Time    `json:"updated_at,omitempty"`
	ConversationID string        `json:"conversation_id,omitempty"`
	User           *UserResponse `json:"user,omitempty"`
	Type           string        `json:"type,omitempty"`
	DeletedAt      *time.Time    `json:"deleted_at,omitempty"`
	ReplyTo        string        `json:"reply_to,omitempty"`
}

type GetListConversationResponse struct {
	ConversationID string                        `json:"conversation_id,omitempty"`
	Title          string                        `json:"title,omitempty"`
	Avatar         string                        `json:"avatar,omitempty"`
	LastMessageID  string                        `json:"last_message_id,omitempty"`
	CreatedAt      *time.Time                    `json:"created_at,omitempty"`
	UpdatedAt      *time.Time                    `json:"updated_at,omitempty"`
	Type           string                        `json:"type,omitempty"`
	LastMessage    *MessageResponse              `json:"last_message,omitempty"`
	Members        []*ConversationMemberResponse `json:"members,omitempty"`
}

type SeenMessageResponse struct {
	MessageID      string     `json:"message_id,omitempty"`
	UserID         string     `json:"user_id,omitempty"`
	ConversationID string     `json:"conversation_id,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type SeenMessageRequest struct {
	ConversationID string `json:"conversation_id,omitempty"`
	MessageID      string `json:"message_id,omitempty"`
}
