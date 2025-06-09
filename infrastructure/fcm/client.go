package fcm

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/chat-socio/backend/configuration"
	"google.golang.org/api/option"
)

func NewFCMClient(ctx context.Context, config *configuration.FCMConfig) (*messaging.Client, error) {
	opts := []option.ClientOption{
		option.WithCredentialsFile(config.CredentialsFile),
	}
	firebaseApp, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: config.ProjectID,
	}, opts...)
	if err != nil {
		return nil, err
	}

	fcmClient, err := firebaseApp.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return fcmClient, nil
}
