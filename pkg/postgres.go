package pkg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"runtime"
)

type PostgresConfig struct {
	Host          string
	Port          int
	Username      string
	Password      string
	Database      string
	MaxConnection int
}

func CreateNewPostgresConnection(config PostgresConfig) *pgxpool.Pool {
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	conn, err := pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		Logger().Error(err)
		os.Exit(1)
	}
	conn.Config().MaxConns = int32(runtime.NumCPU())
	return conn
}
