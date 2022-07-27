package config

import (
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
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
		log.Warn().Msg("Using unencrypted Token from env, consider switching to DISCORD_BOT_TOKEN_FILE")
		return token, true
	}
	return "", false
}
