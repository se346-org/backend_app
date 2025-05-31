package postgresql

import (
	"context"
	"time"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type sessionRepository struct {
	db *pgxpool.Pool
}

// CreateSession implements domain.SessionRepository.
func (s *sessionRepository) CreateSession(ctx context.Context, session *domain.Session) error {
	query := `INSERT INTO session (session_token, account_id, created_at, updated_at, expired_at, is_active, user_agent, ip_address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := s.db.Exec(ctx, query, session.SessionToken, session.AccountID, session.CreatedAt, session.UpdatedAt, session.ExpiredAt, session.IsActive, session.UserAgent, session.IPAddress)
	if err != nil {
		return err
	}

	return nil
}

// DeactivateSession implements domain.SessionRepository.
func (s *sessionRepository) DeactivateSession(ctx context.Context, token string) error {
	query := `UPDATE session SET is_active = false WHERE session_token = $1`
	_, err := s.db.Exec(ctx, query, token)
	if err != nil {
		return err
	}

	return nil
}

// DeactiveAllSessionByAccountID implements domain.SessionRepository.
func (s *sessionRepository) DeactiveAllSessionByAccountID(ctx context.Context, accountID string) error {
	query := `UPDATE session SET is_active = false WHERE account_id = $1`
	_, err := s.db.Exec(ctx, query, accountID)
	if err != nil {
		return err
	}

	return nil
}

// GetListSessionByAccountID implements domain.SessionRepository.
func (s *sessionRepository) GetListSessionByAccountID(ctx context.Context, accountID string) ([]*domain.Session, error) {
	var sessions []*domain.Session
	query := `SELECT session_token, account_id, created_at, updated_at, expired_at, is_active, user_agent, ip_address FROM session WHERE account_id = $1 AND is_active = true`
	rows, err := s.db.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var session domain.Session
		err = rows.Scan(&session.SessionToken, &session.AccountID, &session.CreatedAt, &session.UpdatedAt, &session.ExpiredAt, &session.IsActive, &session.UserAgent, &session.IPAddress)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// GetSessionByToken implements domain.SessionRepository.
func (s *sessionRepository) GetSessionByToken(ctx context.Context, token string) (*domain.Session, error) {
	query := `SELECT session_token, account_id, created_at, updated_at, expired_at, is_active, user_agent, ip_address FROM session WHERE session_token = $1`
	row := s.db.QueryRow(ctx, query, token)
	var session domain.Session
	err := row.Scan(&session.SessionToken, &session.AccountID, &session.CreatedAt, &session.UpdatedAt, &session.ExpiredAt, &session.IsActive, &session.UserAgent, &session.IPAddress)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &session, nil
}

// UpdateExpiredAt implements domain.SessionRepository.
func (s *sessionRepository) UpdateExpiredAt(ctx context.Context, token string, newExpiredAt *time.Time) error {
	query := `UPDATE session SET expired_at = $1 WHERE session_token = $2`
	_, err := s.db.Exec(ctx, query, newExpiredAt, token)
	if err != nil {
		return err
	}

	return nil
}

func NewSessionRepository(db *pgxpool.Pool) *sessionRepository {
	return &sessionRepository{
		db: db,
	}
}

var _ domain.SessionRepository = (*sessionRepository)(nil)
