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
	IPositionRepository interface {
		// CREATE / POST
		Create(ctx context.Context, tx *gorm.DB, position *entity.Position) error

		// READ / GET
		GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.Position, bool, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Position, error)
		GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.PositionPaginationRepositoryResponse, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Position, bool, error)

		// UPDATE / PATCH
		Update(ctx context.Context, tx *gorm.DB, position *entity.Position) error

		// DELETE / DELETE
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
	}

	positionRepository struct {
		db *gorm.DB
	}
)

func NewPositionRepository(db *gorm.DB) *positionRepository {
	return &positionRepository{
		db: db,
	}
}

// CREATE / POST
func (pr *positionRepository) Create(ctx context.Context, tx *gorm.DB, position *entity.Position) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&position).Error
}

// READ / GET
func (pr *positionRepository) GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.Position, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var position *entity.Position
	err := tx.WithContext(ctx).Where("name = ?", name).Take(&position).Error
	if err != nil {
		return &entity.Position{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Position{}, false, nil
	}

	return position, true, nil
}
func (pr *positionRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Position, error) {
	if tx == nil {
		tx = pr.db
	}

	var (
		positions []*entity.Position
		err       error
	)

	query := tx.WithContext(ctx).Model(&entity.Position{})
	if err := query.Order(`"created_at" DESC`).Find(&positions).Error; err != nil {
		return []*entity.Position{}, err
	}

	return positions, err
}
func (pr *positionRepository) GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.PositionPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var positions []entity.Position
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Position{})

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ?", searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.PositionPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(Paginate(req.Page, req.PerPage)).Find(&positions).Error; err != nil {
		return dto.PositionPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.PositionPaginationRepositoryResponse{
		Positions: positions,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (pr *positionRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Position, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var position *entity.Position
	err := tx.WithContext(ctx).Where("id = ?", id).Take(&position).Error
	if err != nil {
		return &entity.Position{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Position{}, false, nil
	}

	return position, true, nil
}

// UPDATE / PATCH
func (pr *positionRepository) Update(ctx context.Context, tx *gorm.DB, position *entity.Position) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Model(&entity.Position{}).
		Where("id = ?", position.ID).
		Select("Name", "IsTech").
		Updates(position).Error
}

// DELETE / DELETE
func (pr *positionRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.Position{}).Error
}
