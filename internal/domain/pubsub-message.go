package domain

type UpdateLastMessageID struct {
	ConversationID string `json:"conversation_id,omitempty"`
	MessageID      string `json:"message_id,omitempty"`
}
