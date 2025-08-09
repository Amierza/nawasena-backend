package repository

import (
	"context"
	"errors"

	"github.com/Amierza/nawasena-backend/entity"
	"gorm.io/gorm"
)

type (
	IAuthRepository interface {
		GetAdminByEmail(ctx context.Context, tx *gorm.DB, email string) (*entity.Admin, bool, error)
	}

	authRepository struct {
		db *gorm.DB
	}
)

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{
		db: db,
	}
}

func (ar *authRepository) GetAdminByEmail(ctx context.Context, tx *gorm.DB, email string) (*entity.Admin, bool, error) {
	if tx == nil {
		tx = ar.db
	}

	var admin *entity.Admin
	err := tx.WithContext(ctx).Where("email = ?", email).Take(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Admin{}, false, nil
	}
	if err != nil {
		return &entity.Admin{}, false, err
	}

	return admin, true, nil
}
