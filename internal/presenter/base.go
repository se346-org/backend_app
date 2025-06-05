package presenter

// EmptyResponse represents an empty response.
type EmptyResponse struct{}

// BaseEmptyResponse is a concrete type for an empty base response.
type BaseEmptyResponse struct {
	Message string        `json:"message"`
	Data    EmptyResponse `json:"data,omitempty"`
}

// BaseFriendListResponse is a concrete type for a list of friends.
type BaseFriendListResponse struct {
	Message string           `json:"message"`
	Data    []FriendResponse `json:"data,omitempty"`
}

// BaseFriendResponse is a concrete type for a single friend response.
type BaseFriendResponse struct {
	Message string         `json:"message"`
	Data    FriendResponse `json:"data,omitempty"`
} 