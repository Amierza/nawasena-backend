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
	IPartnerService interface {
		Create(ctx context.Context, req dto.CreatePartnerRequest) (dto.PartnerResponse, error)
		GetAll(ctx context.Context) ([]dto.PartnerResponse, error)
		GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.PartnerPaginationResponse, error)
		GetDetail(ctx context.Context, id string) (dto.PartnerResponse, error)
		Update(ctx context.Context, req dto.UpdatePartnerRequest) (dto.PartnerResponse, error)
		Delete(ctx context.Context, id string) (dto.PartnerResponse, error)
	}

	partnerService struct {
		partnerRepo repository.IPartnerRepository
		jwt         jwt.IJWT
	}
)

func NewPartnerService(partnerRepo repository.IPartnerRepository, jwt jwt.IJWT) *partnerService {
	return &partnerService{
		partnerRepo: partnerRepo,
		jwt:         jwt,
	}
}

func (as *partnerService) Create(ctx context.Context, req dto.CreatePartnerRequest) (dto.PartnerResponse, error) {
	// handle name request
	if req.Name == "" {
		return dto.PartnerResponse{}, dto.ErrEmptyName
	}
	if len(req.Name) < 3 {
		return dto.PartnerResponse{}, dto.ErrNameTooShort
	}

	// handle image request
	if req.Image == "" {
		return dto.PartnerResponse{}, dto.ErrEmptyImage
	}

	// handle double data
	_, found, _ := as.partnerRepo.GetByName(ctx, nil, req.Name)
	if found {
		return dto.PartnerResponse{}, dto.ErrPartnerAlreadyExists
	}

	partnerID := uuid.New()
	partner := &entity.Partner{
		ID:    partnerID,
		Name:  req.Name,
		Image: req.Image,
	}

	// create partner
	if err := as.partnerRepo.Create(ctx, nil, partner); err != nil {
		return dto.PartnerResponse{}, dto.ErrCreatePartner
	}

	return dto.PartnerResponse{
		ID:    partner.ID.String(),
		Name:  partner.Name,
		Image: partner.Image,
	}, nil
}

func (as *partnerService) GetAll(ctx context.Context) ([]dto.PartnerResponse, error) {
	partners, err := as.partnerRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, dto.ErrGetAllPartnerNoPagination
	}

	var datas []dto.PartnerResponse
	for _, partner := range partners {
		data := dto.PartnerResponse{
			ID:    partner.ID.String(),
			Name:  partner.Name,
			Image: partner.Image,
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (as *partnerService) GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.PartnerPaginationResponse, error) {
	dataWithPaginate, err := as.partnerRepo.GetAllWithPagination(ctx, nil, req)
	if err != nil {
		return dto.PartnerPaginationResponse{}, dto.ErrGetAllPartnerWithPagination
	}

	var datas []dto.PartnerResponse
	for _, partner := range dataWithPaginate.Partners {
		data := dto.PartnerResponse{
			ID:    partner.ID.String(),
			Name:  partner.Name,
			Image: partner.Image,
		}

		datas = append(datas, data)
	}

	return dto.PartnerPaginationResponse{
		Data: datas,
		PaginationResponse: response.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (as *partnerService) GetDetail(ctx context.Context, id string) (dto.PartnerResponse, error) {
	partner, _, err := as.partnerRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.PartnerResponse{}, dto.ErrPartnerNotFound
	}

	res := dto.PartnerResponse{
		ID:    partner.ID.String(),
		Name:  partner.Name,
		Image: partner.Image,
	}

	return res, nil
}

func (as *partnerService) Update(ctx context.Context, req dto.UpdatePartnerRequest) (dto.PartnerResponse, error) {
	// get partner by id
	partner, found, err := as.partnerRepo.GetByID(ctx, nil, req.ID)
	if err != nil {
		return dto.PartnerResponse{}, dto.ErrGetPartnerByID
	}
	if !found {
		return dto.PartnerResponse{}, dto.ErrPartnerNotFound
	}

	// handle name request
	if req.Name != "" && req.Name != partner.Name {
		if len(req.Name) < 3 {
			return dto.PartnerResponse{}, dto.ErrNameTooShort
		}

		partner.Name = req.Name
	}

	// handle image request
	oldImage := partner.Image
	if req.Image != "" && req.Image != partner.Image {
		partner.Image = req.Image
	}

	err = as.partnerRepo.RunInTransaction(ctx, func(txRepo repository.IPartnerRepository) error {
		// update partner
		if err := txRepo.Update(ctx, nil, partner); err != nil {
			return dto.ErrUpdatePartner
		}

		name := strings.TrimPrefix(oldImage, "assets/")
		path := filepath.Join("assets", name)
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return dto.ErrDeleteOldImage
		}

		return nil
	})
	if err != nil {
		return dto.PartnerResponse{}, err
	}

	return dto.PartnerResponse{
		ID:    partner.ID.String(),
		Name:  partner.Name,
		Image: partner.Image,
	}, nil
}

func (as *partnerService) Delete(ctx context.Context, id string) (dto.PartnerResponse, error) {
	deletedPartner, found, err := as.partnerRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.PartnerResponse{}, dto.ErrPartnerNotFound
	}
	if !found {
		return dto.PartnerResponse{}, dto.ErrPartnerNotFound
	}

	err = as.partnerRepo.RunInTransaction(ctx, func(txRepo repository.IPartnerRepository) error {
		// Delete Partner Image
		name := strings.TrimPrefix(deletedPartner.Image, "assets/")
		path := filepath.Join("assets", name)
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return dto.ErrDeleteOldImage
		}

		// Delete Partner
		err = as.partnerRepo.DeleteByID(ctx, nil, id)
		if err != nil {
			return dto.ErrDeletePartnerByID
		}

		return nil
	})
	if err != nil {
		return dto.PartnerResponse{}, err
	}

	res := dto.PartnerResponse{
		ID:    deletedPartner.ID.String(),
		Name:  deletedPartner.Name,
		Image: deletedPartner.Image,
	}

	return res, nil
}
