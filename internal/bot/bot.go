package bot

import (
	"database/sql"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lus/dgc"
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

	db      *sql.DB
	s3      *minio.Client
	client  *http.Client
	riotapi apiclient.Client
}

func (bot *Bot) Run() { //nolint
	rand.Seed(time.Now().Unix())

	session, err := discordgo.New("Bot " + bot.Token)
	if err != nil {
		log.Fatalf("error creating discord session: %v", err)
	}
	session.Client = httpc.NewHTTPClient()

	router := dgc.Create(&dgc.Router{
		Prefixes:         []string{"!"},
		IgnorePrefixCase: true,
		BotsAllowed:      false,
		PingHandler: func(ctx *dgc.Ctx) {
			ctx.RespondText("PONG!") //nolint
		},
	})

	router.RegisterCmd(&dgc.Command{
		Name:       "nbot",
		IgnoreCase: true,
		Handler:    bot.fraseHandler,
	})
	router.RegisterCmd(&dgc.Command{
		Name:       "ping",
		IgnoreCase: true,
		Handler: func(ctx *dgc.Ctx) {
			ctx.RespondText("PONG!") //nolint
		},
	})
	router.RegisterCmd(&dgc.Command{
		Name:       "frases",
		IgnoreCase: true,
		Handler:    bot.frasesHandler,
	})
	router.RegisterCmd(&dgc.Command{
		Name:       "version",
		IgnoreCase: true,
		Handler:    bot.versionHandler,
	})
	router.RegisterCmd(&dgc.Command{
		Name:       "qr",
		IgnoreCase: true,
		Handler:    bot.qrHandler,
	})
	router.RegisterCmd(&dgc.Command{
		Name:       "putero",
		IgnoreCase: true,
		Handler:    bot.ptHandler,
	})
	router.RegisterCmd(&dgc.Command{
		Name:       "gafas",
		IgnoreCase: true,
		Handler:    bot.gafasHandler,
	})
	router.RegisterCmd(&dgc.Command{
		Name:       "elo",
		IgnoreCase: true,
		Handler:    bot.eloHandler,
	})
	router.RegisterCmd(&dgc.Command{
		Name:       "minita",
		IgnoreCase: true,
		Handler:    bot.minitaHandler,
		SubCommands: []*dgc.Command{{
			Name:       "add",
			IgnoreCase: true,
			Handler:    bot.addMinitaHandler,
		}},
	})
	router.Initialize(session)

	db, err := db.OpenDB()
	if err != nil {
		log.Fatalf("error: failed to connect to db: %v", err)
	}
	bot.db = db

	s3client, err := storage.NewS3Client()
	if err != nil {
		log.Fatalf("error: failed to create s3 client")
	}
	bot.s3 = s3client

	bot.client = httpc.NewHTTPClient()

	riotapi, err := newRiotAPIClient()
	if err != nil {
		log.Fatalf("error: failed to create riot api client: %v", err)
	}
	bot.riotapi = riotapi

	err = session.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	session.Close()
}
