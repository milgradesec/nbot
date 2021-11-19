package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lus/dgc"
	log "github.com/sirupsen/logrus"
)

func (bot *Bot) versionHandler(ctx *dgc.Ctx) {
	ctx.RespondText(bot.Version)
}

func (bot *Bot) gafasHandler(ctx *dgc.Ctx) {
	ctx.RespondText("Con o sin 👓? 👀👀")

	url, err := bot.generatePresignedURL("img/congafas.webp")
	if err != nil {
		log.Errorf("error: failed to generate presigned url: %v\n", err)
		return
	}
	ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title: "Con Gafas",

		Image: &discordgo.MessageEmbedImage{
			URL:    url,
			Width:  500,
			Height: 500,
		},
	})

	url, err = bot.generatePresignedURL("img/singafas.webp")
	if err != nil {
		log.Errorf("error: failed to generate presigned url: %v\n", err)
		return
	}
	ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title: "Sin Gafas",

		Image: &discordgo.MessageEmbedImage{
			URL:    url,
			Width:  500,
			Height: 500,
		},
	})
}

func (bot *Bot) ptHandler(ctx *dgc.Ctx) {
	url, err := bot.generatePresignedURL("clips/putero.mp4")
	if err != nil {
		log.Errorf("error: failed to generate presigned url: %v\n", err)
		return
	}
	ctx.RespondText(url)
}
