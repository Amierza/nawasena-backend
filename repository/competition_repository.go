package repository

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/Amierza/nawasena-backend/dto"
	"github.com/Amierza/nawasena-backend/entity"
	"github.com/Amierza/nawasena-backend/response"
	"gorm.io/gorm"
)

type (
	ICompetitionRepository interface {
		RunInTransaction(ctx context.Context, fn func(txRepo ICompetitionRepository) error) error

		// CREATE / POST
		Create(ctx context.Context, tx *gorm.DB, competition *entity.Competition) error
		CreateImage(ctx context.Context, tx *gorm.DB, image *entity.CompetitionImage) error

		// READ / GET
		GetByNameAndDate(ctx context.Context, tx *gorm.DB, name string, date time.Time) (*entity.Competition, bool, error)
		GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Competition, error)
		GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.CompetitionPaginationRepositoryResponse, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Competition, bool, error)
		GetImagesByID(ctx context.Context, tx *gorm.DB, id string) ([]*entity.CompetitionImage, error)

		// UPDATE / PATCH
		Update(ctx context.Context, tx *gorm.DB, competition *entity.Competition) error

		// DELETE / DELETE
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
		DeleteImagesByID(ctx context.Context, tx *gorm.DB, id string) error
	}

	competitionRepository struct {
		db *gorm.DB
	}
)

func NewCompetitionRepository(db *gorm.DB) *competitionRepository {
	return &competitionRepository{
		db: db,
	}
}

func (ar *competitionRepository) RunInTransaction(ctx context.Context, fn func(txRepo ICompetitionRepository) error) error {
	return ar.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &competitionRepository{db: tx}
		return fn(txRepo)
	})
}

// CREATE / POST
func (pr *competitionRepository) Create(ctx context.Context, tx *gorm.DB, competition *entity.Competition) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&competition).Error
}
func (pr *competitionRepository) CreateImage(ctx context.Context, tx *gorm.DB, image *entity.CompetitionImage) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Create(&image).Error
}

// READ / GET
func (pr *competitionRepository) GetByNameAndDate(ctx context.Context, tx *gorm.DB, name string, date time.Time) (*entity.Competition, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var competition *entity.Competition
	err := tx.WithContext(ctx).Where("name = ? AND DATE(date) = ?", name, date.Format("2006-01-02")).Take(&competition).Error
	if err != nil {
		return &entity.Competition{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Competition{}, false, nil
	}

	return competition, true, nil
}
func (pr *competitionRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Competition, error) {
	if tx == nil {
		tx = pr.db
	}

	var (
		competitions []*entity.Competition
		err          error
	)

	query := tx.WithContext(ctx).Model(&entity.Competition{}).Preload("Images")
	if err := query.Order(`"created_at" DESC`).Find(&competitions).Error; err != nil {
		return []*entity.Competition{}, err
	}

	return competitions, err
}
func (pr *competitionRepository) GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.CompetitionPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = pr.db
	}

	var competitions []entity.Competition
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Competition{}).Preload("Images")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchValue, searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.CompetitionPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(Paginate(req.Page, req.PerPage)).Find(&competitions).Error; err != nil {
		return dto.CompetitionPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.CompetitionPaginationRepositoryResponse{
		Competitions: competitions,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (pr *competitionRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Competition, bool, error) {
	if tx == nil {
		tx = pr.db
	}

	var competition *entity.Competition
	err := tx.WithContext(ctx).Preload("Images").Where("id = ?", id).Take(&competition).Error
	if err != nil {
		return &entity.Competition{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Competition{}, false, nil
	}

	return competition, true, nil
}
func (ar *competitionRepository) GetImagesByID(ctx context.Context, tx *gorm.DB, id string) ([]*entity.CompetitionImage, error) {
	if tx == nil {
		tx = ar.db
	}

	query := tx.WithContext(ctx).
		Model(&entity.CompetitionImage{}).
		Preload("Competition").
		Where("competition_id = ?", id)

	var competitionImages []*entity.CompetitionImage
	if err := query.Find(&competitionImages).Error; err != nil {
		return []*entity.CompetitionImage{}, err
	}

	return competitionImages, nil
}

// UPDATE / PATCH
func (pr *competitionRepository) Update(ctx context.Context, tx *gorm.DB, competition *entity.Competition) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", competition.ID).Updates(&competition).Error
}

// DELETE / DELETE
func (pr *competitionRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = pr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.Competition{}).Error
}
func (ar *competitionRepository) DeleteImagesByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = ar.db
	}

	return tx.WithContext(ctx).Where("competition_id = ?", id).Delete(&entity.CompetitionImage{}).Error
}
