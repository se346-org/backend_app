package domain

import "errors"

var (
	ErrNoRows = errors.New("no rows in result set")

	ErrNotFoundMemberOfConversation = errors.New("user is not a member of conversation")

	ErrConversationAlreadyExist = errors.New("conversation already exist")
)
