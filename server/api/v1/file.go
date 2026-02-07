package v1

import (
	"errors"
	"fmt"
	"server/config"
	"server/middleware"
	"server/models"
	"server/models/response"
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

}

type fileApi struct{}

// uploadFile godoc
// @Summary      File Upload
// @Description  Upload a file to the server and associate it with a dataset. The file is stored in MinIO object storage.
// @Tags         File
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "File to upload"
// @Param        id formData int true "Dataset ID to associate the uploaded file with"
// @Success      200 {object} response.ResponseBase[models.FileUploadResp] "File uploaded successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid file, missing parameters, or upload failed"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      403 {object} response.ResponseBase[any] "Forbidden, user not authorized to access the specified dataset"
// @Router       /file/upload [post]
func (this *fileApi) uploadFile(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		utils.Logger.Error(err)
		return response.BadRequestWithMsg("Failed to get uploaded file")
	}
	datasetID, err := echo.FormValue[uint](ctx, "id")
	if err != nil {
		utils.Logger.Error(err)
		return response.BadRequestWithMsg("Failed to get dataset ID")
	}

	// Open the uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		utils.Logger.Errorf("Failed to open uploaded file: %v", err)
		return response.BadRequestWithMsg("Failed to open uploaded file")
	}
	defer src.Close()

	// Get file type/MIME type
	fileType := fileHeader.Header.Get("Content-Type")
	if fileType == "" {
		fileType = "application/octet-stream"
	}

	// Call CreateFile service to upload and save
	switch err := fileService.CreateFile(
		ctx.Request().Context(),
		currentUser.ID,
		datasetID,
		fileHeader.Filename,
		fileType,
		src,
		fileHeader.Size,
	); {
	case err == nil:
		// ok
	case errors.Is(err, service.ErrNotFound) || errors.Is(err, service.ErrDuplicatedKey):
		return response.ForbiddenWithMsg(fmt.Sprintf("Unauthorized access to the dataset: %d", datasetID))
	default:
		utils.Logger.Errorf("Failed to create file: %v", err)
		return response.FailWithMsg(ctx, "Failed to upload file")
	}

	return response.OkWithData(ctx, models.FileUploadResp{
		OwnerID:   currentUser.ID,
		DatasetID: datasetID,
		Name:      fileHeader.Filename,
		Type:      fileType,
		Size:      fileHeader.Size,
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
		return response.NoAuthWithMsg(err.Error())
	}

	fileInfoList, err := utils.BindAndValidate[models.ListFilesReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	total, files, err := fileService.GetFileListByUserID(
		ctx.Request().Context(),
		currentUser.ID,
		fileInfoList.DatasetID,
		fileInfoList.Page,
		fileInfoList.PageSize,
	)
	if err != nil {
		return response.FailWithMsg(ctx, "Failed to get file list")
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
		return response.NoAuthWithMsg(err.Error())
	}

	fileInfoReq, err := utils.BindAndValidate[models.DetailedFileInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg("Missing file_id parameter")
	}

	fileInfo, err := fileService.GetFileInfoByFileID(ctx.Request().Context(), fileInfoReq.ID, currentUser.ID)

	switch err {
	case nil:
		return response.OkWithData(ctx, fileInfo)
	case service.ErrNotFound:
		return response.ForbiddenWithMsg(fmt.Sprintf("Unauthorized access to the file: %d", fileInfoReq.ID))
	default:
		utils.Logger.Errorf("Failed to get file with ID %d: %v", fileInfoReq.ID, err)
		return response.FailWithMsg(ctx, "Unknown error")
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
		return response.NoAuthWithMsg(err.Error())
	}

	fileInfoReq, err := utils.BindAndValidate[models.DetailedFileInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg("Missing file_id parameter")
	}

	// Get file path from database
	filePath, err := fileService.GetFilePathByFileID(ctx.Request().Context(), fileInfoReq.ID, currentUser.ID)
	if err != nil {
		utils.Logger.Errorf("Failed to get file path for file ID %d: %v", fileInfoReq.ID, err)
		return response.NotFoundWithMsg(err.Error())
	}

	// Get presigned download URL from MinIO
	downloadURL, err := fileService.GetDownloadURLByFilePath(ctx.Request().Context(), filePath)
	if err != nil {
		utils.Logger.Errorf("Failed to get download URL for file: %v", err)
		return response.FailWithMsg(ctx, "Failed to get download URL")
	}

	return response.OkWithData(ctx, models.FileDownloadResp{
		URL: downloadURL,
	})
}
