package bot

import (
	"database/sql"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq" // psql driver
	httpc "github.com/milgradesec/go-libs/http"
	"github.com/minio/minio-go/v7"
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

func (bot *Bot) Run() {
	rand.Seed(time.Now().Unix())

	session, err := discordgo.New("Bot " + bot.Token)
	if err != nil {
		log.Fatalf("error creating discord session: %v", err)
	}
	session.Client = httpc.NewHTTPClient()
	session.AddHandler(bot.messageHandler)

	db, err := db.OpenDB()
	if err != nil {
		log.Fatalf("error: failed to connect to db: %v", err)
	}
	bot.db = db

	s3client, err := newS3Client()
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

func (bot *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	switch m.Content {
	case "!nbot":
		bot.nbotHandler(s, m)
		return
	case "!version":
		bot.versionHandler(s, m)
		return
	case "!frases":
		bot.frasesHandler(s, m)
		return
	case "!ping":
		bot.pingHandler(s, m)
		return
	case "!fact":
		bot.factHandler(s, m)
		return
	case "!joke":
		bot.jokeHandler(s, m)
		return
	case "!putero":
		bot.ptHandler(s, m)
		return
	case "!meme":
		return
	}

	if strings.HasPrefix(m.Content, "!elo") {
		bot.eloHandler(s, m)
		return
	}

	if strings.HasPrefix(m.Content, "!qr") {
		bot.qrHandler(s, m)
		return
	}

	if strings.Contains(m.Content, "nbot") {
		bot.fraseHandler(s, m)
		return
	}
}

func (bot *Bot) nbotHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "A su servicio")
	if err != nil {
		log.Errorf("error: failed to send message: %v\n", err)
	}
}

func (bot *Bot) fraseHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, bot.getRandomQuote())
	if err != nil {
		log.Errorf("error: failed to send message: %v\n", err)
	}
}

func (bot *Bot) frasesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, bot.getAllQuotes())
	if err != nil {
		log.Errorf("error: failed to send message: %v\n", err)
	}
}

func (bot *Bot) versionHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, bot.Version)
	if err != nil {
		log.Errorf("error: failed to send message: %v\n", err)
	}
}

func (bot *Bot) pingHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "PONG!")
	if err != nil {
		log.Errorf("error: failed to send message: %v\n", err)
	}
}

func (bot *Bot) factHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, bot.getRandomFact())
	if err != nil {
		log.Errorf("error: failed to send message: %v\n", err)
	}
}

func (bot *Bot) jokeHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, bot.getRandomJoke())
	if err != nil {
		log.Errorf("error: failed to send message: %v\n", err)
	}
}

func (bot *Bot) qrHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := strings.TrimPrefix(m.Content, "!qr")

	url, err := bot.getQRCodeURL(msg)
	if err != nil {
		log.Errorf("error: failed to get QR code from message '%s': %v", msg, err)
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Toma QR",

		Image: &discordgo.MessageEmbedImage{
			URL:    url,
			Width:  400,
			Height: 400,
		},
	})
	if err != nil {
		log.Errorf("error: failed to send message: %v\n", err)
	}
}

func (bot *Bot) ptHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://s3.paesa.es/nbot-data/clips/putero.mp4")
	if err != nil {
		log.Errorf("error: failed to send message: %v\n", err)
	}
}

func (bot *Bot) eloHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!elo" {
		msg, err := bot.getLeagueElo("PEIN PACKER")
		if err != nil {
			log.Errorf("error: failed to get league data: %v\n", err)
			return
		}
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		if err != nil {
			log.Errorf("error: failed to send message: %v\n", err)
		}
		return
	}

	name := strings.TrimPrefix(m.Content, "!elo")
	msg, err := bot.getLeagueElo(name)
	if err != nil {
		log.Errorf("error: failed to get league data for '%s': %v\n", name, err)
		return
	}
	_, err = s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		log.Errorf("error: failed to send message: %v\n", err)
	}
}

func (bot *Bot) memeHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
}
