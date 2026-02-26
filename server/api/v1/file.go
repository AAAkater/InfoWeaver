package v1

import (
	"errors"
	"server/config"
	"server/middleware"
	"server/models"
	"server/models/common/response"
	"server/service"
	"server/utils"

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
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters or no files provided"
// @Failure      401 {object} response.ResponseBase[any] "Invalid or expired token"
// @Failure      403 {object} response.ResponseBase[any] "Dataset not found or access denied"
// @Failure      404 {object} response.ResponseBase[any] "Dataset not found"
// @Failure      500 {object} response.ResponseBase[any] "Internal server error"
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

	uploadedFiles, errors := fileService.UploadFile(ctx.Request().Context(), fileHeaders, fileNumber, currentUser.ID, datasetID)
	// If all files failed, return error
	if len(uploadedFiles) == 0 && len(errors) > 0 {
		Logger.Errorf("All file uploads failed: %v", errors)
		return response.ErrUnknownError()
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
// @Param        dataset_id query int false "Filter by dataset ID"
// @Success      200 {object} response.ResponseBase[models.SimpleFileInfoListResp] "File list retrieved successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      401 {object} response.ResponseBase[any] "Invalid or expired token"
// @Failure      500 {object} response.ResponseBase[any] "Internal server error"
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
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      401 {object} response.ResponseBase[any] "Invalid or expired token"
// @Failure      404 {object} response.ResponseBase[any] "File not found"
// @Failure      500 {object} response.ResponseBase[any] "Internal server error"
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
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      401 {object} response.ResponseBase[any] "Invalid or expired token"
// @Failure      404 {object} response.ResponseBase[any] "File not found"
// @Failure      500 {object} response.ResponseBase[any] "Internal server error"
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
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      401 {object} response.ResponseBase[any] "Invalid or expired token"
// @Failure      404 {object} response.ResponseBase[any] "File not found"
// @Failure      500 {object} response.ResponseBase[any] "Internal server error"
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
