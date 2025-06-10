package usecase

import (
	"context"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/chat-socio/backend/pkg/uuid"
)

type FCMUseCase interface {
	CreateFCMToken(ctx context.Context, fcmToken *domain.FcmToken) error
	DeleteFCMToken(ctx context.Context, userID, token string) error
}

type fcmUseCase struct {
	fcmTokenRepository domain.FcmTokenRepository
	obs                *observability.Observability
}

func NewFCMUseCase(fcmTokenRepository domain.FcmTokenRepository, obs *observability.Observability) *fcmUseCase {
	return &fcmUseCase{fcmTokenRepository: fcmTokenRepository, obs: obs}
}

// CreateFCMToken implements FCMUseCase.
func (f *fcmUseCase) CreateFCMToken(ctx context.Context, fcmToken *domain.FcmToken) error {
	ctx, span := f.obs.StartSpan(ctx, "FcmUseCase.CreateFCMToken")
	defer span()
	logger := f.obs.Logger.WithContext(ctx)
	id, err := uuid.NewID()
	if err != nil {
		logger.Error("failed to generate id", err)
		return err
	}
	fcmToken.ID = id
	err = f.fcmTokenRepository.CreateFcmToken(ctx, fcmToken)
	if err != nil {
		logger.Error("failed to create fcm token", err)
		return err
	}
	return nil
}

// DeleteFCMToken implements FCMUseCase.
func (f *fcmUseCase) DeleteFCMToken(ctx context.Context, userID, token string) error {
	ctx, span := f.obs.StartSpan(ctx, "FcmUseCase.DeleteFCMToken")
	defer span()
	logger := f.obs.Logger.WithContext(ctx)
	err := f.fcmTokenRepository.DeleteFcmTokenByUserIDAndToken(ctx, userID, token)
	if err != nil {
		logger.Error("failed to delete fcm token by user id and token", err)
		return err
	}
	return nil
}

var _ FCMUseCase = &fcmUseCase{}
