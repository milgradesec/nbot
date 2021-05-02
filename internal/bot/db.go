package bot

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq" // psql driver
	log "github.com/sirupsen/logrus"
)

func (bot *Bot) openDB() error {
	var (
		dbHost     string
		dbName     string
		dbUser     string
		dbPassword string
		found      bool
	)

	dbHost, found = os.LookupEnv("POSTGRES_HOST")
	if !found {
		return errors.New("POSTGRES_HOST env variable not set")
	}

	dbName, found = os.LookupEnv("POSTGRES_DB")
	if !found {
		return errors.New("POSTGRES_DB env variable not set")
	}

	dbUser, found = os.LookupEnv("POSTGRES_USER")
	if !found {
		return errors.New("POSTGRES_USER env variable not set")
	}

	dbPassFile, found := os.LookupEnv("POSTGRES_DB_PASSWORD_FILE")
	if found {
		buf, err := ioutil.ReadFile(dbPassFile)
		if err != nil {
			return err
		}
		dbPassword = string(buf)
	} else {
		dbPassword, found = os.LookupEnv("POSTGRES_DB_PASSWORD")
		if !found {
			return errors.New("POSTGRES_DB_PASSWORD env variable not set")
		}
		log.Warnln("Using unencrypted DB password from Env, consider switching to POSTGRES_DB_PASSWORD_FILE")
	}

	connStr := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + "/" + dbName + "?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(3)
	db.SetMaxIdleConns(2)

	err = db.Ping()
	if err != nil {
		return err
	}

	bot.db = db
	return nil
}
