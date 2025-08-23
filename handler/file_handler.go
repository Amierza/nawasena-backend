package handler

import (
	"mime/multipart"
	"net/http"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/Amierza/nawasena-backend/service"
	"github.com/gin-gonic/gin"
)

type (
	IFileHandler interface {
		Upload(ctx *gin.Context)
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

func (fh *fileHandler) Upload(ctx *gin.Context) {
	// coba ambil semua files dari multipart form
	form, err := ctx.MultipartForm()
	var files []*multipart.FileHeader

	if err == nil && form.File != nil && len(form.File["files"]) > 0 {
		// kalau user upload banyak file (key = "files")
		files = form.File["files"]
	} else {
		// kalau user upload single file (key = "file")
		file, err := ctx.FormFile("file")
		if err != nil {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_NO_FILES_UPLOADED, "no file(s) uploaded", nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}
		files = []*multipart.FileHeader{file}
	}

	// ambil folder dari query param (optional, default = "general")
	folder := ctx.DefaultQuery("folder", "")

	// call service
	uploadedURLs, err := fh.fileService.Upload(ctx, files, folder)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_UPLOAD_FILES, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// kalau hanya 1 file → balikin string saja
	if len(uploadedURLs) == 1 {
		res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPLOAD_FILE, uploadedURLs[0])
		ctx.JSON(http.StatusOK, res)
		return
	}

	// kalau banyak file → balikin array
	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPLOAD_FILES, uploadedURLs)
	ctx.JSON(http.StatusOK, res)
}
