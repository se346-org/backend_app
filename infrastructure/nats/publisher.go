package nats

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type NatsPublisher struct {
	js nats.JetStreamContext
}

func (np *NatsPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = np.js.Publish(subject, jsonData)
	return err
}

func NewPublisher(js nats.JetStreamContext) *NatsPublisher {
	return &NatsPublisher{js: js}
}
