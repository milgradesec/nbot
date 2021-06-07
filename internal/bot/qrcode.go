package bot

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lus/dgc"
	log "github.com/sirupsen/logrus"
)

func (bot *Bot) qrHandler(ctx *dgc.Ctx) {
	args := ctx.Arguments
	msg := args.Raw()

	url, err := bot.getQRCodeURL(msg)
	if err != nil {
		log.Errorf("error: failed to get QR code from message '%s': %v", msg, err)
	}

	ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title: "Toma QR",
		Image: &discordgo.MessageEmbedImage{
			URL:    url,
			Width:  400,
			Height: 400,
		},
	})
}

func (bot *Bot) getQRCodeURL(msg string) (string, error) {
	u := "https://qrcode.paesa.es/qr?data=" + url.QueryEscape(msg)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return "", err
	}

	resp, err := bot.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("status code != 200")
	}
	return u, nil
}
