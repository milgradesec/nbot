package bot

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/lus/dgc"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

var minitasObjectKeys []string

func (bot *Bot) minitaHandler(ctx *dgc.Ctx) {
	key := pickRandomMinita()
	ctx.RespondText("https://s3.paesa.es/nbot-data/" + key) //nolint
}

func (bot *Bot) addMinitaHandler(ctx *dgc.Ctx) {
	msg := ctx.Event.Message
	if msg.Author.Username != "MILGRADESEC" {
		ctx.RespondText("Tu no tienes permiso para añadir nada") //nolint
		return
	}

	args := ctx.Arguments
	if args.Amount() != 1 {
		ctx.RespondText("Aprendete el puto comando: !minita add URL") //nolint
		return
	}
	log.Infof("addMinitaHandler called: args = %s", args.Raw())

	urlArg := args.Get(0)
	u, err := url.Parse(urlArg.Raw())
	if err != nil {
		log.Error(err)
		return
	}

	reqctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqctx, "GET", u.String(), nil)
	if err != nil {
		log.Error(err)
		return
	}
	req.Header.Set("Accept", "image/png")

	resp, err := bot.client.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("http status code != 200: %s", resp.Status)
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	h := computeImageMD5(buf)
	log.Infof("Image md5: %s", h)

	uploadMinitaIMG(bot.s3, h, resp.Body)
}

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

func computeImageMD5(buf []byte) string {
	h := md5.New() //nolint
	return hex.EncodeToString(h.Sum(buf))
}

func uploadMinitaIMG(client *minio.Client, key string, src io.Reader) {
	opts := minio.PutObjectOptions{
		ContentType: "image/png",
	}

	uploadInfo, err := client.PutObject(context.Background(), "nbot-data", "minitas/"+key+".png", src, -1, opts)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Info("Successfully uploaded bytes: ", uploadInfo)
}
