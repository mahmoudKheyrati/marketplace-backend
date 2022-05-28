package pkg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
)

type PostgresConfig struct {
	Host          string
	Port          int
	Username      string
	Password      string
	Database      string
	MaxConnection int
}

func CreateNewPostgresConnection(config PostgresConfig) *pgx.Conn {
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}
