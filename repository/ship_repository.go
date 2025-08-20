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
	IShipRepository interface {
		RunInTransaction(ctx context.Context, fn func(txRepo IShipRepository) error) error

		// CREATE / POST
		Create(ctx context.Context, tx *gorm.DB, ship *entity.Ship) error
		CreateImage(ctx context.Context, tx *gorm.DB, image *entity.ShipImage) error

		// READ / GET
		GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.Ship, bool, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Ship, error)
		GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.ShipPaginationRepositoryResponse, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Ship, bool, error)
		GetImagesByID(ctx context.Context, tx *gorm.DB, id string) ([]*entity.ShipImage, error)

		// UPDATE / PATCH
		Update(ctx context.Context, tx *gorm.DB, ship *entity.Ship) error

		// DELETE / DELETE
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
		DeleteImagesByID(ctx context.Context, tx *gorm.DB, id string) error
	}

	shipRepository struct {
		db *gorm.DB
	}
)

func NewShipRepository(db *gorm.DB) *shipRepository {
	return &shipRepository{
		db: db,
	}
}

func (ar *shipRepository) RunInTransaction(ctx context.Context, fn func(txRepo IShipRepository) error) error {
	return ar.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &shipRepository{db: tx}
		return fn(txRepo)
	})
}

// CREATE / POST
func (pr *shipRepository) Create(ctx context.Context, tx *gorm.DB, ship *entity.Ship) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&ship).Error
}
func (pr *shipRepository) CreateImage(ctx context.Context, tx *gorm.DB, image *entity.ShipImage) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Model(&entity.ShipImage{}).Create(&image).Error
}

// READ / GET
func (pr *shipRepository) GetByName(ctx context.Context, tx *gorm.DB, name string) (*entity.Ship, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var ship *entity.Ship
	err := tx.WithContext(ctx).Where("name = ?", name).Take(&ship).Error
	if err != nil {
		return &entity.Ship{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Ship{}, false, nil
	}

	return ship, true, nil
}
func (pr *shipRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Ship, error) {
	if tx == nil {
		tx = pr.db
	}

	var (
		ships []*entity.Ship
		err   error
	)

	query := tx.WithContext(ctx).Model(&entity.Ship{}).Preload("Images")
	if err := query.Order(`"created_at" DESC`).Find(&ships).Error; err != nil {
		return []*entity.Ship{}, err
	}

	return ships, err
}
func (pr *shipRepository) GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.ShipPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var ships []entity.Ship
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Ship{}).Preload("Images")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.ShipPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(Paginate(req.Page, req.PerPage)).Find(&ships).Error; err != nil {
		return dto.ShipPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.ShipPaginationRepositoryResponse{
		Ships: ships,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (pr *shipRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Ship, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var ship *entity.Ship
	err := tx.WithContext(ctx).Preload("Images").Where("id = ?", id).Take(&ship).Error
	if err != nil {
		return &entity.Ship{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Ship{}, false, nil
	}

	return ship, true, nil
}
func (ar *shipRepository) GetImagesByID(ctx context.Context, tx *gorm.DB, id string) ([]*entity.ShipImage, error) {
	if tx == nil {
		tx = ar.db
	}

	query := tx.WithContext(ctx).
		Model(&entity.ShipImage{}).
		Preload("Ship").
		Where("ship_id = ?", id)

	var shipImages []*entity.ShipImage
	if err := query.Find(&shipImages).Error; err != nil {
		return []*entity.ShipImage{}, err
	}

	return shipImages, nil
}

// UPDATE / PATCH
func (pr *shipRepository) Update(ctx context.Context, tx *gorm.DB, ship *entity.Ship) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", ship.ID).Updates(&ship).Error
}

// DELETE / DELETE
func (pr *shipRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.Ship{}).Error
}
func (ar *shipRepository) DeleteImagesByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = ar.db
	}

	return tx.WithContext(ctx).Where("ship_id = ?", id).Delete(&entity.ShipImage{}).Error
}
