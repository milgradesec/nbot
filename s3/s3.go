package s3

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	defaultS3Endpoint = "s3.paesa.es"
	defaultS3Region   = "eu-west-1"
	bucketName        = "nbot"
)

var Client *minio.Client

func NewClient() (*minio.Client, error) {
	var (
		endpoint  string
		region    string
		accessKey string
		secretKey string
	)

	if !viper.IsSet("AWS_S3_ENDPOINT") {
		log.Warn().Msg("AWS_S3_ENDPOINT not set, using default endpoint: '" + defaultS3Endpoint + "'")
		endpoint = defaultS3Endpoint
	} else {
		endpoint = viper.GetString("AWS_S3_ENDPOINT")
	}

	if !viper.IsSet("AWS_DEFAULT_REGION") {
		log.Warn().Msg("AWS_DEFAULT_REGION not set, using default region: '" + defaultS3Region + "'")
		region = defaultS3Region
	} else {
		region = viper.GetString("AWS_DEFAULT_REGION")
	}

	if !viper.IsSet("AWS_ACCESS_KEY_ID") {
		return nil, errors.New("AWS_ACCESS_KEY_ID not set")
	}
	accessKey = viper.GetString("AWS_ACCESS_KEY_ID")

	if !viper.IsSet("AWS_SECRET_ACCESS_KEY") {
		return nil, errors.New("AWS_SECRET_ACCESS_KEY not set")
	}
	secretKey = viper.GetString("AWS_SECRET_ACCESS_KEY")

	client, err := minio.New(endpoint, &minio.Options{
		Region: region,
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GeneratePresignedURL(objectKey string) (string, error) {
	presignedURL, err := Client.PresignedGetObject(context.Background(), bucketName, objectKey, time.Hour*8, make(url.Values))
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
