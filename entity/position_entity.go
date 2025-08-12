package entity

import (
	"github.com/google/uuid"
)

type Position struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(150);unique;not null" json:"name"`

	Members []Member `gorm:"foreignKey:PositionID;constraint:OnDelete:CASCADE" json:"-"`

	TimeStamp
}
