package bot

import (
	"github.com/bwmarrin/discordgo"
)

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
			// 		Name:        "add",
			// 		Description: "Añadir una frase",
			// 		Type:        discordgo.ApplicationCommandOptionSubCommand,
			// 		Options: []*discordgo.ApplicationCommandOption{

			// 			{
			// 				Name:        "frase",
			// 				Description: "La frase en cuestion",
			// 				Type:        discordgo.ApplicationCommandOptionString,
			// 				Required:    true,
			// 			},
			// 		},
			// 		Required: false,
			// 	},
			// 	{
			// 		Name:        "del",
			// 		Description: "Eliminar una frase",
			// 		Type:        discordgo.ApplicationCommandOptionSubCommand,
			// 		Required:    false,
			// 	},
			// },
		},
		{
			Name:        "minita",
			Description: "Te mando la foto de una minita",
		},
		{
			Name:        "lol",
			Description: "Muestra información de las rankeds de un usuario",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "user",
					Description: "Nombre del usuario en League of Leguends",
					Required:    true,
				},
			},
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
	commandHandlers["lol"] = bot.lolHandler
}
