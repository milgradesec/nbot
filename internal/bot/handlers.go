package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lus/dgc"
)

func (bot *Bot) versionHandler(ctx *dgc.Ctx) {
	ctx.RespondText(bot.Version)
}

func (bot *Bot) gafasHandler(ctx *dgc.Ctx) {
	ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title: "Con Gafas",

		Image: &discordgo.MessageEmbedImage{
			URL:    "https://s3.paesa.es/nbot/img/congafas.png",
			Width:  400,
			Height: 400,
		},
	})

	ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title: "Sin Gafas",

		Image: &discordgo.MessageEmbedImage{
			URL:    "https://s3.paesa.es/nbot/img/singafas.png",
			Width:  400,
			Height: 400,
		},
	})
}

func (bot *Bot) ptHandler(ctx *dgc.Ctx) {
	ctx.RespondText("https://s3.paesa.es/nbot/clips/putero.mp4")
}
