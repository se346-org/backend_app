package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/internal/presenter"
	"github.com/chat-socio/backend/internal/usecase"
	"github.com/chat-socio/backend/internal/utils"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/chat-socio/backend/pkg/pointer"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/jackc/pgx/v5/pgconn"
)

type FCMHandler struct {
	FcmUseCase  usecase.FCMUseCase
	UserUseCase usecase.UserUseCase
	Obs         *observability.Observability
}

func (h *FCMHandler) CreateFCMToken(ctx context.Context, c *app.RequestContext) {
	ctx, span := h.Obs.StartSpan(ctx, "FCMHandler.CreateFCMToken")
	defer span()
	logger := h.Obs.Logger.WithContext(ctx)

	accountID := ctx.Value(utils.AccountIDKey).(string)

	userID, err := h.UserUseCase.GetUserIDByAccountID(ctx, accountID)
	if err != nil {
		logger.Error("failed to get user id by account id", err)
		c.JSON(http.StatusInternalServerError, &presenter.BaseResponse[any]{
			Message: err.Error(),
		})
		return
	}

	var fcmToken domain.FcmToken
	err = c.Bind(&fcmToken)
	if err != nil {
		logger.Error("failed to bind json", err)
		c.JSON(http.StatusBadRequest, &presenter.BaseResponse[any]{
			Message: err.Error(),
		})
		return
	}

	if fcmToken.Token == "" {
		c.JSON(http.StatusBadRequest, &presenter.BaseResponse[any]{
			Message: "token is required",
		})
		return
	}

	fcmToken.UserID = userID
	fcmToken.CreatedAt = pointer.ToPtr(time.Now())
	fcmToken.UpdatedAt = pointer.ToPtr(time.Now())

	err = h.FcmUseCase.CreateFCMToken(ctx, &fcmToken)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			c.JSON(http.StatusOK, &presenter.BaseResponse[any]{
				Message: "success",
			})
			return
		}
	}

	if err != nil {
		logger.Error("failed to create fcm token", err)
		c.JSON(http.StatusInternalServerError, &presenter.BaseResponse[any]{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &presenter.BaseResponse[any]{
		Message: "success",
	})
}

func (h *FCMHandler) DeleteFCMToken(ctx context.Context, c *app.RequestContext) {
	ctx, span := h.Obs.StartSpan(ctx, "FCMHandler.DeleteFCMToken")
	defer span()
	logger := h.Obs.Logger.WithContext(ctx)

	accountID := ctx.Value(utils.AccountIDKey).(string)
	if accountID == "" {
		c.JSON(http.StatusUnauthorized, &presenter.BaseResponse[any]{
			Message: "unauthorized",
		})
		return
	}

	userID, err := h.UserUseCase.GetUserIDByAccountID(ctx, accountID)
	if err != nil {
		logger.Error("failed to get user id by account id", err)
		c.JSON(http.StatusInternalServerError, &presenter.BaseResponse[any]{
			Message: err.Error(),
		})
		return
	}

	token := c.Query("fcm_token")
	if token == "" {
		c.JSON(http.StatusBadRequest, &presenter.BaseResponse[any]{
			Message: "fcm_token is required",
		})
		return
	}

	err = h.FcmUseCase.DeleteFCMToken(ctx, userID, token)
	if err != nil {
		logger.Error("failed to delete fcm token", err)
		c.JSON(http.StatusInternalServerError, &presenter.BaseResponse[any]{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &presenter.BaseResponse[any]{
		Message: "success",
	})
}
