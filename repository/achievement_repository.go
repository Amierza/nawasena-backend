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
	IAchievementRepository interface {
		RunInTransaction(ctx context.Context, fn func(txRepo IAchievementRepository) error) error

		// CREATE / POST
		Create(ctx context.Context, tx *gorm.DB, achievement *entity.Achievement) error
		CreateImage(ctx context.Context, tx *gorm.DB, image *entity.AchievementImage) error

		// READ / GET
		GetByNameAndYear(ctx context.Context, tx *gorm.DB, name string, year int) (*entity.Achievement, bool, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Achievement, error)
		GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.AchievementPaginationRepositoryResponse, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Achievement, bool, error)
		GetCategoryByCategoryID(ctx context.Context, tx *gorm.DB, categoryID string) (*entity.AchievementCategory, bool, error)
		GetFeatured(ctx context.Context, tx *gorm.DB, limit *int) ([]*entity.Achievement, error)
		GetImagesByID(ctx context.Context, tx *gorm.DB, id string) ([]*entity.AchievementImage, error)

		// UPDATE / PATCH
		Update(ctx context.Context, tx *gorm.DB, achievement *entity.Achievement) error

		// DELETE / DELETE
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
		DeleteImagesByID(ctx context.Context, tx *gorm.DB, id string) error
	}

	achievementRepository struct {
		db *gorm.DB
	}
)

func NewAchievementRepository(db *gorm.DB) *achievementRepository {
	return &achievementRepository{
		db: db,
	}
}

func (ar *achievementRepository) RunInTransaction(ctx context.Context, fn func(txRepo IAchievementRepository) error) error {
	return ar.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &achievementRepository{db: tx}
		return fn(txRepo)
	})
}

// CREATE / POST
func (pr *achievementRepository) Create(ctx context.Context, tx *gorm.DB, achievement *entity.Achievement) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&achievement).Error
}
func (pr *achievementRepository) CreateImage(ctx context.Context, tx *gorm.DB, image *entity.AchievementImage) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&image).Error
}

// READ / GET
func (pr *achievementRepository) GetByNameAndYear(ctx context.Context, tx *gorm.DB, name string, year int) (*entity.Achievement, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var achievement *entity.Achievement
	err := tx.WithContext(ctx).Where("name = ? AND year = ?", name, year).Take(&achievement).Error
	if err != nil {
		return &entity.Achievement{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Achievement{}, false, nil
	}

	return achievement, true, nil
}
func (pr *achievementRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Achievement, error) {
	if tx == nil {
		tx = pr.db
	}

	var (
		achievements []*entity.Achievement
		err          error
	)

	query := tx.WithContext(ctx).Model(&entity.Achievement{}).Preload("Images").Preload("AchievementCategory")
	if err := query.Order(`"created_at" DESC`).Find(&achievements).Error; err != nil {
		return []*entity.Achievement{}, err
	}

	return achievements, err
}
func (pr *achievementRepository) GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.AchievementPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var achievements []entity.Achievement
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Achievement{}).Preload("Images").Preload("AchievementCategory")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.AchievementPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(Paginate(req.Page, req.PerPage)).Find(&achievements).Error; err != nil {
		return dto.AchievementPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.AchievementPaginationRepositoryResponse{
		Achievements: achievements,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (pr *achievementRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Achievement, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var achievement *entity.Achievement
	err := tx.WithContext(ctx).Preload("Images").Preload("AchievementCategory").Where("id = ?", id).Take(&achievement).Error
	if err != nil {
		return &entity.Achievement{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Achievement{}, false, nil
	}

	return achievement, true, nil
}
func (ar *achievementRepository) GetCategoryByCategoryID(ctx context.Context, tx *gorm.DB, categoryID string) (*entity.AchievementCategory, bool, error) {
	if tx == nil {
		tx = ar.db
	}

	var achievement *entity.AchievementCategory
	err := tx.WithContext(ctx).Where("id = ?", categoryID).Take(&achievement).Error
	if err != nil {
		return &entity.AchievementCategory{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.AchievementCategory{}, false, nil
	}

	return achievement, true, nil
}
func (ar *achievementRepository) GetFeatured(ctx context.Context, tx *gorm.DB, limit *int) ([]*entity.Achievement, error) {
	if tx == nil {
		tx = ar.db
	}

	var achievement []*entity.Achievement
	query := tx.WithContext(ctx).Model(&entity.Achievement{}).
		Preload("Images").
		Preload("AchievementCategory").
		Where("featured = ?", true).
		Order("created_at DESC")

	if limit != nil {
		query = query.Limit(*limit)
	}

	if err := query.Find(&achievement).Error; err != nil {
		return nil, err
	}

	// Kalau hasil kurang dari limit, ambil tambahan (fallback)
	if limit != nil && len(achievement) < *limit {
		remaining := *limit - len(achievement)

		var fallback []*entity.Achievement
		err := tx.WithContext(ctx).Model(&entity.Achievement{}).
			Preload("Images").
			Preload("AchievementCategory").
			Where("featured = ?", false). // jangan ambil yang udah featured
			Order("created_at DESC").
			Limit(remaining).
			Find(&fallback).Error
		if err != nil {
			return nil, err
		}

		achievement = append(achievement, fallback...)
	}

	return achievement, nil
}
func (ar *achievementRepository) GetImagesByID(ctx context.Context, tx *gorm.DB, id string) ([]*entity.AchievementImage, error) {
	if tx == nil {
		tx = ar.db
	}

	query := tx.WithContext(ctx).
		Model(&entity.AchievementImage{}).
		Preload("Achievement").
		Where("achievement_id = ?", id)

	var achievementImages []*entity.AchievementImage
	if err := query.Find(&achievementImages).Error; err != nil {
		return []*entity.AchievementImage{}, err
	}

	return achievementImages, nil
}

// UPDATE / PATCH
func (ar *achievementRepository) Update(ctx context.Context, tx *gorm.DB, achievement *entity.Achievement) error {
	if tx == nil {
		tx = ar.db
	}

	return tx.WithContext(ctx).Model(&entity.Achievement{}).Where("id = ?", achievement.ID).Updates(achievement).Error
}

// DELETE / DELETE
func (ar *achievementRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = ar.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.Achievement{}).Error
}
func (ar *achievementRepository) DeleteImagesByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = ar.db
	}

	return tx.WithContext(ctx).Where("achievement_id = ?", id).Delete(&entity.AchievementImage{}).Error
}
