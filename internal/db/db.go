package db

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"
)

func OpenDB() (*pgxpool.Pool, error) { //nolint
	var (
		dbHost     string
		dbName     string
		dbUser     string
		dbPassword string
		rootCA     string
		sslMode    string
		found      bool
	)

	dbHost, found = os.LookupEnv("POSTGRES_HOST")
	if !found {
		return nil, errors.New("POSTGRES_HOST env variable not set")
	}

	dbName, found = os.LookupEnv("POSTGRES_DB")
	if !found {
		return nil, errors.New("POSTGRES_DB env variable not set")
	}

	dbUser, found = os.LookupEnv("POSTGRES_USER")
	if !found {
		return nil, errors.New("POSTGRES_USER env variable not set")
	}

	dbPassFile, found := os.LookupEnv("POSTGRES_DB_PASSWORD_FILE")
	if found {
		buf, err := os.ReadFile(dbPassFile)
		if err != nil {
			return nil, err
		}
		dbPassword = string(buf)
	} else {
		dbPassword, found = os.LookupEnv("POSTGRES_DB_PASSWORD")
		if !found {
			return nil, errors.New("POSTGRES_DB_PASSWORD env variable not set")
		}
		log.Warn().Msg("Using unencrypted DB password from env, consider switching to POSTGRES_DB_PASSWORD_FILE")
	}

	rootCA, found = os.LookupEnv("POSTGRES_SSL_ROOT_CERT")
	if !found {
		return nil, errors.New("POSTGRES_SSL_ROOT_CERT env variable not set")
	}

	sslMode, found = os.LookupEnv("POSTGRES_SSL_MODE")
	if !found {
		sslMode = "verify-full"
	}

	connStr := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + "/" + dbName + "?sslmode=" + sslMode + "&sslrootcert=" + rootCA
	dbpool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}
