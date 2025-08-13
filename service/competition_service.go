package service

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/entity"
	"github.com/Amierza/nawasena-backend/helper"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/repository"
	"github.com/Amierza/nawasena-backend/response"
	"github.com/google/uuid"
)

type (
	ICompetitionService interface {
		Create(ctx context.Context, req dto.CreateCompetitionRequest) (dto.CompetitionResponse, error)
		GetAll(ctx context.Context) ([]dto.CompetitionResponse, error)
		GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.CompetitionPaginationResponse, error)
		GetDetail(ctx context.Context, id string) (dto.CompetitionResponse, error)
		Update(ctx context.Context, req dto.UpdateCompetitionRequest) (dto.CompetitionResponse, error)
		Delete(ctx context.Context, id string) (dto.CompetitionResponse, error)
	}

	competitionService struct {
		competitionRepo repository.ICompetitionRepository
		jwt             jwt.IJWT
	}
)

func NewCompetitionService(competitionRepo repository.ICompetitionRepository, jwt jwt.IJWT) *competitionService {
	return &competitionService{
		competitionRepo: competitionRepo,
		jwt:             jwt,
	}
}

func (as *competitionService) Create(ctx context.Context, req dto.CreateCompetitionRequest) (dto.CompetitionResponse, error) {
	// handle name request
	if req.Name == "" {
		return dto.CompetitionResponse{}, dto.ErrEmptyName
	}
	if len(req.Name) < 3 {
		return dto.CompetitionResponse{}, dto.ErrNameTooShort
	}

	// handle date request
	if req.Date == "" {
		return dto.CompetitionResponse{}, dto.ErrEmptyDate
	}
	date, err := helper.StringToTime(req.Date)
	if err != nil {
		return dto.CompetitionResponse{}, dto.ErrParseTimeFromStringToTime
	}

	// handle description request
	if req.Description == "" {
		return dto.CompetitionResponse{}, dto.ErrEmptyDescription
	}
	if len(req.Description) < 5 {
		return dto.CompetitionResponse{}, dto.ErrDescriptionTooShort
	}

	// handle double data
	_, found, _ := as.competitionRepo.GetByNameAndDate(ctx, nil, req.Name, date)
	if found {
		return dto.CompetitionResponse{}, dto.ErrCompetitionAlreadyExists
	}

	competitionID := uuid.New()
	competition := &entity.Competition{
		ID:          competitionID,
		Name:        req.Name,
		Date:        date,
		Description: req.Description,
	}

	// handle image url
	var (
		competitionImages         []*entity.CompetitionImage
		competitionImageResponses []dto.CompetitionImageResponse
	)
	if len(req.Images) == 0 {
		return dto.CompetitionResponse{}, dto.ErrEmptyImages
	}
	for _, imgName := range req.Images {
		imgID := uuid.New()
		// handle entity
		competitionImages = append(competitionImages, &entity.CompetitionImage{
			ID:            imgID,
			Name:          imgName,
			CompetitionID: &competitionID,
		})

		// handle response
		competitionImageResponses = append(competitionImageResponses, dto.CompetitionImageResponse{
			ID:   imgID.String(),
			Name: imgName,
		})
	}

	err = as.competitionRepo.RunInTransaction(ctx, func(txRepo repository.ICompetitionRepository) error {
		// create competition
		if err := txRepo.Create(ctx, nil, competition); err != nil {
			return dto.ErrCreateCompetition
		}

		// create competition images
		for _, img := range competitionImages {
			if err := txRepo.CreateImage(ctx, nil, img); err != nil {
				return dto.ErrCreateCompetitionImage
			}
		}

		return nil
	})
	if err != nil {
		return dto.CompetitionResponse{}, err
	}

	return dto.CompetitionResponse{
		ID:          competition.ID.String(),
		Name:        competition.Name,
		Date:        competition.Date.Format("2006-01-02"),
		Description: competition.Description,
		Images:      competitionImageResponses,
	}, nil
}

func (as *competitionService) GetAll(ctx context.Context) ([]dto.CompetitionResponse, error) {
	competitions, err := as.competitionRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, dto.ErrGetAllCompetitionNoPagination
	}

	var datas []dto.CompetitionResponse
	for _, competition := range competitions {
		data := dto.CompetitionResponse{
			ID:          competition.ID.String(),
			Name:        competition.Name,
			Date:        competition.Date.Format("2006-01-02"),
			Description: competition.Description,
		}

		for _, a := range competition.Images {
			data.Images = append(data.Images, dto.CompetitionImageResponse{
				ID:   a.ID.String(),
				Name: a.Name,
			})
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (as *competitionService) GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.CompetitionPaginationResponse, error) {
	dataWithPaginate, err := as.competitionRepo.GetAllWithPagination(ctx, nil, req)
	if err != nil {
		return dto.CompetitionPaginationResponse{}, dto.ErrGetAllCompetitionWithPagination
	}

	var datas []dto.CompetitionResponse
	for _, competition := range dataWithPaginate.Competitions {
		data := dto.CompetitionResponse{
			ID:          competition.ID.String(),
			Name:        competition.Name,
			Date:        competition.Date.Format("2006-01-02"),
			Description: competition.Description,
		}

		for _, a := range competition.Images {
			data.Images = append(data.Images, dto.CompetitionImageResponse{
				ID:   a.ID.String(),
				Name: a.Name,
			})
		}

		datas = append(datas, data)
	}

	return dto.CompetitionPaginationResponse{
		Data: datas,
		PaginationResponse: response.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (as *competitionService) GetDetail(ctx context.Context, id string) (dto.CompetitionResponse, error) {
	competition, _, err := as.competitionRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.CompetitionResponse{}, dto.ErrCompetitionNotFound
	}

	res := dto.CompetitionResponse{
		ID:          competition.ID.String(),
		Name:        competition.Name,
		Date:        competition.Date.Format("2006-01-02"),
		Description: competition.Description,
	}

	for _, a := range competition.Images {
		res.Images = append(res.Images, dto.CompetitionImageResponse{
			ID:   a.ID.String(),
			Name: a.Name,
		})
	}

	return res, nil
}

func (as *competitionService) Update(ctx context.Context, req dto.UpdateCompetitionRequest) (dto.CompetitionResponse, error) {
	// get competition by id
	competition, found, err := as.competitionRepo.GetByID(ctx, nil, req.ID)
	if err != nil {
		return dto.CompetitionResponse{}, dto.ErrGetCompetitionByID
	}
	if !found {
		return dto.CompetitionResponse{}, dto.ErrCompetitionNotFound
	}

	// handle name request
	if req.Name != "" && req.Name != competition.Name {
		if len(req.Name) < 3 {
			return dto.CompetitionResponse{}, dto.ErrNameTooShort
		}

		competition.Name = req.Name
	}

	// handle date request
	var date time.Time
	if req.Date != "" && req.Date != competition.Date.Format("2006-01-02") {
		d, err := helper.StringToTime(req.Date)
		if err != nil {
			return dto.CompetitionResponse{}, dto.ErrParseTimeFromStringToTime
		}

		date = d
		competition.Date = date
	}

	// handle description request
	if req.Description != "" && req.Description != competition.Description {
		if len(req.Description) < 5 {
			return dto.CompetitionResponse{}, dto.ErrDescriptionTooShort
		}

		competition.Description = req.Description
	}

	// handle double data
	_, found, _ = as.competitionRepo.GetByNameAndDate(ctx, nil, req.Name, date)
	if found {
		return dto.CompetitionResponse{}, dto.ErrCompetitionAlreadyExists
	}

	// handle image url
	var (
		competitionImages         []*entity.CompetitionImage
		competitionImageResponses []dto.CompetitionImageResponse
	)
	if len(req.Images) > 0 {
		for _, imgName := range req.Images {
			imgID := uuid.New()
			// handle entity
			competitionImages = append(competitionImages, &entity.CompetitionImage{
				ID:            imgID,
				Name:          imgName,
				CompetitionID: &competition.ID,
			})

			// handle response
			competitionImageResponses = append(competitionImageResponses, dto.CompetitionImageResponse{
				ID:   imgID.String(),
				Name: imgName,
			})
		}
	}

	err = as.competitionRepo.RunInTransaction(ctx, func(txRepo repository.ICompetitionRepository) error {
		// update competition
		if err := txRepo.Update(ctx, nil, competition); err != nil {
			return dto.ErrUpdateCompetition
		}

		// handle new image
		if len(req.Name) > 0 {
			// check request images
			oldImages, err := txRepo.GetImagesByID(ctx, nil, competition.ID.String())
			if err != nil {
				return dto.ErrGetCompetitionImages
			}

			// Delete Existing Competition Image
			// in assets
			for _, img := range oldImages {
				name := strings.TrimPrefix(img.Name, "assets/")
				path := filepath.Join("assets", name)
				if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
					return dto.ErrDeleteOldImage
				}
			}
			// in db
			if err := txRepo.DeleteImagesByID(ctx, nil, competition.ID.String()); err != nil {
				return dto.ErrDeleteCompetitionImageByCompetitionID
			}

			// Create new competition images
			for _, img := range competitionImages {
				if err := txRepo.CreateImage(ctx, nil, img); err != nil {
					return dto.ErrCreateCompetitionImage
				}
			}
		}

		return nil
	})
	if err != nil {
		return dto.CompetitionResponse{}, err
	}

	return dto.CompetitionResponse{
		ID:          competition.ID.String(),
		Name:        competition.Name,
		Date:        competition.Date.Format("2006-01-02"),
		Description: competition.Description,
		Images:      competitionImageResponses,
	}, nil
}

func (as *competitionService) Delete(ctx context.Context, id string) (dto.CompetitionResponse, error) {
	deletedCompetition, found, err := as.competitionRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.CompetitionResponse{}, dto.ErrCompetitionNotFound
	}
	if !found {
		return dto.CompetitionResponse{}, dto.ErrCompetitionNotFound
	}

	err = as.competitionRepo.RunInTransaction(ctx, func(txRepo repository.ICompetitionRepository) error {
		// Delete Competition Images
		oldCompetitionImages, err := txRepo.GetImagesByID(ctx, nil, id)
		if err != nil {
			return dto.ErrGetCompetitionImages
		}
		for _, img := range oldCompetitionImages {
			name := strings.TrimPrefix(img.Name, "assets/")
			path := filepath.Join("assets", name)
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return dto.ErrDeleteOldImage
			}
		}
		if err := txRepo.DeleteImagesByID(ctx, nil, id); err != nil {
			return dto.ErrDeleteCompetitionImageByCompetitionID
		}

		// Delete Competition
		err = as.competitionRepo.DeleteByID(ctx, nil, id)
		if err != nil {
			return dto.ErrDeleteCompetitionByID
		}

		return nil
	})
	if err != nil {
		return dto.CompetitionResponse{}, err
	}

	res := dto.CompetitionResponse{
		ID:          deletedCompetition.ID.String(),
		Name:        deletedCompetition.Name,
		Date:        deletedCompetition.Date.Format("2006-01-02"),
		Description: deletedCompetition.Description,
	}

	for _, a := range deletedCompetition.Images {
		res.Images = append(res.Images, dto.CompetitionImageResponse{
			ID:   a.ID.String(),
			Name: a.Name,
		})
	}

	return res, nil
}
