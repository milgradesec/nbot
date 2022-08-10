package bot

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/milgradesec/nbot/bot/modules/league"
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

func (bot *Bot) messageErrorRespond(i *discordgo.InteractionCreate, err error) {
	bot.messageRespond(i, fmt.Sprintf("❌ Se ha producido un error ==> %v", err))
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
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	if _, ok := optionMap["add"]; ok {
		bot.messageRespond(i, "Este comando no esta implementado.")
		return
	}

	if _, ok := optionMap["del"]; ok {
		bot.messageRespond(i, "Este comando no esta implementado.")
		return
	}

	bot.messageRespond(i, quotes.GetRandom(context.TODO()))
}

func (bot *Bot) minitaHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	url, err := minitas.GetRandom()
	if err != nil {
		log.Error().Err(err).Msg("failed to retrieve a random minita")
		bot.messageErrorRespond(i, err)
		return
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

func (bot *Bot) lolHandler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	if opt, ok := optionMap["user"]; ok {
		msg, err := league.GetRankedSummary(context.TODO(), opt.StringValue())
		if err != nil {
			bot.messageErrorRespond(i, err)
			return
		}
		bot.messageRespond(i, msg)
	}
}
