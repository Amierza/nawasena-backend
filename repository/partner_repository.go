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
	IPartnerRepository interface {
		RunInTransaction(ctx context.Context, fn func(txRepo IPartnerRepository) error) error

		// CREATE / POST
		Create(ctx context.Context, tx *gorm.DB, partner *entity.Partner) error

		// READ / GET
		GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.Partner, bool, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Partner, error)
		GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.PartnerPaginationRepositoryResponse, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Partner, bool, error)

		// UPDATE / PATCH
		Update(ctx context.Context, tx *gorm.DB, partner *entity.Partner) error

		// DELETE / DELETE
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
	}

	partnerRepository struct {
		db *gorm.DB
	}
)

func NewPartnerRepository(db *gorm.DB) *partnerRepository {
	return &partnerRepository{
		db: db,
	}
}

func (ar *partnerRepository) RunInTransaction(ctx context.Context, fn func(txRepo IPartnerRepository) error) error {
	return ar.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &partnerRepository{db: tx}
		return fn(txRepo)
	})
}

// CREATE / POST
func (pr *partnerRepository) Create(ctx context.Context, tx *gorm.DB, partner *entity.Partner) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&partner).Error
}

// READ / GET
func (pr *partnerRepository) GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.Partner, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var partner *entity.Partner
	err := tx.WithContext(ctx).Where("name = ?", name).Take(&partner).Error
	if err != nil {
		return &entity.Partner{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Partner{}, false, nil
	}

	return partner, true, nil
}
func (pr *partnerRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Partner, error) {
	if tx == nil {
		tx = pr.db
	}

	var (
		partners []*entity.Partner
		err      error
	)

	query := tx.WithContext(ctx).Model(&entity.Partner{})
	if err := query.Order(`"created_at" DESC`).Find(&partners).Error; err != nil {
		return []*entity.Partner{}, err
	}

	return partners, err
}
func (pr *partnerRepository) GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.PartnerPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var partners []entity.Partner
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Partner{})

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ?", searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.PartnerPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(Paginate(req.Page, req.PerPage)).Find(&partners).Error; err != nil {
		return dto.PartnerPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.PartnerPaginationRepositoryResponse{
		Partners: partners,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (pr *partnerRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Partner, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var partner *entity.Partner
	err := tx.WithContext(ctx).Where("id = ?", id).Take(&partner).Error
	if err != nil {
		return &entity.Partner{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Partner{}, false, nil
	}

	return partner, true, nil
}

// UPDATE / PATCH
func (pr *partnerRepository) Update(ctx context.Context, tx *gorm.DB, partner *entity.Partner) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", partner.ID).Updates(&partner).Error
}

// DELETE / DELETE
func (pr *partnerRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.Partner{}).Error
}
