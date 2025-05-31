package nats

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type NatsSubscriber struct {
	js           nats.JetStreamContext
	subs         *nats.Subscription
	consumerName string
}

func (ns *NatsSubscriber) Subscribe(ctx context.Context, subject string, handler func(ctx context.Context, data interface{}) error) error {
	var err error
	ns.subs, err = ns.js.Subscribe(subject, func(msg *nats.Msg) {
		var data interface{}
		err := json.Unmarshal(msg.Data, &data)
		if err != nil {
			msg.Nak()
			return
		}
		err = handler(ctx, data)
		if err != nil {
			msg.Nak()
			return
		}
		msg.Ack()
	}, nats.Durable(ns.consumerName), nats.ManualAck())
	if err != nil {
		return err
	}
	return nil
}

func (ns *NatsSubscriber) Unsubscribe() error {
	return ns.subs.Unsubscribe()
}

func NewSubscriber(js nats.JetStreamContext, consumerName string) *NatsSubscriber {
	return &NatsSubscriber{js: js, consumerName: consumerName}
}
