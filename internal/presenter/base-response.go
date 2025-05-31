package presenter

type BaseResponse[T any] struct {
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}
