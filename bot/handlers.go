package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/milgradesec/nbot/s3"
	"github.com/rs/zerolog/log"
)

func (bot *Bot) gafasHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	url, err := s3.GeneratePresignedURL("img/congafas.webp")
	if err != nil {
		log.Error().Err(err).Msgf("failed to generate presigned url")
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

	url, err = s3.GeneratePresignedURL("img/singafas.webp")
	if err != nil {
		log.Error().Err(err).Msgf("failed to generate presigned url")
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
	url, err := s3.GeneratePresignedURL("clips/putero.mp4")
	if err != nil {
		log.Error().Err(err).Msgf("failed to generate presigned url")
		return
	}
	s.ChannelMessageSend(m.ChannelID, url)
}
