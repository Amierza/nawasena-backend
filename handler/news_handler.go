package handler

import (
	"net/http"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/Amierza/nawasena-backend/service"
	"github.com/gin-gonic/gin"
)

type (
	INewsHandler interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetFeatured(ctx *gin.Context)
		GetDetail(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	newsHandler struct {
		newsService service.INewsService
	}
)

func NewNewsHandler(newsService service.INewsService) *newsHandler {
	return &newsHandler{
		newsService: newsService,
	}
}

func (ah *newsHandler) Create(ctx *gin.Context) {
	var payload dto.CreateNewsRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.newsService.Create(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_NEWS, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_NEWS, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *newsHandler) GetAll(ctx *gin.Context) {
	paginationParam := ctx.DefaultQuery("pagination", "true")
	usePagination := paginationParam != "false"

	if !usePagination {
		// Tanpa pagination
		result, err := ah.newsService.GetAll(ctx)
		if err != nil {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_NEWS, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}

		res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_NEWS, result)
		ctx.JSON(http.StatusOK, res)
		return
	}

	var payload response.PaginationRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.newsService.GetAllWithPagination(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_NEWS, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: dto.MESSAGE_SUCCESS_GET_LIST_NEWS,
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, res)
}

func (ah *newsHandler) GetFeatured(ctx *gin.Context) {
	limit := ctx.Query("limit")
	result, err := ah.newsService.GetFeatured(ctx, limit)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DETAIL_NEWS, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DETAIL_NEWS, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *newsHandler) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ah.newsService.GetDetail(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DETAIL_NEWS, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DETAIL_NEWS, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *newsHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var payload dto.UpdateNewsRequest
	payload.ID = idStr
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.newsService.Update(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_NEWS, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_NEWS, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *newsHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ah.newsService.Delete(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_NEWS, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_NEWS, result)
	ctx.JSON(http.StatusOK, res)
}
