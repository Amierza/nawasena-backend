package entity

import (
	"github.com/google/uuid"
)

type Ship struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name"`
	Description string    `json:"description"`

	Images []ShipImage `gorm:"foreignKey:ShipID;constraint:OnDelete:CASCADE" json:"-"`

	TimeStamp
}
