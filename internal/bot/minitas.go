package bot

import (
	"bytes"
	"context"
	"crypto/md5" //nolint
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/lus/dgc"
	"github.com/minio/minio-go/v7"
	log "github.com/sirupsen/logrus"
)

func (bot *Bot) minitaHandler(ctx *dgc.Ctx) {
	key := bot.pickRandomMinitaID()
	ctx.RespondText("https://s3.paesa.es/nbot-data/minitas/" + key) //nolint
}

func (bot *Bot) addMinitaHandler(ctx *dgc.Ctx) {
	msg := ctx.Event.Message
	if msg.Author.Username != "MILGRADESEC" && msg.Author.Username != "L. L." {
		ctx.RespondText("Tu no tienes permiso para añadir nada") //nolint
		return
	}

	args := ctx.Arguments
	if args.Amount() != 1 {
		ctx.RespondText("Aprendete el puto comando: !minita add URL") //nolint
		return
	}

	urlArg := args.Get(0)
	u, err := url.Parse(urlArg.Raw())
	if err != nil {
		log.Error(err)
		return
	}

	resp, err := fetchImage(bot.client, u.String())
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}
	h := computeMD5(buf)

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		log.Error("error: Content-Type header not set")
		return
	}
	key := addContentTypeToKey(contentType, h)

	err = bot.insertMinitaID(key)
	if err != nil {
		log.Error(err)
		return
	}

	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}
	err = bot.uploadMinitaIMG(key, bytes.NewReader(buf), int64(len(buf)), opts)
	if err != nil {
		log.Errorf("error uploading img: %v", err)
		return
	}

	ctx.RespondText("Nueva minita añadida correctamente.\nMinita ID: " + key) //nolint
}

func (bot *Bot) pickRandomMinitaID() string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id string
	row := bot.db.QueryRowContext(ctx, `SELECT id FROM minitas ORDER BY RANDOM() LIMIT 1`)
	if err := row.Scan(&id); err != nil {
		log.Error(err)
	}

	if err := row.Err(); err != nil {
		log.Errorf("error: failed to handle db response: %v\n", err)
	}
	return id
}

func (bot *Bot) uploadMinitaIMG(key string, src io.Reader, size int64, opts minio.PutObjectOptions) error {
	_, err := bot.s3.PutObject(context.Background(), "nbot-data", "minitas/"+key, src, size, opts)
	if err != nil {
		return err
	}
	return nil
}

func (bot *Bot) insertMinitaID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := bot.db.QueryContext(ctx, `INSERT INTO minitas VALUES ($1)`, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

func fetchImage(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "image/png,image/jpeg")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("http status code != 200: " + resp.Status)
	}
	return resp, nil
}

func computeMD5(buf []byte) string {
	h := md5.New() //nolint
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func addContentTypeToKey(contentType, key string) string {
	switch contentType {
	case "image/png":
		return key + ".png"
	case "image/jpeg":
		return key + ".jpeg"
	}
	return key
}
