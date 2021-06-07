package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lus/dgc"
)

func (bot *Bot) versionHandler(ctx *dgc.Ctx) {
	ctx.RespondText(bot.Version) //nolint
}

func (bot *Bot) gafasHandler(ctx *dgc.Ctx) {
	ctx.RespondEmbed(&discordgo.MessageEmbed{ //nolint
		Title: "Con Gafas",

		Image: &discordgo.MessageEmbedImage{
			URL:    "https://s3.paesa.es/nbot-data/img/congafas.png",
			Width:  400,
			Height: 400,
		},
	})

	ctx.RespondEmbed(&discordgo.MessageEmbed{ //nolint
		Title: "Sin Gafas",

		Image: &discordgo.MessageEmbedImage{
			URL:    "https://s3.paesa.es/nbot-data/img/singafas.png",
			Width:  400,
			Height: 400,
		},
	})
}

func (bot *Bot) ptHandler(ctx *dgc.Ctx) {
	ctx.RespondText("https://s3.paesa.es/nbot-data/clips/putero.mp4") //nolint
}
