package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/chat-socio/backend/configuration"
	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/internal/presenter"
	"github.com/chat-socio/backend/internal/utils"
	"github.com/chat-socio/backend/pkg/jwt"
	"github.com/cloudwego/hertz/pkg/app"
)

type Middleware struct {
	sessionCacheRepository domain.SessionCacheRepository
	sessionRepository      domain.SessionRepository
}

func NewMiddleware(sessionCacheRepository domain.SessionCacheRepository, sessionRepository domain.SessionRepository) *Middleware {
	return &Middleware{
		sessionCacheRepository: sessionCacheRepository,
		sessionRepository:      sessionRepository,
	}
}

func (m *Middleware) AuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var accessToken string
		// Extract the session token from the request header
		bearToken := c.Request.Header.Get("Authorization")
		if bearToken == "" {
			// Extract the session token from the cookie if not found in the header
			b := c.Cookie("access_token")
			if len(b) == 0 {
				c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{
					Message: "Unauthorized",
				})
				c.Abort()
				return
			}

			accessToken = string(b)
		} else {
			// Extract the session token from the header (e.g., "Bearer <token>")
			strings := strings.Split(bearToken, " ")
			if len(strings) != 2 || strings[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{
					Message: "Unauthorized",
				})
				c.Abort()
				return
			}
			accessToken = strings[1]
		}

		// Check if the access token is valid
		jwtToken, err := jwt.ParseHS256JWT(accessToken, configuration.ConfigInstance.JWT.SecretKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		// Extract claims from the token
		claims, err := jwt.ExtractClaims(jwtToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		// validate expired token
		if claims.Exp < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		// Check if the session exists in the cache
		session, err := m.sessionCacheRepository.GetSessionByToken(ctx, claims.Jit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{
				Message: "Internal server error",
			})
			c.Abort()
			return
		}

		if session == nil {
			// If the session is not found in the cache, check the database
			session, err = m.sessionRepository.GetSessionByToken(ctx, claims.Jit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{
					Message: "Internal server error",
				})
				c.Abort()
				return
			}
			if session == nil {
				c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{
					Message: "Unauthorized",
				})
				c.Abort()
				return
			}
			// Check if the session is expired
			if session.ExpiredAt != nil && session.ExpiredAt.Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{
					Message: "Session expired",
				})
				c.Abort()
				return
			}
			// Check if the session is active
			if session.IsActive != nil && !*session.IsActive {
				c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{
					Message: "Session is not active",
				})
				c.Abort()
				return
			}
			// Cache the session for future use
			err = m.sessionCacheRepository.CreateSessionWithExpireTime(ctx, session)
			if err != nil {
				c.JSON(http.StatusInternalServerError, presenter.BaseResponse[any]{
					Message: "Internal server error",
				})
				c.Abort()
				return
			}
		}

		// Check session is active
		if session.IsActive != nil && !*session.IsActive {
			c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{
				Message: "Session is not active",
			})
			c.Abort()
			return
		}

		// Check user agent and IP address
		userAgent := c.Request.Header.Get("User-Agent")
		ipAddress := c.ClientIP()
		if session.UserAgent != userAgent || session.IPAddress != ipAddress {
			c.JSON(http.StatusUnauthorized, presenter.BaseResponse[any]{
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		// Set the account ID in the context for later use
		ctx = context.WithValue(ctx, utils.AccountIDKey, claims.Sub)
		ctx = context.WithValue(ctx, utils.SessionTokenKey, claims.Jit)
		c.Next(ctx)
	}
}
