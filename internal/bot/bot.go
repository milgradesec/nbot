package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Token  string
	quotes []string
}

func (bot *Bot) Run() {
	session, err := discordgo.New("Bot " + bot.Token)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	rand.Seed(time.Now().Unix())
	session.AddHandler(bot.messageHandler)

	err = session.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}

	fmt.Println("Bot is running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	session.Close()
}

func (bot *Bot) getRandomQuote() string {
	return bot.quotes[rand.Intn(len(bot.quotes))] //nolint
}

func (bot *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!frases" {
		var msg string
		for _, frase := range bot.quotes {
			msg += frase
			msg += "\n"
		}
		_, err := s.ChannelMessageSend(m.ChannelID, msg)
		if err != nil {
			log.Printf("error: failed to send message; %v", err)
		}
		return
	}

	if strings.Contains(m.Content, "nbot") {
		_, err := s.ChannelMessageSend(m.ChannelID, bot.getRandomQuote())
		if err != nil {
			log.Printf("error: failed to send message; %v", err)
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
