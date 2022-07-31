package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/milgradesec/nbot/bot/modules/minitas"
	"github.com/milgradesec/nbot/bot/modules/quotes"
	"github.com/milgradesec/nbot/s3"
	"github.com/rs/zerolog/log"
)

func (bot *Bot) versionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: bot.Version,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to return response for command 'version'")
	}
}

func (bot *Bot) pingHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "PONG!",
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to return response for command 'ping'")
	}
}

func (bot *Bot) gafasHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	key := "img/congafas.webp"
	congafas, err := s3.PresignedURL(key)
	if err != nil {
		log.Error().Err(err).Msgf("failed to generate presigned url for '%s'", key)
		return
	}

	key = "img/singafas.webp"
	singafas, err := s3.PresignedURL(key)
	if err != nil {
		log.Error().Err(err).Msgf("failed to generate presigned url for '%s'", key)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "¿Con Gafas?",
					Image: &discordgo.MessageEmbedImage{
						URL: congafas,
					},
				},
				{
					Title: "¿O mejor sin gafas?",
					Image: &discordgo.MessageEmbedImage{
						URL: singafas,
					},
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to return response for command 'gafas'")
	}
}

func (bot *Bot) puteroHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	key := "clips/putero.mp4"
	url, err := s3.PresignedURL(key)
	if err != nil {
		log.Error().Err(err).Msgf("failed to generate presigned url for '%s'", key)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: url,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to return response for command 'putero'")
	}
}

func (bot *Bot) fraseHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: quotes.GetRandom(),
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to return response for command 'frase'")
	}
}

func (bot *Bot) minitaHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	url, err := minitas.GetRandom()
	if err != nil {
		log.Error().Err(err).Msg("failed to get a random minita")
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Image: &discordgo.MessageEmbedImage{
						URL: url,
					},
				},
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to return response for command 'minita'")
	}
}
