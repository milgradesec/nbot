package bot

import (
	"database/sql"
	"errors"
	"os"
)

func (bot *Bot) openDB() error {
	dbHost, found := os.LookupEnv("POSTGRES_HOST")
	if !found {
		return errors.New("POSTGRES_HOST env variable not set")
	}

	dbName, found := os.LookupEnv("POSTGRES_DB")
	if !found {
		return errors.New("POSTGRES_DB env variable not set")
	}

	dbUser, found := os.LookupEnv("POSTGRES_USER")
	if !found {
		return errors.New("POSTGRES_USER env variable not set")
	}

	dbPassword, found := os.LookupEnv("POSTGRES_DB_PASSWORD")
	if !found {
		return errors.New("POSTGRES_DB_PASSWORD env variable not set")
	}

	connStr := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + "/" + dbName

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(2)

	bot.db = db
	return nil
}
