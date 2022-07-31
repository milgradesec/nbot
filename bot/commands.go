package bot

import (
	"github.com/bwmarrin/discordgo"
)

type commandHandler func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)

/*func (bot *Bot) registerCommands() {
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
	bot.cmd = commands
}

func (bot *Bot) commandDispatcher(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Message.Content, "!") {
		args := strings.Split(m.Message.Content, " ")

		handler, found := bot.cmd[args[0]]
		if found {
			handler(s, m, args)
		}
	}
}*/

///////////////////////////////////////////////////////////////////////////////////////////////

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "version",
			Description: "Muestra la version de Nbot",
		},
		{
			Name:        "ping",
			Description: "Responde PONG!",
		},
		{
			Name:        "gafas",
			Description: "¿Con o sin gafas?",
		},
		{
			Name:        "putero",
			Description: "Eres un putero y lo sabes",
		},
		{
			Name:        "frase",
			Description: "Una frase aleatoria de L.L",
			// Options: []*discordgo.ApplicationCommandOption{
			// 	{
			// 		Type:        discordgo.ApplicationCommandOptionString,
			// 		Name:        "add",
			// 		Description: "Añade una nueva frase",
			// 	},
			// },
		},
		{
			Name:        "minita",
			Description: "Te mando la foto de una minita",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
)

func (bot *Bot) registerCommandHandlers() {
	commandHandlers["version"] = bot.versionHandler
	commandHandlers["ping"] = bot.pingHandler
	commandHandlers["gafas"] = bot.gafasHandler
	commandHandlers["putero"] = bot.puteroHandler
	commandHandlers["frase"] = bot.fraseHandler
	commandHandlers["minita"] = bot.minitaHandler
}
