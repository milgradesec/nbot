package bot

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Fact struct {
	Text string `json:"text"`
}

func (bot *Bot) getRandomFact() string {
	const url = "https://uselessfacts.jsph.pl/random.json?language=en"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err.Error()
	}

	resp, err := bot.client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}

	var fact Fact
	err = json.Unmarshal(buf, &fact)
	if err != nil {
		return err.Error()
	}

	return fact.Text
}
