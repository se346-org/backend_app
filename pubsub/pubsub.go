package pubsub

import (
	"context"
)

// type HandleFunc func(ctx context.Context, data interface{}) error

type Publisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, subject string, handler func(ctx context.Context, data interface{}) error) error
	Unsubscribe() error
}

type QueueSubscriber interface {
	Subscribe(ctx context.Context, subject string, handler func(ctx context.Context, data interface{}) error) error
	Unsubscribe() error
}
