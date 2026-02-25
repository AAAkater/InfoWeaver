package v1

import (
	"errors"
	"fmt"
	"server/config"
	"server/middleware"
	"server/models"
	"server/models/common/response"
	"server/service"
	"server/utils"
	"sync"

	"github.com/labstack/echo/v5"
)

func SetFileRouter(e *echo.Echo) {

	fileRouterGroup := e.Group(config.API_V1+"/file", middleware.TokenMiddleware())

	fileHandler := &fileApi{}
	fileRouterGroup.POST("/upload", fileHandler.uploadFile)
	fileRouterGroup.GET("/list", fileHandler.ListFiles)
	fileRouterGroup.GET("/info/:file_id", fileHandler.getSingleDetailedFileInfo)
	fileRouterGroup.GET("/download/:file_id", fileHandler.getDownloadFileURL)
	fileRouterGroup.POST("/delete/:file_id", fileHandler.deleteFile)

}

type fileApi struct{}

// uploadFile godoc
// @Summary      Multi-File Upload
// @Description  Upload multiple files to the server and associate them with a dataset. Files are stored in MinIO object storage.
// @Tags         File
// @Accept       multipart/form-data
// @Produce      json
// @Param        files formData []file true "Files to upload" format(binary)
// @Param        id formData uint true "Dataset ID to associate the uploaded files with"
// @Success      200 {object} response.ResponseBase[models.MultiFileUploadResp] "Files uploaded successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid files, missing parameters, or upload failed"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      403 {object} response.ResponseBase[any] "Forbidden, user not authorized to access the specified dataset"
// @Router       /file/upload [post]
func (this *fileApi) uploadFile(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}

	// Get dataset ID
	datasetID, err := echo.FormValue[uint](ctx, "id")
	if err != nil {
		return response.ErrMissDatasetID()
	}

	// Validate dataset ownership
	if _, err := datasetService.GetDatasetInfoByID(ctx.Request().Context(), datasetID, currentUser.ID); err != nil {
		return response.ErrDatasetNotFound()
	}

	// Get all uploaded files
	form, err := ctx.MultipartForm()
	if err != nil {
		return response.ErrMissFile()
	}

	fileHeaders := form.File["files"]
	fileNumber := len(fileHeaders)
	if fileNumber == 0 {
		return response.ErrNoFileUploaded()
	} else if fileNumber > 5 { // Limit maximum number of files to 5
		return response.ErrFileNumberLimited()
	}

	Logger.Infof("Received %d files for upload", fileNumber)

	// Process files in parallel
	var wg sync.WaitGroup
	resultChan := make(chan models.FileUploadInfo, fileNumber)
	errChan := make(chan error, fileNumber)

	for _, fileHeader := range fileHeaders {
		wg.Go(func() {
			fh := fileHeader

			// Open the uploaded file
			src, err := fh.Open()
			if err != nil {
				Logger.Errorf("Failed to open uploaded file %s: %v", fh.Filename, err)
				errChan <- fmt.Errorf("failed to open file %s: %w", fh.Filename, err)
				return
			}
			defer src.Close()

			// Get file type/MIME type
			fileType := fh.Header.Get("Content-Type")
			if fileType == "" {
				fileType = "application/octet-stream"
			}

			// Use channels to coordinate parallel MinIO upload and DB record creation
			minioErrChan := make(chan error, 1)
			dbFileChan := make(chan *models.File, 1)
			dbErrChan := make(chan error, 1)

			var fileWg sync.WaitGroup

			// Goroutine 1: Upload to MinIO
			fileWg.Go(func() {
				if err := fileService.UploadFileToMinio(
					ctx.Request().Context(), currentUser.ID, datasetID, fh.Filename, src, fh.Size); err != nil {
					Logger.Errorf("Failed to upload file %s to Minio: %v", fh.Filename, err)
					minioErrChan <- fmt.Errorf("failed to upload file %s to Minio: %w", fh.Filename, err)
				} else {
					minioErrChan <- nil
				}
			})

			// Goroutine 2: Create database record
			fileWg.Go(func() {
				// Create database record
				dbFile, err := fileService.CreateFileInfo(
					ctx.Request().Context(),
					currentUser.ID,
					datasetID,
					fh.Filename,
					fileType,
					fh.Size,
				)
				switch {
				case errors.Is(err, service.ErrDuplicatedKey):
					Logger.Errorf("Failed to save file record %s to database: %v", fh.Filename, err)
					dbErrChan <- response.ForbiddenWithMsg(fmt.Sprintf("Unauthorized access to the dataset: %d", datasetID))
					return
				case err != nil:
					Logger.Errorf("Failed to create file record %s: %v", fh.Filename, err)
					dbErrChan <- fmt.Errorf("failed to create file record %s: %w", fh.Filename, err)
					return
				}
				dbFileChan <- dbFile
				dbErrChan <- nil
			})

			fileWg.Wait()
			close(minioErrChan)
			close(dbFileChan)
			close(dbErrChan)

			// Check for MinIO upload errors
			if minioErr := <-minioErrChan; minioErr != nil {
				errChan <- minioErr
				return
			}

			// Check for database errors
			if dbErr := <-dbErrChan; dbErr != nil {
				errChan <- dbErr
				return
			}

			// Get dbFile from channel
			dbFile := <-dbFileChan

			// Publish file upload event
			if err := fileService.PublishFileUploadEvent(ctx.Request().Context(), dbFile); err != nil {
				Logger.Errorf("Failed to publish file upload event for %s: %v", fh.Filename, err)
				// Continue even if event publishing fails
			}

			Logger.Infof("File uploaded successfully: %s", fh.Filename)
			resultChan <- models.FileUploadInfo{
				OwnerID:   currentUser.ID,
				DatasetID: datasetID,
				Name:      fh.Filename,
				Type:      fileType,
				Size:      fh.Size,
			}
		})
	}

	wg.Wait()
	close(resultChan)
	close(errChan)

	// Collect results
	var uploadedFiles []models.FileUploadInfo
	var errors []error

	for resp := range resultChan {
		uploadedFiles = append(uploadedFiles, resp)
	}

	for err := range errChan {
		errors = append(errors, err)
	}

	// If all files failed, return error
	if len(uploadedFiles) == 0 && len(errors) > 0 {
		Logger.Errorf("All file uploads failed: %v", errors)
		return response.FailWithMsg(ctx, fmt.Sprintf("All file uploads failed: %v", errors))
	}

	// If some files failed, log warnings but return success for successful uploads
	if len(errors) > 0 {
		Logger.Warnf("Some file uploads failed: %v", errors)
	}

	Logger.Infof("Successfully uploaded %d out of %d files", len(uploadedFiles), len(fileHeaders))
	return response.OkWithData(ctx, models.MultiFileUploadResp{
		Files: uploadedFiles,
	})
}

// ListFiles godoc
// @Summary      Get File List
// @Description  Retrieve a paginated list of files for the current user
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        page query int true "Page number" minimum(1)
// @Param        page_size query int true "Number of files per page" minimum(1) maximum(100)
// @Success      200 {object} response.ResponseBase[models.SimpleFileInfoListResp] "File list retrieved successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Router       /file/list [get]
func (this *fileApi) ListFiles(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}

	args, err := utils.BindAndValidate[models.ListFilesReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	total, files, err := fileService.GetFileListByUserID(
		ctx.Request().Context(),
		currentUser.ID,
		args.DatasetID,
		args.Page,
		args.PageSize,
	)
	if err != nil {
		Logger.Errorf("Failed to get file list for user %d: %v", currentUser.ID, err)
		return response.ErrUnknownError()
	}

	return response.OkWithData(ctx, models.SimpleFileInfoListResp{
		Total: total,
		Files: files,
	})
}

// getSingleDetailedFileInfo godoc
// @Summary      Get File Info
// @Description  Get detailed information about a specific file
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        file_id path int true "File ID"
// @Success      200 {object} response.ResponseBase[models.DetailedFileInfo] "File info retrieved successfully"
// @Failure      400 {object} response.ResponseBase[any] "Missing or invalid file_id parameter"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      404 {object} response.ResponseBase[any] "File not found"
// @Router       /file/info/{file_id} [get]
func (this *fileApi) getSingleDetailedFileInfo(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}

	args, err := utils.BindAndValidate[models.DetailedFileInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	switch fileInfo, err := fileService.GetFileInfoByFileID(ctx.Request().Context(), args.ID, currentUser.ID); {
	case err == nil:
		return response.OkWithData(ctx, fileInfo)
	case errors.Is(err, service.ErrNotFound):
		return response.ErrFileNotFound()
	default:
		Logger.Errorf("Failed to get file with ID %d: %v", args.ID, err)
		return response.ErrUnknownError()
	}
}

// getDownloadFileURL godoc
// @Summary      Get File Download URL
// @Description  Get a presigned download URL for a file from MinIO
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        file_id path int true "File ID"
// @Success      200 {object} response.ResponseBase[models.FileDownloadResp] "Download URL generated successfully"
// @Failure      400 {object} response.ResponseBase[any] "Missing or invalid file_id parameter"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      404 {object} response.ResponseBase[any] "File not found"
// @Router       /file/download/{file_id} [get]
func (this *fileApi) getDownloadFileURL(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}

	args, err := utils.BindAndValidate[models.DetailedFileInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	// Get file path from database
	filePath, err := fileService.GetFilePathByFileID(ctx.Request().Context(), args.ID, currentUser.ID)
	if err != nil {
		Logger.Errorf("Failed to get file path for file ID %d: %v", args.ID, err)
		return response.ErrFileNotFound()
	}

	// Get presigned download URL from MinIO
	downloadURL, err := fileService.GetDownloadURLByFilePath(ctx.Request().Context(), filePath)
	if err != nil {
		Logger.Errorf("Failed to get download URL for file: %v", err)
		return response.ErrUnknownError()
	}

	return response.OkWithData(ctx, models.FileDownloadResp{
		URL: downloadURL,
	})
}

// deleteFile godoc
// @Summary      Delete File
// @Description  Delete a file from both MinIO storage and database
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        file_id path int true "File ID"
// @Success      200 {object} response.ResponseBase[any] "File deleted successfully"
// @Failure      400 {object} response.ResponseBase[any] "Missing or invalid file_id parameter"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      403 {object} response.ResponseBase[any] "Forbidden, user not authorized to access the file"
// @Failure      404 {object} response.ResponseBase[any] "File not found"
// @Router       /file/delete/{file_id} [post]
func (this *fileApi) deleteFile(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.ErrInvalidToken()
	}

	args, err := utils.BindAndValidate[models.DetailedFileInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	// Get file path from database (also validates ownership)
	filePath, err := fileService.GetFilePathByFileID(ctx.Request().Context(), args.ID, currentUser.ID)
	switch err {
	case nil:
		//ok
	case service.ErrNotFound:
		return response.ErrFileNotFound()
	default:
		Logger.Errorf("Failed to get file path for file ID %d: %v", args.ID, err)
		return response.ErrUnknownError()
	}

	// Delete file from MinIO and database
	if err := fileService.DeleteFileByFileID(ctx.Request().Context(), args.ID, filePath); err != nil {
		Logger.Errorf("Failed to delete file with ID %d: %v", args.ID, err)
		return response.ErrUnknownError()
	}

	return response.Ok(ctx)
}
