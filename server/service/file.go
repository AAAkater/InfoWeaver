package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"server/db"
	"server/models"
	"server/utils"
	"sync"
	"time"

	"gorm.io/gorm"
)

var FileServiceApp = new(FileService)

type FileService struct{}

// CreateFile uploads a file to Minio and creates a database record
func (this *FileService) CreateFile(ctx context.Context, ownerID uint, datasetID uint, filename string, fileType string, fileReader io.Reader, fileSize int64) error {

	objectName := fmt.Sprintf("%d/%s", ownerID, filename)

	// Prepare database record
	dbFile := models.File{
		UserID:    ownerID,
		Name:      filename,
		MinioPath: objectName,
		Size:      fileSize,
		Type:      fileType,
		DatasetID: datasetID,
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// Upload file to Minio in parallel
	wg.Go(func() {
		if err := db.MinioClient.UploadFile(ctx, objectName, fileReader, fileSize); err != nil {
			utils.Logger.Errorf("Failed to upload file to Minio: %v", err)
			errChan <- err
		}
	})

	// Save file record to database in parallel
	wg.Go(func() {
		if err := gorm.G[models.File](db.PgSqlDB).Create(ctx, &dbFile); err != nil {
			utils.Logger.Errorf("Failed to save file record to database: %v", err)
			errChan <- err
		}
	})

	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	utils.Logger.Infof("File created successfully: %s (ID: %d)", dbFile.Name, dbFile.ID)
	return nil
}

// GetFileListByUserID retrieves files for a specific user with pagination
// page: page number (1-indexed), pageSize: number of items per page
func (this *FileService) GetFileListByUserID(ctx context.Context, userID uint, page int, pageSize int) (total int64, files []models.SimpleFileInfo, e error) {

	page = max(page, 1)
	pageSize = max(pageSize, 10)

	offset := (page - 1) * pageSize

	result := db.PgSqlDB.Model(&models.File{}).
		Where("user_id= ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&files)

	switch result.Error {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		utils.Logger.Errorf("No files found for user %d: %v", userID, result.Error)
		return 0, files, nil
	default:
		utils.Logger.Errorf("Failed to get file list for user %d: %v", userID, result.Error)
		return 0, nil, errors.New("Unknown error occurred while retrieving file list")
	}
	return result.RowsAffected, files, nil
}

// GetFileInfoByFileID retrieves a file by fileID
func (this *FileService) GetFileInfoByFileID(ctx context.Context, fileID uint, ownerID uint) (fileInfo *models.DetailedFileInfo, e error) {
	result := db.PgSqlDB.Model(&models.File{}).
		Where("ID = ? AND user_id = ?", fileID, ownerID).
		Find(&fileInfo)
	if result.Error != nil {
		utils.Logger.Errorf("Failed to get file with ID %d: %v", fileID, result.Error)
		return nil, result.Error
	}
	return fileInfo, nil
}

func (this *FileService) GetFilePathByFileID(ctx context.Context, fileID uint, ownerID uint) (string, error) {
	dbFile, err := gorm.G[models.File](db.PgSqlDB).
		Select("minio_path").
		Where("id = ? AND user_id = ?", fileID, ownerID).
		First(ctx)
	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		return "", errors.New("file not found")
	default:
		utils.Logger.Errorf("Failed to get file path for file ID %d: %v", fileID, err)
		return "", err
	}
	return dbFile.MinioPath, nil
}

// GetDownloadURLByFilePath generates a presigned download URL for a file by file path
func (this *FileService) GetDownloadURLByFilePath(ctx context.Context, filePath string) (string, error) {

	// Generate presigned URL with 1 hour expiration (3600 seconds)
	downloadURL, err := db.MinioClient.GetPresignedDownloadURL(ctx, filePath, 3600)
	if err != nil {
		utils.Logger.Errorf("Failed to generate download URL for file %s: %v", filePath, err)
		return "", err
	}

	utils.Logger.Infof("Download URL generated successfully: %s", filePath)
	return downloadURL, nil
}

// DeleteFileByFileID deletes a file from both Minio and database
func (this *FileService) DeleteFileByFileID(ctx context.Context, fileID uint, filePath string) error {

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	// Delete from Minio in parallel
	wg.Go(func() {
		if err := db.MinioClient.DeleteFile(ctx, filePath); err != nil {
			utils.Logger.Errorf("Failed to delete file from Minio: %v", err)
			errChan <- err
		}
	})

	// Delete from database in parallel
	wg.Go(func() {
		if _, err := gorm.G[models.File](db.PgSqlDB).
			Where("ID = ?", fileID).
			Update(ctx, "DeletedAt", time.Now()); err != nil {
			utils.Logger.Errorf("Failed to delete file record from database: %v", err)
			errChan <- err
		}
	})

	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	utils.Logger.Infof("File deleted successfully: %s (ID: %d)", filePath, fileID)
	return nil
}

// UpdateFileInfo updates file info in the database
func (this *FileService) UpdateFileInfo(ctx context.Context, fileID uint, userID uint, newFileInfo models.FileInfoUpdate) error {

	if _, err := gorm.G[models.File](db.PgSqlDB).
		Where("ID = ? AND UserID = ?", fileID, userID).
		Updates(ctx, models.File{
			Size:      newFileInfo.Size,
			Name:      newFileInfo.Name,
			Type:      newFileInfo.Type,
			MinioPath: newFileInfo.MinioPath,
			UserID:    newFileInfo.UserID,
		}); err != nil {
		utils.Logger.Errorf("Failed to update file metadata: %v", err)
		return err
	}
	utils.Logger.Infof("File metadata updated successfully: %s (ID: %d)", newFileInfo.Name, fileID)
	return nil
}

// CheckFileExistsByFileID checks if a file exists in Minio
func (this *FileService) CheckFileExistsByFileID(ctx context.Context, fileID uint) (bool, error) {
	fileInfo, err := gorm.G[models.File](db.PgSqlDB).
		Where("ID = ?", fileID).
		First(ctx)

	if err != nil {
		utils.Logger.Errorf("Failed to get file path from database: %v", err)
		return false, err
	}
	exists, err := db.MinioClient.FileExists(ctx, fileInfo.MinioPath)
	if err != nil {
		utils.Logger.Errorf("Failed to check file existence: %v", err)
		return false, err
	}

	return exists, nil
}
