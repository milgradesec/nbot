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

	"github.com/bwmarrin/discordgo"
	"github.com/milgradesec/nbot/internal/db"
)

type Bot struct {
	Version string
	Token   string
	db      *sql.DB
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
		_, err := s.ChannelMessageSend(m.ChannelID, "A su servicio")
		if err != nil {
			log.Errorf("error: failed to send message: %v\n", err)
		}
		return
	case "!version":
		_, err := s.ChannelMessageSend(m.ChannelID, bot.Version)
		if err != nil {
			log.Errorf("error: failed to send message: %v\n", err)
		}
		return
	case "!frases":
		_, err := s.ChannelMessageSend(m.ChannelID, bot.getAllQuotes())
		if err != nil {
			log.Errorf("error: failed to send message: %v\n", err)
		}
		return
	case "!ping":
		_, err := s.ChannelMessageSend(m.ChannelID, "PONG!")
		if err != nil {
			log.Errorf("error: failed to send message: %v\n", err)
		}
		return
	case "!elo":
		msg, err := bot.getLeagueElo("PEIN PACKER")
		if err != nil {
			log.Errorf("error: failed to get league data: %v\n", err)
		}
		_, err = s.ChannelMessageSend(m.ChannelID, msg)
		if err != nil {
			log.Errorf("error: failed to send message: %v\n", err)
		}
		return
	}

	if strings.Contains(m.Content, "nbot") {
		_, err := s.ChannelMessageSend(m.ChannelID, bot.getRandomQuote())
		if err != nil {
			log.Errorf("error: failed to send message: %v\n", err)
		}
	}
}
