package storage

import (
	"errors"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

const (
	defaultS3Endpoint = "s3.paesa.es"
)

func NewS3Client() (*minio.Client, error) {
	var (
		endpoint  string
		accessKey string
		secretKey string
		found     bool
	)

	endpoint, found = os.LookupEnv("S3_ENDPOINT_URL")
	if !found {
		endpoint = defaultS3Endpoint
		log.Warn().Msg("S3_ENDPOINT_URL not set, using default endpoint: '" + defaultS3Endpoint + "'")
	}

	accessKeyFile, found := os.LookupEnv("S3_ACCESS_KEY_FILE")
	if found {
		buf, err := os.ReadFile(accessKeyFile)
		if err != nil {
			return nil, err
		}
		accessKey = string(buf)
	} else {
		accessKey, found = os.LookupEnv("S3_ACCESS_KEY")
		if !found {
			return nil, errors.New("S3_ACCESS_KEY env variable not set")
		}
		log.Warn().Msg("Using unencrypted S3 access key from env, consider switching to S3_ACCESS_KEY_FILE")
	}

	secretKeyFile, found := os.LookupEnv("S3_SECRET_KEY_FILE")
	if found {
		buf, err := os.ReadFile(secretKeyFile)
		if err != nil {
			return nil, err
		}
		secretKey = string(buf)
	} else {
		secretKey, found = os.LookupEnv("S3_SECRET_KEY")
		if !found {
			return nil, errors.New("S3_SECRET_KEY env variable not set")
		}
		log.Warn().Msg("Using unencrypted S3 secret key from env, consider switching to S3_SECRET_KEY_FILE")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Region: "eu-west-1",
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
