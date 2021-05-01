package bot

import (
	"testing"
)

func TestLoadQuotes(t *testing.T) {
	bot := &Bot{}
	err := bot.loadQuotes()
	if err != nil {
		t.Fatalf("error: failed to load quotes from quotes.json: %v", err)
	}
}
