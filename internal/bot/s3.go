package bot

import (
	"errors"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func newS3Client() (*minio.Client, error) {
	accessKey, found := os.LookupEnv("S3_ACCESS_KEY")
	if !found {
		return nil, errors.New("")
	}

	secretKey, found := os.LookupEnv("S3_SECRET_KEY")
	if !found {
		return nil, errors.New("")
	}

	client, err := minio.New("s3.paesa.es", &minio.Options{
		Region: "eu-west-1",
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
