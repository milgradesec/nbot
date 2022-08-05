package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/viper"
)

var (
	Conn *pgxpool.Pool
)

func Open() error {
	var (
		host     string
		database string
		username string
		password string
		sslMode  string
	)

	if !viper.IsSet("POSTGRES_HOST") {
		return errors.New("NBOT_POSTGRES_HOST not set")
	}
	host = viper.GetString("POSTGRES_HOST")

	if !viper.IsSet("POSTGRES_DB") {
		return errors.New("NBOT_POSTGRES_DB not set")
	}
	database = viper.GetString("POSTGRES_DB")

	if !viper.IsSet("POSTGRES_USER") {
		return errors.New("NBOT_POSTGRES_USER not set")
	}
	username = viper.GetString("POSTGRES_USER")

	if !viper.IsSet("POSTGRES_DB_PASSWORD") {
		return errors.New("NBOT_POSTGRES_DB_PASSWORD not set")
	}
	password = viper.GetString("POSTGRES_DB_PASSWORD")

	if !viper.IsSet("POSTGRES_SSL_MODE") {
		sslMode = "require"
	} else {
		sslMode = viper.GetString("POSTGRES_SSL_MODE")
	}

	connStr := "postgres://" + username + ":" + password + "@" + host + "/" + database + "?sslmode=" + sslMode
	dbpool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return err
	}
	Conn = dbpool

	err = dbpool.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}
