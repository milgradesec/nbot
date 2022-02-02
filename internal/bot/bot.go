package bot

import (
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	httpc "github.com/milgradesec/go-libs/http"
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog/log"
	"github.com/yuhanfang/riot/apiclient"

	"github.com/bwmarrin/discordgo"
	"github.com/milgradesec/nbot/internal/db"
	"github.com/milgradesec/nbot/internal/storage"
)

type Bot struct {
	Version string
	Token   string

	commands map[string]commandHandler
	dbpool   *pgxpool.Pool
	s3       *minio.Client
	client   *http.Client
	riotapi  apiclient.Client
}

func (bot *Bot) Run() {
	rand.Seed(time.Now().Unix())

	session, err := discordgo.New("Bot " + bot.Token)
	if err != nil {
		log.Fatal().Err(err).Msgf("error creating discord session")
	}
	session.Client = httpc.NewHTTPClient()

	bot.registerCommands()
	session.AddHandler(bot.commandDispatcher)

	dbpool, err := db.OpenDB()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to connect to db")
	}
	bot.dbpool = dbpool

	s3client, err := storage.NewS3Client()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create s3 client")
	}
	bot.s3 = s3client

	bot.client = httpc.NewHTTPClient()

	riotapi, err := newRiotAPIClient()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create riot api client")
	}
	bot.riotapi = riotapi

	err = session.Open()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to open websocket connection to discord")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	session.Close()
}

const superUser = "MILGRADESEC"
