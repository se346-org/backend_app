package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/chat-socio/backend/internal/presenter"
	"github.com/chat-socio/backend/internal/usecase"
	"github.com/chat-socio/backend/internal/utils"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/cloudwego/hertz/pkg/app"
)

type ConversationHandler struct {
	ConversationUseCase usecase.ConversationUseCase
	UserUseCase         usecase.UserUseCase
	Obs                 *observability.Observability
}

// @Summary Send a message
// @Description Send a message in a conversation
// @Tags conversation
// @Accept json
// @Produce json
// @Param request body presenter.SendMessageRequest true "Send Message Request"
// @Success 200 {object} presenter.MessageResponse "Message sent successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /conversation/message [post]
func (ch *ConversationHandler) SendMessage(ctx context.Context, c *app.RequestContext) {
	ctx, span := ch.Obs.StartSpan(ctx, "ConversationHandler.SendMessage")
	defer span()

	// get account id from context
	accountID := ctx.Value(utils.AccountIDKey)
	if accountID == nil {
		c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{
			Message: "Unauthorized",
		})
		return
	}

	userID, err := ch.UserUseCase.GetUserIDByAccountID(ctx, accountID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{
			Message: err.Error(),
		})
		return
	}

	var request presenter.SendMessageRequest
	if err := c.BindAndValidate(&request); err != nil {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[any]{
			Message: err.Error(),
		})
		return
	}

	request.UserID = userID

	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[any]{
			Message: err.Error(),
		})
		return
	}

	sendMessageResponse, err := ch.ConversationUseCase.SendMessage(ctx, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, presenter.BaseResponse[*presenter.MessageResponse]{
		Data:    sendMessageResponse,
		Message: "Message sent successfully",
	})
}

// @Summary Create a conversation
// @Description Create a new conversation
// @Tags conversation
// @Accept json
// @Produce json
// @Param request body presenter.CreateConversationRequest true "Create Conversation Request"
// @Success 200 {object} presenter.ConversationResponse "Conversation created successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /conversation/create [post]
func (ch *ConversationHandler) CreateConversation(ctx context.Context, c *app.RequestContext) {
	ctx, span := ch.Obs.StartSpan(ctx, "ConversationHandler.CreateConversation")
	defer span()

	var request presenter.CreateConversationRequest
	if err := c.BindAndValidate(&request); err != nil {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[*presenter.ConversationResponse]{
			Message: err.Error(),
		})
		return
	}

	if err := request.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[*presenter.ConversationResponse]{
			Message: err.Error(),
		})
		return
	}

	createConversationResponse, err := ch.ConversationUseCase.CreateConversation(ctx, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[*presenter.ConversationResponse]{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, presenter.BaseResponse[*presenter.ConversationResponse]{
		Data:    createConversationResponse,
		Message: "Conversation created successfully",
	})
}

// @Summary Get list of conversations
// @Description Get a list of conversations for the current user
// @Tags conversation
// @Accept json
// @Produce json
// @Param last_message_id query string false "Last message ID for pagination"
// @Param limit query int false "Limit number of conversations"
// @Success 200 {array} presenter.GetListConversationResponse "List conversation fetched successfully"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /conversation/list [get]
func (ch *ConversationHandler) GetListConversation(ctx context.Context, c *app.RequestContext) {
	ctx, span := ch.Obs.StartSpan(ctx, "ConversationHandler.GetListConversation")
	defer span()

	accountID := ctx.Value(utils.AccountIDKey)
	if accountID == nil {
		c.JSON(http.StatusUnauthorized, presenter.BaseResponse[[]*presenter.GetListConversationResponse]{
			Message: "Unauthorized",
		})
		return
	}

	userID, err := ch.UserUseCase.GetUserIDByAccountID(ctx, accountID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[[]*presenter.GetListConversationResponse]{
			Message: err.Error(),
		})
		return
	}

	conversationID := c.Query("conversation_id")
	if conversationID != "" {
		conversation, err := ch.ConversationUseCase.GetConversationByID(ctx, conversationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, presenter.BaseResponse[[]*presenter.GetListConversationResponse]{
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, presenter.BaseResponse[*presenter.ConversationResponse]{
			Data:    conversation,
			Message: "Conversation fetched successfully",
		})
		return
	}

	lastMessageID := c.Query("last_message_id")
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 20
	}

	listConversation, err := ch.ConversationUseCase.GetListConversationByUserID(ctx, userID, lastMessageID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[[]*presenter.GetListConversationResponse]{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, presenter.BaseResponse[[]*presenter.GetListConversationResponse]{
		Data:    listConversation,
		Message: "List conversation fetched successfully",
	})
}

// @Summary Get list of messages
// @Description Get a list of messages in a conversation
// @Tags conversation
// @Accept json
// @Produce json
// @Param conversation_id query string true "Conversation ID"
// @Param last_message_id query string false "Last message ID for pagination"
// @Param limit query int false "Limit number of messages"
// @Success 200 {array} presenter.MessageResponse "List message fetched successfully"
// @Failure 400 {string} string "Conversation ID is required"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /conversation/messages [get]
func (ch *ConversationHandler) GetListMessage(ctx context.Context, c *app.RequestContext) {
	ctx, span := ch.Obs.StartSpan(ctx, "ConversationHandler.GetListMessage")
	defer span()

	accountID := ctx.Value(utils.AccountIDKey)
	if accountID == nil {
		c.JSON(http.StatusUnauthorized, presenter.BaseResponse[[]*presenter.MessageResponse]{
			Message: "Unauthorized",
		})
		return
	}

	userID, err := ch.UserUseCase.GetUserIDByAccountID(ctx, accountID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, presenter.BaseResponse[[]*presenter.MessageResponse]{
			Message: err.Error(),
		})
		return
	}

	conversationID := c.Query("conversation_id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[[]*presenter.MessageResponse]{
			Message: "Conversation ID is required",
		})
		return
	}
	lastMessageID := c.Query("last_message_id")
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 20
	}

	listMessage, err := ch.ConversationUseCase.GetListMessageByConversationID(ctx, userID, conversationID, lastMessageID, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, presenter.BaseResponse[[]*presenter.MessageResponse]{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, presenter.BaseResponse[[]*presenter.MessageResponse]{
		Data:    listMessage,
		Message: "List message fetched successfully",
	})
}

// @Summary Get conversation by ID
// @Description Get a conversation by its ID
// @Tags conversation
// @Accept json
// @Produce json
// @Param conversation_id query string true "Conversation ID"
// @Success 200 {object} presenter.ConversationResponse "Conversation fetched successfully"
// @Failure 500 {string} string "Internal server error"
// @Router /conversation [get]
func (ch *ConversationHandler) GetConversationByID(ctx context.Context, c *app.RequestContext) {
	ctx, span := ch.Obs.StartSpan(ctx, "ConversationHandler.GetConversationByID")
	defer span()

	conversationID := c.Query("conversation_id")

	conversation, err := ch.ConversationUseCase.GetConversationByID(ctx, conversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[*presenter.ConversationResponse]{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, presenter.BaseResponse[*presenter.ConversationResponse]{
		Data:    conversation,
		Message: "Conversation fetched successfully",
	})
}

// @Summary Seen message
// @Description Mark a message as seen
// @Tags conversation
// @Accept json
// @Produce json
// @Param request body presenter.SeenMessageRequest true "Seen Message Request"
// @Success 200 {string} string "Seen message successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /conversation/seen [post]
func (ch *ConversationHandler) SeenMessage(ctx context.Context, c *app.RequestContext) {
	ctx, span := ch.Obs.StartSpan(ctx, "ConversationHandler.SeenMessage")
	defer span()

	var request presenter.SeenMessageRequest
	if err := c.BindAndValidate(&request); err != nil {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[any]{
			Message: err.Error(),
		})
		return
	}

	accountID := ctx.Value(utils.AccountIDKey)
	if accountID == nil {
		c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{
			Message: "Unauthorized",
		})
		return
	}

	userID, err := ch.UserUseCase.GetUserIDByAccountID(ctx, accountID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{
			Message: err.Error(),
		})
		return
	}

	err = ch.ConversationUseCase.SeenMessage(ctx, request.MessageID, userID, request.ConversationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, presenter.BaseResponse[any]{
		Message: "Seen message successfully",
	})
}
