package domain

import "time"

type Account struct {
	ID        string     `json:"id,omitempty"`
	Username  string     `json:"username,omitempty"`
	Password  string     `json:"-"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (a *Account) TableName() string {
	return "account"
}

func (a *Account) MapFields() ([]string, []any) {
	return []string{
			"id",
			"username",
			"password",
			"created_at",
			"updated_at",
		}, []any{
			&a.ID,
			&a.Username,
			&a.Password,
			&a.CreatedAt,
			&a.UpdatedAt,
		}
}
