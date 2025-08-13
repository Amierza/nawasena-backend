package handler

import (
	"net/http"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/Amierza/nawasena-backend/service"
	"github.com/gin-gonic/gin"
)

type (
	ICompetitionHandler interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetDetail(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	competitionHandler struct {
		competitionService service.ICompetitionService
	}
)

func NewCompetitionHandler(competitionService service.ICompetitionService) *competitionHandler {
	return &competitionHandler{
		competitionService: competitionService,
	}
}

func (ah *competitionHandler) Create(ctx *gin.Context) {
	var payload dto.CreateCompetitionRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.competitionService.Create(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_COMPETITION, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_COMPETITION, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *competitionHandler) GetAll(ctx *gin.Context) {
	paginationParam := ctx.DefaultQuery("pagination", "true")
	usePagination := paginationParam != "false"

	if !usePagination {
		// Tanpa pagination
		result, err := ah.competitionService.GetAll(ctx)
		if err != nil {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_COMPETITION, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}

		res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_COMPETITION, result)
		ctx.JSON(http.StatusOK, res)
		return
	}

	var payload response.PaginationRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.competitionService.GetAllWithPagination(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_COMPETITION, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: dto.MESSAGE_SUCCESS_GET_LIST_COMPETITION,
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, res)
}

func (ah *competitionHandler) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ah.competitionService.GetDetail(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DETAIL_COMPETITION, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DETAIL_COMPETITION, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *competitionHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var payload dto.UpdateCompetitionRequest
	payload.ID = idStr
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.competitionService.Update(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_COMPETITION, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_COMPETITION, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *competitionHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ah.competitionService.Delete(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_COMPETITION, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_COMPETITION, result)
	ctx.JSON(http.StatusOK, res)
}
