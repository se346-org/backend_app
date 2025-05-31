package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/chat-socio/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type accountRepository struct {
	db *pgxpool.Pool
}

// CreateAccount implements domain.AccountRepository.
func (a *accountRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	query := fmt.Sprintf(`INSERT INTO %s (id,username, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`, account.TableName())
	_, err := a.db.Exec(ctx, query, account.Username, account.Password, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

// CreateAccountUser implements domain.AccountRepository.
func (a *accountRepository) CreateAccountUser(ctx context.Context, account *domain.Account, user *domain.UserInfo) error {
	tx, err := a.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	query := fmt.Sprintf(`INSERT INTO %s (id, username, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`, account.TableName())
	_, err = tx.Exec(ctx, query, account.ID, account.Username, account.Password, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`INSERT INTO %s (id,account_id, type, email, full_name, avatar, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, user.TableName())
	_, err = tx.Exec(ctx, query, user.ID, account.ID, user.Type, user.Email, user.FullName, user.Avatar, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// GetAccountByID implements domain.AccountRepository.
func (a *accountRepository) GetAccountByID(ctx context.Context, id string) (*domain.Account, error) {
	account := &domain.Account{}
	fields, values := account.MapFields()
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE id = $1`, strings.Join(fields, ","), account.TableName())
	row := a.db.QueryRow(ctx, query, id)
	err := row.Scan(values...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return account, nil
}

// GetAccountByUsername implements domain.AccountRepository.
func (a *accountRepository) GetAccountByUsername(ctx context.Context, username string) (*domain.Account, error) {
	account := &domain.Account{}
	fields, values := account.MapFields()
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE username = $1`, strings.Join(fields, ","), account.TableName())
	row := a.db.QueryRow(ctx, query, username)
	err := row.Scan(values...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return account, nil
}

// UpdatePassword implements domain.AccountRepository.
func (a *accountRepository) UpdatePassword(ctx context.Context, id string, password string) error {
	var temp domain.Account
	query := fmt.Sprintf(`UPDATE %s SET password = $1 WHERE id = $2`, temp.TableName())
	_, err := a.db.Exec(ctx, query, password, id)
	if err != nil {
		return err
	}
	return nil
}

var _ domain.AccountRepository = (*accountRepository)(nil)

// NewAccountRepository creates a new instance of AccountRepository.
func NewAccountRepository(db *pgxpool.Pool) domain.AccountRepository {
	return &accountRepository{
		db: db,
	}
}
