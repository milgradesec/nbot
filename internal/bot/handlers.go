package bot

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (bot *Bot) gafasHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	url, err := bot.generatePresignedURL("img/congafas.webp")
	if err != nil {
		log.Errorf("error: failed to generate presigned url: %v\n", err)
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
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
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Sin Gafas",

		Image: &discordgo.MessageEmbedImage{
			URL:    url,
			Width:  500,
			Height: 500,
		},
	})

	s.ChannelMessageSend(m.ChannelID, "Con o sin gafas?")
}

func (bot *Bot) ptHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	url, err := bot.generatePresignedURL("clips/putero.mp4")
	if err != nil {
		log.Errorf("error: failed to generate presigned url: %v\n", err)
		return
	}
	s.ChannelMessageSend(m.ChannelID, url)
}
