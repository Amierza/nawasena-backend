package handler

import (
	"net/http"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/Amierza/nawasena-backend/service"
	"github.com/gin-gonic/gin"
)

type (
	IAchievementHandler interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetDetail(ctx *gin.Context)
		GetFeatured(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	achievementHandler struct {
		achievementService service.IAchievementService
	}
)

func NewAchievementHandler(achievementService service.IAchievementService) *achievementHandler {
	return &achievementHandler{
		achievementService: achievementService,
	}
}

func (ah *achievementHandler) Create(ctx *gin.Context) {
	var payload dto.CreateAchievementRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.achievementService.Create(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_ACHIEVEMENT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_ACHIEVEMENT, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *achievementHandler) GetAll(ctx *gin.Context) {
	paginationParam := ctx.DefaultQuery("pagination", "true")
	usePagination := paginationParam != "false"

	if !usePagination {
		// Tanpa pagination
		result, err := ah.achievementService.GetAll(ctx)
		if err != nil {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_ACHIEVEMENT, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}

		res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_ACHIEVEMENT, result)
		ctx.JSON(http.StatusOK, res)
		return
	}

	var payload response.PaginationRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.achievementService.GetAllWithPagination(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_ACHIEVEMENT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: dto.MESSAGE_SUCCESS_GET_LIST_ACHIEVEMENT,
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, res)
}

func (ah *achievementHandler) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ah.achievementService.GetDetail(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DETAIL_ACHIEVEMENT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DETAIL_ACHIEVEMENT, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *achievementHandler) GetFeatured(ctx *gin.Context) {
	limit := ctx.Query("limit")
	result, err := ah.achievementService.GetFeatured(ctx, limit)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DETAIL_ACHIEVEMENT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DETAIL_ACHIEVEMENT, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *achievementHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var payload dto.UpdateAchievementRequest
	payload.ID = idStr
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.achievementService.Update(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ACHIEVEMENT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_ACHIEVEMENT, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *achievementHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ah.achievementService.Delete(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ACHIEVEMENT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_ACHIEVEMENT, result)
	ctx.JSON(http.StatusOK, res)
}
