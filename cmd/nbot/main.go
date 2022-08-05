package main

import (
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/milgradesec/nbot/bot"
)

var (
	Version = "DEV"
)

func main() {
	log.Info().Msgf("Nbot %s -- %s/%s %s", Version, runtime.GOOS, runtime.GOARCH, runtime.Version())

	if !viper.IsSet("DISCORD_TOKEN") {
		log.Fatal().Msg("NBOT_DISCORD_TOKEN not set")
	}
	token := viper.GetString("DISCORD_TOKEN")

	bot, err := bot.NewBot(token, Version)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create bot")
	}
	bot.Run()
}
