package handler

import (
	"net/http"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/Amierza/nawasena-backend/service"
	"github.com/gin-gonic/gin"
)

type (
	IFileHandler interface {
		UploadFiles(ctx *gin.Context)
	}

	fileHandler struct {
		fileService service.IFileService
	}
)

func NewFileHandler(fileService service.IFileService) *fileHandler {
	return &fileHandler{
		fileService: fileService,
	}
}

// UploadFiles godoc
//
//	@Summary		Upload one or multiple files
//	@Description	Upload single or multiple files
//	@Tags			Upload
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			files	formData	file										true	"Files to upload"
//	@Success		200		{object}	response.SwaggerResponseSuccess[[]string]	"List of uploaded file URLs"
//	@Failure		400		{object}	response.SwaggerResponseError				"Invalid input"
//	@Failure		401		{object}	response.SwaggerResponseError				"Unauthorized"
//	@Router			/uploads [post]
func (fh *fileHandler) UploadFiles(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PARSE_MULTIPART_FORM, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_FILES_IS_EMPTY, dto.MESSAGE_FAILED_NO_FILES_UPLOADED, nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	uploadedURLs, err := fh.fileService.UploadFiles(ctx, files)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_UPLOAD_FILES, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPLOAD_FILES, uploadedURLs)
	ctx.JSON(http.StatusOK, res)
}
