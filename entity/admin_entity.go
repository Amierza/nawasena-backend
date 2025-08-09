package entity

import (
	"github.com/Amierza/nawasena-backend/helper"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name"`
	Email       string    `gorm:"unique;not null" json:"email"`
	Password    string    `gorm:"not null" json:"password"`
	PhoneNumber string    `gorm:"type:varchar(20)" json:"phone_number"`

	TimeStamp
}

func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	var err error

	a.Password, err = helper.HashPassword(a.Password)
	if err != nil {
		return err
	}

	a.PhoneNumber, err = helper.StandardizePhoneNumber(a.PhoneNumber)
	if err != nil {
		return err
	}

	return nil
}
