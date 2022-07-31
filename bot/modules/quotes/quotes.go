package quotes

import (
	"context"
	"time"

	"github.com/milgradesec/nbot/db"
	"github.com/rs/zerolog/log"
)

/*func (bot *Bot) quoteHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
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

func (bot *Bot) insertNewQuote(quote string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := db.Conn.Query(ctx, `INSERT INTO quotes VALUES ($1)`, quote)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to handle db response: %w", err)
	}
	return nil
}

*/
func GetRandom() string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var quote string
	err := db.Conn.QueryRow(ctx, `SELECT quote FROM quotes ORDER BY RANDOM() LIMIT 1`).Scan(&quote)
	if err != nil {
		log.Error().Err(err).Msg("failed to handle db response")
	}
	return quote
}

func GetAll() string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := db.Conn.Query(ctx, `SELECT * FROM quotes`)
	if err != nil {
		log.Error().Err(err).Msgf("failed to query db")
	}
	defer rows.Close()

	var msg string
	for rows.Next() {
		var quote string
		err = rows.Scan(&quote)
		if err != nil {
			log.Error().Err(err).Msg("failed to handle db response")
			break
		}
		msg += quote
		msg += "\n"
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("failed to handle db response")
	}
	return msg
}
