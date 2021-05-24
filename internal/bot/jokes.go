package bot

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type joke struct {
	Joke string `json:"joke"`
}

func (bot *Bot) getRandomJoke() string {
	const url = "https://icanhazdadjoke.com/"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err.Error()
	}
	req.Header.Add("Accept", "application/json")

	resp, err := bot.client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}

	var joke joke
	err = json.Unmarshal(buf, &joke)
	if err != nil {
		return err.Error()
	}

	return joke.Joke
}
