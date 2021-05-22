package bot

import "testing"

func TestQR(t *testing.T) {
	bot := &Bot{}
	if _, err := bot.getQRCodeURL("Me ahogo"); err != nil {
		t.Fatal(err)
	}
}
