package bot

import (
	"bytes"
	"context"
	"crypto/md5" //nolint
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
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
	key, err := bot.pickRandomMinitaID()
	if err != nil {
		ctx.RespondText("Se ha producido un error interno.")
		return
	}
	ctx.RespondText("https://s3.paesa.es/nbot/minitas/" + key)
}

func (bot *Bot) addMinitaHandler(ctx *dgc.Ctx) {
	args := ctx.Arguments
	if args.Amount() != 1 {
		ctx.RespondText("Aprendete el comando: !minita add URL")
		return
	}

	msg := ctx.Event.Message
	if msg.Author.Username != "MILGRADESEC" {
		ctx.RespondText("Tu no tienes permiso para añadir nada. Putero.")
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
		log.Errorf("error: failed to fetch image from source: %v", err)
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
		ContentType:  contentType,
		CacheControl: "public, max-age=604800",
	}
	err = bot.uploadMinitaIMG(key, bytes.NewReader(buf), int64(len(buf)), opts)
	if err != nil {
		log.Errorf("error uploading img: %v", err)
		return
	}

	ctx.RespondText("Nueva minita añadida correctamente.\nMinita ID: " + key)
}

func (bot *Bot) deleteMinitaHandler(ctx *dgc.Ctx) {
	args := ctx.Arguments
	if args.Amount() != 1 {
		ctx.RespondText("El comando es !minita delete MinitaID")
		return
	}

	msg := ctx.Event.Message
	if msg.Author.Username != "MILGRADESEC" {
		ctx.RespondText("Tu no tienes permiso para quitar ninguna minita.")
		return
	}

	id := args.Get(0).Raw()
	found, err := bot.minitaExists(id)
	if err != nil {
		log.Error(err)
		return
	}
	if !found {
		ctx.RespondText("No existe ninguna minita con ese ID.")
		return
	}

	if err = bot.deleteMinitaID(id); err != nil {
		log.Error(err)
	}

	if err = bot.deleteMinitaIMG(id); err != nil {
		log.Error(err)
	}

	ctx.RespondText("Minita eliminada correctamente.")
}

func (bot *Bot) minitaExists(id string) (bool, error) {
	var result int
	row := bot.db.QueryRow(`SELECT 1 FROM minitas WHERE id = $1`, id)
	if err := row.Scan(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	if err := row.Err(); err != nil {
		return false, fmt.Errorf("error: failed to handle db response: %w", err)
	}
	return true, nil
}

func (bot *Bot) pickRandomMinitaID() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var id string
	row := bot.db.QueryRowContext(ctx, `SELECT id FROM minitas ORDER BY RANDOM() LIMIT 1`)
	if err := row.Scan(&id); err != nil {
		return id, err
	}

	if err := row.Err(); err != nil {
		return id, fmt.Errorf("error: failed to handle db response: %w", err)
	}
	return id, nil
}

func (bot *Bot) uploadMinitaIMG(key string, src io.Reader, size int64, opts minio.PutObjectOptions) error {
	_, err := bot.s3.PutObject(context.Background(), "nbot", "minitas/"+key, src, size, opts)
	if err != nil {
		return err
	}
	return nil
}

func (bot *Bot) insertMinitaID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := bot.db.QueryContext(ctx, `INSERT INTO minitas VALUES ($1)`, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error: failed to handle db response: %w", err)
	}
	return nil
}

func (bot *Bot) deleteMinitaID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := bot.db.ExecContext(ctx, `DELETE FROM minitas WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("no rows affected")
	}
	return nil
}

func (bot *Bot) deleteMinitaIMG(key string) error {
	return bot.s3.RemoveObject(context.Background(), "nbot", "minitas/"+key, minio.RemoveObjectOptions{})
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
