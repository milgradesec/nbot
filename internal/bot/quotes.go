package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/lus/dgc"
	log "github.com/sirupsen/logrus"
)

func (bot *Bot) fraseHandler(ctx *dgc.Ctx) {
	ctx.RespondText(bot.getRandomQuote())
}

func (bot *Bot) frasesHandler(ctx *dgc.Ctx) {
	ctx.RespondText(bot.getAllQuotes())
}

func (bot *Bot) addQuoteHandler(ctx *dgc.Ctx) {
	args := ctx.Arguments

	msg := ctx.Event.Message
	if msg.Author.Username != "MILGRADESEC" {
		ctx.RespondText("Tu no tienes permiso para añadir nada. Putero.")
		return
	}

	err := bot.insertNewQuote(args.Raw())
	if err != nil {
		ctx.RespondText("Se ha producido un error al añadir la frase.")
		return
	}
	ctx.RespondText("Frase añadida correctamente.")
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

func (bot *Bot) getRandomQuote() string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := bot.db.QueryRowContext(ctx, `SELECT quote FROM quotes ORDER BY RANDOM() LIMIT 1`)
	var quote string
	if err := row.Scan(&quote); err != nil {
		log.Error(err)
	}

	if err := row.Err(); err != nil {
		log.Errorf("error: failed to handle db response: %v\n", err)
	}
	return quote
}

func (bot *Bot) getAllQuotes() string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
