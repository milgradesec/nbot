package bot

import "testing"

func TestOpenDB(t *testing.T) {
	bot := &Bot{}
	err := bot.openDB()
	if err != nil {
		t.Error(err)
	}
}
