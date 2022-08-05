package bot

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/milgradesec/nbot/bot/modules/minitas"
	"github.com/milgradesec/nbot/bot/modules/quotes"
	"github.com/milgradesec/nbot/storage"
	"github.com/rs/zerolog/log"
)

func (bot *Bot) messageRespond(i *discordgo.InteractionCreate, msg string) {
	bot.s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

func (bot *Bot) versionHandler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	bot.messageRespond(i, bot.Version)
}

func (bot *Bot) pingHandler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	bot.messageRespond(i, "PONG!")
}

func (bot *Bot) gafasHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	key := "img/congafas.webp"
	congafas, err := storage.PresignedGet(context.TODO(), key)
	if err != nil {
		log.Error().Err(err).Msgf("failed to generate presigned url for '%s'", key)
		return
	}

	key = "img/singafas.webp"
	singafas, err := storage.PresignedGet(context.TODO(), key)
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
	url, err := storage.PresignedGet(context.TODO(), key)
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

func (bot *Bot) fraseHandler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	if _, ok := optionMap["add"]; ok {
		// log.Info().Msgf("Recibido comando /frase add %s", option.StringValue())
		bot.messageRespond(i, "Este comando aun no esta implementado.")
		return
	}

	if _, ok := optionMap["del"]; ok {
		bot.messageRespond(i, "Este comando aun no esta implementado.")
		return
	}

	bot.messageRespond(i, quotes.GetRandom())
}

func (bot *Bot) minitaHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	url, err := minitas.GetRandom()
	if err != nil {
		log.Error().Err(err).Msg("failed to get a random minita")
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
}
