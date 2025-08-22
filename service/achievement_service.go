package service

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
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
		GetFeatured(ctx context.Context, limit string) ([]dto.AchievementResponse, error)
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
	// validate name
	if len(req.Name) < 3 {
		return dto.AchievementResponse{}, dto.ErrNameTooShort
	}

	// validate description
	if len(req.Description) < 5 {
		return dto.AchievementResponse{}, dto.ErrDescriptionTooShort
	}

	// validate required arrays
	if len(req.Team) == 0 {
		return dto.AchievementResponse{}, dto.ErrEmptyTeam
	}
	if len(req.Tags) == 0 {
		return dto.AchievementResponse{}, dto.ErrEmptyTags
	}

	// handle double data
	_, found, _ := as.achievementRepo.GetByNameAndYear(ctx, nil, req.Name, req.Year)
	if found {
		return dto.AchievementResponse{}, dto.ErrAchievementAlreadyExists
	}

	achievementID := uuid.New()
	achievement := &entity.Achievement{
		ID:          achievementID,
		Name:        req.Name,
		Year:        req.Year,
		Description: req.Description,
		Location:    req.Location,
		Rank:        req.Rank,
		Competition: req.Competition,
		Team:        req.Team,
		Impact:      req.Impact,
		VideoURL:    req.VideoURL,
		Featured:    req.Featured,
		Tags:        req.Tags,
	}

	// handle category
	category, found, _ := as.achievementRepo.GetCategoryByCategoryID(ctx, nil, req.CategoryID)
	if !found {
		return dto.AchievementResponse{}, dto.ErrAchievementCategoryNotFound
	}
	categoryUUID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return dto.AchievementResponse{}, dto.ErrParseUUID
	}
	achievement.AchievementCategoryID = &categoryUUID

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

	err = as.achievementRepo.RunInTransaction(ctx, func(txRepo repository.IAchievementRepository) error {
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
		Location:    achievement.Location,
		Rank:        achievement.Rank,
		Competition: achievement.Competition,
		Team:        achievement.Team,
		Impact:      achievement.Impact,
		VideoURL:    achievement.VideoURL,
		Featured:    achievement.Featured,
		Tags:        achievement.Tags,
		Images:      achievementImageResponses,
		Category: dto.AchievementCategoryResponse{
			ID:   achievement.AchievementCategoryID.String(),
			Name: category.Name,
		},
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
			Location:    achievement.Location,
			Rank:        achievement.Rank,
			Competition: achievement.Competition,
			Team:        achievement.Team,
			Impact:      achievement.Impact,
			VideoURL:    achievement.VideoURL,
			Featured:    achievement.Featured,
			Tags:        achievement.Tags,
			Category: dto.AchievementCategoryResponse{
				ID:   achievement.AchievementCategoryID.String(),
				Name: achievement.AchievementCategory.Name,
			},
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
			Location:    achievement.Location,
			Rank:        achievement.Rank,
			Competition: achievement.Competition,
			Team:        achievement.Team,
			Impact:      achievement.Impact,
			VideoURL:    achievement.VideoURL,
			Featured:    achievement.Featured,
			Tags:        achievement.Tags,
			Category: dto.AchievementCategoryResponse{
				ID:   achievement.AchievementCategoryID.String(),
				Name: achievement.AchievementCategory.Name,
			},
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
		Location:    achievement.Location,
		Rank:        achievement.Rank,
		Competition: achievement.Competition,
		Team:        achievement.Team,
		Impact:      achievement.Impact,
		VideoURL:    achievement.VideoURL,
		Featured:    achievement.Featured,
		Tags:        achievement.Tags,
		Category: dto.AchievementCategoryResponse{
			ID:   achievement.AchievementCategoryID.String(),
			Name: achievement.AchievementCategory.Name,
		},
	}

	for _, a := range achievement.Images {
		res.Images = append(res.Images, dto.AchievementImageResponse{
			ID:   a.ID.String(),
			Name: a.Name,
		})
	}

	return res, nil
}

func (ns *achievementService) GetFeatured(ctx context.Context, limit string) ([]dto.AchievementResponse, error) {
	lim := 1 // default
	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			lim = l
		} else {
			return nil, dto.ErrParseLimit
		}
	}

	featuredAchievement, err := ns.achievementRepo.GetFeatured(ctx, nil, &lim)
	if err != nil {
		return nil, dto.ErrGetAllFeaturedAchievement
	}

	var datas []dto.AchievementResponse
	for _, achievement := range featuredAchievement {
		data := dto.AchievementResponse{
			ID:          achievement.ID.String(),
			Name:        achievement.Name,
			Year:        achievement.Year,
			Description: achievement.Description,
			Location:    achievement.Location,
			Rank:        achievement.Rank,
			Competition: achievement.Competition,
			Team:        achievement.Team,
			Impact:      achievement.Impact,
			VideoURL:    achievement.VideoURL,
			Featured:    achievement.Featured,
			Tags:        achievement.Tags,
			Category: dto.AchievementCategoryResponse{
				ID:   achievement.AchievementCategoryID.String(),
				Name: achievement.AchievementCategory.Name,
			},
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
	if req.Name != "" {
		achievement.Name = req.Name
	}

	// handle year request
	if req.Year != nil && req.Year != &achievement.Year {
		achievement.Year = *req.Year
	}

	// handle description request
	if req.Description != "" {
		achievement.Description = req.Description
	}

	// handle other optional fields
	if req.Location != "" {
		achievement.Location = req.Location
	}

	if req.Rank != "" {
		achievement.Rank = req.Rank
	}

	if req.Competition != "" {
		achievement.Competition = req.Competition
	}

	if len(req.Team) > 0 {
		achievement.Team = req.Team
	}

	if req.Impact != "" {
		achievement.Impact = req.Impact
	}

	if req.VideoURL != "" {
		achievement.VideoURL = req.VideoURL
	}

	achievement.Featured = req.Featured

	if len(req.Tags) > 0 {
		achievement.Tags = req.Tags
	}

	// handle category
	if req.CategoryID != "" {
		_, found, _ := as.achievementRepo.GetCategoryByCategoryID(ctx, nil, req.CategoryID)
		if !found {
			return dto.AchievementResponse{}, dto.ErrAchievementCategoryNotFound
		}

		categoryUUID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return dto.AchievementResponse{}, dto.ErrParseUUID
		}
		achievement.AchievementCategoryID = &categoryUUID
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
		Location:    achievement.Location,
		Rank:        achievement.Rank,
		Competition: achievement.Competition,
		Team:        achievement.Team,
		Impact:      achievement.Impact,
		VideoURL:    achievement.VideoURL,
		Featured:    achievement.Featured,
		Tags:        achievement.Tags,
		Images:      achievementImageResponses,
		Category: dto.AchievementCategoryResponse{
			ID:   achievement.AchievementCategoryID.String(),
			Name: achievement.AchievementCategory.Name,
		},
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
		Location:    deletedAchievement.Location,
		Rank:        deletedAchievement.Rank,
		Competition: deletedAchievement.Competition,
		Team:        deletedAchievement.Team,
		Impact:      deletedAchievement.Impact,
		VideoURL:    deletedAchievement.VideoURL,
		Featured:    deletedAchievement.Featured,
		Tags:        deletedAchievement.Tags,
		Category: dto.AchievementCategoryResponse{
			ID:   deletedAchievement.AchievementCategoryID.String(),
			Name: deletedAchievement.AchievementCategory.Name,
		},
	}

	for _, a := range deletedAchievement.Images {
		res.Images = append(res.Images, dto.AchievementImageResponse{
			ID:   a.ID.String(),
			Name: a.Name,
		})
	}

	return res, nil
}
