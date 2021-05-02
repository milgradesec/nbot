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
	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Token string
	db    *sql.DB
}

func (bot *Bot) Run() {
	rand.Seed(time.Now().Unix())

	session, err := discordgo.New("Bot " + bot.Token)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}
	session.AddHandler(bot.messageHandler)

	err = bot.openDB()
	if err != nil {
		log.Fatalf("error: failed to connect to db: %v", err)
	}

	err = session.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}

	log.Infoln("Bot is running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	session.Close()
}

func (bot *Bot) getRandomQuote() string {
	row := bot.db.QueryRow("SELECT quote FROM quotes ORDER BY RANDOM() LIMIT 1")

	var name string
	if err := row.Scan(&name); err != nil {
		log.Error(err)
	}

	if err := row.Err(); err != nil {
		log.Errorf("error: failed to handle db response: %v\n", err)
	}
	return name
}

func (bot *Bot) getAllQuotes() string {
	rows, err := bot.db.Query("SELECT * FROM quotes")
	if err != nil {
		log.Errorf("error: failed to query db: %v\n", err)
	}
	defer rows.Close()

	var msg string
	for rows.Next() {
		var quote string
		err = rows.Scan(&quote)
		if err != nil {
			log.Errorf("error: failed to handle db response: %v\n", err)
			break
		}
		msg += quote
		msg += "\n"
	}

	if err := rows.Err(); err != nil {
		log.Errorf("error: failed to handle db response: %v\n", err)
	}
	return msg
}

func (bot *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!frases" {
		_, err := s.ChannelMessageSend(m.ChannelID, bot.getAllQuotes())
		if err != nil {
			log.Errorf("error: failed to send message; %v\n", err)
		}
		return
	}

	if strings.Contains(m.Content, "nbot") {
		_, err := s.ChannelMessageSend(m.ChannelID, bot.getRandomQuote())
		if err != nil {
			log.Errorf("error: failed to send message; %v\n", err)
		}
	}
}
