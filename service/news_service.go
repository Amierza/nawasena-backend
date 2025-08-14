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
	INewsService interface {
		Create(ctx context.Context, req dto.CreateNewsRequest) (dto.NewsResponse, error)
		GetAll(ctx context.Context) ([]dto.NewsResponse, error)
		GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.NewsPaginationResponse, error)
		GetDetail(ctx context.Context, id string) (dto.NewsResponse, error)
		Update(ctx context.Context, req dto.UpdateNewsRequest) (dto.NewsResponse, error)
		Delete(ctx context.Context, id string) (dto.NewsResponse, error)
	}

	newsService struct {
		newsRepo repository.INewsRepository
		jwt      jwt.IJWT
	}
)

func NewNewsService(newsRepo repository.INewsRepository, jwt jwt.IJWT) *newsService {
	return &newsService{
		newsRepo: newsRepo,
		jwt:      jwt,
	}
}

func (as *newsService) Create(ctx context.Context, req dto.CreateNewsRequest) (dto.NewsResponse, error) {
	// handle name request
	if req.Name == "" {
		return dto.NewsResponse{}, dto.ErrEmptyName
	}
	if len(req.Name) < 3 {
		return dto.NewsResponse{}, dto.ErrNameTooShort
	}

	// handle description request
	if req.Description == "" {
		return dto.NewsResponse{}, dto.ErrEmptyDescription
	}
	if len(req.Description) < 5 {
		return dto.NewsResponse{}, dto.ErrDescriptionTooShort
	}

	// handle double data
	_, found, _ := as.newsRepo.GetByName(ctx, nil, req.Name)
	if found {
		return dto.NewsResponse{}, dto.ErrNewsAlreadyExists
	}

	newsID := uuid.New()
	news := &entity.News{
		ID:          newsID,
		Name:        req.Name,
		Description: req.Description,
	}

	// handle image url
	var (
		newsImages         []*entity.NewsImage
		newsImageResponses []dto.NewsImageResponse
	)
	if len(req.Images) == 0 {
		return dto.NewsResponse{}, dto.ErrEmptyImages
	}
	for _, imgName := range req.Images {
		imgID := uuid.New()
		// handle entity
		newsImages = append(newsImages, &entity.NewsImage{
			ID:     imgID,
			Name:   imgName,
			NewsID: &newsID,
		})

		// handle response
		newsImageResponses = append(newsImageResponses, dto.NewsImageResponse{
			ID:   imgID.String(),
			Name: imgName,
		})
	}

	err := as.newsRepo.RunInTransaction(ctx, func(txRepo repository.INewsRepository) error {
		// create news
		if err := txRepo.Create(ctx, nil, news); err != nil {
			return dto.ErrCreateNews
		}

		// handle new image
		if len(req.Name) > 0 {
			// check request images
			oldImages, err := txRepo.GetImagesByID(ctx, nil, news.ID.String())
			if err != nil {
				return dto.ErrGetNewsImages
			}

			// Delete Existing News Image
			// in assets
			for _, img := range oldImages {
				name := strings.TrimPrefix(img.Name, "assets/")
				path := filepath.Join("assets", name)
				if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
					return dto.ErrDeleteOldImage
				}
			}
			// in db
			if err := txRepo.DeleteImagesByID(ctx, nil, news.ID.String()); err != nil {
				return dto.ErrDeleteNewsImageByNewsID
			}

			// Create new news images
			for _, img := range newsImages {
				if err := txRepo.CreateImage(ctx, nil, img); err != nil {
					return dto.ErrCreateNewsImage
				}
			}
		}

		return nil
	})
	if err != nil {
		return dto.NewsResponse{}, err
	}

	return dto.NewsResponse{
		ID:          news.ID.String(),
		Name:        news.Name,
		Description: news.Description,
		Images:      newsImageResponses,
	}, nil
}

func (as *newsService) GetAll(ctx context.Context) ([]dto.NewsResponse, error) {
	newss, err := as.newsRepo.GetAll(ctx, nil)
	if err != nil {
		return nil, dto.ErrGetAllNewsNoPagination
	}

	var datas []dto.NewsResponse
	for _, news := range newss {
		data := dto.NewsResponse{
			ID:          news.ID.String(),
			Name:        news.Name,
			Description: news.Description,
		}

		for _, a := range news.Images {
			data.Images = append(data.Images, dto.NewsImageResponse{
				ID:   a.ID.String(),
				Name: a.Name,
			})
		}

		datas = append(datas, data)
	}

	return datas, nil
}

func (as *newsService) GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.NewsPaginationResponse, error) {
	dataWithPaginate, err := as.newsRepo.GetAllWithPagination(ctx, nil, req)
	if err != nil {
		return dto.NewsPaginationResponse{}, dto.ErrGetAllNewsWithPagination
	}

	var datas []dto.NewsResponse
	for _, news := range dataWithPaginate.Newss {
		data := dto.NewsResponse{
			ID:          news.ID.String(),
			Name:        news.Name,
			Description: news.Description,
		}

		for _, a := range news.Images {
			data.Images = append(data.Images, dto.NewsImageResponse{
				ID:   a.ID.String(),
				Name: a.Name,
			})
		}

		datas = append(datas, data)
	}

	return dto.NewsPaginationResponse{
		Data: datas,
		PaginationResponse: response.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (as *newsService) GetDetail(ctx context.Context, id string) (dto.NewsResponse, error) {
	news, _, err := as.newsRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.NewsResponse{}, dto.ErrNewsNotFound
	}

	res := dto.NewsResponse{
		ID:          news.ID.String(),
		Name:        news.Name,
		Description: news.Description,
	}

	for _, a := range news.Images {
		res.Images = append(res.Images, dto.NewsImageResponse{
			ID:   a.ID.String(),
			Name: a.Name,
		})
	}

	return res, nil
}

func (as *newsService) Update(ctx context.Context, req dto.UpdateNewsRequest) (dto.NewsResponse, error) {
	// get news by id
	news, found, err := as.newsRepo.GetByID(ctx, nil, req.ID)
	if err != nil {
		return dto.NewsResponse{}, dto.ErrGetNewsByID
	}
	if !found {
		return dto.NewsResponse{}, dto.ErrNewsNotFound
	}

	// handle name request
	if req.Name != "" && req.Name != news.Name {
		if len(req.Name) < 3 {
			return dto.NewsResponse{}, dto.ErrNameTooShort
		}

		news.Name = req.Name
	}

	// handle description request
	if req.Description != "" && req.Description != news.Description {
		if len(req.Description) < 5 {
			return dto.NewsResponse{}, dto.ErrDescriptionTooShort
		}

		news.Description = req.Description
	}

	// handle image url
	var (
		newsImages         []*entity.NewsImage
		newsImageResponses []dto.NewsImageResponse
	)
	if len(req.Images) > 0 {
		for _, imgName := range req.Images {
			imgID := uuid.New()
			// handle entity
			newsImages = append(newsImages, &entity.NewsImage{
				ID:     imgID,
				Name:   imgName,
				NewsID: &news.ID,
			})

			// handle response
			newsImageResponses = append(newsImageResponses, dto.NewsImageResponse{
				ID:   imgID.String(),
				Name: imgName,
			})
		}
	}

	err = as.newsRepo.RunInTransaction(ctx, func(txRepo repository.INewsRepository) error {
		// update news
		if err := txRepo.Update(ctx, nil, news); err != nil {
			return dto.ErrUpdateNews
		}

		// handle new image
		if len(req.Name) > 0 {
			// check request images
			oldImages, err := txRepo.GetImagesByID(ctx, nil, news.ID.String())
			if err != nil {
				return dto.ErrGetNewsImages
			}

			// Delete Existing News Image
			// in assets
			for _, img := range oldImages {
				if err := os.Remove(img.Name); err != nil && !os.IsNotExist(err) {
					return dto.ErrDeleteOldImage
				}
			}
			// in db
			if err := txRepo.DeleteImagesByID(ctx, nil, news.ID.String()); err != nil {
				return dto.ErrDeleteNewsImageByNewsID
			}

			// Create new news images
			for _, img := range newsImages {
				if err := txRepo.CreateImage(ctx, nil, img); err != nil {
					return dto.ErrCreateNewsImage
				}
			}
		}

		return nil
	})
	if err != nil {
		return dto.NewsResponse{}, err
	}

	return dto.NewsResponse{
		ID:          news.ID.String(),
		Name:        news.Name,
		Description: news.Description,
		Images:      newsImageResponses,
	}, nil
}

func (as *newsService) Delete(ctx context.Context, id string) (dto.NewsResponse, error) {
	deletedNews, found, err := as.newsRepo.GetByID(ctx, nil, id)
	if err != nil {
		return dto.NewsResponse{}, dto.ErrNewsNotFound
	}
	if !found {
		return dto.NewsResponse{}, dto.ErrNewsNotFound
	}

	err = as.newsRepo.RunInTransaction(ctx, func(txRepo repository.INewsRepository) error {
		// Delete News Images
		oldNewsImages, err := txRepo.GetImagesByID(ctx, nil, id)
		if err != nil {
			return dto.ErrGetNewsImages
		}
		for _, img := range oldNewsImages {
			name := strings.TrimPrefix(img.Name, "assets/")
			path := filepath.Join("assets", name)
			if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
				return dto.ErrDeleteOldImage
			}
		}
		if err := txRepo.DeleteImagesByID(ctx, nil, id); err != nil {
			return dto.ErrDeleteNewsImageByNewsID
		}

		// Delete News
		err = as.newsRepo.DeleteByID(ctx, nil, id)
		if err != nil {
			return dto.ErrDeleteNewsByID
		}

		return nil
	})
	if err != nil {
		return dto.NewsResponse{}, err
	}

	res := dto.NewsResponse{
		ID:          deletedNews.ID.String(),
		Name:        deletedNews.Name,
		Description: deletedNews.Description,
	}

	for _, a := range deletedNews.Images {
		res.Images = append(res.Images, dto.NewsImageResponse{
			ID:   a.ID.String(),
			Name: a.Name,
		})
	}

	return res, nil
}
