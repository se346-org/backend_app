package postgresql

import (
	"context"
	"fmt"

	"github.com/chat-socio/backend/configuration"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, postgresConfig *configuration.PostgresConfig) (*pgxpool.Pool, error) {
	// Build the connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.Username,
		postgresConfig.Password,
		postgresConfig.Database,
		postgresConfig.SSLMode,
	)
	// Open a connection to the database
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
