package presenter

import (
	"fmt"
	"regexp"
	"time"
)

var (
	emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regexEmail = regexp.MustCompile(emailRegex)
)

type RegisterRequest struct {
	Email    string `json:"email,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Password string `json:"password,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

func (r RegisterRequest) Validate() error {
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	if r.FullName == "" {
		return fmt.Errorf("full_name is required")
	}
	if r.Password == "" {
		return fmt.Errorf("password is required")
	}

	// validate email format
	if !isValidEmail(r.Email) {
		return fmt.Errorf("invalid email format")
	}

	// validate password length
	if len(r.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	return nil
}

func isValidEmail(email string) bool {
	return regexEmail.MatchString(email)
}

type RegisterResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (r LoginRequest) Validate() error {
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	if r.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

type LoginResponse struct {
	AccessToken string `json:"access_token,omitempty"`
}

type GetUserInfoResponse struct {
	AccountID      string     `json:"account_id,omitempty"`
	UserID         string     `json:"user_id,omitempty"`
	Type           string     `json:"type,omitempty"`
	Email          string     `json:"email,omitempty"`
	FullName       string     `json:"full_name,omitempty"`
	Avatar         string     `json:"avatar,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	ConversationID *string    `json:"conversation_id,omitempty"`
}

type UserResponse struct {
	UserID    string     `json:"user_id,omitempty"`
	FullName  string     `json:"full_name,omitempty"`
	Avatar    string     `json:"avatar,omitempty"`
	UserType  string     `json:"user_type,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
