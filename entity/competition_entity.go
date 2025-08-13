package entity

import (
	"time"

	"github.com/google/uuid"
)

type Competition struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`

	Images []CompetitionImage `gorm:"foreignKey:CompetitionID;constraint:OnDelete:CASCADE" json:"-"`

	TimeStamp
}
