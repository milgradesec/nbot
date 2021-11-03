package bot

import (
	"context"
	"net/url"
	"time"
)

const bucketName = "nbot"

func (bot *Bot) generatePresignedURL(objectKey string) (string, error) {
	presignedURL, err := bot.s3.PresignedGetObject(context.Background(), bucketName, objectKey, time.Hour*8, make(url.Values))
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
