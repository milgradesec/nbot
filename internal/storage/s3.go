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
	defaultS3Region   = "eu-west-1"
)

func NewS3Client() (*minio.Client, error) { //nolint
	var (
		endpoint  string
		region    string
		useTLS    bool
		accessKey string
		secretKey string
		found     bool
	)

	endpoint, found = os.LookupEnv("S3_ENDPOINT_URL")
	if !found {
		endpoint = defaultS3Endpoint
		log.Warn().Msg("S3_ENDPOINT_URL not set, using default endpoint: '" + defaultS3Endpoint + "'")
	}

	region, found = os.LookupEnv("S3_DEFAULT_REGION")
	if !found {
		region = defaultS3Region
		log.Warn().Msg("S3_DEFAULT_REGION not set, using: '" + defaultS3Region + "'")
	}

	if os.Getenv("S3_ENDPOINT_INSECURE") != "" {
		useTLS = true
	} else {
		useTLS = false
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
		Region: region,
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useTLS,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
