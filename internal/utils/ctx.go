package utils

type contextKey string

const (
	UserAgentKey    contextKey = "user_agent"
	IpAddressKey    contextKey = "ip_address"
	AccountIDKey    contextKey = "account_id"
	SessionTokenKey contextKey = "session_token"
)
