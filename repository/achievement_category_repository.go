package repository

import (
	"context"
	"errors"

	"github.com/Amierza/nawasena-backend/entity"
	"gorm.io/gorm"
)

type (
	IAchievementCategoryRepository interface {
		RunInTransaction(ctx context.Context, fn func(txRepo IAchievementCategoryRepository) error) error

		// CREATE / POST
		Create(ctx context.Context, tx *gorm.DB, achievementCategory *entity.AchievementCategory) error

		// READ / GET
		GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.AchievementCategory, bool, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.AchievementCategory, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.AchievementCategory, bool, error)

		// UPDATE / PATCH
		Update(ctx context.Context, tx *gorm.DB, achievementCategory *entity.AchievementCategory) error

		// DELETE / DELETE
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
	}

	achievementCategoryRepository struct {
		db *gorm.DB
	}
)

func NewAchievementCategoryRepository(db *gorm.DB) *achievementCategoryRepository {
	return &achievementCategoryRepository{
		db: db,
	}
}

func (ar *achievementCategoryRepository) RunInTransaction(ctx context.Context, fn func(txRepo IAchievementCategoryRepository) error) error {
	return ar.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &achievementCategoryRepository{db: tx}
		return fn(txRepo)
	})
}

// CREATE / POST
func (pr *achievementCategoryRepository) Create(ctx context.Context, tx *gorm.DB, achievementCategory *entity.AchievementCategory) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&achievementCategory).Error
}

// READ / GET
func (pr *achievementCategoryRepository) GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.AchievementCategory, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var achievementCategory *entity.AchievementCategory
	err := tx.WithContext(ctx).Where("name = ?", name).Take(&achievementCategory).Error
	if err != nil {
		return &entity.AchievementCategory{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.AchievementCategory{}, false, nil
	}

	return achievementCategory, true, nil
}
func (pr *achievementCategoryRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.AchievementCategory, error) {
	if tx == nil {
		tx = pr.db
	}

	var (
		achievementCategorys []*entity.AchievementCategory
		err                  error
	)

	query := tx.WithContext(ctx).Model(&entity.AchievementCategory{})
	if err := query.Order(`"created_at" DESC`).Find(&achievementCategorys).Error; err != nil {
		return []*entity.AchievementCategory{}, err
	}

	return achievementCategorys, err
}
func (pr *achievementCategoryRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.AchievementCategory, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var achievementCategory *entity.AchievementCategory
	err := tx.WithContext(ctx).Where("id = ?", id).Take(&achievementCategory).Error
	if err != nil {
		return &entity.AchievementCategory{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.AchievementCategory{}, false, nil
	}

	return achievementCategory, true, nil
}

// UPDATE / PATCH
func (pr *achievementCategoryRepository) Update(ctx context.Context, tx *gorm.DB, achievementCategory *entity.AchievementCategory) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", achievementCategory.ID).Updates(&achievementCategory).Error
}

// DELETE / DELETE
func (pr *achievementCategoryRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.AchievementCategory{}).Error
}
