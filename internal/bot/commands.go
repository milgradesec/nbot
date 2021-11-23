package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type commandHandler func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)

func (bot *Bot) registerCommands() {
	var commands = map[string]commandHandler{
		"!version": func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			s.ChannelMessageSend(m.ChannelID, bot.Version)
		},
		"!ping": func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			s.ChannelMessageSend(m.ChannelID, "PONG!")
		},
		"!putero": bot.ptHandler,
		"!nbot":   bot.quoteHandler,
		"!frase":  bot.quoteHandler,
		"!quote":  bot.quoteHandler,
		"!frases": bot.quotesHandler,
		"!gafas":  bot.gafasHandler,
		"!elo":    bot.eloHandler,
		"!minita": bot.minitaHandler,
	}
	bot.commands = commands
}

func (bot *Bot) commandDispatcher(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Message.Content, "!") {
		args := strings.Split(m.Message.Content, " ")

		handler, found := bot.commands[args[0]]
		if found {
			handler(s, m, args)
		}
	}
}
