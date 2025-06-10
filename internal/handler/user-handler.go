package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chat-socio/backend/configuration"
	"github.com/chat-socio/backend/internal/presenter"
	"github.com/chat-socio/backend/internal/usecase"
	"github.com/chat-socio/backend/internal/utils"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
)

type UserHandler struct {
	UserUseCase usecase.UserUseCase
	Obs         *observability.Observability
}

func (uh *UserHandler) Register(ctx context.Context, c *app.RequestContext) {
	ctx, span := uh.Obs.StartSpan(ctx, "UserHandler.Register")
	defer span()

	var registerRequest presenter.RegisterRequest
	if err := c.Bind(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[any]{
			Message: fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	err := registerRequest.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[any]{
			Message: fmt.Sprintf("Validation error: %v", err),
		})
		return
	}

	response, err := uh.UserUseCase.Register(ctx, &registerRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{Message: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, presenter.BaseResponse[*presenter.RegisterResponse]{
		Message: "User registered successfully",
		Data:    response,
	})
}

func (uh *UserHandler) Login(ctx context.Context, c *app.RequestContext) {
	ctx, span := uh.Obs.StartSpan(ctx, "UserHandler.Login")
	defer span()

	// Extract user agent and IP address from context
	userAgent := c.Request.Header.Get("User-Agent")
	ipAddress := c.ClientIP()
	ctx = context.WithValue(ctx, utils.UserAgentKey, userAgent)
	ctx = context.WithValue(ctx, utils.IpAddressKey, ipAddress)
	// Bind the request body to the LoginRequest struct
	var loginRequest presenter.LoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[any]{
			Message: fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	err := loginRequest.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[any]{
			Message: fmt.Sprintf("Validation error: %v", err),
		})
		return
	}

	response, err := uh.UserUseCase.Login(ctx, &loginRequest)
	if err != nil && err != usecase.ErrNotFoundAccount && err != usecase.ErrWrongPassword {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{Message: fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	if err == usecase.ErrNotFoundAccount {
		c.JSON(http.StatusNotFound, presenter.BaseResponse[any]{Message: "Account not found"})
		return
	}

	if err == usecase.ErrWrongPassword {
		c.JSON(http.StatusNotFound, presenter.BaseResponse[any]{Message: "Wrong password"})
		return
	}

	c.SetCookie("access_token", response.AccessToken, 3600*24*6*30, "/", configuration.ConfigInstance.Server.Origin, protocol.CookieSameSiteNoneMode, true, true)

	c.JSON(http.StatusOK, presenter.BaseResponse[*presenter.LoginResponse]{
		Message: "User logged in successfully",
		Data:    response,
	})
}

func (uh *UserHandler) GetMyInfo(ctx context.Context, c *app.RequestContext) {
	ctx, span := uh.Obs.StartSpan(ctx, "UserHandler.GetMyInfo")
	defer span()

	// Extract account ID from the request context
	accountID := ctx.Value(utils.AccountIDKey).(string)
	if accountID == "" {
		c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{Message: "Could not get account ID"})
		return
	}

	userInfo, err := uh.UserUseCase.GetMyInfo(ctx)
	if err != nil && err != usecase.ErrNotFoundAccount {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{Message: fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	if err == usecase.ErrNotFoundAccount {
		c.JSON(http.StatusNotFound, presenter.BaseResponse[any]{Message: "Account not found"})
		return
	}

	c.JSON(http.StatusOK, presenter.BaseResponse[*presenter.GetUserInfoResponse]{
		Message: "User info retrieved successfully",
		Data:    userInfo,
	})
}

func (uh *UserHandler) GetListUser(ctx context.Context, c *app.RequestContext) {
	ctx, span := uh.Obs.StartSpan(ctx, "UserHandler.GetListUser")
	defer span()

	// Extract params from request
	accountID := ctx.Value(utils.AccountIDKey).(string)
	userID, err := uh.UserUseCase.GetUserIDByAccountID(ctx, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{Message: fmt.Sprintf("Internal server error: %v", err)})
		return
	}
	keyword := c.Query("keyword")
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	lastID := c.Query("last_id")
	response, err := uh.UserUseCase.GetListUser(ctx, userID, keyword, limit, lastID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{Message: fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	c.JSON(http.StatusOK, presenter.BaseResponse[[]*presenter.GetUserInfoResponse]{
		Message: "User list retrieved successfully",
		Data:    response,
	})
}

func (uh *UserHandler) UpdateUser(ctx context.Context, c *app.RequestContext) {
	ctx, span := uh.Obs.StartSpan(ctx, "UserHandler.UpdateUser")
	defer span()

	// Extract account ID from the request context
	accountID := ctx.Value(utils.AccountIDKey).(string)
	if accountID == "" {
		c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{Message: "Could not get account ID"})
		return
	}

	var updateRequest presenter.UpdateUserRequest
	if err := c.Bind(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[any]{
			Message: fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	err := updateRequest.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, presenter.BaseResponse[any]{
			Message: fmt.Sprintf("Validation error: %v", err),
		})
		return
	}

	response, err := uh.UserUseCase.UpdateUser(ctx, &updateRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{Message: fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	c.JSON(http.StatusOK, presenter.BaseResponse[*presenter.GetUserInfoResponse]{
		Message: "User updated successfully",
		Data:    response,
	})
}
