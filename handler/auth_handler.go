package handler

import (
	"net/http"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/Amierza/nawasena-backend/service"
	"github.com/gin-gonic/gin"
)

type (
	IAuthHandler interface {
		Login(ctx *gin.Context)
		RefreshToken(ctx *gin.Context)
	}

	authHandler struct {
		authService service.IAuthService
	}
)

func NewAuthHandler(authService service.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}

// Login godoc
//
//	@Summary		Login admin
//	@Description	Authenticate admin and return access & refresh tokens
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.LoginRequest									true	"Login credentials"
//	@Success		200		{object}	response.SwaggerResponseSuccess[dto.LoginResponse]	"Success login"
//	@Failure		400		{object}	response.SwaggerResponseError						"Invalid input"
//	@Failure		404		{object}	response.SwaggerResponseError						"Admin account not found"
//	@Router			/login [post]
func (ah *authHandler) Login(ctx *gin.Context) {
	var payload dto.LoginRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.authService.Login(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_LOGIN_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_LOGIN_USER, result)
	ctx.JSON(http.StatusOK, res)
}

// RefreshToken godoc
//
//	@Summary		Refresh access token
//	@Description	Get a new access token using refresh token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.RefreshTokenRequest										true	"Refresh token payload"
//	@Success		200		{object}	response.SwaggerResponseSuccess[dto.RefreshTokenResponse]	"Success refresh access token"
//	@Failure		400		{object}	response.SwaggerResponseError								"Invalid <access_token>"
//	@Router			/refresh-token [post]
func (ah *authHandler) RefreshToken(ctx *gin.Context) {
	var payload dto.RefreshTokenRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := ah.authService.RefreshToken(ctx, payload)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_REFRESH_TOKEN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REFRESH_TOKEN, result)
	ctx.AbortWithStatusJSON(http.StatusOK, res)
}
