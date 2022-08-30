package minitas

import (
	"context"
	"errors"
	"fmt"

	db "github.com/milgradesec/nbot/database"
	"github.com/milgradesec/nbot/storage"
	"github.com/minio/minio-go/v7"
)

func Get(ctx context.Context, id string) (string, error) {
	return generateURLFromID(ctx, id)
}

func GetRandom(ctx context.Context) (string, error) {
	id, err := pickRandomID(ctx)
	if err != nil {
		return "", err
	}
	return generateURLFromID(ctx, id)
}

func Delete(ctx context.Context, id string) error {
	result, err := db.Conn.Exec(ctx, `DELETE FROM minitas WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if rows := result.RowsAffected(); rows == 0 {
		return errors.New("no rows affected")
	}
	return storage.Client.RemoveObject(ctx, "nbot", "minitas/"+id, minio.RemoveObjectOptions{})
}

func Insert(ctx context.Context) error {
	return errors.New("not implemented yet")
}

func generateURLFromID(ctx context.Context, id string) (string, error) {
	url, err := storage.PresignedGet(ctx, "minitas/"+id)
	if err != nil {
		return "", err
	}
	return url, nil
}

/*func (bot *Bot) minitaHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) > 2 {
		switch args[1] {
		case "add":
			bot.addMinitaHandler(s, m, args)
		case "delete":
			bot.deleteMinitaHandler(s, m, args)
		}
	}

	if len(args) == 1 {
		key, err := bot.pickRandomMinitaID()
		if err != nil {
			log.Error().Err(err).Msg("failed pick a random minita")
			s.ChannelMessageSend(m.ChannelID, "Se ha producido un error interno.")
			return
		}

		url, err := s3.PresignedURL("minitas/" + key)
		if err != nil {
			log.Error().Err(err).Msg("failed to generate presigned")
			return
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: "Minita",
			Image: &discordgo.MessageEmbedImage{
				URL: url,
			},
		})
	}
}*/

/*func (bot *Bot) addMinitaHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) <= 2 {
		s.ChannelMessageSend(m.ChannelID, "Aprendete el comando: !minita add URL")
		return
	}

	if m.Author.Username != superUser {
		s.ChannelMessageSend(m.ChannelID, "Tu no tienes permiso para añadir nada. Putero.")
		return
	}

	srcURL := strings.Join(args[2:], " ")
	u, err := url.Parse(srcURL)
	if err != nil {
		log.Error().Err(err).Msgf("failed to parse url '%s'", srcURL)
		return
	}

	resp, err := fetchImage(bot.client, u.String())
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch image from source")
		return
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return
	}
	h := computeMD5(buf)

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		log.Error().Msg("Content-Type header not set")
		return
	}
	key := addContentTypeToKey(contentType, h)

	err = bot.insertMinitaID(key)
	if err != nil {
		log.Error().Err(err).Msgf("failed to insert minita with id '%s'", key)
		return
	}

	opts := minio.PutObjectOptions{
		ContentType:  contentType,
		CacheControl: "public, max-age=604800",
	}
	err = bot.uploadMinitaIMG(key, bytes.NewReader(buf), int64(len(buf)), opts)
	if err != nil {
		log.Error().Err(err).Msgf("error uploading image to s3")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Nueva minita añadida correctamente.\nMinita ID: "+key)
}*/

/*func (bot *Bot) deleteMinitaHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) <= 2 {
		s.ChannelMessageSend(m.ChannelID, "El comando es !minita delete MinitaID")
		return
	}

	if m.Author.Username != superUser {
		s.ChannelMessageSend(m.ChannelID, "Tu no tienes permiso para quitar ninguna minita.")
		return
	}

	id := strings.Join(args[2:], " ")
	found, err := bot.minitaExists(id)
	if err != nil {
		log.Error().Err(err).Msgf("failed check if minita with id '%s' already exists", id)
		return
	}
	if !found {
		s.ChannelMessageSend(m.ChannelID, "No existe ninguna minita con ese ID.")
		return
	}

	if err = bot.deleteMinitaID(id); err != nil {
		log.Error().Err(err).Msgf("failed to delete minita with id '%s'", id)
	}

	if err = bot.deleteMinitaIMG(id); err != nil {
		log.Error().Err(err).Msgf("failed to delete minita image from s3")
	}

	s.ChannelMessageSend(m.ChannelID, "Minita eliminada correctamente.")
}

func checkAlreadyExists(id string) (bool, error) {
	var result int
	err := db.Conn.QueryRow(context.Background(), `SELECT 1 FROM minitas WHERE id = $1`, id).Scan(&result)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to handle db response: %w", err)
	}
	return true, nil
}*/

func pickRandomID(ctx context.Context) (string, error) {
	var id string
	err := db.Conn.QueryRow(ctx, `SELECT id FROM minitas ORDER BY RANDOM() LIMIT 1`).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to handle db response: %w", err)
	}
	return id, nil
}

/*func uploadMinitaIMG(key string, src io.Reader, size int64, opts minio.PutObjectOptions) error {
	_, err := s3.Client.PutObject(context.Background(), "nbot", "minitas/"+key, src, size, opts)
	if err != nil {
		return err
	}
	return nil
}

func insertMinitaID(id string) error {
	rows, err := db.Conn.Query(context.Background(), `INSERT INTO minitas VALUES ($1)`, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to handle db response: %w", err)
	}
	return nil
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
	case "image/webp":
		return key + ".webp"
	}
	return key
}

func fetchImage(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "image/webp,image/jpeg,image/png")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("http status code != 200: " + resp.Status)
	}
	return resp, nil
}*/
