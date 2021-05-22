package bot

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"
)

func (bot *Bot) getQRCodeURL(msg string) (string, error) {
	u := "https://qrcode.paesa.es/qr?data=" + url.QueryEscape(msg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("status code != 200")
	}
	return u, nil
}
