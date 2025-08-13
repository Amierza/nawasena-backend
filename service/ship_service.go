package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/entity"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/repository"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/google/uuid"
)

type (
	IShipService interface {
		Create(ctx context.Context, req dto.CreateShipRequest) (dto.ShipResponse, error)
		GetAll(ctx context.Context) ([]dto.ShipResponse, error)
		GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.ShipPaginationResponse, error)
		GetDetail(ctx context.Context, id string) (dto.ShipResponse, error)
		Update(ctx context.Context, req dto.UpdateShipRequest) (dto.ShipResponse, error)
		Delete(ctx context.Context, id string) (dto.ShipResponse, error)
	}

	shipService struct {
		shipRepo repository.IShipRepository
		jwt      jwt.IJWT
	}
)

func NewShipService(shipRepo repository.IShipRepository, jwt jwt.IJWT) *shipService {
	return &shipService{
		shipRepo: shipRepo,
		jwt:      jwt,
	}
}

func (as *shipService) Create(ctx context.Context, req dto.CreateShipRequest) (dto.ShipResponse, error) {
	// handle name request
	if req.Name == "" {
		return dto.ShipResponse{}, dto.ErrEmptyName
	}
	if len(req.Name) < 3 {
		return dto.ShipResponse{}, dto.ErrNameTooShort
	}

	// handle description request
	if req.Description == "" {
		return dto.ShipResponse{}, dto.ErrEmptyDescription
	}
	if len(req.Description) < 5 {
		return dto.ShipResponse{}, dto.ErrDescriptionTooShort
	}

	// handle double data
	_, found, _ := as.shipRepo.GetByName(ctx, nil, req.Name)
	if found {
		return dto.ShipResponse{}, dto.ErrShipAlreadyExists
	}

	shipID := uuid.New()
	ship := &entity.Ship{
		ID:          shipID,
		Name:        req.Name,
		Description: req.Description,
	}

	// handle image url
	var (
		shipImages         []*entity.ShipImage
		shipImageResponses []dto.ShipImageResponse
	)
	if len(req.Images) == 0 {
		return dto.ShipResponse{}, dto.ErrEmptyImages
	}
	for _, imgName := range req.Images {
		imgID := uuid.New()
		// handle entity
		shipImages = append(shipImages, &entity.ShipImage{
			ID:     imgID,
			Name:   imgName,
			ShipID: &shipID,
		})

		// handle response
		shipImageResponses = append(shipImageResponses, dto.ShipImageResponse{
			ID:   imgID.String(),
			Name: imgName,
		})
	}

	err := as.shipRepo.RunInTransaction(ctx, func(txRepo repository.IShipRepository) error {
		// create ship
		if err := txRepo.Create(ctx, nil, ship); err != nil {
			return dto.ErrCreateShip
		}

		// handle new image
		if len(req.Name) > 0 {
			// check request images
			oldImages, err := txRepo.GetImagesByID(ctx, nil, ship.ID.String())
			if err != nil {
				return dto.ErrGetShipImages
			}

			// Delete Existing Ship Image
			// in assets
			for _, img := range oldImages {
				name := strings.TrimPrefix(img.Name, "assets/")
				path := filepath.Join("assets", name)
				if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
					return dto.ErrDeleteOldImage
				}
			}
			// in db
			if err := txRepo.DeleteImagesByID(ctx, nil, ship.ID.String()); err != nil {
				return dto.ErrDeleteShipImageByShipID
			}

			// Create new ship images
			for _, img := range shipImages {
				if err := txRepo.CreateImage(ctx, nil, img); err != nil {
					return dto.ErrCreateShipImage
				}
			}
		}

		return nil
	})
	if err != nil {
		return dto.ShipResponse{}, err
	}

	return dto.ShipResponse{
		ID:          ship.ID.String(),
		Name:        ship.Name,
		Description: ship.Description,
		Images:      shipImageResponses,
	}, nil
}

func (as *shipService) GetAll(ctx context.Context) ([]dto.ShipResponse, error) {
	ships, err := as.shipRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, dto.ErrGetAllShipNoPagination
	}

	var datas []dto.ShipResponse
	for _, ship := range ships {
		data := dto.ShipResponse{
			ID:          ship.ID.String(),
			Name:        ship.Name,
			Description: ship.Description,
		}

		for _, a := range ship.Images {
			data.Images = append(data.Images, dto.ShipImageResponse{
				ID:   a.ID.String(),
				Name: a.Name,
			})
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (as *shipService) GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.ShipPaginationResponse, error) {
	dataWithPaginate, err := as.shipRepo.GetAllWithPagination(ctx, nil, req)
	if err != nil {
		return dto.ShipPaginationResponse{}, dto.ErrGetAllShipWithPagination
	}

	var datas []dto.ShipResponse
	for _, ship := range dataWithPaginate.Ships {
		data := dto.ShipResponse{
			ID:          ship.ID.String(),
			Name:        ship.Name,
			Description: ship.Description,
		}

		for _, a := range ship.Images {
			data.Images = append(data.Images, dto.ShipImageResponse{
				ID:   a.ID.String(),
				Name: a.Name,
			})
		}

		datas = append(datas, data)
	}

	return dto.ShipPaginationResponse{
		Data: datas,
		PaginationResponse: response.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (as *shipService) GetDetail(ctx context.Context, id string) (dto.ShipResponse, error) {
	ship, _, err := as.shipRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.ShipResponse{}, dto.ErrShipNotFound
	}

	res := dto.ShipResponse{
		ID:          ship.ID.String(),
		Name:        ship.Name,
		Description: ship.Description,
	}

	for _, a := range ship.Images {
		res.Images = append(res.Images, dto.ShipImageResponse{
			ID:   a.ID.String(),
			Name: a.Name,
		})
	}

	return res, nil
}

func (as *shipService) Update(ctx context.Context, req dto.UpdateShipRequest) (dto.ShipResponse, error) {
	// get ship by id
	ship, found, err := as.shipRepo.GetByID(ctx, nil, req.ID)
	if err != nil {
		return dto.ShipResponse{}, dto.ErrGetShipByID
	}
	if !found {
		return dto.ShipResponse{}, dto.ErrShipNotFound
	}

	// handle name request
	if req.Name != "" && req.Name != ship.Name {
		if len(req.Name) < 3 {
			return dto.ShipResponse{}, dto.ErrNameTooShort
		}

		ship.Name = req.Name
	}

	// handle description request
	if req.Description != "" && req.Description != ship.Description {
		if len(req.Description) < 5 {
			return dto.ShipResponse{}, dto.ErrDescriptionTooShort
		}

		ship.Description = req.Description
	}

	// handle image url
	var (
		shipImages         []*entity.ShipImage
		shipImageResponses []dto.ShipImageResponse
	)
	if len(req.Images) > 0 {
		for _, imgName := range req.Images {
			imgID := uuid.New()
			// handle entity
			shipImages = append(shipImages, &entity.ShipImage{
				ID:     imgID,
				Name:   imgName,
				ShipID: &ship.ID,
			})

			// handle response
			shipImageResponses = append(shipImageResponses, dto.ShipImageResponse{
				ID:   imgID.String(),
				Name: imgName,
			})
		}
	}

	err = as.shipRepo.RunInTransaction(ctx, func(txRepo repository.IShipRepository) error {
		// update ship
		if err := txRepo.Update(ctx, nil, ship); err != nil {
			return dto.ErrUpdateShip
		}

		// handle new image
		if len(req.Name) > 0 {
			// check request images
			oldImages, err := txRepo.GetImagesByID(ctx, nil, ship.ID.String())
			if err != nil {
				return dto.ErrGetShipImages
			}

			// Delete Existing Ship Image
			// in assets
			for _, img := range oldImages {
				if err := os.Remove(img.Name); err != nil && !os.IsNotExist(err) {
					return dto.ErrDeleteOldImage
				}
			}
			// in db
			if err := txRepo.DeleteImagesByID(ctx, nil, ship.ID.String()); err != nil {
				return dto.ErrDeleteShipImageByShipID
			}

			// Create new ship images
			for _, img := range shipImages {
				if err := txRepo.CreateImage(ctx, nil, img); err != nil {
					return dto.ErrCreateShipImage
				}
			}
		}

		return nil
	})
	if err != nil {
		return dto.ShipResponse{}, err
	}

	return dto.ShipResponse{
		ID:          ship.ID.String(),
		Name:        ship.Name,
		Description: ship.Description,
		Images:      shipImageResponses,
	}, nil
}

func (as *shipService) Delete(ctx context.Context, id string) (dto.ShipResponse, error) {
	deletedShip, found, err := as.shipRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.ShipResponse{}, dto.ErrShipNotFound
	}
	if !found {
		return dto.ShipResponse{}, dto.ErrShipNotFound
	}

	err = as.shipRepo.RunInTransaction(ctx, func(txRepo repository.IShipRepository) error {
		// Delete Ship Images
		oldShipImages, err := txRepo.GetImagesByID(ctx, nil, id)
		if err != nil {
			return dto.ErrGetShipImages
		}
		for _, img := range oldShipImages {
			name := strings.TrimPrefix(img.Name, "assets/")
			path := filepath.Join("assets", name)
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return dto.ErrDeleteOldImage
			}
		}
		if err := txRepo.DeleteImagesByID(ctx, nil, id); err != nil {
			return dto.ErrDeleteShipImageByShipID
		}

		// Delete Ship
		err = as.shipRepo.DeleteByID(ctx, nil, id)
		if err != nil {
			return dto.ErrDeleteShipByID
		}

		return nil
	})
	if err != nil {
		return dto.ShipResponse{}, err
	}

	res := dto.ShipResponse{
		ID:          deletedShip.ID.String(),
		Name:        deletedShip.Name,
		Description: deletedShip.Description,
	}

	for _, a := range deletedShip.Images {
		res.Images = append(res.Images, dto.ShipImageResponse{
			ID:   a.ID.String(),
			Name: a.Name,
		})
	}

	return res, nil
}
