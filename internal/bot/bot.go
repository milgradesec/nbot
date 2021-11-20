package bot

import (
	"database/sql"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpc "github.com/milgradesec/go-libs/http"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
	"github.com/yuhanfang/riot/apiclient"

	"github.com/bwmarrin/discordgo"
	"github.com/milgradesec/nbot/internal/db"
	"github.com/milgradesec/nbot/internal/storage"
)

type Bot struct {
	Version string
	Token   string

	commands map[string]commandHandler
	db       *sql.DB
	s3       *minio.Client
	client   *http.Client
	riotapi  apiclient.Client
}

func (bot *Bot) Run() {
	rand.Seed(time.Now().Unix())

	session, err := discordgo.New("Bot " + bot.Token)
	if err != nil {
		log.Fatalf("error: error creating discord session: %v\n", err)
	}
	session.Client = httpc.NewHTTPClient()

	bot.registerCommands()
	session.AddHandler(bot.commandDispatcher)

	db, err := db.OpenDB()
	if err != nil {
		log.Fatalf("error: failed to connect to db: %v\n", err)
	}
	bot.db = db

	s3client, err := storage.NewS3Client()
	if err != nil {
		log.Fatalf("error: failed to create s3 client: %v\n", err)
	}
	bot.s3 = s3client

	bot.client = httpc.NewHTTPClient()

	riotapi, err := newRiotAPIClient()
	if err != nil {
		log.Fatalf("error: failed to create riot api client: %v\n", err)
	}
	bot.riotapi = riotapi

	err = session.Open()
	if err != nil {
		log.Fatalf("error: failed to open websocket connection to discord: %v\n", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	session.Close()
}

const superUser = "MILGRADESEC"
