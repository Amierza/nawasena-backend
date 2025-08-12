package entity

import (
	"github.com/google/uuid"
)

type Member struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(150);not null" json:"name"`
	Image      string    `json:"image"`
	Major      string    `gorm:"not null" json:"major"`
	Generation string    `gorm:"not null" json:"generation"`

	PositionID *uuid.UUID `gorm:"type:uuid" json:"position_id,omitempty"`
	Position   Position   `gorm:"foreignKey:PositionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"position,omitempty"`

	TimeStamp
}
