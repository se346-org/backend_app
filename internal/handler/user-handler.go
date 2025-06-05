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

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags user
// @Accept json
// @Produce json
// @Param request body presenter.RegisterRequest true "Register Request"
// @Success 200 {object} presenter.RegisterResponse "User registered successfully"
// @Failure 400 {string} string "Invalid request or validation error"
// @Failure 500 {string} string "Internal server error"
// @Router /user/register [post]
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

// @Summary Login user
// @Description Login with email and password
// @Tags user
// @Accept json
// @Produce json
// @Param request body presenter.LoginRequest true "Login Request"
// @Success 200 {object} presenter.LoginResponse "User logged in successfully"
// @Failure 400 {string} string "Invalid request or validation error"
// @Failure 404 {string} string "Account not found or wrong password"
// @Failure 500 {string} string "Internal server error"
// @Router /user/login [post]
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

// @Summary Get my info
// @Description Get information about the current user
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} presenter.GetUserInfoResponse "User info retrieved successfully"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Account not found"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/user/info [get]
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

// @Summary Get list of users
// @Description Get a list of users with optional keyword, limit, and last_id
// @Tags user
// @Accept json
// @Produce json
// @Param keyword query string false "Keyword to search"
// @Param limit query int false "Limit number of users"
// @Param last_id query string false "Last user ID for pagination"
// @Success 200 {array} presenter.GetUserInfoResponse "User list retrieved successfully"
// @Failure 500 {string} string "Internal server error"
// @Router /user/list [get]
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
