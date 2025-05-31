package observability

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

type contextKey string

const (
	CorrelationIDKey contextKey = "correlation_id"
	UserIDKey        contextKey = "user_id"
)

func NewLogger() *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	return &Logger{Logger: log}
}

func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	entry := l.WithFields(logrus.Fields{})
	
	if correlationID := ctx.Value(CorrelationIDKey); correlationID != nil {
		entry = entry.WithField("correlation_id", correlationID)
	}
	
	if userID := ctx.Value(UserIDKey); userID != nil {
		entry = entry.WithField("user_id", userID)
	}
	
	return entry
}

func (l *Logger) WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

func (l *Logger) WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
} 