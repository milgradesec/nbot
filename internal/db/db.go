package db

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
)

func OpenDB() (*sql.DB, error) {
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
		log.Warnln("Using unencrypted DB password from env, consider switching to POSTGRES_DB_PASSWORD_FILE")
	}

	rootCA, found = os.LookupEnv("POSTGRES_SSL_ROOT_CERT")
	if !found {
		return nil, errors.New("POSTGRES_SSL_ROOT_CERT env variable not set")
	}

	sslMode, found = os.LookupEnv("POSTGRES_SSL_MODE")
	if !found {
		sslMode = "verify-full"
	}

	connStr := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + "/" + dbName + "?sslmode=" + sslMode + "&sslrootcert=" + rootCA //nolint
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(2)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
