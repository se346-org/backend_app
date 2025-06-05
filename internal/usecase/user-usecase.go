package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/chat-socio/backend/configuration"
	"github.com/chat-socio/backend/pkg/hash"
	"github.com/chat-socio/backend/pkg/jwt"
	"github.com/chat-socio/backend/pkg/pointer"
	"github.com/chat-socio/backend/pkg/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/internal/presenter"
	"github.com/chat-socio/backend/internal/utils"
	"github.com/chat-socio/backend/pkg/observability"
)

var (
	ErrNotFoundAccount = errors.New("account not found")
	ErrWrongPassword   = errors.New("wrong password")
)

type UserUseCase interface {
	Register(ctx context.Context, registerRequest *presenter.RegisterRequest) (*presenter.RegisterResponse, error)
	Login(ctx context.Context, loginRequest *presenter.LoginRequest) (*presenter.LoginResponse, error)
	GetUserInfo(ctx context.Context, userID string) (*presenter.GetUserInfoResponse, error)
	GetMyInfo(ctx context.Context) (*presenter.GetUserInfoResponse, error)
	GetUserInfoByEmail(ctx context.Context, email string) (*presenter.GetUserInfoResponse, error)
	GetListUser(ctx context.Context, userID string, keyword string, limit int, lastID string) ([]*presenter.GetUserInfoResponse, error)
	GetUserIDByAccountID(ctx context.Context, accountID string) (string, error)
}

type userUseCase struct {
	accountRepository      domain.AccountRepository
	userRepository         domain.UserRepository
	sessionRepository      domain.SessionRepository
	sessionCacheRepository domain.SessionCacheRepository
	userCacheRepository    domain.UserCacheRepository
	obs                    *observability.Observability
}

// GetUserIDByAccountID implements UserUseCase.
func (u *userUseCase) GetUserIDByAccountID(ctx context.Context, accountID string) (string, error) {
	ctx, span := u.obs.StartSpan(ctx, "UserUsecase.GetUserIDByAccountID")
	defer span()

	userID, err := u.userCacheRepository.GetUserIDByAccountID(ctx, accountID)
	if err != nil && err != redis.Nil {
		return "", err
	}

	if userID != "" {
		return userID, nil
	}

	user, err := u.userRepository.GetUserByAccountID(ctx, accountID)
	if err != nil {
		return "", err
	}

	err = u.userCacheRepository.SetUserIDByAccountID(ctx, accountID, user.ID)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

// GetListUser implements UserUseCase.
func (u *userUseCase) GetListUser(ctx context.Context, userID string, keyword string, limit int, lastID string) ([]*presenter.GetUserInfoResponse, error) {
	ctx, span := u.obs.StartSpan(ctx, "UserUsecase.GetListUser")
	defer span()

	users, err := u.userRepository.GetListUserWithConversation(ctx, userID, keyword, limit, lastID)
	if err != nil {
		return nil, err
	}

	userResponses := make([]*presenter.GetUserInfoResponse, 0)
	for _, user := range users {
		userResponses = append(userResponses, &presenter.GetUserInfoResponse{
			UserID:         user.ID,
			Email:          user.Email,
			FullName:       user.FullName,
			Avatar:         user.Avatar,
			Type:           user.Type,
			AccountID:      user.AccountID,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			ConversationID: user.ConversationID,
		})
	}

	return userResponses, nil
}

// GetUserInfoByEmail implements UserUseCase.
func (u *userUseCase) GetUserInfoByEmail(ctx context.Context, email string) (*presenter.GetUserInfoResponse, error) {
	ctx, span := u.obs.StartSpan(ctx, "UserUsecase.GetUserInfoByEmail")
	defer span()

	user, err := u.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &presenter.GetUserInfoResponse{
		UserID:    user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Avatar:    user.Avatar,
		Type:      user.Type,
		AccountID: user.AccountID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetMyInfo implements UserUseCase.
func (u *userUseCase) GetMyInfo(ctx context.Context) (*presenter.GetUserInfoResponse, error) {
	ctx, span := u.obs.StartSpan(ctx, "UserUsecase.GetMyInfo")
	defer span()

	accountID := ctx.Value(utils.AccountIDKey).(string)
	user, err := u.userRepository.GetUserByAccountID(ctx, accountID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	if err == pgx.ErrNoRows {
		return nil, ErrNotFoundAccount
	}

	return &presenter.GetUserInfoResponse{
		UserID:    user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Avatar:    user.Avatar,
		Type:      user.Type,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		AccountID: user.AccountID,
	}, nil
}

// GetUserInfo implements UserUseCase.
func (u *userUseCase) GetUserInfo(ctx context.Context, userID string) (*presenter.GetUserInfoResponse, error) {
	ctx, span := u.obs.StartSpan(ctx, "UserUsecase.GetUserInfo")
	defer span()

	user, err := u.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &presenter.GetUserInfoResponse{
		UserID:    user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Avatar:    user.Avatar,
		Type:      user.Type,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		AccountID: user.AccountID,
	}, nil
}

// Login implements UserUseCase.
func (u *userUseCase) Login(ctx context.Context, loginRequest *presenter.LoginRequest) (*presenter.LoginResponse, error) {
	ctx, span := u.obs.StartSpan(ctx, "UserUsecase.Login")
	defer span()

	userAgent := ctx.Value(utils.UserAgentKey).(string)
	ipAddress := ctx.Value(utils.IpAddressKey).(string)

	account, err := u.accountRepository.GetAccountByUsername(ctx, loginRequest.Email)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return nil, err
	}

	if err == pgx.ErrNoRows {
		return nil, ErrNotFoundAccount
	}

	if !hash.CheckPasswordHash(loginRequest.Password, account.Password) {
		return nil, ErrWrongPassword
	}

	sessionToken, err := uuid.NewID()
	if err != nil {
		return nil, err
	}

	// Create session
	now := time.Now()
	active := true
	expiredAt := now.Add(24 * time.Hour * 30 * 6) // 180 days
	session := &domain.Session{
		SessionToken: sessionToken,
		AccountID:    account.ID,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		IsActive:     &active,
		ExpiredAt:    &expiredAt,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}

	err = u.sessionRepository.CreateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	// Cache the session
	err = u.sessionCacheRepository.CreateSessionWithExpireTime(ctx, session)
	if err != nil {
		return nil, err
	}

	// Get user info and cache the user ID
	user, err := u.userRepository.GetUserByAccountID(ctx, account.ID)
	if err != nil {
		return nil, err
	}

	// Cache the user ID
	err = u.userCacheRepository.SetUserIDByAccountID(ctx, account.ID, user.ID)
	if err != nil {
		return nil, err
	}

	// return the access token
	var claims = jwt.JWTClaims{
		Jit: sessionToken,
		Sub: account.ID,
		Iat: now.Unix(),
		Exp: expiredAt.Unix(),
	}

	jwtToken, err := jwt.GenerateHS256JWT(claims, configuration.ConfigInstance.JWT.SecretKey)
	if err != nil {
		return nil, err
	}

	loginResponse := &presenter.LoginResponse{
		AccessToken: jwtToken,
	}

	return loginResponse, nil
}

// Register implements UserUseCase.
func (u *userUseCase) Register(ctx context.Context, registerRequest *presenter.RegisterRequest) (*presenter.RegisterResponse, error) {
	ctx, span := u.obs.StartSpan(ctx, "UserUsecase.Register")
	defer span()

	logger := u.obs.Logger.WithContext(ctx)
	logger.Info("Starting user registration", map[string]interface{}{
		"email": registerRequest.Email,
	})

	accountID, err := uuid.NewID()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := hash.HashPassword(registerRequest.Password)
	if err != nil {
		return nil, err
	}

	account := &domain.Account{
		ID:        accountID,
		Username:  registerRequest.Email,
		Password:  hashedPassword,
		CreatedAt: pointer.ToPtr(time.Now()),
		UpdatedAt: pointer.ToPtr(time.Now()),
	}

	userID, err := uuid.NewID()
	if err != nil {
		return nil, err
	}

	user := &domain.UserInfo{
		ID:        userID,
		AccountID: account.ID,
		Type:      domain.ExternalUserType,
		Email:     registerRequest.Email,
		FullName:  registerRequest.FullName,
		Avatar:    registerRequest.Avatar,
		CreatedAt: pointer.ToPtr(time.Now()),
		UpdatedAt: pointer.ToPtr(time.Now()),
	}

	err = u.accountRepository.CreateAccountUser(ctx, account, user)
	if err != nil {
		logger.WithError(err).Error("User registration failed")
		return nil, err
	}

	logger.Info("User registration completed successfully")
	registerResponse := &presenter.RegisterResponse{
		Success: true,
		Message: "Register success",
	}

	return registerResponse, nil
}

var _ UserUseCase = (*userUseCase)(nil)

func NewUserUseCase(accountRepository domain.AccountRepository, userRepository domain.UserRepository, sessionRepository domain.SessionRepository, sessionCacheRepository domain.SessionCacheRepository, userCacheRepository domain.UserCacheRepository, obs *observability.Observability) UserUseCase {
	return &userUseCase{
		accountRepository:      accountRepository,
		userRepository:         userRepository,
		sessionRepository:      sessionRepository,
		sessionCacheRepository: sessionCacheRepository,
		userCacheRepository:    userCacheRepository,
		obs:                    obs,
	}
}
