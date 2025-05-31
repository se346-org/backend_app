package migrate

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/chat-socio/backend/configuration"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Migrate() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		configuration.ConfigInstance.Postgres.Host,
		configuration.ConfigInstance.Postgres.Port,
		configuration.ConfigInstance.Postgres.Username,
		configuration.ConfigInstance.Postgres.Password,
		configuration.ConfigInstance.Postgres.Database,
		configuration.ConfigInstance.Postgres.SSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Error connecting to database:", err)
		os.Exit(1)
	}

	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Println("Error creating migration driver:", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Println("Error creating new migration instance:", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Println("Error applying migrations:", err)
			os.Exit(1)
		}
	}

	fmt.Println("Migrations applied successfully")
}
