package entity

import (
	"github.com/google/uuid"
)

type NewsCategory struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(150);not null" json:"name"`

	News []News `gorm:"foreignKey:NewsCategoryID;constraint:OnDelete:CASCADE" json:"-"`

	TimeStamp
}
