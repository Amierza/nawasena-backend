package repository

import (
	"context"
	"errors"

	"github.com/Amierza/nawasena-backend/entity"
	"gorm.io/gorm"
)

type (
	INewsCategoryRepository interface {
		RunInTransaction(ctx context.Context, fn func(txRepo INewsCategoryRepository) error) error

		// CREATE / POST
		Create(ctx context.Context, tx *gorm.DB, newsCategory *entity.NewsCategory) error

		// READ / GET
		GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.NewsCategory, bool, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.NewsCategory, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.NewsCategory, bool, error)

		// UPDATE / PATCH
		Update(ctx context.Context, tx *gorm.DB, newsCategory *entity.NewsCategory) error

		// DELETE / DELETE
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
	}

	newsCategoryRepository struct {
		db *gorm.DB
	}
)

func NewNewsCategoryRepository(db *gorm.DB) *newsCategoryRepository {
	return &newsCategoryRepository{
		db: db,
	}
}

func (ar *newsCategoryRepository) RunInTransaction(ctx context.Context, fn func(txRepo INewsCategoryRepository) error) error {
	return ar.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &newsCategoryRepository{db: tx}
		return fn(txRepo)
	})
}

// CREATE / POST
func (pr *newsCategoryRepository) Create(ctx context.Context, tx *gorm.DB, newsCategory *entity.NewsCategory) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&newsCategory).Error
}

// READ / GET
func (pr *newsCategoryRepository) GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.NewsCategory, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var newsCategory *entity.NewsCategory
	err := tx.WithContext(ctx).Where("name = ?", name).Take(&newsCategory).Error
	if err != nil {
		return &entity.NewsCategory{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.NewsCategory{}, false, nil
	}

	return newsCategory, true, nil
}
func (pr *newsCategoryRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.NewsCategory, error) {
	if tx == nil {
		tx = pr.db
	}

	var (
		newsCategorys []*entity.NewsCategory
		err           error
	)

	query := tx.WithContext(ctx).Model(&entity.NewsCategory{})
	if err := query.Order(`"created_at" DESC`).Find(&newsCategorys).Error; err != nil {
		return []*entity.NewsCategory{}, err
	}

	return newsCategorys, err
}
func (pr *newsCategoryRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.NewsCategory, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var newsCategory *entity.NewsCategory
	err := tx.WithContext(ctx).Where("id = ?", id).Take(&newsCategory).Error
	if err != nil {
		return &entity.NewsCategory{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.NewsCategory{}, false, nil
	}

	return newsCategory, true, nil
}

// UPDATE / PATCH
func (pr *newsCategoryRepository) Update(ctx context.Context, tx *gorm.DB, newsCategory *entity.NewsCategory) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", newsCategory.ID).Updates(&newsCategory).Error
}

// DELETE / DELETE
func (pr *newsCategoryRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.NewsCategory{}).Error
}
