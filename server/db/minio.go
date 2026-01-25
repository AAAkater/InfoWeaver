package db

import (
	"context"
	"server/config"
	"server/utils"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func connectMinioClient(cfg *config.Config) (*minio.Client, error) {
	endpoint := cfg.GetMinioDSN()
	utils.Logger.Infof("use Minio endpoint:%s", endpoint)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MINIO_ACCESS_KEY, cfg.MINIO_SECRET_KEY, ""),
		Secure: cfg.MINIO_USE_SSL,
	})
	if err != nil {
		return nil, err
	}

	// Test connection
	exists, err := client.BucketExists(context.Background(), cfg.MINIO_BUCKET_NAME)
	if err != nil {
		return nil, err
	}

	// Create bucket if it doesn't exist
	if !exists {
		err = client.MakeBucket(context.Background(), cfg.MINIO_BUCKET_NAME, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		utils.Logger.Infof("Bucket '%s' created successfully", cfg.MINIO_BUCKET_NAME)
	}

	return client, nil
}
