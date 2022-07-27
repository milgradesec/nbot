package bot

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/milgradesec/nbot/db"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"github.com/yuhanfang/riot/apiclient"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Version string

	s       *discordgo.Session
	cmd     map[string]commandHandler
	s3      *minio.Client
	client  *http.Client
	riotapi apiclient.Client
}

func NewBot(token string, version string) (*Bot, error) {
	bot := &Bot{
		Version: version,
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("error creating discord session: %w", err)
	}
	bot.s = session

	bot.registerCommands()
	session.AddHandler(bot.commandDispatcher)

	err = db.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// bot.s3, err = storage.NewS3Client()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create s3 client: %w", err)
	// }

	// bot.client = httpc.NewHTTPClient()

	// bot.riotapi, err = newRiotAPIClient()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create riot api client: %w", err)
	// }

	return bot, nil
}

func (bot *Bot) Run() {
	rand.Seed(time.Now().Unix())

	if err := bot.s.Open(); err != nil {
		log.Fatal().Err(err).Msgf("failed to create a websocket connection to Discord")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	bot.s.Close()
}

const superUser = "MILGRADESEC"
