package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/milgradesec/nbot/bot/modules/league"
	db "github.com/milgradesec/nbot/database"
	"github.com/milgradesec/nbot/storage"
	"github.com/rs/zerolog/log"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Version string
	s       *discordgo.Session
}

func NewBot(token string, version string) (*Bot, error) {
	bot := &Bot{
		Version: version,
	}

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating discord session: %w", err)
	}
	bot.s = s

	err = db.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	storage.Client, err = storage.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 client: %w", err)
	}

	league.Client, err = league.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create riot api client: %w", err)
	}

	return bot, nil
}

func (bot *Bot) Run() {
	bot.registerCommandHandlers()
	bot.s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	if err := bot.s.Open(); err != nil {
		log.Fatal().Err(err).Msgf("Cannot establish a connection to Discord")
	}
	defer bot.s.Close()

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := bot.s.ApplicationCommandCreate(bot.s.State.User.ID, "", v)
		if err != nil {
			log.Error().Err(err).Msgf("Cannot create '%v' command", v.Name)
		}

		registeredCommands[i] = cmd
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	for _, v := range registeredCommands {
		err := bot.s.ApplicationCommandDelete(bot.s.State.User.ID, "", v.ID)
		if err != nil {
			log.Error().Err(err).Msgf("Cannot delete '%v' command", v.Name)
		}
	}
}
