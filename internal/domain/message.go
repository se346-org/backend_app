package domain

import "time"

const (
	MessageTypeText     = "text"
	MessageTypeImage    = "image"
	MessageTypeVideo    = "video"
	MessageTypeAudio    = "audio"
	MessageTypeFile     = "file"
	MessageTypeLocation = "location"
	MessageTypeContact  = "contact"
	MessageTypeSticker  = "sticker"
	MessageTypePoll     = "poll"
	MessageTypeSystem   = "system"
)

type Message struct {
	ID             string     `json:"id,omitempty"`
	ConversationID string     `json:"conversation_id,omitempty"`
	UserID         string     `json:"user_id,omitempty"`
	Type           string     `json:"type,omitempty"`
	Body           string     `json:"body,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	ReplyTo        string     `json:"reply_to,omitempty"`
	User           *UserInfo  `json:"-"`
	IgnoreSend     string     `json:"-"`
}

func (m *Message) TableName() string {
	return "message"
}

func (m *Message) MapFields() ([]string, []any) {
	return []string{
			"id",
			"conversation_id",
			"user_id",
			"type",
			"body",
			"created_at",
			"updated_at",
			"deleted_at",
			"reply_to",
		}, []any{
			&m.ID,
			&m.ConversationID,
			&m.UserID,
			&m.Type,
			&m.Body,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.DeletedAt,
			&m.ReplyTo,
		}
}
