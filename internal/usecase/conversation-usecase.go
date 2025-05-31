package usecase

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/internal/presenter"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/chat-socio/backend/pkg/pointer"
	"github.com/chat-socio/backend/pkg/uuid"
	"github.com/chat-socio/backend/pubsub"
	"github.com/jackc/pgx/v5"
)

type ConversationUseCase interface {
	GetListConversationByUserID(ctx context.Context, userID string, lastMessageID string, limit int) ([]*presenter.GetListConversationResponse, error)
	GetConversationByID(ctx context.Context, conversationID string) (*presenter.ConversationResponse, error)
	GetListMessageByConversationID(ctx context.Context, userID string, conversationID string, lastMessageID string, limit int) ([]*presenter.MessageResponse, error)
	CreateConversation(ctx context.Context, conversation *presenter.CreateConversationRequest) (*presenter.ConversationResponse, error)
	SendMessage(ctx context.Context, message *presenter.SendMessageRequest) (*presenter.MessageResponse, error)
	HandleNewMessage(ctx context.Context, message *domain.WebSocketMessage) error
	HandleUpdateLastMessageID(ctx context.Context, data domain.UpdateLastMessageID) error
	SeenMessage(ctx context.Context, messageID string, userID string, conversationID string) error
	GetListSeenMessageByConversationID(ctx context.Context, conversationID string) ([]*presenter.SeenMessageResponse, error)
	HandleSeenMessage(ctx context.Context, message *domain.SeenMessage) error
}

type conversationUseCase struct {
	conversationRepository domain.ConversationRepository
	messageRepository      domain.MessageRepository
	messagePublisher       pubsub.Publisher
	userOnlineRepository   domain.UserOnlineRepository
	userRepository         domain.UserRepository
	seenMessageRepository  domain.SeenMessageRepository
	obs                    *observability.Observability
}

func (c *conversationUseCase) HandleSeenMessage(ctx context.Context, message *domain.SeenMessage) error {
	logger := c.obs.Logger.WithContext(ctx)
	id, err := uuid.NewID()
	if err != nil {
		logger.Error("failed to generate id", err)
		return err
	}
	message.ID = id
	err = c.seenMessageRepository.CreateSeenMessage(ctx, message)
	if err != nil {
		logger.Error("failed to upsert seen message", err, message)
		return err
	}
	wsMessage := &domain.WebSocketMessage{
		Type: domain.WsSeenMessage,
		Payload: map[string]any{
			"conversation_id": message.ConversationID,
			"message_id":      message.MessageID,
			"user_id":         message.UserID,
		},
	}
	err = c.messagePublisher.Publish(ctx, domain.SUBJECT_NEW_MESSAGE, wsMessage)
	if err != nil {
		logger.Error("failed to publish seen message", err, message)
		return err
	}
	return nil
}

// GetListSeenMessageByConversationID implements ConversationUseCase.
func (c *conversationUseCase) GetListSeenMessageByConversationID(ctx context.Context, conversationID string) ([]*presenter.SeenMessageResponse, error) {
	ctx, span := c.obs.StartSpan(ctx, "ConversationUsecase.GetListSeenMessageByConversationID")
	defer span()
	logger := c.obs.Logger.WithContext(ctx)
	seenMessages, err := c.seenMessageRepository.GetListSeenMessageByConversationID(ctx, conversationID)
	if err != nil {
		logger.Error("failed to get list seen message by conversation id", err, conversationID)
		return nil, err
	}
	seenMessageResponses := make([]*presenter.SeenMessageResponse, 0)
	for _, seenMessage := range seenMessages {
		seenMessageResponses = append(seenMessageResponses, &presenter.SeenMessageResponse{
			MessageID:      seenMessage.MessageID,
			UserID:         seenMessage.UserID,
			ConversationID: seenMessage.ConversationID,
			CreatedAt:      seenMessage.CreatedAt,
			UpdatedAt:      seenMessage.UpdatedAt,
		})
	}
	return seenMessageResponses, nil
}

// SeenMessage implements ConversationUseCase.
func (c *conversationUseCase) SeenMessage(ctx context.Context, messageID string, userID string, conversationID string) error {
	ctx, span := c.obs.StartSpan(ctx, "ConversationUsecase.SeenMessage")
	defer span()
	logger := c.obs.Logger.WithContext(ctx)
	err := c.messagePublisher.Publish(ctx, domain.SUBJECT_SEEN_MESSAGE, domain.SeenMessage{
		MessageID:      messageID,
		UserID:         userID,
		ConversationID: conversationID,
	})
	if err != nil {
		logger.Error("failed to publish seen message", err, messageID, userID, conversationID)
		return err
	}
	return nil
}

// HandleUpdateLastMessageID implements ConversationUseCase.
func (c *conversationUseCase) HandleUpdateLastMessageID(ctx context.Context, data domain.UpdateLastMessageID) error {
	logger := c.obs.Logger.WithContext(ctx)
	err := c.conversationRepository.UpdateLastMessageID(ctx, data.ConversationID, data.MessageID)
	if err != nil {
		logger.Error("error update last message id", err, data)
		return err
	}
	conversation, _, err := c.conversationRepository.GetConversationByID(ctx, data.ConversationID)
	if err != nil {
		logger.Error("error get conversation by id", err, data)
		return err
	}

	message, err := c.messageRepository.GetMessageByID(ctx, data.MessageID)
	if err != nil {
		logger.Error("error get message by id", err, data)
		return err
	}

	userMap, err := pointer.ToMap(message.User)
	if err != nil {
		logger.Error("error convert message to map", err, message)
		return err
	}
	messageMap, err := pointer.ToMap(message)
	if err != nil {
		logger.Error("error convert message to map", err, message)
		return err
	}

	messageMap["user"] = userMap

	// conversation.LastMessage = message
	conversationMap, err := pointer.ToMap(conversation)
	if err != nil {
		logger.Error("error convert conversation to map", err, conversation)
		return err
	}
	conversationMap["last_message"] = messageMap

	wsMessage := &domain.WebSocketMessage{
		Type:    domain.WsUpdateLastMessage,
		Payload: conversationMap,
	}

	return c.messagePublisher.Publish(ctx, domain.SUBJECT_NEW_MESSAGE, wsMessage)
}

func (c *conversationUseCase) getUserOnlineByConversationID(ctx context.Context, conversationID string) ([]*domain.UserOnline, error) {
	return c.userOnlineRepository.GetUserOnlineByConversationID(ctx, conversationID)
}

func (c *conversationUseCase) handleSendEventNewMessage(ctx context.Context, message *domain.WebSocketMessage) error {
	logger := c.obs.Logger.WithContext(ctx)
	// get user online by conversation id
	userOnlines, err := c.getUserOnlineByConversationID(ctx, message.Payload["conversation_id"].(string))
	if err != nil {
		logger.Error("error get user online by conversation id", err, message)
		return err
	}

	mapIgnoreUserOnlines := make(map[string]bool)
	for _, uo := range message.IgnoreUserOnlines {
		mapIgnoreUserOnlines[uo] = true
	}

	// send message to websocket
	for _, userOnline := range userOnlines {
		// exclude user who send message
		if _, ok := mapIgnoreUserOnlines[userOnline.ID]; ok {
			continue
		}
		wsConn, ok := domain.WebSocket.GetConnection(userOnline.ConnectionID)
		if !ok {
			continue
		}

		b, err := json.Marshal(message)
		if err != nil {
			logger.Error("failed to marshal message to json", err, message)
			continue
		}
		err = wsConn.SendMessage(b)
		if err != nil {
			logger.Error("failed to send message to websocket", err, message)
			continue
		}
	}
	return nil
}

func (c *conversationUseCase) handleSendEventUpdateLastMessageID(ctx context.Context, message *domain.WebSocketMessage) error {
	// get user online by conversation id
	logger := c.obs.Logger.WithContext(ctx)
	userOnlines, err := c.getUserOnlineByConversationID(ctx, message.Payload["id"].(string))
	if err != nil {
		logger.Error("error get user online by conversation id", err, message)
		return err
	}

	mapIgnoreUserOnlines := make(map[string]bool)
	for _, uo := range message.IgnoreUserOnlines {
		mapIgnoreUserOnlines[uo] = true
	}
	// send message to websocket
	for _, userOnline := range userOnlines {
		// exclude user who send message
		if _, ok := mapIgnoreUserOnlines[userOnline.ID]; ok {
			continue
		}
		wsConn, ok := domain.WebSocket.GetConnection(userOnline.ConnectionID)
		if !ok {
			continue
		}
		b, err := json.Marshal(message)
		if err != nil {
			logger.Error("failed to marshal message to json", err, message)
			continue
		}
		err = wsConn.SendMessage(b)
		if err != nil {
			logger.Error("failed to send message to websocket", err, message)
			continue
		}
	}
	return nil
}

// HandleNewMessage implements ConversationUseCase.
func (c *conversationUseCase) HandleNewMessage(ctx context.Context, message *domain.WebSocketMessage) error {
	switch message.Type {
	case domain.WsMessage:
		return c.handleSendEventNewMessage(ctx, message)
	case domain.WsUpdateLastMessage:
		return c.handleSendEventUpdateLastMessageID(ctx, message)
	case domain.WsSeenMessage:
		return c.handleSendEventNewMessage(ctx, message)
	}
	return nil
}

func NewConversationUseCase(conversationRepository domain.ConversationRepository, messageRepository domain.MessageRepository, messagePublisher pubsub.Publisher, userOnlineRepository domain.UserOnlineRepository, userRepository domain.UserRepository, seenMessageRepository domain.SeenMessageRepository, obs *observability.Observability) ConversationUseCase {
	return &conversationUseCase{
		conversationRepository: conversationRepository,
		messageRepository:      messageRepository,
		messagePublisher:       messagePublisher,
		userOnlineRepository:   userOnlineRepository,
		userRepository:         userRepository,
		seenMessageRepository:  seenMessageRepository,
		obs:                    obs,
	}
}

// CreateConversation implements ConversationUseCase.
func (c *conversationUseCase) CreateConversation(ctx context.Context, conversation *presenter.CreateConversationRequest) (*presenter.ConversationResponse, error) {
	conversationID, err := uuid.NewID()
	if err != nil {
		return nil, err
	}
	conversationDomain := &domain.Conversation{
		ID:        conversationID,
		Type:      conversation.Type,
		Title:     conversation.Title,
		Avatar:    conversation.Avatar,
		CreatedAt: pointer.ToPtr(time.Now()),
		UpdatedAt: pointer.ToPtr(time.Now()),
	}
	conversationMembers := make([]*domain.ConversationMember, 0)
	conversationMemberResponses := make([]*presenter.ConversationMemberResponse, 0)
	for _, userID := range conversation.Members {
		conversationMemberID, err := uuid.NewID()
		if err != nil {
			return nil, err
		}
		conversationMembers = append(conversationMembers, &domain.ConversationMember{
			ID:             conversationMemberID,
			ConversationID: conversationID,
			UserID:         userID,
			CreatedAt:      pointer.ToPtr(time.Now()),
			UpdatedAt:      pointer.ToPtr(time.Now()),
		})
		conversationMemberResponses = append(conversationMemberResponses, &presenter.ConversationMemberResponse{
			UserID: userID,
		})
	}
	conversationDomain, err = c.conversationRepository.CreateConversation(ctx, conversationDomain, conversationMembers)
	if err != nil {
		return nil, err
	}
	return &presenter.ConversationResponse{
		ConversationID: conversationDomain.ID,
		Type:           conversationDomain.Type,
		Title:          conversationDomain.Title,
		Avatar:         conversationDomain.Avatar,
		Members:        conversationMemberResponses,
	}, nil
}

// GetConversationByID implements ConversationUseCase.
func (c *conversationUseCase) GetConversationByID(ctx context.Context, conversationID string) (*presenter.ConversationResponse, error) {
	conversation, conversationMembers, err := c.conversationRepository.GetConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	conversationMemberResponses := make([]*presenter.ConversationMemberResponse, 0)
	for _, conversationMember := range conversationMembers {
		conversationMemberResponses = append(conversationMemberResponses, &presenter.ConversationMemberResponse{
			UserID:   conversationMember.UserID,
			FullName: conversationMember.FullName,
			Avatar:   conversationMember.Avatar,
			UserType: conversationMember.UserType,
		})
	}
	return &presenter.ConversationResponse{
		ConversationID: conversation.ID,
		Type:           conversation.Type,
		Title:          conversation.Title,
		Avatar:         conversation.Avatar,
		Members:        conversationMemberResponses,
	}, nil
}

// GetListConversationByUserID implements ConversationUseCase.
func (c *conversationUseCase) GetListConversationByUserID(ctx context.Context, userID string, lastMessageID string, limit int) ([]*presenter.GetListConversationResponse, error) {
	conversations, err := c.conversationRepository.GetListConversationByUserID(ctx, userID, lastMessageID, limit)
	if err != nil {
		return nil, err
	}
	conversationResponses := make([]*presenter.GetListConversationResponse, 0)
	for _, conversation := range conversations {
		conversationResponse := &presenter.GetListConversationResponse{
			ConversationID: conversation.ID,
			Type:           conversation.Type,
			Title:          conversation.Title,
			Avatar:         conversation.Avatar,
			LastMessageID:  conversation.LastMessageID,
			CreatedAt:      conversation.CreatedAt,
			UpdatedAt:      conversation.UpdatedAt,
		}

		if conversation.LastMessage != nil {
			conversationResponse.LastMessage = &presenter.MessageResponse{
				MessageID:      conversation.LastMessage.ID,
				Body:           conversation.LastMessage.Body,
				CreatedAt:      conversation.LastMessage.CreatedAt,
				UpdatedAt:      conversation.LastMessage.UpdatedAt,
				Type:           conversation.LastMessage.Type,
				DeletedAt:      conversation.LastMessage.DeletedAt,
				ReplyTo:        conversation.LastMessage.ReplyTo,
				ConversationID: conversation.LastMessage.ConversationID,
				User: &presenter.UserResponse{
					UserID:   conversation.LastMessage.User.ID,
					FullName: conversation.LastMessage.User.FullName,
					Avatar:   conversation.LastMessage.User.Avatar,
					UserType: conversation.LastMessage.User.Type,
				},
			}
		}
		if len(conversation.Members) > 0 {
			for _, member := range conversation.Members {
				conversationResponse.Members = append(conversationResponse.Members, &presenter.ConversationMemberResponse{
					UserID:   member.ID,
					FullName: member.FullName,
					Avatar:   member.Avatar,
					UserType: member.Type,
				})
			}
		}
		conversationResponses = append(conversationResponses, conversationResponse)
	}
	return conversationResponses, nil
}

// GetListMessageByConversationID implements ConversationUseCase.
func (c *conversationUseCase) GetListMessageByConversationID(ctx context.Context, userID string, conversationID string, lastMessageID string, limit int) ([]*presenter.MessageResponse, error) {
	ctx, span := c.obs.StartSpan(ctx, "ConversationUsecase.GetListMessageByConversationID")
	defer span()
	// check is member of conversation
	isMember, err := c.conversationRepository.CheckIsMemberOfConversation(ctx, userID, conversationID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	if !isMember {
		return nil, domain.ErrNotFoundMemberOfConversation
	}
	// get list message by conversation id
	messages, err := c.messageRepository.GetListMessageByConversationID(ctx, conversationID, lastMessageID, limit)
	if err != nil {
		return nil, err
	}
	if err == pgx.ErrNoRows {
		return []*presenter.MessageResponse{}, nil
	}
	messageResponses := make([]*presenter.MessageResponse, 0)
	for _, message := range messages {
		messageResponses = append(messageResponses, &presenter.MessageResponse{
			MessageID:      message.ID,
			Body:           message.Body,
			CreatedAt:      message.CreatedAt,
			UpdatedAt:      message.UpdatedAt,
			Type:           message.Type,
			DeletedAt:      message.DeletedAt,
			ReplyTo:        message.ReplyTo,
			ConversationID: message.ConversationID,
			User: &presenter.UserResponse{
				UserID:   message.User.ID,
				FullName: message.User.FullName,
				Avatar:   message.User.Avatar,
				UserType: message.User.Type,
			},
		})
	}
	return messageResponses, nil
}

// SendMessage implements ConversationUseCase.
func (c *conversationUseCase) SendMessage(ctx context.Context, message *presenter.SendMessageRequest) (*presenter.MessageResponse, error) {
	ctx, span := c.obs.StartSpan(ctx, "ConversationUsecase.SendMessage")
	defer span()
	logger := c.obs.Logger.WithContext(ctx)
	messageID, err := uuid.NewID()
	if err != nil {
		return nil, err
	}
	messageDomain := &domain.Message{
		ID:             messageID,
		ConversationID: message.ConversationID,
		UserID:         message.UserID,
		Type:           message.Type,
		Body:           message.Body,
		CreatedAt:      pointer.ToPtr(time.Now()),
		UpdatedAt:      pointer.ToPtr(time.Now()),
		ReplyTo:        message.ReplyTo,
	}
	messageDomain, err = c.messageRepository.CreateMessage(ctx, messageDomain)
	if err != nil {
		logger.Error("error create message", err, message)
		return nil, err
	}

	user, err := c.userRepository.GetUserByID(ctx, message.UserID)
	if err != nil {
		logger.Error("error get user by id", err, message)
		return nil, err
	}
	userMap, err := pointer.ToMap(user)
	if err != nil {
		logger.Error("error convert user to map", err, user)
		return nil, err
	}

	messageMap, err := pointer.ToMap(messageDomain)
	if err != nil {
		logger.Error("error convert message to map", err, message)
		return nil, err
	}
	messageMap["user"] = userMap
	wsMessage := &domain.WebSocketMessage{
		Type:              domain.WsMessage,
		Payload:           messageMap,
		IgnoreUserOnlines: []string{message.UserID},
	}
	// send message to websocket
	err = c.messagePublisher.Publish(ctx, domain.SUBJECT_NEW_MESSAGE, wsMessage)
	if err != nil {
		logger.Error("failed to publish message to websocket", err, message)
		// return nil, fmt.Errorf("failed to publish message to websocket: %w", err)
	}
	// update last message id of conversation
	err = c.messagePublisher.Publish(ctx, domain.SUBJECT_UPDATE_LAST_MESSAGE_ID, domain.UpdateLastMessageID{
		ConversationID: message.ConversationID,
		MessageID:      messageDomain.ID,
	})

	if err != nil {
		logger.Error("failed to publish message to websocket", err, message)
		// return nil, fmt.Errorf("failed to publish message to websocket: %w", err)
	}

	return &presenter.MessageResponse{
		MessageID:      messageDomain.ID,
		Body:           messageDomain.Body,
		CreatedAt:      messageDomain.CreatedAt,
		UpdatedAt:      messageDomain.UpdatedAt,
		Type:           messageDomain.Type,
		DeletedAt:      messageDomain.DeletedAt,
		ReplyTo:        messageDomain.ReplyTo,
		ConversationID: messageDomain.ConversationID,
	}, nil
}

var _ ConversationUseCase = &conversationUseCase{}
