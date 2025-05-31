package nats

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

func Connect(address string, opts ...nats.Option) *nats.Conn {
	nc, err := nats.Connect(address, opts...)
	if err != nil {
		log.Fatal(err)
	}
	return nc
}

func ConnectJetStream(address string, opts ...nats.Option) nats.JetStreamContext {
	nc, err := nats.Connect(address, opts...)
	if err != nil {
		log.Fatal(err)
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	return js
}

// func Publish[T any](js nats.JetStreamContext, subject string, data T) error {
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = js.Publish(subject, jsonData)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// type HandlerFunc[T any] func(ctx context.Context, data T) error

// func Subscribe[T any](ctx context.Context, js nats.JetStreamContext, subject string, handler HandlerFunc[T]) error {
// 	_, err := js.Subscribe(subject, func(msg *nats.Msg) {
// 		var data T
// 		err := json.Unmarshal(msg.Data, &data)
// 		if err != nil {
// 			msg.Nak()
// 			return
// 		}
// 		err = handler(ctx, data)
// 		if err != nil {
// 			msg.Nak()
// 			return
// 		}
// 		msg.Ack()
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func WrapHandler[T any](handler func(ctx context.Context, data T) error) func(ctx context.Context, data interface{}) error {
	return func(ctx context.Context, data interface{}) error {
		var typedData T
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, &typedData)
		if err != nil {
			return err
		}
		return handler(ctx, typedData)
	}
}
