package v1

import (
	"server/config"
	"server/middleware"
	"server/models"
	"server/models/response"
	"server/utils"

	"github.com/labstack/echo/v5"
)

func SetFileRouter(e *echo.Echo) {

	fileRouterGroup := e.Group(config.API_V1+"/file", middleware.TokenMiddleware())

	fileHandler := &fileApi{}
	fileRouterGroup.POST("/upload", fileHandler.uploadFile)
	fileRouterGroup.GET("/list", fileHandler.getFileList)
	fileRouterGroup.GET("/info/:file_id", fileHandler.getFileInfo)
	fileRouterGroup.GET("/download/:file_id", fileHandler.getDownloadFileURL)

}

type fileApi struct{}

// uploadFile godoc
// @Summary      File Upload
// @Description  Upload a file to the server
// @Tags         File
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "File to upload"
// @Success      200 {object} response.ResponseBase[models.FileUploadResp] "File uploaded successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid file or upload failed"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Router       /file/upload [post]
func (this *fileApi) uploadFile(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	newFileInfo, err := utils.BindAndValidate[models.FileUploadReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg("Failed to get uploaded file")
	}

	// Open the uploaded file
	src, err := newFileInfo.File.Open()
	if err != nil {
		utils.Logger.Errorf("Failed to open uploaded file: %v", err)
		return response.BadRequestWithMsg("Failed to open uploaded file")
	}
	defer src.Close()

	// Get file type/MIME type
	fileType := newFileInfo.File.Header.Get("Content-Type")
	if fileType == "" {
		fileType = "application/octet-stream"
	}

	// Call CreateFile service to upload and save
	if err := fileService.CreateFile(
		ctx.Request().Context(),
		currentUser.ID,
		newFileInfo.File.Filename,
		fileType,
		src,
		newFileInfo.File.Size,
	); err != nil {
		utils.Logger.Errorf("Failed to create file: %v", err)
		return response.FailWithMsg(ctx, "Failed to upload file")
	}

	return response.OkWithData(ctx, models.FileUploadResp{
		OwnerID: currentUser.ID,
		Name:    newFileInfo.File.Filename,
		Type:    fileType,
		Size:    newFileInfo.File.Size,
	})
}

// getFileList godoc
// @Summary      Get File List
// @Description  Retrieve a paginated list of files for the current user
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        page query int true "Page number" minimum(1)
// @Param        page_size query int true "Number of files per page" minimum(1) maximum(100)
// @Success      200 {object} response.ResponseBase[models.FileInfoListResp] "File list retrieved successfully"
// @Failure      400 {object} response.ResponseBase[any] "Invalid request parameters"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Router       /file/list [get]
func (this *fileApi) getFileList(ctx *echo.Context) error {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	fileInfoList, err := utils.BindAndValidate[models.FileInfoListReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg(err.Error())
	}

	files, err := fileService.GetFileListByUserID(
		ctx.Request().Context(),
		currentUser.ID,
		fileInfoList.Page,
		fileInfoList.PageSize,
	)
	if err != nil {
		utils.Logger.Errorf("Failed to get file list: %v", err)
		return response.FailWithMsg(ctx, "Failed to get file list")
	}

	return response.OkWithData(ctx, models.FileInfoListResp{
		Total: int64(len(files)),
		Files: files,
	})
}

// getFileInfo godoc
// @Summary      Get File Info
// @Description  Get detailed information about a specific file
// @Tags         File
// @Accept       json
// @Produce      json
// @Param        file_id path int true "File ID"
// @Success      200 {object} response.ResponseBase[models.FileInfo] "File info retrieved successfully"
// @Failure      400 {object} response.ResponseBase[any] "Missing or invalid file_id parameter"
// @Failure      401 {object} response.ResponseBase[any] "Unauthorized, authentication token required"
// @Failure      404 {object} response.ResponseBase[any] "File not found"
// @Router       /file/info/{file_id} [get]
func (this *fileApi) getFileInfo(ctx *echo.Context) error {
	_, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	fileInfoReq, err := utils.BindAndValidate[models.FileInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg("Missing file_id parameter")
	}

	fileInfo, err := fileService.GetFileInfoByFileID(ctx.Request().Context(), fileInfoReq.ID)
	if err != nil {
		utils.Logger.Errorf("Failed to get file info: %v", err)
		return response.FailWithMsg(ctx, "Failed to get file info")
	}

	return response.OkWithData(ctx, fileInfo)
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
	_, err := utils.GetCurrentUser(ctx)
	if err != nil {
		return response.NoAuthWithMsg(err.Error())
	}

	fileInfoReq, err := utils.BindAndValidate[models.FileInfoReq](ctx)
	if err != nil {
		return response.BadRequestWithMsg("Missing file_id parameter")
	}

	// Get presigned download URL from MinIO
	downloadURL, err := fileService.GetDownloadURLByFileID(ctx.Request().Context(), fileInfoReq.ID)
	if err != nil {
		utils.Logger.Errorf("Failed to get download URL for file: %v", err)
		return response.FailWithMsg(ctx, "Failed to get download URL")
	}

	return response.OkWithData(ctx, models.FileDownloadResp{
		URL: downloadURL,
	})
}
