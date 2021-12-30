package main

import (
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/milgradesec/nbot/internal/bot"
	"github.com/milgradesec/nbot/internal/config"
)

var (
	Version = "DEV"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	log.Info().Msgf("Nbot %s", Version)
	log.Info().Msgf("%s/%s %s", runtime.GOOS, runtime.GOARCH, runtime.Version())

	token, found := config.GetToken()
	if !found {
		log.Fatal().Msgf("error: Discord token not found")
	}

	bot := &bot.Bot{
		Version: Version,
		Token:   token,
	}
	bot.Run()
}
