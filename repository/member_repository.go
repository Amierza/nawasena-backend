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
	IMemberRepository interface {
		// CREATE / POST
		Create(ctx context.Context, tx *gorm.DB, member *entity.Member) error

		// READ / GET
		GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Member, error)
		GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.MemberPaginationRepositoryResponse, error)
		GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Member, bool, error)
		GetByNameMajorGenerationAndPositionID(ctx context.Context, tx *gorm.DB, name, major, positionID string, generation int) (*entity.Member, bool, error)
		GetPositionByPositionID(ctx context.Context, tx *gorm.DB, positionID string) (*entity.Position, bool, error)

		// UPDATE / PATCH
		Update(ctx context.Context, tx *gorm.DB, member *entity.Member) error

		// DELETE / DELETE
		DeleteByID(ctx context.Context, tx *gorm.DB, id string) error
	}

	memberRepository struct {
		db *gorm.DB
	}
)

func NewMemberRepository(db *gorm.DB) *memberRepository {
	return &memberRepository{
		db: db,
	}
}

// CREATE / POST
func (mr *memberRepository) Create(ctx context.Context, tx *gorm.DB, member *entity.Member) error {
	if tx == nil {
		tx = mr.db
	}

	return tx.WithContext(ctx).Create(&member).Error
}

// READ / GET
func (mr *memberRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]*entity.Member, error) {
	if tx == nil {
		tx = mr.db
	}

	var (
		members []*entity.Member
		err     error
	)

	query := tx.WithContext(ctx).Preload("Position").Model(&entity.Member{})
	if err := query.Order(`"created_at" DESC`).Find(&members).Error; err != nil {
		return []*entity.Member{}, err
	}

	return members, err
}
func (mr *memberRepository) GetAllWithPagination(ctx context.Context, tx *gorm.DB, req response.PaginationRequest) (dto.MemberPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = mr.db
	}

	var members []entity.Member
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	query := tx.WithContext(ctx).Model(&entity.Member{}).Preload("Position")

	if req.Search != "" {
		searchValue := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ?", searchValue)
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.MemberPaginationRepositoryResponse{}, err
	}

	if err := query.Order(`"created_at" DESC`).Scopes(Paginate(req.Page, req.PerPage)).Find(&members).Error; err != nil {
		return dto.MemberPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return dto.MemberPaginationRepositoryResponse{
		Members: members,
		PaginationResponse: response.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (mr *memberRepository) GetByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Member, bool, error) {
	if tx == nil {
		tx = mr.db
	}

	var member *entity.Member
	err := tx.WithContext(ctx).Preload("Position").Where("id = ?", id).Take(&member).Error
	if err != nil {
		return &entity.Member{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Member{}, false, nil
	}

	return member, true, nil
}
func (mr *memberRepository) GetByNameMajorGenerationAndPositionID(ctx context.Context, tx *gorm.DB, name, major, positionID string, generation int) (*entity.Member, bool, error) {
	if tx == nil {
		tx = mr.db
	}

	var member *entity.Member
	err := tx.WithContext(ctx).
		Where("name = ? AND major = ? AND generation = ? AND position_id = ?", name, major, generation, positionID).
		Take(&member).Error
	if err != nil {
		return &entity.Member{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Member{}, false, nil
	}

	return member, true, nil
}
func (mr *memberRepository) GetPositionByPositionID(ctx context.Context, tx *gorm.DB, positionID string) (*entity.Position, bool, error) {
	if tx == nil {
		tx = mr.db
	}

	var position *entity.Position
	err := tx.WithContext(ctx).Where("id = ?", positionID).Take(&position).Error
	if err != nil {
		return &entity.Position{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Position{}, false, nil
	}

	return position, true, nil
}

// UPDATE / PATCH
func (mr *memberRepository) Update(ctx context.Context, tx *gorm.DB, member *entity.Member) error {
	if tx == nil {
		tx = mr.db
	}

	return tx.WithContext(ctx).Where("id = ?", member.ID).Updates(&member).Error
}

// DELETE / DELETE
func (mr *memberRepository) DeleteByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = mr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.Member{}).Error
}
