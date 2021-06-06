package bot

import (
	"context"
	"math/rand"
	"time"

	"github.com/lus/dgc"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

var minitasObjectKeys []string

func loadMinitasKeys(client *minio.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectCh := client.ListObjects(ctx, "nbot-data", minio.ListObjectsOptions{
		Prefix:    "img/minitas",
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			log.Errorln(object.Err)
		}
		minitasObjectKeys = append(minitasObjectKeys, object.Key)
	}
}

func pickRandomMinita() string {
	randomIndex := rand.Intn(len(minitasObjectKeys)) //nolint
	return minitasObjectKeys[randomIndex]
}

func (bot *Bot) minitaHandler(ctx *dgc.Ctx) {
	key := pickRandomMinita()
	ctx.RespondText("https://s3.paesa.es/nbot-data/" + key) //nolint
}
