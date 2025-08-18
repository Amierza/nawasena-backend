package handler

import (
	"net/http"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/Amierza/nawasena-backend/service"
	"github.com/gin-gonic/gin"
)

type (
	INewsCategoryHandler interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetDetail(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	newsCategoryHandler struct {
		newsCategoryService service.INewsCategoryService
	}
)

func NewNewsCategoryHandler(newsCategoryService service.INewsCategoryService) *newsCategoryHandler {
	return &newsCategoryHandler{
		newsCategoryService: newsCategoryService,
	}
}

func (ph *newsCategoryHandler) Create(ctx *gin.Context) {
	var payload dto.CreateNewsCategoryRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ph.newsCategoryService.Create(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_NEWS_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_NEWS_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *newsCategoryHandler) GetAll(ctx *gin.Context) {
	result, err := ph.newsCategoryService.GetAll(ctx)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_NEWS_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_NEWS_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *newsCategoryHandler) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ph.newsCategoryService.GetDetail(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DETAIL_NEWS_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DETAIL_NEWS_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *newsCategoryHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var payload dto.UpdateNewsCategoryRequest
	payload.ID = idStr
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ph.newsCategoryService.Update(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_NEWS_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_NEWS_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *newsCategoryHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ph.newsCategoryService.Delete(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_NEWS_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_NEWS_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}
