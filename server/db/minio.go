package db

import (
	"context"
	"io"
	"server/config"
	"server/utils"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioService wraps minio client operations

var MinioClient *MinioService

type MinioService struct {
	client     *minio.Client
	bucketName string
}

// NewMinioService creates and initializes a new MinioService instance
func NewMinioService(cfg *config.Config, bucketName string) (*MinioService, error) {
	endpoint := cfg.GetMinioDSN()
	utils.Logger.Infof("use Minio endpoint:%s", endpoint)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MINIO_ACCESS_KEY, cfg.MINIO_SECRET_KEY, ""),
		Secure: cfg.MINIO_USE_SSL,
	})
	if err != nil {
		utils.Logger.Errorf("Failed to create Minio client: %v", err)
		return nil, err
	}

	// Test connection
	if exists, err := client.BucketExists(context.Background(), cfg.MINIO_BUCKET_NAME); err != nil {
		utils.Logger.Errorf("Failed to connect to Minio: %v", err)
		return nil, err
	} else {
		// Create bucket if it doesn't exist
		if !exists {
			if err := client.MakeBucket(context.Background(), cfg.MINIO_BUCKET_NAME, minio.MakeBucketOptions{}); err != nil {
				utils.Logger.Errorf("Failed to create bucket '%s': %v", cfg.MINIO_BUCKET_NAME, err)
				return nil, err
			}
			utils.Logger.Infof("Bucket '%s' created successfully", cfg.MINIO_BUCKET_NAME)
		}
	}

	return &MinioService{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// UploadFile uploads a file to Minio
func (ms *MinioService) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64) error {
	_, err := ms.client.PutObject(ctx, ms.bucketName, objectName, reader, objectSize, minio.PutObjectOptions{})
	if err != nil {
		utils.Logger.Errorf("Failed to upload file %s: %v", objectName, err)
		return err
	}
	utils.Logger.Infof("File %s uploaded successfully to bucket %s", objectName, ms.bucketName)
	return nil
}

// DownloadFile downloads a file from Minio
func (ms *MinioService) DownloadFile(ctx context.Context, objectName string) (io.Reader, error) {
	object, err := ms.client.GetObject(ctx, ms.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		utils.Logger.Errorf("Failed to download file %s: %v", objectName, err)
		return nil, err
	}
	return object, nil
}

// DeleteFile deletes a file from Minio
func (ms *MinioService) DeleteFile(ctx context.Context, objectName string) error {
	err := ms.client.RemoveObject(ctx, ms.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		utils.Logger.Errorf("Failed to delete file %s: %v", objectName, err)
		return err
	}
	utils.Logger.Infof("File %s deleted successfully from bucket %s", objectName, ms.bucketName)
	return nil
}

// ListFiles lists all files in a bucket
func (ms *MinioService) ListFiles(ctx context.Context, prefix string) ([]string, error) {
	var files []string
	objectsCh := ms.client.ListObjects(ctx, ms.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectsCh {
		if object.Err != nil {
			utils.Logger.Errorf("Error listing objects: %v", object.Err)
			return nil, object.Err
		}
		files = append(files, object.Key)
	}
	return files, nil
}

// FileExists checks if a file exists in Minio
func (ms *MinioService) FileExists(ctx context.Context, objectName string) (bool, error) {
	_, err := ms.client.StatObject(ctx, ms.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		utils.Logger.Errorf("Failed to check file existence %s: %v", objectName, err)
		return false, err
	}
	return true, nil
}

// GetFileMetadata gets the metadata of a file
func (ms *MinioService) GetFileMetadata(ctx context.Context, objectName string) (*minio.ObjectInfo, error) {
	objInfo, err := ms.client.StatObject(ctx, ms.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		utils.Logger.Errorf("Failed to get file metadata %s: %v", objectName, err)
		return nil, err
	}
	return &objInfo, nil
}

// CopyFile copies a file within Minio (same or different buckets)
func (ms *MinioService) CopyFile(ctx context.Context, sourceObject string, destBucket string, destObject string) error {
	source := minio.CopySrcOptions{
		Bucket: ms.bucketName,
		Object: sourceObject,
	}
	dest := minio.CopyDestOptions{
		Bucket: destBucket,
		Object: destObject,
	}

	_, err := ms.client.CopyObject(ctx, dest, source)
	if err != nil {
		utils.Logger.Errorf("Failed to copy file from %s/%s to %s/%s: %v", ms.bucketName, sourceObject, destBucket, destObject, err)
		return err
	}
	utils.Logger.Infof("File copied successfully from %s/%s to %s/%s", ms.bucketName, sourceObject, destBucket, destObject)
	return nil
}
