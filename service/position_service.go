package service

import (
	"context"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/entity"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/repository"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/google/uuid"
)

type (
	IPositionService interface {
		Create(ctx context.Context, req dto.CreatePositionRequest) (dto.PositionResponse, error)
		GetAll(ctx context.Context) ([]dto.PositionResponse, error)
		GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.PositionPaginationResponse, error)
		GetDetail(ctx context.Context, id string) (dto.PositionResponse, error)
		Update(ctx context.Context, req dto.UpdatePositionRequest) (dto.PositionResponse, error)
		Delete(ctx context.Context, id string) (dto.PositionResponse, error)
	}

	positionService struct {
		positionRepo repository.IPositionRepository
		jwt          jwt.IJWT
	}
)

func NewPositionService(positionRepo repository.IPositionRepository, jwt jwt.IJWT) *positionService {
	return &positionService{
		positionRepo: positionRepo,
		jwt:          jwt,
	}
}

func (ps *positionService) Create(ctx context.Context, req dto.CreatePositionRequest) (dto.PositionResponse, error) {
	// handle name request
	if req.Name == "" {
		return dto.PositionResponse{}, dto.ErrEmptyName
	}
	if len(req.Name) < 3 {
		return dto.PositionResponse{}, dto.ErrNameTooShort
	}
	_, found, _ := ps.positionRepo.GetByName(ctx, nil, req.Name)
	if found {
		return dto.PositionResponse{}, dto.ErrPositionAlreadyExists
	}

	id := uuid.New()
	position := &entity.Position{
		ID:     id,
		Name:   req.Name,
		IsTech: req.IsTech,
	}

	err := ps.positionRepo.Create(ctx, nil, position)
	if err != nil {
		return dto.PositionResponse{}, dto.ErrCreatePosition
	}

	return dto.PositionResponse{
		ID:     position.ID.String(),
		Name:   position.Name,
		IsTech: position.IsTech,
	}, nil
}

func (ps *positionService) GetAll(ctx context.Context) ([]dto.PositionResponse, error) {
	positions, err := ps.positionRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, dto.ErrGetAllPositionNoPagination
	}

	var datas []dto.PositionResponse
	for _, position := range positions {
		data := dto.PositionResponse{
			ID:     position.ID.String(),
			Name:   position.Name,
			IsTech: position.IsTech,
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (ps *positionService) GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.PositionPaginationResponse, error) {
	dataWithPaginate, err := ps.positionRepo.GetAllWithPagination(ctx, nil, req)
	if err != nil {
		return dto.PositionPaginationResponse{}, dto.ErrGetAllPositionWithPagination
	}

	var datas []dto.PositionResponse
	for _, position := range dataWithPaginate.Positions {
		data := dto.PositionResponse{
			ID:     position.ID.String(),
			Name:   position.Name,
			IsTech: position.IsTech,
		}

		datas = append(datas, data)
	}

	return dto.PositionPaginationResponse{
		Data: datas,
		PaginationResponse: response.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (ps *positionService) GetDetail(ctx context.Context, id string) (dto.PositionResponse, error) {
	position, _, err := ps.positionRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.PositionResponse{}, dto.ErrPositionNotFound
	}

	return dto.PositionResponse{
		ID:     position.ID.String(),
		Name:   position.Name,
		IsTech: position.IsTech,
	}, nil
}

func (ps *positionService) Update(ctx context.Context, req dto.UpdatePositionRequest) (dto.PositionResponse, error) {
	// get position from db
	position, flag, err := ps.positionRepo.GetByID(ctx, nil, req.ID)
	if err != nil || !flag {
		return dto.PositionResponse{}, dto.ErrPositionNotFound
	}

	// handle name request
	if req.Name != "" && req.Name != position.Name {
		if len(req.Name) < 3 {
			return dto.PositionResponse{}, dto.ErrNameTooShort
		}

		_, found, _ := ps.positionRepo.GetByName(ctx, nil, req.Name)
		if found {
			return dto.PositionResponse{}, dto.ErrPositionAlreadyExists
		}

		position.Name = req.Name
	}

	if req.IsTech != nil {
		position.IsTech = *req.IsTech
	}

	if err := ps.positionRepo.Update(ctx, nil, position); err != nil {
		return dto.PositionResponse{}, dto.ErrUpdatePosition
	}

	res := dto.PositionResponse{
		ID:     position.ID.String(),
		Name:   position.Name,
		IsTech: position.IsTech,
	}

	return res, nil
}

func (ps *positionService) Delete(ctx context.Context, id string) (dto.PositionResponse, error) {
	deletedPosition, found, err := ps.positionRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.PositionResponse{}, dto.ErrPositionNotFound
	}
	if !found {
		return dto.PositionResponse{}, dto.ErrPositionNotFound
	}

	err = ps.positionRepo.DeleteByID(ctx, nil, id)
	if err != nil {
		return dto.PositionResponse{}, dto.ErrDeletePositionByID
	}

	res := dto.PositionResponse{
		ID:     deletedPosition.ID.String(),
		Name:   deletedPosition.Name,
		IsTech: deletedPosition.IsTech,
	}

	return res, nil
}
