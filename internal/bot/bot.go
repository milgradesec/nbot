package bot

import (
	"database/sql"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" // psql driver
	"github.com/lus/dgc"
	httpc "github.com/milgradesec/go-libs/http"
	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
	"github.com/yuhanfang/riot/apiclient"

	"github.com/bwmarrin/discordgo"
	"github.com/milgradesec/nbot/internal/db"
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
		Name:       "fact",
		IgnoreCase: true,
		Handler:    bot.factHandler,
	})
	router.RegisterCmd(&dgc.Command{
		Name:       "joke",
		IgnoreCase: true,
		Handler:    bot.jokeHandler,
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
	router.Initialize(session)

	db, err := db.OpenDB()
	if err != nil {
		log.Fatalf("error: failed to connect to db: %v", err)
	}
	bot.db = db

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

func (bot *Bot) versionHandler(ctx *dgc.Ctx) {
	ctx.RespondText(bot.Version) //nolint
}

func (bot *Bot) gafasHandler(ctx *dgc.Ctx) {
	ctx.RespondEmbed(&discordgo.MessageEmbed{ //nolint
		Title: "Con Gafas",

		Image: &discordgo.MessageEmbedImage{
			URL:    "https://s3.paesa.es/nbot-data/img/congafas.png",
			Width:  400,
			Height: 400,
		},
	})

	ctx.RespondEmbed(&discordgo.MessageEmbed{ //nolint
		Title: "Sin Gafas",

		Image: &discordgo.MessageEmbedImage{
			URL:    "https://s3.paesa.es/nbot-data/img/singafas.png",
			Width:  400,
			Height: 400,
		},
	})
}

func (bot *Bot) ptHandler(ctx *dgc.Ctx) {
	ctx.RespondText("https://s3.paesa.es/nbot-data/clips/putero.mp4") //nolint
}
