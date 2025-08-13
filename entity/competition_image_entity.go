package entity

import (
	"github.com/google/uuid"
)

type CompetitionImage struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(150);not null" json:"name"`

	CompetitionID *uuid.UUID  `gorm:"type:uuid" json:"competition_id,omitempty"`
	Competition   Competition `gorm:"foreignKey:CompetitionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"competition,omitempty"`

	TimeStamp
}
