package handler

import (
	"net/http"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/Amierza/nawasena-backend/service"
	"github.com/gin-gonic/gin"
)

type (
	IAchievementCategoryHandler interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetDetail(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	achievementCategoryHandler struct {
		achievementCategoryService service.IAchievementCategoryService
	}
)

func NewAchievementCategoryHandler(achievementCategoryService service.IAchievementCategoryService) *achievementCategoryHandler {
	return &achievementCategoryHandler{
		achievementCategoryService: achievementCategoryService,
	}
}

func (ph *achievementCategoryHandler) Create(ctx *gin.Context) {
	var payload dto.CreateAchievementCategoryRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ph.achievementCategoryService.Create(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_ACHIEVEMENT_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_ACHIEVEMENT_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *achievementCategoryHandler) GetAll(ctx *gin.Context) {
	result, err := ph.achievementCategoryService.GetAll(ctx)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_ACHIEVEMENT_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_ACHIEVEMENT_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *achievementCategoryHandler) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ph.achievementCategoryService.GetDetail(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DETAIL_ACHIEVEMENT_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DETAIL_ACHIEVEMENT_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *achievementCategoryHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var payload dto.UpdateAchievementCategoryRequest
	payload.ID = idStr
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ph.achievementCategoryService.Update(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ACHIEVEMENT_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_ACHIEVEMENT_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}

func (ph *achievementCategoryHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ph.achievementCategoryService.Delete(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ACHIEVEMENT_CATEGORY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_ACHIEVEMENT_CATEGORY, result)
	ctx.JSON(http.StatusOK, res)
}
