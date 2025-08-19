package service

import (
	"context"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/entity"
	"github.com/Amierza/nawasena-backend/jwt"
	"github.com/Amierza/nawasena-backend/repository"
	"github.com/google/uuid"
)

type (
	IAchievementCategoryService interface {
		Create(ctx context.Context, req dto.CreateAchievementCategoryRequest) (dto.AchievementCategoryResponse, error)
		GetAll(ctx context.Context) ([]dto.AchievementCategoryResponse, error)
		GetDetail(ctx context.Context, id string) (dto.AchievementCategoryResponse, error)
		Update(ctx context.Context, req dto.UpdateAchievementCategoryRequest) (dto.AchievementCategoryResponse, error)
		Delete(ctx context.Context, id string) (dto.AchievementCategoryResponse, error)
	}

	achievementCategoryService struct {
		achievementCategoryRepo repository.IAchievementCategoryRepository
		jwt                     jwt.IJWT
	}
)

func NewAchievementCategoryService(achievementCategoryRepo repository.IAchievementCategoryRepository, jwt jwt.IJWT) *achievementCategoryService {
	return &achievementCategoryService{
		achievementCategoryRepo: achievementCategoryRepo,
		jwt:                     jwt,
	}
}

func (as *achievementCategoryService) Create(ctx context.Context, req dto.CreateAchievementCategoryRequest) (dto.AchievementCategoryResponse, error) {
	// handle name request
	if req.Name == "" {
		return dto.AchievementCategoryResponse{}, dto.ErrEmptyName
	}
	if len(req.Name) < 3 {
		return dto.AchievementCategoryResponse{}, dto.ErrNameTooShort
	}

	// handle double data
	_, found, _ := as.achievementCategoryRepo.GetByName(ctx, nil, req.Name)
	if found {
		return dto.AchievementCategoryResponse{}, dto.ErrAchievementCategoryAlreadyExists
	}

	achievementCategoryID := uuid.New()
	achievementCategory := &entity.AchievementCategory{
		ID:   achievementCategoryID,
		Name: req.Name,
	}

	// create achievementCategory
	if err := as.achievementCategoryRepo.Create(ctx, nil, achievementCategory); err != nil {
		return dto.AchievementCategoryResponse{}, dto.ErrCreateAchievementCategory
	}

	return dto.AchievementCategoryResponse{
		ID:   achievementCategory.ID.String(),
		Name: achievementCategory.Name,
	}, nil
}

func (as *achievementCategoryService) GetAll(ctx context.Context) ([]dto.AchievementCategoryResponse, error) {
	achievementCategorys, err := as.achievementCategoryRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, dto.ErrGetAllAchievementCategory
	}

	var datas []dto.AchievementCategoryResponse
	for _, achievementCategory := range achievementCategorys {
		data := dto.AchievementCategoryResponse{
			ID:   achievementCategory.ID.String(),
			Name: achievementCategory.Name,
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (as *achievementCategoryService) GetDetail(ctx context.Context, id string) (dto.AchievementCategoryResponse, error) {
	achievementCategory, _, err := as.achievementCategoryRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.AchievementCategoryResponse{}, dto.ErrAchievementCategoryNotFound
	}

	res := dto.AchievementCategoryResponse{
		ID:   achievementCategory.ID.String(),
		Name: achievementCategory.Name,
	}

	return res, nil
}

func (as *achievementCategoryService) Update(ctx context.Context, req dto.UpdateAchievementCategoryRequest) (dto.AchievementCategoryResponse, error) {
	// get achievementCategory by id
	achievementCategory, found, err := as.achievementCategoryRepo.GetByID(ctx, nil, req.ID)
	if err != nil {
		return dto.AchievementCategoryResponse{}, dto.ErrGetAchievementCategoryByID
	}
	if !found {
		return dto.AchievementCategoryResponse{}, dto.ErrAchievementCategoryNotFound
	}

	// handle name request
	if req.Name != "" && req.Name != achievementCategory.Name {
		if len(req.Name) < 3 {
			return dto.AchievementCategoryResponse{}, dto.ErrNameTooShort
		}

		achievementCategory.Name = req.Name
	}

	// update achievementCategory
	if err := as.achievementCategoryRepo.Update(ctx, nil, achievementCategory); err != nil {
		return dto.AchievementCategoryResponse{}, dto.ErrUpdateAchievementCategory
	}

	return dto.AchievementCategoryResponse{
		ID:   achievementCategory.ID.String(),
		Name: achievementCategory.Name,
	}, nil
}

func (as *achievementCategoryService) Delete(ctx context.Context, id string) (dto.AchievementCategoryResponse, error) {
	deletedAchievementCategory, found, err := as.achievementCategoryRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.AchievementCategoryResponse{}, dto.ErrAchievementCategoryNotFound
	}
	if !found {
		return dto.AchievementCategoryResponse{}, dto.ErrAchievementCategoryNotFound
	}

	// Delete AchievementCategory
	err = as.achievementCategoryRepo.DeleteByID(ctx, nil, id)
	if err != nil {
		return dto.AchievementCategoryResponse{}, dto.ErrDeleteAchievementCategoryByID
	}

	res := dto.AchievementCategoryResponse{
		ID:   deletedAchievementCategory.ID.String(),
		Name: deletedAchievementCategory.Name,
	}

	return res, nil
}
