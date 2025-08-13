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
	IAchievementService interface {
		Create(ctx context.Context, req dto.CreateAchievementRequest) (dto.AchievementResponse, error)
		GetAll(ctx context.Context) ([]dto.AchievementResponse, error)
		GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.AchievementPaginationResponse, error)
		GetDetail(ctx context.Context, id string) (dto.AchievementResponse, error)
		Update(ctx context.Context, req dto.UpdateAchievementRequest) (dto.AchievementResponse, error)
		Delete(ctx context.Context, id string) (dto.AchievementResponse, error)
	}

	achievementService struct {
		achievementRepo repository.IAchievementRepository
		jwt             jwt.IJWT
	}
)

func NewAchievementService(achievementRepo repository.IAchievementRepository, jwt jwt.IJWT) *achievementService {
	return &achievementService{
		achievementRepo: achievementRepo,
		jwt:             jwt,
	}
}

func (as *achievementService) Create(ctx context.Context, req dto.CreateAchievementRequest) (dto.AchievementResponse, error) {
	// handle name request
	if req.Name == "" {
		return dto.AchievementResponse{}, dto.ErrEmptyName
	}
	if len(req.Name) < 3 {
		return dto.AchievementResponse{}, dto.ErrNameTooShort
	}

	// handle year request
	if req.Year == nil {
		return dto.AchievementResponse{}, dto.ErrEmptyYear
	}

	// handle description request
	if req.Description == "" {
		return dto.AchievementResponse{}, dto.ErrEmptyDescription
	}
	if len(req.Description) < 5 {
		return dto.AchievementResponse{}, dto.ErrDescriptionTooShort
	}

	// handle double data
	_, found, _ := as.achievementRepo.GetByNameAndYear(ctx, nil, req.Name, *req.Year)
	if found {
		return dto.AchievementResponse{}, dto.ErrAchievementAlreadyExists
	}

	achievementID := uuid.New()
	achievement := &entity.Achievement{
		ID:          achievementID,
		Name:        req.Name,
		Year:        *req.Year,
		Description: req.Description,
	}

	// handle image url
	var (
		achievementImages         []*entity.AchievementImage
		achievementImageResponses []dto.AchievementImageResponse
	)
	if len(req.Images) == 0 {
		return dto.AchievementResponse{}, dto.ErrEmptyImages
	}
	for _, imgName := range req.Images {
		imgID := uuid.New()
		// handle entity
		achievementImages = append(achievementImages, &entity.AchievementImage{
			ID:            imgID,
			Name:          imgName,
			AchievementID: &achievementID,
		})

		// handle response
		achievementImageResponses = append(achievementImageResponses, dto.AchievementImageResponse{
			ID:   imgID.String(),
			Name: imgName,
		})
	}

	err := as.achievementRepo.RunInTransaction(ctx, func(txRepo repository.IAchievementRepository) error {
		// create achievement
		if err := txRepo.Create(ctx, nil, achievement); err != nil {
			return dto.ErrCreateAchievement
		}

		// create achievement images
		for _, img := range achievementImages {
			if err := txRepo.CreateImage(ctx, nil, img); err != nil {
				return dto.ErrCreateAchievementImage
			}
		}

		return nil
	})
	if err != nil {
		return dto.AchievementResponse{}, err
	}

	return dto.AchievementResponse{
		ID:          achievement.ID.String(),
		Name:        achievement.Name,
		Year:        achievement.Year,
		Description: achievement.Description,
		Images:      achievementImageResponses,
	}, nil
}

func (as *achievementService) GetAll(ctx context.Context) ([]dto.AchievementResponse, error) {
	achievements, err := as.achievementRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, dto.ErrGetAllAchievementNoPagination
	}

	var datas []dto.AchievementResponse
	for _, achievement := range achievements {
		data := dto.AchievementResponse{
			ID:          achievement.ID.String(),
			Name:        achievement.Name,
			Year:        achievement.Year,
			Description: achievement.Description,
		}

		for _, a := range achievement.Images {
			data.Images = append(data.Images, dto.AchievementImageResponse{
				ID:   a.ID.String(),
				Name: a.Name,
			})
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (as *achievementService) GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.AchievementPaginationResponse, error) {
	dataWithPaginate, err := as.achievementRepo.GetAllWithPagination(ctx, nil, req)
	if err != nil {
		return dto.AchievementPaginationResponse{}, dto.ErrGetAllAchievementWithPagination
	}

	var datas []dto.AchievementResponse
	for _, achievement := range dataWithPaginate.Achievements {
		data := dto.AchievementResponse{
			ID:          achievement.ID.String(),
			Name:        achievement.Name,
			Year:        achievement.Year,
			Description: achievement.Description,
		}

		for _, a := range achievement.Images {
			data.Images = append(data.Images, dto.AchievementImageResponse{
				ID:   a.ID.String(),
				Name: a.Name,
			})
		}

		datas = append(datas, data)
	}

	return dto.AchievementPaginationResponse{
		Data: datas,
		PaginationResponse: response.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (as *achievementService) GetDetail(ctx context.Context, id string) (dto.AchievementResponse, error) {
	achievement, _, err := as.achievementRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.AchievementResponse{}, dto.ErrAchievementNotFound
	}

	res := dto.AchievementResponse{
		ID:          achievement.ID.String(),
		Name:        achievement.Name,
		Year:        achievement.Year,
		Description: achievement.Description,
	}

	for _, a := range achievement.Images {
		res.Images = append(res.Images, dto.AchievementImageResponse{
			ID:   a.ID.String(),
			Name: a.Name,
		})
	}

	return res, nil
}

func (as *achievementService) Update(ctx context.Context, req dto.UpdateAchievementRequest) (dto.AchievementResponse, error) {
	// get achievement by id
	achievement, found, err := as.achievementRepo.GetByID(ctx, nil, req.ID)
	if err != nil {
		return dto.AchievementResponse{}, dto.ErrGetAchievementByID
	}
	if !found {
		return dto.AchievementResponse{}, dto.ErrAchievementNotFound
	}

	// handle name request
	if req.Name != "" && req.Name != achievement.Name {
		if len(req.Name) < 3 {
			return dto.AchievementResponse{}, dto.ErrNameTooShort
		}

		achievement.Name = req.Name
	}

	// handle year request
	if req.Year != nil && req.Year != &achievement.Year {
		achievement.Year = *req.Year
	}

	// handle description request
	if req.Description != "" && req.Description != achievement.Description {
		if len(req.Description) < 5 {
			return dto.AchievementResponse{}, dto.ErrDescriptionTooShort
		}

		achievement.Description = req.Description
	}

	// handle image url
	var (
		achievementImages         []*entity.AchievementImage
		achievementImageResponses []dto.AchievementImageResponse
	)
	if len(req.Images) > 0 {
		for _, imgName := range req.Images {
			imgID := uuid.New()
			// handle entity
			achievementImages = append(achievementImages, &entity.AchievementImage{
				ID:            imgID,
				Name:          imgName,
				AchievementID: &achievement.ID,
			})

			// handle response
			achievementImageResponses = append(achievementImageResponses, dto.AchievementImageResponse{
				ID:   imgID.String(),
				Name: imgName,
			})
		}
	}

	err = as.achievementRepo.RunInTransaction(ctx, func(txRepo repository.IAchievementRepository) error {
		// update achievement
		if err := txRepo.Update(ctx, nil, achievement); err != nil {
			return dto.ErrUpdateAchievement
		}

		// handle new image
		if len(req.Name) > 0 {
			// check request images
			oldImages, err := txRepo.GetImagesByID(ctx, nil, achievement.ID.String())
			if err != nil {
				return dto.ErrGetAchievementImages
			}

			// Delete Existing Achievement Image
			// in assets
			for _, img := range oldImages {
				if err := os.Remove(img.Name); err != nil && !os.IsNotExist(err) {
					return dto.ErrDeleteOldImage
				}
			}
			// in db
			if err := txRepo.DeleteImagesByID(ctx, nil, achievement.ID.String()); err != nil {
				return dto.ErrDeleteAchievementImageByAchievementID
			}

			// Create new achievement images
			for _, img := range achievementImages {
				if err := txRepo.CreateImage(ctx, nil, img); err != nil {
					return dto.ErrCreateAchievementImage
				}
			}
		}

		return nil
	})
	if err != nil {
		return dto.AchievementResponse{}, err
	}

	return dto.AchievementResponse{
		ID:          achievement.ID.String(),
		Name:        achievement.Name,
		Year:        achievement.Year,
		Description: achievement.Description,
		Images:      achievementImageResponses,
	}, nil
}

func (as *achievementService) Delete(ctx context.Context, id string) (dto.AchievementResponse, error) {
	deletedAchievement, found, err := as.achievementRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.AchievementResponse{}, dto.ErrAchievementNotFound
	}
	if !found {
		return dto.AchievementResponse{}, dto.ErrAchievementNotFound
	}

	err = as.achievementRepo.RunInTransaction(ctx, func(txRepo repository.IAchievementRepository) error {
		// Delete Achievement Images
		oldAchievementImages, err := txRepo.GetImagesByID(ctx, nil, id)
		if err != nil {
			return dto.ErrGetAchievementImages
		}
		for _, img := range oldAchievementImages {
			name := strings.TrimPrefix(img.Name, "assets/")
			path := filepath.Join("assets", name)
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return dto.ErrDeleteOldImage
			}
		}
		if err := txRepo.DeleteImagesByID(ctx, nil, id); err != nil {
			return dto.ErrDeleteAchievementImageByAchievementID
		}

		// Delete Achievement
		err = as.achievementRepo.DeleteByID(ctx, nil, id)
		if err != nil {
			return dto.ErrDeleteAchievementByID
		}

		return nil
	})
	if err != nil {
		return dto.AchievementResponse{}, err
	}

	res := dto.AchievementResponse{
		ID:          deletedAchievement.ID.String(),
		Name:        deletedAchievement.Name,
		Year:        deletedAchievement.Year,
		Description: deletedAchievement.Description,
	}

	for _, a := range deletedAchievement.Images {
		res.Images = append(res.Images, dto.AchievementImageResponse{
			ID:   a.ID.String(),
			Name: a.Name,
		})
	}

	return res, nil
}
