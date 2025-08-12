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
	IAdminRepository interface {
		// CREATE / POST
		Create(ctx context.Context, tx *gorm.DB, admin *entity.Admin) error

		// READ / GET
		GetByEmail(ctx context.Context, tx *gorm.DB, email string) (*entity.Admin, bool, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Admin, error)
		GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.AdminPaginationRepositoryResponse, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Admin, bool, error)

		// UPDATE / PATCH
		Update(ctx context.Context, tx *gorm.DB, admin *entity.Admin) error

		// DELETE / DELETE
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
	}

	adminRepository struct {
		db *gorm.DB
	}
)

func NewAdminRepository(db *gorm.DB) *adminRepository {
	return &adminRepository{
		db: db,
	}
}

// CREATE / POST
func (ar *adminRepository) Create(ctx context.Context, tx *gorm.DB, admin *entity.Admin) error {
	if tx == nil {
		tx = ar.db
	}

	return tx.WithContext(ctx).Create(&admin).Error
}

// READ / GET
func (ar *adminRepository) GetByEmail(ctx context.Context, tx *gorm.DB, email string) (*entity.Admin, bool, error) {
	if tx == nil {
		tx = ar.db
	}

	var admin *entity.Admin
	err := tx.WithContext(ctx).Where("email = ?", email).Take(&admin).Error
	if err != nil {
		return &entity.Admin{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Admin{}, false, nil
	}

	return admin, true, nil
}
func (ar *adminRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Admin, error) {
	if tx == nil {
		tx = ar.db
	}

	var (
		admins []*entity.Admin
		err    error
	)

	query := tx.WithContext(ctx).Model(&entity.Admin{}).Where(`role != 'super admin'`)
	if err := query.Order(`"created_at" DESC`).Find(&admins).Error; err != nil {
		return []*entity.Admin{}, err
	}

	return admins, err
}
func (ar *adminRepository) GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.AdminPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = ar.db
	}

	var admins []entity.Admin
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Admin{}).Where(`role != 'super admin'`)

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.AdminPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(Paginate(req.Page, req.PerPage)).Find(&admins).Error; err != nil {
		return dto.AdminPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.AdminPaginationRepositoryResponse{
		Admins: admins,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (ar *adminRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Admin, bool, error) {
	if tx == nil {
		tx = ar.db
	}

	var admin *entity.Admin
	err := tx.WithContext(ctx).Where("id = ?", id).Take(&admin).Error
	if err != nil {
		return &entity.Admin{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Admin{}, false, nil
	}

	return admin, true, nil
}

// UPDATE / PATCH
func (ar *adminRepository) Update(ctx context.Context, tx *gorm.DB, admin *entity.Admin) error {
	if tx == nil {
		tx = ar.db
	}

	return tx.WithContext(ctx).Where("id = ?", admin.ID).Updates(&admin).Error
}

// DELETE / DELETE
func (ar *adminRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = ar.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.Admin{}).Error
}
