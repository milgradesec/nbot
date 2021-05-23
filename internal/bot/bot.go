package bot

import (
	"database/sql"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq" // psql driver
	httpc "github.com/milgradesec/go-libs/http"
	log "github.com/sirupsen/logrus"
	"github.com/yuhanfang/riot/apiclient"

	"github.com/bwmarrin/discordgo"
	"github.com/milgradesec/nbot/internal/db"
)

type Bot struct {
	Version string
	Token   string

	db      *sql.DB
	riotapi apiclient.Client
}

func (bot *Bot) Run() {
	rand.Seed(time.Now().Unix())

	session, err := discordgo.New("Bot " + bot.Token)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}
	session.Client = httpc.NewHTTPClient()
	session.AddHandler(bot.messageHandler)

	db, err := db.OpenDB()
	if err != nil {
		log.Fatalf("error: failed to connect to db: %v", err)
	}
	bot.db = db

	riotapi, err := newRiotAPIClient()
	if err != nil {
		log.Fatalf("error: failed to create Riot API Client: %v", err)
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
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Soy Putero",
		Video: &discordgo.MessageEmbedVideo{
			URL:    "https://s3.paesa.es/nbot-data/clips/putero.mp4?Content-Disposition=attachment%3B%20filename%3D%22clips%2Fputero.mp4%22&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=4E93nyrmHS3WzURmVzb%2F20210523%2Feu-west-1%2Fs3%2Faws4_request&X-Amz-Date=20210523T223531Z&X-Amz-Expires=432000&X-Amz-SignedHeaders=host&X-Amz-Signature=2b60adfd4d4c9c02997f249bd21b3a1fdecff6cb3d8d7e2ecac76850f6269bcf",
			Width:  400,
			Height: 400,
		},
	})
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
