package storage

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
	defaultEndpoint   = "s3.paesa.es"
	defaultRegion     = "eu-west-1"
	defaultBucketName = "nbot"
	defaultExpiration = time.Hour * 6
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
		log.Warn().Msg("AWS_S3_ENDPOINT not set, using default endpoint: '" + defaultEndpoint + "'")
		endpoint = defaultEndpoint
	} else {
		endpoint = viper.GetString("AWS_S3_ENDPOINT")
	}

	if !viper.IsSet("AWS_DEFAULT_REGION") {
		log.Warn().Msg("AWS_DEFAULT_REGION not set, using default region: '" + defaultRegion + "'")
		region = defaultRegion
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

func Get(ctx context.Context, objectName string) error {
	return errors.New("not implemented")
}

func Put(ctx context.Context, objectName string) error {
	return errors.New("not implemented")
}

func Delete(ctx context.Context, objectName string) error {
	return errors.New("not implemented")
}

func PresignedGet(ctx context.Context, objectName string) (string, error) {
	u, err := Client.PresignedGetObject(ctx, defaultBucketName, objectName, defaultExpiration, make(url.Values))
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
