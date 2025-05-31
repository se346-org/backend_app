package postgresql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/chat-socio/backend/pkg/observability"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db  *pgxpool.Pool
	obs *observability.Observability
}

// GetListUser implements domain.UserRepository.
func (u *userRepository) GetListUser(ctx context.Context, keyword string, limit int, lastID string) ([]*domain.UserInfo, error) {
	var users []*domain.UserInfo
	var user domain.UserInfo
	fields, _ := user.MapFields()
	query := fmt.Sprintf(`SELECT %s FROM %s`, strings.Join(fields, ","), user.TableName())

	args := []any{}
	conditions := []string{}

	if keyword != "" {
		conditions = append(conditions, fmt.Sprintf("(full_name ILIKE $%d OR email ILIKE $%d)", len(args)+1, len(args)+2))
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}

	if lastID != "" {
		conditions = append(conditions, fmt.Sprintf("id < $%d", len(args)+1))
		args = append(args, lastID)
	}

	conditions = append(conditions, "type = 'EXTERNAL'")

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	query = fmt.Sprintf("%s ORDER BY id DESC LIMIT $%d", query, len(args)+1)
	args = append(args, limit)

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.UserInfo
		_, values := user.MapFields()
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (u *userRepository) GetListUserWithConversation(ctx context.Context, userID string, keyword string, limit int, lastID string) ([]*domain.UserInfo, error) {

	var users []*domain.UserInfo
	var user domain.UserInfo
	fields, _ := user.MapFields()
	query := fmt.Sprintf(`
		WITH user_conversations AS (
			SELECT DISTINCT cm.user_id, cm.conversation_id 
			FROM conversation_member cm
			WHERE cm.conversation_id IN (
				SELECT conversation_id 
				FROM conversation_member 
				WHERE user_id = $1
			)
		)
		SELECT %s, uc.conversation_id
		FROM %s u 
		LEFT JOIN user_conversations uc ON u.id = uc.user_id
		WHERE u.id != $1`, strings.Join(fields, ","), user.TableName())

	args := []any{userID}
	conditions := []string{}

	if keyword != "" {
		conditions = append(conditions, fmt.Sprintf("(u.full_name ILIKE $%d OR u.email ILIKE $%d)", len(args)+1, len(args)+2))
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}

	if lastID != "" {
		conditions = append(conditions, fmt.Sprintf("u.id < $%d", len(args)+1))
		args = append(args, lastID)
	}

	conditions = append(conditions, "u.type = 'EXTERNAL'")

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s AND %s", query, strings.Join(conditions, " AND "))
	}

	query = fmt.Sprintf("%s ORDER BY u.id DESC LIMIT $%d", query, len(args)+1)
	args = append(args, limit)

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.UserInfo
		_, values := user.MapFields()
		values = append(values, &user.ConversationID)
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

// CreateUser implements domain.UserRepository.
func (u *userRepository) CreateUser(ctx context.Context, user *domain.UserInfo) error {
	ctx, span := u.obs.StartSpan(ctx, "UserRepository.Create")
	defer span()

	start := time.Now()
	logger := u.obs.Logger.WithContext(ctx)

	logger.Info("Creating user", map[string]interface{}{
		"email": user.Email,
	})

	query := `INSERT INTO user_info (id ,account_id, type, email, full_name, avatar, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := u.db.Exec(ctx, query, user.AccountID, user.Type, user.Email, user.FullName, user.Avatar, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		logger.WithError(err).Error("Failed to create user")
		return err
	}

	logger.Info("User created successfully")
	u.obs.Metrics.UsersRegistered.Inc()
	u.obs.Metrics.RecordDBQuery("INSERT", "users", "success", time.Since(start))
	return nil
}

// GetUserByAccountID implements domain.UserRepository.
func (u *userRepository) GetUserByAccountID(ctx context.Context, accountID string) (*domain.UserInfo, error) {
	var user domain.UserInfo
	fields, values := user.MapFields()
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE account_id = $1`, strings.Join(fields, ","), user.TableName())
	row := u.db.QueryRow(ctx, query, accountID)
	err := row.Scan(values...)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail implements domain.UserRepository.
func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.UserInfo, error) {
	ctx, span := u.obs.StartSpan(ctx, "UserRepository.GetByEmail")
	defer span()

	start := time.Now()
	logger := u.obs.Logger.WithContext(ctx)

	logger.Info("Getting user by email", map[string]interface{}{
		"email": email,
	})

	var user domain.UserInfo
	fields, values := user.MapFields()
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE email = $1`, strings.Join(fields, ","), user.TableName())
	row := u.db.QueryRow(ctx, query, email)
	err := row.Scan(values...)
	if err != nil {
		logger.WithError(err).Error("Failed to get user by email")
		return nil, err
	}

	logger.Info("User retrieved successfully")
	u.obs.Metrics.RecordDBQuery("SELECT", "users", "success", time.Since(start))
	return &user, nil
}

// GetUserByID implements domain.UserRepository.
func (u *userRepository) GetUserByID(ctx context.Context, id string) (*domain.UserInfo, error) {
	var user domain.UserInfo
	fields, values := user.MapFields()
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE id = $1`, strings.Join(fields, ","), user.TableName())
	row := u.db.QueryRow(ctx, query, id)
	err := row.Scan(values...)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser implements domain.UserRepository.
func (u *userRepository) UpdateUser(ctx context.Context, user *domain.UserInfo) error {
	query := `UPDATE user_info SET full_name = $1, avatar = $2, updated_at = NOW() WHERE id = $3`
	_, err := u.db.Exec(ctx, query, user.FullName, user.Avatar, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepository(db *pgxpool.Pool, obs *observability.Observability) domain.UserRepository {
	return &userRepository{
		db:  db,
		obs: obs,
	}
}

var _ domain.UserRepository = (*userRepository)(nil)
