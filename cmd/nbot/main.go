package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/milgradesec/nbot/internal/bot"
	"github.com/milgradesec/nbot/internal/config"
)

var (
	Version = "DEV"
)

func main() {
	log.Infoln("Nbot is running.")

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
