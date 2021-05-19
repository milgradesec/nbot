package main

import (
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/milgradesec/nbot/internal/bot"
	"github.com/milgradesec/nbot/internal/config"
)

var (
	Version = "DEV"
)

func main() {
	log.Infof("Nbot %s", Version)
	log.Infof("%s/%s %s", runtime.GOOS, runtime.GOARCH, runtime.Version())

	token, found := config.GetToken()
	if !found {
		log.Fatal("error: Discord token not found")
	}

	bot := &bot.Bot{
		Version: Version,
		Token:   token,
	}
	bot.Run()
}
