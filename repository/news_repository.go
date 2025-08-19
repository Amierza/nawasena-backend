package repository

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/entity"
	"github.com/Amierza/nawasena-backend/response"
	"gorm.io/gorm"
)

type (
	INewsRepository interface {
		RunInTransaction(ctx context.Context, fn func(txRepo INewsRepository) error) error

		// CREATE / POST
		Create(ctx context.Context, tx *gorm.DB, news *entity.News) error
		CreateImage(ctx context.Context, tx *gorm.DB, image *entity.NewsImage) error

		// READ / GET
		GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.News, bool, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.News, error)
		GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.NewsPaginationRepositoryResponse, error)
		GetFeatured(ctx context.Context, tx *gorm.DB, limit *int) ([]*entity.News, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.News, bool, error)
		GetCategoryByCategoryID(ctx context.Context, tx *gorm.DB, categoryID string) (*entity.NewsCategory, bool, error)
		GetImagesByID(ctx context.Context, tx *gorm.DB, id string) ([]*entity.NewsImage, error)

		// UPDATE / PATCH
		Update(ctx context.Context, tx *gorm.DB, news *entity.News) error
		IncrementViews(ctx context.Context, tx *gorm.DB, id string) error

		// DELETE / DELETE
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
		DeleteImagesByID(ctx context.Context, tx *gorm.DB, id string) error
	}

	newsRepository struct {
		db *gorm.DB
	}
)

func NewNewsRepository(db *gorm.DB) *newsRepository {
	return &newsRepository{
		db: db,
	}
}

func (nr *newsRepository) RunInTransaction(ctx context.Context, fn func(txRepo INewsRepository) error) error {
	return nr.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &newsRepository{db: tx}
		return fn(txRepo)
	})
}

// CREATE / POST
func (nr *newsRepository) Create(ctx context.Context, tx *gorm.DB, news *entity.News) error {
	if tx == nil {
		tx = nr.db
	}

	return tx.WithContext(ctx).Create(&news).Error
}
func (nr *newsRepository) CreateImage(ctx context.Context, tx *gorm.DB, image *entity.NewsImage) error {
	if tx == nil {
		tx = nr.db
	}

	return tx.WithContext(ctx).Create(&image).Error
}

// READ / GET
func (nr *newsRepository) GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.News, bool, error) {
	if tx == nil {
		tx = nr.db
	}

	var news *entity.News
	err := tx.WithContext(ctx).Where("name = ?", name).Take(&news).Error
	if err != nil {
		return &entity.News{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.News{}, false, nil
	}

	return news, true, nil
}
func (nr *newsRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.News, error) {
	if tx == nil {
		tx = nr.db
	}

	var (
		newss []*entity.News
		err   error
	)

	query := tx.WithContext(ctx).Model(&entity.News{}).Preload("Images").Preload("NewsCategory")
	if err := query.Order(`"created_at" DESC`).Find(&newss).Error; err != nil {
		return []*entity.News{}, err
	}

	return newss, err
}
func (nr *newsRepository) GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.NewsPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = nr.db
	}

	var newss []entity.News
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.News{}).Preload("Images").Preload("NewsCategory")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.NewsPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(Paginate(req.Page, req.PerPage)).Find(&newss).Error; err != nil {
		return dto.NewsPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.NewsPaginationRepositoryResponse{
		Newss: newss,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (nr *newsRepository) GetFeatured(ctx context.Context, tx *gorm.DB, limit *int) ([]*entity.News, error) {
	if tx == nil {
		tx = nr.db
	}

	var news []*entity.News
	query := tx.WithContext(ctx).Model(&entity.News{}).
		Preload("Images").
		Preload("NewsCategory").
		Where("featured = ?", true).
		Order("published_at DESC")

	if limit != nil {
		query = query.Limit(*limit)
	}

	if err := query.Find(&news).Error; err != nil {
		return nil, err
	}

	// Kalau hasil kurang dari limit, ambil tambahan (fallback)
	if limit != nil && len(news) < *limit {
		remaining := *limit - len(news)

		var fallback []*entity.News
		err := tx.WithContext(ctx).Model(&entity.News{}).
			Preload("Images").
			Preload("NewsCategory").
			Where("featured = ?", false). // jangan ambil yang udah featured
			Order("views DESC").
			Limit(remaining).
			Find(&fallback).Error
		if err != nil {
			return nil, err
		}

		news = append(news, fallback...)
	}

	return news, nil
}
func (nr *newsRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.News, bool, error) {
	if tx == nil {
		tx = nr.db
	}

	var news *entity.News
	err := tx.WithContext(ctx).Preload("Images").Preload("NewsCategory").Where("id = ?", id).Take(&news).Error
	if err != nil {
		return &entity.News{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.News{}, false, nil
	}

	return news, true, nil
}
func (nr *newsRepository) GetCategoryByCategoryID(ctx context.Context, tx *gorm.DB, categoryID string) (*entity.NewsCategory, bool, error) {
	if tx == nil {
		tx = nr.db
	}

	var news *entity.NewsCategory
	err := tx.WithContext(ctx).Where("id = ?", categoryID).Take(&news).Error
	if err != nil {
		return &entity.NewsCategory{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.NewsCategory{}, false, nil
	}

	return news, true, nil
}
func (nr *newsRepository) GetImagesByID(ctx context.Context, tx *gorm.DB, id string) ([]*entity.NewsImage, error) {
	if tx == nil {
		tx = nr.db
	}

	query := tx.WithContext(ctx).
		Model(&entity.NewsImage{}).
		Preload("News").
		Where("news_id = ?", id)

	var newsImages []*entity.NewsImage
	if err := query.Find(&newsImages).Error; err != nil {
		return []*entity.NewsImage{}, err
	}

	return newsImages, nil
}

// UPDATE / PATCH
func (nr *newsRepository) Update(ctx context.Context, tx *gorm.DB, news *entity.News) error {
	if tx == nil {
		tx = nr.db
	}

	return tx.WithContext(ctx).Model(&entity.News{}).Where("id = ?", news.ID).Updates(news).Error
}
func (nr *newsRepository) IncrementViews(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = nr.db
	}

	return tx.WithContext(ctx).
		Model(&entity.News{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

// DELETE / DELETE
func (nr *newsRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = nr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.News{}).Error
}
func (nr *newsRepository) DeleteImagesByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = nr.db
	}

	return tx.WithContext(ctx).Where("news_id = ?", id).Delete(&entity.NewsImage{}).Error
}
