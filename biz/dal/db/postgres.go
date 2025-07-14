package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wnnce/fserv-template/config"
)

var (
	defaultDB *pgxpool.Pool
)

// InitPostgres initializes a Postgresql connection pool using configuration values
// retrieved via Viper. It establishes a connection, verifies it via Ping, and stores
// the pool in the defaultDB global variable.
//
// It returns a cleanup function that closes the connection pool when invoked,
// or an error if the initialization fails.
func InitPostgres(ctx context.Context) (func(), error) {
	host := config.ViperGet[string]("database.host", "127.0.0.1")
	port := config.ViperGet[int]("database.port", 3456)
	username := config.ViperGet[string]("database.username")
	password := config.ViperGet[string]("database.password")
	dbName := config.ViperGet[string]("database.db")
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbName,
	)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		slog.Error("failed to create postgres pool", slog.Group("data",
			slog.String("host", host),
			slog.Int("port", port),
			slog.String("db", dbName),
		), slog.String("error", err.Error()))
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	slog.Info("Postgresql connected successfully", slog.Group("data",
		slog.String("host", host),
		slog.Int("port", port),
		slog.String("db", dbName),
	))
	defaultDB = pool
	return func() {
		defaultDB.Close()
	}, nil
}

// Postgres returns the connected Postgresql database instance.
// Panics if the database has not been initialized.
func Postgres() *pgxpool.Pool {
	return defaultDB
}
