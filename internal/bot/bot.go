package bot

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Token string

	quotes []string
	db     *sql.DB
}

func (bot *Bot) Run() {
	rand.Seed(time.Now().Unix())

	session, err := discordgo.New("Bot " + bot.Token)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}
	session.AddHandler(bot.messageHandler)

	err = bot.loadQuotes()
	if err != nil {
		log.Fatalf("error: failed to load quotes from quotes.json: %v", err)
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
	return bot.quotes[rand.Intn(len(bot.quotes))] //nolint
}

func (bot *Bot) getAllQuotes() string {
	var msg string
	for _, frase := range bot.quotes {
		msg += frase
		msg += "\n"
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

type Quotes struct {
	Quotes []string
}

func (bot *Bot) loadQuotes() error {
	content, err := ioutil.ReadFile("../../quotes.json")
	if err != nil {
		return err
	}

	var payload Quotes
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return err
	}

	bot.quotes = payload.Quotes
	return nil
}
