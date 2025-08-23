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
	IFlyerService interface {
		Create(ctx context.Context, req dto.CreateFlyerRequest) (dto.FlyerResponse, error)
		GetAll(ctx context.Context) ([]dto.FlyerResponse, error)
		GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.FlyerPaginationResponse, error)
		GetDetail(ctx context.Context, id string) (dto.FlyerResponse, error)
		Update(ctx context.Context, req dto.UpdateFlyerRequest) (dto.FlyerResponse, error)
		Delete(ctx context.Context, id string) (dto.FlyerResponse, error)
	}

	flyerService struct {
		flyerRepo repository.IFlyerRepository
		jwt       jwt.IJWT
	}
)

func NewFlyerService(flyerRepo repository.IFlyerRepository, jwt jwt.IJWT) *flyerService {
	return &flyerService{
		flyerRepo: flyerRepo,
		jwt:       jwt,
	}
}

func (as *flyerService) Create(ctx context.Context, req dto.CreateFlyerRequest) (dto.FlyerResponse, error) {
	// handle name request
	if req.Name == "" {
		return dto.FlyerResponse{}, dto.ErrEmptyName
	}
	if len(req.Name) < 3 {
		return dto.FlyerResponse{}, dto.ErrNameTooShort
	}

	// handle image request
	if req.Image == "" {
		return dto.FlyerResponse{}, dto.ErrEmptyImage
	}

	// handle double data
	_, found, _ := as.flyerRepo.GetByName(ctx, nil, req.Name)
	if found {
		return dto.FlyerResponse{}, dto.ErrFlyerAlreadyExists
	}

	flyerID := uuid.New()
	flyer := &entity.Flyer{
		ID:    flyerID,
		Name:  req.Name,
		Image: req.Image,
	}

	// create flyer
	if err := as.flyerRepo.Create(ctx, nil, flyer); err != nil {
		return dto.FlyerResponse{}, dto.ErrCreateFlyer
	}

	return dto.FlyerResponse{
		ID:    flyer.ID.String(),
		Name:  flyer.Name,
		Image: flyer.Image,
	}, nil
}

func (as *flyerService) GetAll(ctx context.Context) ([]dto.FlyerResponse, error) {
	flyers, err := as.flyerRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, dto.ErrGetAllFlyerNoPagination
	}

	var datas []dto.FlyerResponse
	for _, flyer := range flyers {
		data := dto.FlyerResponse{
			ID:    flyer.ID.String(),
			Name:  flyer.Name,
			Image: flyer.Image,
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (as *flyerService) GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.FlyerPaginationResponse, error) {
	dataWithPaginate, err := as.flyerRepo.GetAllWithPagination(ctx, nil, req)
	if err != nil {
		return dto.FlyerPaginationResponse{}, dto.ErrGetAllFlyerWithPagination
	}

	var datas []dto.FlyerResponse
	for _, flyer := range dataWithPaginate.Flyers {
		data := dto.FlyerResponse{
			ID:    flyer.ID.String(),
			Name:  flyer.Name,
			Image: flyer.Image,
		}

		datas = append(datas, data)
	}

	return dto.FlyerPaginationResponse{
		Data: datas,
		PaginationResponse: response.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (as *flyerService) GetDetail(ctx context.Context, id string) (dto.FlyerResponse, error) {
	flyer, _, err := as.flyerRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.FlyerResponse{}, dto.ErrFlyerNotFound
	}

	res := dto.FlyerResponse{
		ID:    flyer.ID.String(),
		Name:  flyer.Name,
		Image: flyer.Image,
	}

	return res, nil
}

func (as *flyerService) Update(ctx context.Context, req dto.UpdateFlyerRequest) (dto.FlyerResponse, error) {
	// get flyer by id
	flyer, found, err := as.flyerRepo.GetByID(ctx, nil, req.ID)
	if err != nil {
		return dto.FlyerResponse{}, dto.ErrGetFlyerByID
	}
	if !found {
		return dto.FlyerResponse{}, dto.ErrFlyerNotFound
	}

	// handle name request
	if req.Name != "" && req.Name != flyer.Name {
		if len(req.Name) < 3 {
			return dto.FlyerResponse{}, dto.ErrNameTooShort
		}

		flyer.Name = req.Name
	}

	err = as.flyerRepo.RunInTransaction(ctx, func(txRepo repository.IFlyerRepository) error {
		// update flyer
		if err := txRepo.Update(ctx, nil, flyer); err != nil {
			return dto.ErrUpdateFlyer
		}

		return nil
	})
	if err != nil {
		return dto.FlyerResponse{}, err
	}

	return dto.FlyerResponse{
		ID:    flyer.ID.String(),
		Name:  flyer.Name,
		Image: flyer.Image,
	}, nil
}

func (as *flyerService) Delete(ctx context.Context, id string) (dto.FlyerResponse, error) {
	deletedFlyer, found, err := as.flyerRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.FlyerResponse{}, dto.ErrFlyerNotFound
	}
	if !found {
		return dto.FlyerResponse{}, dto.ErrFlyerNotFound
	}

	err = as.flyerRepo.RunInTransaction(ctx, func(txRepo repository.IFlyerRepository) error {
		// Delete Flyer
		err = as.flyerRepo.DeleteByID(ctx, nil, id)
		if err != nil {
			return dto.ErrDeleteFlyerByID
		}

		return nil
	})
	if err != nil {
		return dto.FlyerResponse{}, err
	}

	res := dto.FlyerResponse{
		ID:    deletedFlyer.ID.String(),
		Name:  deletedFlyer.Name,
		Image: deletedFlyer.Image,
	}

	return res, nil
}
