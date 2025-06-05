package handler

import (
	"context"
	"strconv"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/internal/presenter"
	"github.com/chat-socio/backend/internal/usecase"
	"github.com/chat-socio/backend/internal/utils"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type FriendHandler struct {
	FriendUseCase usecase.FriendUseCase
	Obs           *observability.Observability
	UserCacheRepo domain.UserCacheRepository
}

func NewFriendHandler(friendUseCase usecase.FriendUseCase, obs *observability.Observability, userCacheRepo domain.UserCacheRepository) *FriendHandler {
	return &FriendHandler{
		FriendUseCase: friendUseCase,
		Obs:           obs,
		UserCacheRepo: userCacheRepo,
	}
}

// @Summary Send friend request
// @Description Send a friend request to another user
// @Tags friends
// @Accept json
// @Produce json
// @Param friend_id path string true "Friend ID"
// @Success 200 {object} presenter.BaseEmptyResponse
// @Failure 400 {object} presenter.BaseEmptyResponse
// @Failure 500 {object} presenter.BaseEmptyResponse
// @Security BearerAuth
// @Router /auth/friend/{friend_id} [post]
func (h *FriendHandler) SendFriendRequest(ctx context.Context, c *app.RequestContext) {
	friendID := c.Param("friend_id")
	accountID, ok := ctx.Value(utils.AccountIDKey).(string)
	if !ok {
		logger := h.Obs.Logger.WithContext(ctx)
		logger.Error("Failed to get account ID from context")
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	logger := h.Obs.Logger.WithContext(ctx)
	logger.Info("Sending friend request", map[string]interface{}{
		"friend_id":  friendID,
		"account_id": accountID,
	})

	// Get user ID from account ID
	userID, err := h.UserCacheRepo.GetUserIDByAccountID(ctx, accountID)
	if err != nil {
		logger.WithError(err).Error("Failed to get user ID from account ID", map[string]interface{}{
			"account_id": accountID,
		})
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	logger.Info("Got user ID", map[string]interface{}{
		"user_id": userID,
	})

	err = h.FriendUseCase.SendFriendRequest(ctx, userID, friendID)
	if err != nil {
		logger.WithError(err).Error("Failed to send friend request", map[string]interface{}{
			"user_id":   userID,
			"friend_id": friendID,
		})
		c.JSON(consts.StatusBadRequest, presenter.BaseEmptyResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(consts.StatusOK, presenter.BaseEmptyResponse{
		Message: "Friend request sent successfully",
	})
}

// @Summary Accept friend request
// @Description Accept a friend request
// @Tags friends
// @Accept json
// @Produce json
// @Param friend_id path string true "Friend ID"
// @Success 200 {object} presenter.BaseEmptyResponse
// @Failure 400 {object} presenter.BaseEmptyResponse
// @Failure 500 {object} presenter.BaseEmptyResponse
// @Security BearerAuth
// @Router /auth/friend/{friend_id}/accept [post]
func (h *FriendHandler) AcceptFriendRequest(ctx context.Context, c *app.RequestContext) {
	friendID := c.Param("friend_id")
	accountID, ok := ctx.Value(utils.AccountIDKey).(string)
	if !ok {
		logger := h.Obs.Logger.WithContext(ctx)
		logger.Error("Failed to get account ID from context")
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	logger := h.Obs.Logger.WithContext(ctx)
	logger.Info("Accepting friend request", map[string]interface{}{
		"friend_id":  friendID,
		"account_id": accountID,
	})

	// Get user ID from account ID
	userID, err := h.UserCacheRepo.GetUserIDByAccountID(ctx, accountID)
	if err != nil {
		logger.WithError(err).Error("Failed to get user ID from account ID", map[string]interface{}{
			"account_id": accountID,
		})
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	logger.Info("Got user ID", map[string]interface{}{
		"user_id": userID,
	})

	err = h.FriendUseCase.AcceptFriendRequest(ctx, userID, friendID)
	if err != nil {
		logger.WithError(err).Error("Failed to accept friend request", map[string]interface{}{
			"user_id":   userID,
			"friend_id": friendID,
		})
		c.JSON(consts.StatusBadRequest, presenter.BaseEmptyResponse{
			Message: err.Error(),
		})
		return
	}

	logger.Info("Friend request accepted successfully", map[string]interface{}{
		"user_id":   userID,
		"friend_id": friendID,
	})

	c.JSON(consts.StatusOK, presenter.BaseEmptyResponse{
		Message: "Friend request accepted successfully",
	})
}

// @Summary Reject friend request
// @Description Reject a friend request
// @Tags friends
// @Accept json
// @Produce json
// @Param friend_id path string true "Friend ID"
// @Success 200 {object} presenter.BaseEmptyResponse
// @Failure 400 {object} presenter.BaseEmptyResponse
// @Failure 500 {object} presenter.BaseEmptyResponse
// @Security BearerAuth
// @Router /auth/friend/{friend_id}/reject [post]
func (h *FriendHandler) RejectFriendRequest(ctx context.Context, c *app.RequestContext) {
	friendID := c.Param("friend_id")
	accountID, ok := ctx.Value(utils.AccountIDKey).(string)
	if !ok {
		logger := h.Obs.Logger.WithContext(ctx)
		logger.Error("Failed to get account ID from context")
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	logger := h.Obs.Logger.WithContext(ctx)
	logger.Info("Rejecting friend request", map[string]interface{}{
		"friend_id":  friendID,
		"account_id": accountID,
	})

	// Get user ID from account ID
	userID, err := h.UserCacheRepo.GetUserIDByAccountID(ctx, accountID)
	if err != nil {
		logger.WithError(err).Error("Failed to get user ID from account ID", map[string]interface{}{
			"account_id": accountID,
		})
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	logger.Info("Got user ID", map[string]interface{}{
		"user_id": userID,
	})

	err = h.FriendUseCase.RejectFriendRequest(ctx, userID, friendID)
	if err != nil {
		logger.WithError(err).Error("Failed to reject friend request", map[string]interface{}{
			"user_id":   userID,
			"friend_id": friendID,
		})
		c.JSON(consts.StatusBadRequest, presenter.BaseEmptyResponse{
			Message: err.Error(),
		})
		return
	}

	logger.Info("Friend request rejected successfully", map[string]interface{}{
		"user_id":   userID,
		"friend_id": friendID,
	})

	c.JSON(consts.StatusOK, presenter.BaseEmptyResponse{
		Message: "Friend request rejected successfully",
	})
}

// @Summary Get friends
// @Description Get list of friends
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Friend status (pending/accepted)" default(accepted)
// @Param limit query int false "Limit" default(10)
// @Param last_id query string false "Last ID for pagination"
// @Success 200 {object} presenter.BaseFriendListResponse
// @Failure 400 {object} presenter.BaseEmptyResponse
// @Failure 500 {object} presenter.BaseEmptyResponse
// @Router /auth/friend [get]
func (h *FriendHandler) GetFriends(ctx context.Context, c *app.RequestContext) {
	accountID, ok := ctx.Value(utils.AccountIDKey).(string)
	if !ok {
		logger := h.Obs.Logger.WithContext(ctx)
		logger.Error("Failed to get account ID from context")
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	logger := h.Obs.Logger.WithContext(ctx)
	logger.Info("Getting friends", map[string]interface{}{
		"account_id": accountID,
	})

	// Get user ID from account ID
	userID, err := h.UserCacheRepo.GetUserIDByAccountID(ctx, accountID)
	if err != nil {
		logger.WithError(err).Error("Failed to get user ID from account ID", map[string]interface{}{
			"account_id": accountID,
		})
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	logger.Info("Got user ID", map[string]interface{}{
		"user_id": userID,
	})

	status := domain.FriendStatus(c.Query("status"))
	if status == "" {
		status = domain.FriendStatusAccepted
	}
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		logger.Error("Invalid limit parameter", map[string]interface{}{
			"limit": limitStr,
		})
		c.JSON(consts.StatusBadRequest, presenter.BaseEmptyResponse{
			Message: "Invalid limit parameter",
		})
		return
	}
	lastID := c.Query("last_id")

	logger.Info("Getting friends with parameters", map[string]interface{}{
		"user_id": userID,
		"status":  status,
		"limit":   limit,
		"last_id": lastID,
	})

	friends, err := h.FriendUseCase.GetFriends(ctx, userID, status, limit, lastID)
	if err != nil {
		logger.WithError(err).Error("Failed to get friends", map[string]interface{}{
			"user_id": userID,
			"status":  status,
		})
		c.JSON(consts.StatusBadRequest, presenter.BaseEmptyResponse{
			Message: err.Error(),
		})
		return
	}

	// Convert []*domain.FriendWithUser to []domain.FriendWithUser
	friendSlice := make([]domain.FriendWithUser, len(friends))
	for i, friend := range friends {
		friendSlice[i] = *friend
	}

	logger.Info("Friends retrieved successfully", map[string]interface{}{
		"user_id": userID,
		"count":   len(friends),
	})

	c.JSON(consts.StatusOK, presenter.BaseFriendListResponse{
		Message: "Friends retrieved successfully",
		Data:    presenter.ToFriendResponses(friendSlice),
	})
}

// @Summary Unfriend
// @Description Remove a friend
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param friend_id path string true "Friend ID"
// @Success 200 {object} presenter.BaseEmptyResponse
// @Failure 400 {object} presenter.BaseEmptyResponse
// @Failure 500 {object} presenter.BaseEmptyResponse
// @Router /auth/friend/{friend_id} [delete]
func (h *FriendHandler) Unfriend(ctx context.Context, c *app.RequestContext) {
	friendID := c.Param("friend_id")
	accountID, ok := ctx.Value(utils.AccountIDKey).(string)
	if !ok {
		logger := h.Obs.Logger.WithContext(ctx)
		logger.Error("Failed to get account ID from context")
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	logger := h.Obs.Logger.WithContext(ctx)
	logger.Info("Unfriending user", map[string]interface{}{
		"friend_id":  friendID,
		"account_id": accountID,
	})

	// Get user ID from account ID
	userID, err := h.UserCacheRepo.GetUserIDByAccountID(ctx, accountID)
	if err != nil {
		logger.WithError(err).Error("Failed to get user ID from account ID", map[string]interface{}{
			"account_id": accountID,
		})
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	logger.Info("Got user ID", map[string]interface{}{
		"user_id": userID,
	})

	err = h.FriendUseCase.Unfriend(ctx, userID, friendID)
	if err != nil {
		logger.WithError(err).Error("Failed to unfriend", map[string]interface{}{
			"user_id":   userID,
			"friend_id": friendID,
		})
		c.JSON(consts.StatusBadRequest, presenter.BaseEmptyResponse{
			Message: err.Error(),
		})
		return
	}

	logger.Info("Friend removed successfully", map[string]interface{}{
		"user_id":   userID,
		"friend_id": friendID,
	})

	c.JSON(consts.StatusOK, presenter.BaseEmptyResponse{
		Message: "Friend removed successfully",
	})
}

// @Summary Get friend requests sent
// @Description Get list of pending friend requests sent by the user
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param last_id query string false "Last ID for pagination"
// @Success 200 {object} presenter.BaseFriendListResponse
// @Failure 400 {object} presenter.BaseEmptyResponse
// @Failure 500 {object} presenter.BaseEmptyResponse
// @Router /auth/friend/requests [get]
func (h *FriendHandler) GetFriendRequests(ctx context.Context, c *app.RequestContext) {
	accountID, ok := ctx.Value(utils.AccountIDKey).(string)
	if !ok {
		logger := h.Obs.Logger.WithContext(ctx)
		logger.Error("Failed to get account ID from context")
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	logger := h.Obs.Logger.WithContext(ctx)
	logger.Info("Getting friend requests", map[string]interface{}{
		"account_id": accountID,
	})

	// Get user ID from account ID
	userID, err := h.UserCacheRepo.GetUserIDByAccountID(ctx, accountID)
	if err != nil {
		logger.WithError(err).Error("Failed to get user ID from account ID", map[string]interface{}{
			"account_id": accountID,
		})
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(consts.StatusBadRequest, presenter.BaseEmptyResponse{
			Message: "Invalid limit parameter",
		})
		return
	}
	lastID := c.Query("last_id")

	friends, err := h.FriendUseCase.GetFriendRequests(ctx, userID, limit, lastID)
	if err != nil {
		logger.WithError(err).Error("Failed to get friend requests", map[string]interface{}{
			"user_id": userID,
		})
		c.JSON(consts.StatusBadRequest, presenter.BaseEmptyResponse{
			Message: err.Error(),
		})
		return
	}

	// Convert []*domain.FriendWithUser to []domain.FriendWithUser
	friendSlice := make([]domain.FriendWithUser, len(friends))
	for i, friend := range friends {
		friendSlice[i] = *friend
	}

	c.JSON(consts.StatusOK, presenter.BaseFriendListResponse{
		Message: "Friend requests retrieved successfully",
		Data:    presenter.ToFriendResponses(friendSlice),
	})
}

// @Summary Get friend requests received
// @Description Get list of friend requests received by the user
// @Tags friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param last_id query string false "Last ID for pagination"
// @Success 200 {object} presenter.BaseFriendListResponse
// @Failure 400 {object} presenter.BaseEmptyResponse
// @Failure 500 {object} presenter.BaseEmptyResponse
// @Router /auth/friend/received [get]
func (h *FriendHandler) GetFriendRequestsReceived(ctx context.Context, c *app.RequestContext) {
	accountID, ok := ctx.Value(utils.AccountIDKey).(string)
	if !ok {
		logger := h.Obs.Logger.WithContext(ctx)
		logger.Error("Failed to get account ID from context")
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	logger := h.Obs.Logger.WithContext(ctx)
	logger.Info("Getting friend requests received", map[string]interface{}{
		"account_id": accountID,
	})

	// Get user ID from account ID
	userID, err := h.UserCacheRepo.GetUserIDByAccountID(ctx, accountID)
	if err != nil {
		logger.WithError(err).Error("Failed to get user ID from account ID", map[string]interface{}{
			"account_id": accountID,
		})
		c.JSON(consts.StatusInternalServerError, presenter.BaseEmptyResponse{
			Message: "Internal server error",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(consts.StatusBadRequest, presenter.BaseEmptyResponse{
			Message: "Invalid limit parameter",
		})
		return
	}
	lastID := c.Query("last_id")

	friends, err := h.FriendUseCase.GetFriendRequestsReceived(ctx, userID, limit, lastID)
	if err != nil {
		logger.WithError(err).Error("Failed to get friend requests received", map[string]interface{}{
			"user_id": userID,
		})
		c.JSON(consts.StatusBadRequest, presenter.BaseEmptyResponse{
			Message: err.Error(),
		})
		return
	}

	// Convert []*domain.FriendWithUser to []domain.FriendWithUser
	friendSlice := make([]domain.FriendWithUser, len(friends))
	for i, friend := range friends {
		friendSlice[i] = *friend
	}

	c.JSON(consts.StatusOK, presenter.BaseFriendListResponse{
		Message: "Friend requests received retrieved successfully",
		Data:    presenter.ToFriendResponses(friendSlice),
	})
} 