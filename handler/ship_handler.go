package handler

import (
	"net/http"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/Amierza/nawasena-backend/service"
	"github.com/gin-gonic/gin"
)

type (
	IShipHandler interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetDetail(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	shipHandler struct {
		shipService service.IShipService
	}
)

func NewShipHandler(shipService service.IShipService) *shipHandler {
	return &shipHandler{
		shipService: shipService,
	}
}

func (ah *shipHandler) Create(ctx *gin.Context) {
	var payload dto.CreateShipRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.shipService.Create(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_SHIP, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_SHIP, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *shipHandler) GetAll(ctx *gin.Context) {
	paginationParam := ctx.DefaultQuery("pagination", "true")
	usePagination := paginationParam != "false"

	if !usePagination {
		// Tanpa pagination
		result, err := ah.shipService.GetAll(ctx)
		if err != nil {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_SHIP, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}

		res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_SHIP, result)
		ctx.JSON(http.StatusOK, res)
		return
	}

	var payload response.PaginationRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.shipService.GetAllWithPagination(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_SHIP, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: dto.MESSAGE_SUCCESS_GET_LIST_SHIP,
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, res)
}

func (ah *shipHandler) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ah.shipService.GetDetail(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DETAIL_SHIP, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DETAIL_SHIP, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *shipHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var payload dto.UpdateShipRequest
	payload.ID = idStr
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.shipService.Update(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_SHIP, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_SHIP, result)
	ctx.JSON(http.StatusOK, res)
}

func (ah *shipHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ah.shipService.Delete(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_SHIP, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_SHIP, result)
	ctx.JSON(http.StatusOK, res)
}
