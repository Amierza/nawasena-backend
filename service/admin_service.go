package service

import (
	"context"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/entity"
	"github.com/Amierza/nawasena-backend/helper"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/repository"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/google/uuid"
)

type (
	IAdminService interface {
		Create(ctx context.Context, req dto.CreateAdminRequest) (dto.AdminResponse, error)
		GetAll(ctx context.Context) ([]dto.AdminResponse, error)
		GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.AdminPaginationResponse, error)
		GetDetail(ctx context.Context, id string) (dto.AdminResponse, error)
		Update(ctx context.Context, req dto.UpdateAdminRequest) (dto.AdminResponse, error)
		Delete(ctx context.Context, id string) (dto.AdminResponse, error)
	}

	adminService struct {
		adminRepo repository.IAdminRepository
		jwt       jwt.IJWT
	}
)

func NewAdminService(adminRepo repository.IAdminRepository, jwt jwt.IJWT) *adminService {
	return &adminService{
		adminRepo: adminRepo,
		jwt:       jwt,
	}
}

func (as *adminService) Create(ctx context.Context, req dto.CreateAdminRequest) (dto.AdminResponse, error) {
	// handle name request
	if req.Name == "" {
		return dto.AdminResponse{}, dto.ErrEmptyName
	}
	if len(req.Name) < 3 {
		return dto.AdminResponse{}, dto.ErrNameTooShort
	}

	// handle email request
	if req.Email == "" {
		return dto.AdminResponse{}, dto.ErrEmptyEmail
	}
	if !helper.IsValidEmail(req.Email) {
		return dto.AdminResponse{}, dto.ErrInvalidEmail
	}
	_, found, _ := as.adminRepo.GetByEmail(ctx, nil, req.Email)
	if found {
		return dto.AdminResponse{}, dto.ErrEmailAlreadyExists
	}

	// handle password request
	if req.Password == "" {
		return dto.AdminResponse{}, dto.ErrEmptyPassword
	}
	if len(req.Password) < 8 {
		return dto.AdminResponse{}, dto.ErrInvalidPassword
	}

	// handle phone number request
	if req.PhoneNumber == "" {
		return dto.AdminResponse{}, dto.ErrEmptyPhoneNumber
	}

	role := "admin"
	id := uuid.New()
	admin := &entity.Admin{
		ID:          id,
		Name:        req.Name,
		Email:       req.Email,
		Password:    req.Password,
		Role:        entity.Role(role),
		PhoneNumber: req.PhoneNumber,
	}

	err := as.adminRepo.Create(ctx, nil, admin)
	if err != nil {
		return dto.AdminResponse{}, dto.ErrCreateAdmin
	}

	return dto.AdminResponse{
		ID:          admin.ID.String(),
		Name:        admin.Name,
		Email:       admin.Email,
		Password:    admin.Password,
		Role:        admin.Role,
		PhoneNumber: admin.PhoneNumber,
	}, nil
}

func (as *adminService) GetAll(ctx context.Context) ([]dto.AdminResponse, error) {
	admins, err := as.adminRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, dto.ErrGetAllAdminNoPagination
	}

	var datas []dto.AdminResponse
	for _, admin := range admins {
		data := dto.AdminResponse{
			ID:          admin.ID.String(),
			Name:        admin.Name,
			Email:       admin.Email,
			Password:    admin.Password,
			Role:        admin.Role,
			PhoneNumber: admin.PhoneNumber,
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (as *adminService) GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.AdminPaginationResponse, error) {
	dataWithPaginate, err := as.adminRepo.GetAllWithPagination(ctx, nil, req)
	if err != nil {
		return dto.AdminPaginationResponse{}, dto.ErrGetAllAdminWithPagination
	}

	var datas []dto.AdminResponse
	for _, admin := range dataWithPaginate.Admins {
		data := dto.AdminResponse{
			ID:          admin.ID.String(),
			Name:        admin.Name,
			Email:       admin.Email,
			Password:    admin.Password,
			Role:        admin.Role,
			PhoneNumber: admin.PhoneNumber,
		}

		datas = append(datas, data)
	}

	return dto.AdminPaginationResponse{
		Data: datas,
		PaginationResponse: response.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (as *adminService) GetDetail(ctx context.Context, id string) (dto.AdminResponse, error) {
	admin, _, err := as.adminRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.AdminResponse{}, dto.ErrAdminNotFound
	}

	return dto.AdminResponse{
		ID:          admin.ID.String(),
		Name:        admin.Name,
		Email:       admin.Email,
		Password:    admin.Password,
		Role:        admin.Role,
		PhoneNumber: admin.PhoneNumber,
	}, nil
}

func (as *adminService) Update(ctx context.Context, req dto.UpdateAdminRequest) (dto.AdminResponse, error) {
	// get admin from db
	admin, flag, err := as.adminRepo.GetByID(ctx, nil, req.ID)
	if err != nil || !flag {
		return dto.AdminResponse{}, dto.ErrAdminNotFound
	}

	// handle name request
	if req.Name != "" && req.Name != admin.Name {
		if len(req.Name) < 3 {
			return dto.AdminResponse{}, dto.ErrNameTooShort
		}

		admin.Name = req.Name
	}

	// handle email request
	if req.Email != "" && req.Email != admin.Email {
		if !helper.IsValidEmail(req.Email) {
			return dto.AdminResponse{}, dto.ErrInvalidEmail
		}

		_, found, _ := as.adminRepo.GetByEmail(ctx, nil, req.Email)
		if found {
			return dto.AdminResponse{}, dto.ErrEmailAlreadyExists
		}

		admin.Email = req.Email
	}

	// handle password request
	if req.Password != "" {
		if len(req.Password) < 8 {
			return dto.AdminResponse{}, dto.ErrInvalidPassword
		}

		hashP, err := helper.HashPassword(req.Password)
		if err != nil {
			return dto.AdminResponse{}, dto.ErrHashPassword
		}

		admin.Password = hashP
	}

	// handle phone number request
	if req.PhoneNumber != "" && req.PhoneNumber != admin.PhoneNumber {
		formattedPhoneNumber, err := helper.StandardizePhoneNumber(req.PhoneNumber)
		if err != nil {
			return dto.AdminResponse{}, dto.ErrFormatPhoneNumber
		}

		admin.PhoneNumber = formattedPhoneNumber
	}

	err = as.adminRepo.Update(ctx, nil, admin)
	if err != nil {
		return dto.AdminResponse{}, dto.ErrUpdateAdmin
	}

	res := dto.AdminResponse{
		ID:          admin.ID.String(),
		Name:        admin.Name,
		Email:       admin.Email,
		Password:    admin.Password,
		Role:        admin.Role,
		PhoneNumber: admin.PhoneNumber,
	}

	return res, nil
}
func (as *adminService) Delete(ctx context.Context, id string) (dto.AdminResponse, error) {
	deletedAdmin, found, err := as.adminRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.AdminResponse{}, dto.ErrAdminNotFound
	}
	if !found {
		return dto.AdminResponse{}, dto.ErrAdminNotFound
	}

	err = as.adminRepo.DeleteByID(ctx, nil, id)
	if err != nil {
		return dto.AdminResponse{}, dto.ErrDeleteAdminByID
	}

	res := dto.AdminResponse{
		ID:          deletedAdmin.ID.String(),
		Name:        deletedAdmin.Name,
		Email:       deletedAdmin.Email,
		Password:    deletedAdmin.Password,
		Role:        deletedAdmin.Role,
		PhoneNumber: deletedAdmin.PhoneNumber,
	}

	return res, nil
}
