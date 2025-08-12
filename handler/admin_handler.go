package handler

import (
	"net/http"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/Amierza/nawasena-backend/service"
	"github.com/gin-gonic/gin"
)

type (
	IAdminHandler interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetDetail(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	adminHandler struct {
		adminService service.IAdminService
	}
)

func NewAdminHandler(adminService service.IAdminService) *adminHandler {
	return &adminHandler{
		adminService: adminService,
	}
}

// CreateAdminForSuperAdmin godoc
//
//	@Summary		Create new admin (only for super admin)
//	@Description	Create a new admin account only for super admin
//	@Tags			Admin
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreateAdminRequest								true	"Admin data"
//	@Success		200		{object}	response.SwaggerResponseSuccess[dto.AdminResponse]	"Success create admin"
//	@Failure		400		{object}	response.SwaggerResponseError						"Invalid input"
//	@Failure		401		{object}	response.SwaggerResponseError						"Unauthorized"
//	@Router			/admins [post]
func (ah *adminHandler) Create(ctx *gin.Context) {
	var payload dto.CreateAdminRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.adminService.Create(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_ADMIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_ADMIN, result)
	ctx.JSON(http.StatusOK, res)
}

// GetAllAdmin godoc
//
//	@Summary		Get all admins
//	@Description	Get list of all admins (with or without pagination)
//	@Tags			Admin
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			search		query		string															false	"Search keyword if with pagination"
//	@Param			page		query		int																false	"Page number if with pagination"
//	@Param			limit		query		int																false	"Items per page if with pagination"
//	@Param			pagination	query		bool															false	"With pagination or not (default true)"
//	@Success		200			{object}	response.SwaggerResponseSuccess[[]dto.AdminResponse]			"Without pagination"
//	@Success		200			{object}	response.SwaggerResponseSuccess[dto.AdminPaginationResponse]	"With pagination"
//	@Failure		400			{object}	response.SwaggerResponseError									"Invalid input"
//	@Router			/admins [get]
func (ah *adminHandler) GetAll(ctx *gin.Context) {
	paginationParam := ctx.DefaultQuery("pagination", "true")
	usePagination := paginationParam != "false"

	if !usePagination {
		// Tanpa pagination
		result, err := ah.adminService.GetAll(ctx)
		if err != nil {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_ADMIN, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
			return
		}

		res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_LIST_ADMIN, result)
		ctx.JSON(http.StatusOK, res)
		return
	}

	var payload response.PaginationRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.adminService.GetAllWithPagination(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_LIST_ADMIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: dto.MESSAGE_SUCCESS_GET_LIST_ADMIN,
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}

	ctx.JSON(http.StatusOK, res)
}

// GetDetailAdmin godoc
//
//	@Summary		Get Admin Detail
//	@Description	Retrieve admin details by ID
//	@Tags			Admin
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Admin ID"
//	@Success		200	{object}	response.SwaggerResponseSuccess[dto.AdminResponse]
//	@Failure		400	{object}	response.SwaggerResponseError	"Invalid input"
//	@Failure		404	{object}	response.SwaggerResponseError	"Admin not found"
//	@Router			/admins/{id} [get]
func (ah *adminHandler) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ah.adminService.GetDetail(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DETAIL_ADMIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DETAIL_ADMIN, result)
	ctx.JSON(http.StatusOK, res)
}

// UpdateAdmin godoc
//
//	@Summary		Update Admin
//	@Description	Update admin information by ID
//	@Tags			Admin
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Admin ID"
//	@Param			payload	body		dto.UpdateAdminRequest	true	"Update Admin Payload"
//	@Success		200		{object}	response.SwaggerResponseSuccess[dto.AdminResponse]
//	@Failure		400		{object}	response.SwaggerResponseError	"Invalid input"
//	@Failure		404		{object}	response.SwaggerResponseError	"Admin not found"
//	@Router			/admins/{id} [patch]
func (ah *adminHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var payload dto.UpdateAdminRequest
	payload.ID = idStr
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.adminService.Update(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_UPDATE_ADMIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_UPDATE_ADMIN, result)
	ctx.JSON(http.StatusOK, res)
}

// DeleteAdmin godoc
//
//	@Summary		Delete Admin
//	@Description	Delete admin by ID
//	@Tags			Admin
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Admin ID"
//	@Success		200	{object}	response.SwaggerResponseSuccess[dto.AdminResponse]
//	@Failure		400	{object}	response.SwaggerResponseError	"Invalid request"
//	@Failure		404	{object}	response.SwaggerResponseError	"Admin not found"
//	@Router			/admins/{id} [delete]
func (ah *adminHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := ah.adminService.Delete(ctx, idStr)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_DELETE_ADMIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_ADMIN, result)
	ctx.JSON(http.StatusOK, res)
}
