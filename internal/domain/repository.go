package domain

import (
	"context"
	"time"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account *Account) error
	GetAccountByUsername(ctx context.Context, username string) (*Account, error)
	GetAccountByID(ctx context.Context, id string) (*Account, error)
	UpdatePassword(ctx context.Context, id string, password string) error
	CreateAccountUser(ctx context.Context, account *Account, user *UserInfo) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *UserInfo) error
	GetUserByID(ctx context.Context, id string) (*UserInfo, error)
	GetUserByEmail(ctx context.Context, email string) (*UserInfo, error)
	GetUserByAccountID(ctx context.Context, accountID string) (*UserInfo, error)
	UpdateUser(ctx context.Context, user *UserInfo) error
	GetListUser(ctx context.Context, keyword string, limit int, lastID string) ([]*UserInfo, error)
	GetListUserWithConversation(ctx context.Context, userID string, keyword string, limit int, lastID string) ([]*UserInfo, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByToken(ctx context.Context, token string) (*Session, error)
	GetListSessionByAccountID(ctx context.Context, accountID string) ([]*Session, error)
	DeactivateSession(ctx context.Context, token string) error
	DeactiveAllSessionByAccountID(ctx context.Context, accountID string) error
	UpdateExpiredAt(ctx context.Context, token string, newExpiredAt *time.Time) error
}

type SessionCacheRepository interface {
	CreateSessionWithExpireTime(ctx context.Context, session *Session) error
	GetSessionByToken(ctx context.Context, token string) (*Session, error)
	DeleteSession(ctx context.Context, token string) error
}

type UserOnlineRepository interface {
	CreateUserOnline(ctx context.Context, userOnline *UserOnline) error
	DeleteUserOnline(ctx context.Context, id string) error
	GetUserOnlineByConversationID(ctx context.Context, conversationID string) ([]*UserOnline, error)
}

type ConversationRepository interface {
	CreateConversation(ctx context.Context, conversation *Conversation, conversationMembers []*ConversationMember) (*Conversation, error)
	GetListConversationByUserID(ctx context.Context, userID string, lastMessageID string, limit int) ([]*Conversation, error)
	GetConversationByID(ctx context.Context, id string) (*Conversation, []*ConversationMemberWithUser, error)
	UpdateLastMessageID(ctx context.Context, conversationID string, lastMessageID string) error
	CheckIsMemberOfConversation(ctx context.Context, userID string, conversationID string) (bool, error)
	CheckDMConversationExist(ctx context.Context, userID1 string, userID2 string) (*Conversation, error)
}

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *Message) (*Message, error)
	GetListMessageByConversationID(ctx context.Context, conversationID string, lastID string, limit int) ([]*Message, error)
	GetMessageByID(ctx context.Context, id string) (*Message, error)
}

type UserCacheRepository interface {
	GetUserIDByAccountID(ctx context.Context, accountID string) (string, error)
	SetUserIDByAccountID(ctx context.Context, accountID string, userID string) error
}

type SeenMessageRepository interface {
	CreateSeenMessage(ctx context.Context, seenMessage *SeenMessage) error
	GetListSeenMessageByConversationID(ctx context.Context, conversationID string) ([]*SeenMessage, error)
}

type ContactRepository interface {
	CreateContacts(ctx context.Context, contacts []*Contact) error
	GetListContactByUserID(ctx context.Context, userID string, limit int, lastID string) ([]*Contact, error)
	CreateRequestFriend(ctx context.Context, requestFriend *RequestFriend) error
	GetListRequestFriendSentByUserID(ctx context.Context, userID string, limit int, lastID string) ([]*RequestFriend, error)
	GetListRequestFriendReceivedByUserID(ctx context.Context, userID string, limit int, lastID string) ([]*RequestFriend, error)
	AcceptRequestFriend(ctx context.Context, requestFriendID string) error
	UpdateRequestFriendStatus(ctx context.Context, id string, status string) error
}

type FcmTokenRepository interface {
	CreateFcmToken(ctx context.Context, fcmToken *FcmToken) error
	GetFcmTokenByUserID(ctx context.Context, userID string) ([]*FcmToken, error)
	DeleteFcmToken(ctx context.Context, id string) error
	DeleteFcmTokenByUserIDAndToken(ctx context.Context, userID, token string) error
}
