package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func (bot *Bot) quoteHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) == 1 {
		s.ChannelMessageSend(m.ChannelID, bot.getRandomQuote())
	}

	if len(args) > 2 {
		bot.addQuoteHandler(s, m, args)
	}
}

func (bot *Bot) quotesHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	s.ChannelMessageSend(m.ChannelID, bot.getAllQuotes())
}

func (bot *Bot) addQuoteHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.Author.Username != superUser {
		s.ChannelMessageSend(m.ChannelID, "Tu no tienes permiso para añadir nada. Putero.")
		return
	}

	quote := strings.Join(args[2:], " ")
	err := bot.insertNewQuote(quote)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Se ha producido un error al añadir la frase.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Frase añadida correctamente.")
}

func (bot *Bot) getRandomQuote() string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var quote string
	row := bot.db.QueryRowContext(ctx, `SELECT quote FROM quotes ORDER BY RANDOM() LIMIT 1`)
	if err := row.Scan(&quote); err != nil {
		log.Errorf("error: failed to handle db response: %v\n", err)
	}

	if err := row.Err(); err != nil {
		log.Errorf("error: failed to handle db response: %v\n", err)
	}
	return quote
}

func (bot *Bot) getAllQuotes() string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := bot.db.QueryContext(ctx, `SELECT * FROM quotes`)
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

func (bot *Bot) insertNewQuote(quote string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := bot.db.QueryContext(ctx, `INSERT INTO quotes VALUES ($1)`, quote)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error: failed to handle db response: %w", err)
	}
	return nil
}
