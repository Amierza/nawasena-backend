package service

import (
	"context"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/helper"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/repository"
)

type (
	IAuthService interface {
		Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
		RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.RefreshTokenResponse, error)
	}

	authService struct {
		authRepo   repository.IAuthRepository
		jwtService jwt.IJWTService
	}
)

func NewAuthService(authRepo repository.IAuthRepository, jwtService jwt.IJWTService) *authService {
	return &authService{
		authRepo:   authRepo,
		jwtService: jwtService,
	}
}

func (as *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	if req.Email == "" || !helper.IsValidEmail(req.Email) {
		return dto.LoginResponse{}, dto.ErrInvalidEmail
	}

	if req.Password == "" || len(req.Password) < 8 {
		return dto.LoginResponse{}, dto.ErrInvalidPassword
	}

	admin, found, err := as.authRepo.GetAdminByEmail(ctx, nil, req.Email)
	if err != nil {
		return dto.LoginResponse{}, dto.ErrGetAdminByEmail
	}
	if !found {
		return dto.LoginResponse{}, dto.ErrAdminNotFound
	}

	checkPassword, err := helper.CheckPassword(admin.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		return dto.LoginResponse{}, dto.ErrIncorrectPassword
	}

	accessToken, refreshToken, err := as.jwtService.GenerateToken(admin.ID.String())
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (as *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (dto.RefreshTokenResponse, error) {
	_, err := as.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		return dto.RefreshTokenResponse{}, dto.ErrValidateToken
	}

	adminID, err := as.jwtService.GetAdminIDByToken(req.RefreshToken)
	if err != nil {
		return dto.RefreshTokenResponse{}, dto.ErrGetAdminIDFromToken
	}

	accessToken, _, err := as.jwtService.GenerateToken(adminID)
	if err != nil {
		return dto.RefreshTokenResponse{}, dto.ErrGenerateAccessToken
	}

	return dto.RefreshTokenResponse{AccessToken: accessToken}, nil
}
