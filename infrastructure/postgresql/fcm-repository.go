package postgresql

import (
	"context"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/jackc/pgx/v5/pgxpool"
)

type fcmRepository struct {
	db  *pgxpool.Pool
	obs *observability.Observability
}

// DeleteFcmTokenByUserIDAndToken implements domain.FcmTokenRepository.
func (f *fcmRepository) DeleteFcmTokenByUserIDAndToken(ctx context.Context, userID string, token string) error {
	ctx, span := f.obs.StartSpan(ctx, "FcmRepository.DeleteFcmTokenByUserIDAndToken")
	defer span()
	logger := f.obs.Logger.WithContext(ctx)
	query := `
		DELETE FROM fcm_token
		WHERE user_id = $1 AND token = $2
	`
	_, err := f.db.Exec(ctx, query, userID, token)
	if err != nil {
		logger.Error("failed to delete fcm token by user id and token", err)
		return err
	}
	return nil
}

func NewFcmRepository(db *pgxpool.Pool, obs *observability.Observability) domain.FcmTokenRepository {
	return &fcmRepository{db: db, obs: obs}
}

// CreateFcmToken implements domain.FcmTokenRepository.
func (f *fcmRepository) CreateFcmToken(ctx context.Context, fcmToken *domain.FcmToken) error {
	ctx, span := f.obs.StartSpan(ctx, "FcmRepository.CreateFcmToken")
	defer span()
	logger := f.obs.Logger.WithContext(ctx)
	query := `
		INSERT INTO fcm_token (id,user_id, token,created_at,updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := f.db.Exec(ctx, query, fcmToken.ID, fcmToken.UserID, fcmToken.Token, fcmToken.CreatedAt, fcmToken.UpdatedAt)
	if err != nil {
		logger.Error("failed to create fcm token", err)
		return err
	}
	return nil
}

// DeleteFcmToken implements domain.FcmTokenRepository.
func (f *fcmRepository) DeleteFcmToken(ctx context.Context, id string) error {
	ctx, span := f.obs.StartSpan(ctx, "FcmRepository.DeleteFcmToken")
	defer span()
	logger := f.obs.Logger.WithContext(ctx)
	query := `
		DELETE FROM fcm_token
		WHERE id = $1
	`
	_, err := f.db.Exec(ctx, query, id)
	if err != nil {
		logger.Error("failed to delete fcm token", err)
		return err
	}
	return nil
}

// GetFcmTokenByUserID implements domain.FcmTokenRepository.
func (f *fcmRepository) GetFcmTokenByUserID(ctx context.Context, userID string) ([]*domain.FcmToken, error) {
	ctx, span := f.obs.StartSpan(ctx, "FcmRepository.GetFcmTokenByUserID")
	defer span()
	logger := f.obs.Logger.WithContext(ctx)
	query := `
		SELECT id, user_id, token, created_at, updated_at
		FROM fcm_token
		WHERE user_id = $1
	`
	rows, err := f.db.Query(ctx, query, userID)
	if err != nil {
		logger.Error("failed to get fcm token by user id", err)
		return nil, err
	}
	defer rows.Close()

	var fcmTokens []*domain.FcmToken
	for rows.Next() {
		var fcmToken domain.FcmToken
		err := rows.Scan(&fcmToken.ID, &fcmToken.UserID, &fcmToken.Token, &fcmToken.CreatedAt, &fcmToken.UpdatedAt)
		if err != nil {
			return nil, err
		}
		fcmTokens = append(fcmTokens, &fcmToken)
	}
	return fcmTokens, nil
}

var _ domain.FcmTokenRepository = &fcmRepository{}
