package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/milgradesec/nbot/internal/bot"
	"github.com/milgradesec/nbot/internal/config"
)

func main() {
	log.Infoln("NBOT-" + Version)

	token, found := config.GetToken()
	if !found {
		log.Fatal("error: Bot token not found")
	}

	bot := &bot.Bot{
		Token: token,
	}
	bot.Run()
}

var (
	Version string
)
