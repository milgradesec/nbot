package config

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

func GetToken() (string, bool) {
	file, found := os.LookupEnv("DISCORD_BOT_TOKEN_FILE")
	if found {
		token, err := ioutil.ReadFile(file)
		if err != nil {
			return "", false
		}
		return string(token), true
	}

	token, found := os.LookupEnv("DISCORD_BOT_TOKEN")
	if found {
		log.Warnln("Using unencrypted Token from Env, consider switching to DISCORD_BOT_TOKEN_FILE")
		return token, true
	}
	return "", false
}
