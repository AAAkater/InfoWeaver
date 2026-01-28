package service

import (
	"context"
	"fmt"
	"io"
	"server/db"
	"server/models"
	"server/utils"
	"sync"

	"gorm.io/gorm"
)

var FileServiceApp = new(FileService)

type FileService struct{}

// CreateFile uploads a file to Minio and creates a database record
func (this *FileService) CreateFile(ctx context.Context, ownerID uint, filename string, fileType string, fileReader io.Reader, fileSize int64) error {

	objectName := fmt.Sprintf("%d/%s", ownerID, filename)

	// Prepare database record
	db_file := models.File{
		UserID:    ownerID,
		Name:      filename,
		MinIOPath: objectName,
		Size:      fileSize,
		Type:      fileType,
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
		if err := gorm.G[models.File](db.PgSqlDB).Create(ctx, &db_file); err != nil {
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

	utils.Logger.Infof("File created successfully: %s (ID: %d)", db_file.Name, db_file.ID)
	return nil
}

// GetFileListByUserID retrieves files for a specific user with pagination
// page: page number (1-indexed), pageSize: number of items per page
func (this *FileService) GetFileListByUserID(ctx context.Context, userID uint, page int, pageSize int) ([]models.FileInfo, error) {

	page = max(page, 1)
	pageSize = max(pageSize, 10)

	offset := (page - 1) * pageSize

	var files []models.FileInfo
	res := db.PgSqlDB.Model(&models.User{}).
		Where("ID = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(files)

	if res.Error != nil {
		utils.Logger.Errorf("Failed to get file list for user %d: %v", userID, res.Error)
		return nil, res.Error
	}
	return files, nil
}

// GetFileByID retrieves a file by ID
func (this *FileService) GetFileByID(ctx context.Context, fileID uint) (*models.File, error) {
	var file models.File
	result := db.PgSqlDB.First(&file, fileID)
	if result.Error != nil {
		utils.Logger.Errorf("Failed to get file with ID %d: %v", fileID, result.Error)
		return nil, result.Error
	}
	return &file, nil
}

// DownloadFile downloads a file from Minio
func (this *FileService) DownloadFile(ctx context.Context, fileID uint) (io.Reader, error) {
	file, err := this.GetFileByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	reader, err := db.MinioClient.DownloadFile(ctx, file.MinIOPath)
	if err != nil {
		utils.Logger.Errorf("Failed to download file %s: %v", file.MinIOPath, err)
		return nil, err
	}

	utils.Logger.Infof("File downloaded successfully: %s", file.Name)
	return reader, nil
}

// DeleteFileByID deletes a file from both Minio and database
func (this *FileService) DeleteFileByID(ctx context.Context, fileID uint, userID uint) error {
	file, err := this.GetFileByID(ctx, fileID)
	if err != nil {
		return err
	}

	// Verify the file belongs to the user
	if file.UserID != userID {
		return fmt.Errorf("user %d is not authorized to delete file %d", userID, fileID)
	}

	// Delete from Minio
	if err := db.MinioClient.DeleteFile(ctx, file.MinIOPath); err != nil {
		utils.Logger.Errorf("Failed to delete file from Minio: %v", err)
		return err
	}

	// Delete from database
	if result := db.PgSqlDB.Delete(file); result.Error != nil {
		utils.Logger.Errorf("Failed to delete file record from database: %v", result.Error)
		return result.Error
	}

	utils.Logger.Infof("File deleted successfully: %s (ID: %d)", file.Name, file.ID)
	return nil
}

// UpdateFileMetadata updates file metadata in the database
func (this *FileService) UpdateFileMetadata(ctx context.Context, fileID uint, userID uint, updates *models.File) error {
	file, err := this.GetFileByID(ctx, fileID)
	if err != nil {
		return err
	}

	// Verify the file belongs to the user
	if file.UserID != userID {
		return fmt.Errorf("user %d is not authorized to update file %d", userID, fileID)
	}

	// Update only specific fields
	if updates.Name != "" {
		file.Name = updates.Name
	}
	if updates.Type != "" {
		file.Type = updates.Type
	}

	result := db.PgSqlDB.Save(file)
	if result.Error != nil {
		utils.Logger.Errorf("Failed to update file metadata: %v", result.Error)
		return result.Error
	}

	utils.Logger.Infof("File metadata updated successfully: %s (ID: %d)", file.Name, file.ID)
	return nil
}

// CheckFileExists checks if a file exists in Minio
func (this *FileService) CheckFileExists(ctx context.Context, fileID uint) (bool, error) {
	file, err := this.GetFileByID(ctx, fileID)
	if err != nil {
		return false, err
	}

	exists, err := db.MinioClient.FileExists(ctx, file.MinIOPath)
	if err != nil {
		utils.Logger.Errorf("Failed to check file existence: %v", err)
		return false, err
	}

	return exists, nil
}
