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
		GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.News, bool, error)
		GetImagesByID(ctx context.Context, tx *gorm.DB, id string) ([]*entity.NewsImage, error)

		// UPDATE / PATCH
		Update(ctx context.Context, tx *gorm.DB, news *entity.News) error

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

func (ar *newsRepository) RunInTransaction(ctx context.Context, fn func(txRepo INewsRepository) error) error {
	return ar.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &newsRepository{db: tx}
		return fn(txRepo)
	})
}

// CREATE / POST
func (pr *newsRepository) Create(ctx context.Context, tx *gorm.DB, news *entity.News) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&news).Error
}
func (pr *newsRepository) CreateImage(ctx context.Context, tx *gorm.DB, image *entity.NewsImage) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&image).Error
}

// READ / GET
func (pr *newsRepository) GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.News, bool, error) {
	if tx == nil {
		tx = pr.db
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
func (pr *newsRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.News, error) {
	if tx == nil {
		tx = pr.db
	}

	var (
		newss []*entity.News
		err   error
	)

	query := tx.WithContext(ctx).Model(&entity.News{}).Preload("Images")
	if err := query.Order(`"created_at" DESC`).Find(&newss).Error; err != nil {
		return []*entity.News{}, err
	}

	return newss, err
}
func (pr *newsRepository) GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.NewsPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = pr.db
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

	query := tx.WithContext(ctx).Model(&entity.News{}).Preload("Images")

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
func (pr *newsRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.News, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var news *entity.News
	err := tx.WithContext(ctx).Preload("Images").Where("id = ?", id).Take(&news).Error
	if err != nil {
		return &entity.News{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.News{}, false, nil
	}

	return news, true, nil
}
func (ar *newsRepository) GetImagesByID(ctx context.Context, tx *gorm.DB, id string) ([]*entity.NewsImage, error) {
	if tx == nil {
		tx = ar.db
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
func (pr *newsRepository) Update(ctx context.Context, tx *gorm.DB, news *entity.News) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", news.ID).Updates(&news).Error
}

// DELETE / DELETE
func (pr *newsRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.News{}).Error
}
func (ar *newsRepository) DeleteImagesByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = ar.db
	}

	return tx.WithContext(ctx).Where("news_id = ?", id).Delete(&entity.NewsImage{}).Error
}
