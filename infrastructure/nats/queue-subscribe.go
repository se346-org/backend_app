package nats

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type QueueSubscriber struct {
	js           nats.JetStreamContext
	subs         *nats.Subscription
	queueName    string
	consumerName string
}

func (qs *QueueSubscriber) Subscribe(ctx context.Context, subject string, handler func(ctx context.Context, data interface{}) error) error {
	var err error
	qs.subs, err = qs.js.QueueSubscribe(subject, qs.queueName, func(msg *nats.Msg) {
		// fmt.Println("msg", string(msg.Data))
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
	}, nats.Durable(qs.consumerName), nats.ManualAck())
	if err != nil {
		return err
	}
	return nil
}

func (qs *QueueSubscriber) Unsubscribe() error {
	return qs.subs.Unsubscribe()
}

func NewQueueSubscriber(js nats.JetStreamContext, queueName string, consumerName string) *QueueSubscriber {
	return &QueueSubscriber{js: js, queueName: queueName, consumerName: consumerName}
}
