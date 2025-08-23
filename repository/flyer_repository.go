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
	IFlyerRepository interface {
		RunInTransaction(ctx context.Context, fn func(txRepo IFlyerRepository) error) error

		// CREATE / POST
		Create(ctx context.Context, tx *gorm.DB, flyer *entity.Flyer) error

		// READ / GET
		GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.Flyer, bool, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Flyer, error)
		GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.FlyerPaginationRepositoryResponse, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Flyer, bool, error)

		// UPDATE / PATCH
		Update(ctx context.Context, tx *gorm.DB, flyer *entity.Flyer) error

		// DELETE / DELETE
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
	}

	flyerRepository struct {
		db *gorm.DB
	}
)

func NewFlyerRepository(db *gorm.DB) *flyerRepository {
	return &flyerRepository{
		db: db,
	}
}

func (ar *flyerRepository) RunInTransaction(ctx context.Context, fn func(txRepo IFlyerRepository) error) error {
	return ar.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &flyerRepository{db: tx}
		return fn(txRepo)
	})
}

// CREATE / POST
func (pr *flyerRepository) Create(ctx context.Context, tx *gorm.DB, flyer *entity.Flyer) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&flyer).Error
}

// READ / GET
func (pr *flyerRepository) GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.Flyer, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var flyer *entity.Flyer
	err := tx.WithContext(ctx).Where("name = ?", name).Take(&flyer).Error
	if err != nil {
		return &entity.Flyer{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Flyer{}, false, nil
	}

	return flyer, true, nil
}
func (pr *flyerRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Flyer, error) {
	if tx == nil {
		tx = pr.db
	}

	var (
		flyers []*entity.Flyer
		err    error
	)

	query := tx.WithContext(ctx).Model(&entity.Flyer{})
	if err := query.Order(`"created_at" DESC`).Find(&flyers).Error; err != nil {
		return []*entity.Flyer{}, err
	}

	return flyers, err
}
func (pr *flyerRepository) GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.FlyerPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var flyers []entity.Flyer
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Flyer{})

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ?", searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.FlyerPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(Paginate(req.Page, req.PerPage)).Find(&flyers).Error; err != nil {
		return dto.FlyerPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.FlyerPaginationRepositoryResponse{
		Flyers: flyers,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (pr *flyerRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Flyer, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var flyer *entity.Flyer
	err := tx.WithContext(ctx).Where("id = ?", id).Take(&flyer).Error
	if err != nil {
		return &entity.Flyer{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Flyer{}, false, nil
	}

	return flyer, true, nil
}

// UPDATE / PATCH
func (pr *flyerRepository) Update(ctx context.Context, tx *gorm.DB, flyer *entity.Flyer) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", flyer.ID).Updates(&flyer).Error
}

// DELETE / DELETE
func (pr *flyerRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.Flyer{}).Error
}
