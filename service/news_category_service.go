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
	INewsCategoryService interface {
		Create(ctx context.Context, req dto.CreateNewsCategoryRequest) (dto.NewsCategoryResponse, error)
		GetAll(ctx context.Context) ([]dto.NewsCategoryResponse, error)
		GetDetail(ctx context.Context, id string) (dto.NewsCategoryResponse, error)
		Update(ctx context.Context, req dto.UpdateNewsCategoryRequest) (dto.NewsCategoryResponse, error)
		Delete(ctx context.Context, id string) (dto.NewsCategoryResponse, error)
	}

	newsCategoryService struct {
		newsCategoryRepo repository.INewsCategoryRepository
		jwt              jwt.IJWT
	}
)

func NewNewsCategoryService(newsCategoryRepo repository.INewsCategoryRepository, jwt jwt.IJWT) *newsCategoryService {
	return &newsCategoryService{
		newsCategoryRepo: newsCategoryRepo,
		jwt:              jwt,
	}
}

func (as *newsCategoryService) Create(ctx context.Context, req dto.CreateNewsCategoryRequest) (dto.NewsCategoryResponse, error) {
	// handle name request
	if req.Name == "" {
		return dto.NewsCategoryResponse{}, dto.ErrEmptyName
	}
	if len(req.Name) < 3 {
		return dto.NewsCategoryResponse{}, dto.ErrNameTooShort
	}

	// handle double data
	_, found, _ := as.newsCategoryRepo.GetByName(ctx, nil, req.Name)
	if found {
		return dto.NewsCategoryResponse{}, dto.ErrNewsCategoryAlreadyExists
	}

	newsCategoryID := uuid.New()
	newsCategory := &entity.NewsCategory{
		ID:   newsCategoryID,
		Name: req.Name,
	}

	// create newsCategory
	if err := as.newsCategoryRepo.Create(ctx, nil, newsCategory); err != nil {
		return dto.NewsCategoryResponse{}, dto.ErrCreateNewsCategory
	}

	return dto.NewsCategoryResponse{
		ID:   newsCategory.ID.String(),
		Name: newsCategory.Name,
	}, nil
}

func (as *newsCategoryService) GetAll(ctx context.Context) ([]dto.NewsCategoryResponse, error) {
	newsCategorys, err := as.newsCategoryRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, dto.ErrGetAllNewsCategory
	}

	var datas []dto.NewsCategoryResponse
	for _, newsCategory := range newsCategorys {
		data := dto.NewsCategoryResponse{
			ID:   newsCategory.ID.String(),
			Name: newsCategory.Name,
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (as *newsCategoryService) GetDetail(ctx context.Context, id string) (dto.NewsCategoryResponse, error) {
	newsCategory, _, err := as.newsCategoryRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.NewsCategoryResponse{}, dto.ErrNewsCategoryNotFound
	}

	res := dto.NewsCategoryResponse{
		ID:   newsCategory.ID.String(),
		Name: newsCategory.Name,
	}

	return res, nil
}

func (as *newsCategoryService) Update(ctx context.Context, req dto.UpdateNewsCategoryRequest) (dto.NewsCategoryResponse, error) {
	// get newsCategory by id
	newsCategory, found, err := as.newsCategoryRepo.GetByID(ctx, nil, req.ID)
	if err != nil {
		return dto.NewsCategoryResponse{}, dto.ErrGetNewsCategoryByID
	}
	if !found {
		return dto.NewsCategoryResponse{}, dto.ErrNewsCategoryNotFound
	}

	// handle name request
	if req.Name != "" && req.Name != newsCategory.Name {
		if len(req.Name) < 3 {
			return dto.NewsCategoryResponse{}, dto.ErrNameTooShort
		}

		newsCategory.Name = req.Name
	}

	// update newsCategory
	if err := as.newsCategoryRepo.Update(ctx, nil, newsCategory); err != nil {
		return dto.NewsCategoryResponse{}, dto.ErrUpdateNewsCategory
	}

	return dto.NewsCategoryResponse{
		ID:   newsCategory.ID.String(),
		Name: newsCategory.Name,
	}, nil
}

func (as *newsCategoryService) Delete(ctx context.Context, id string) (dto.NewsCategoryResponse, error) {
	deletedNewsCategory, found, err := as.newsCategoryRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.NewsCategoryResponse{}, dto.ErrNewsCategoryNotFound
	}
	if !found {
		return dto.NewsCategoryResponse{}, dto.ErrNewsCategoryNotFound
	}

	// Delete NewsCategory
	err = as.newsCategoryRepo.DeleteByID(ctx, nil, id)
	if err != nil {
		return dto.NewsCategoryResponse{}, dto.ErrDeleteNewsCategoryByID
	}

	res := dto.NewsCategoryResponse{
		ID:   deletedNewsCategory.ID.String(),
		Name: deletedNewsCategory.Name,
	}

	return res, nil
}
