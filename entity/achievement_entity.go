package entity

import (
	"github.com/google/uuid"
)

type Achievement struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name"`
	Year        int       `gorm:"not null" json:"year"`
	Description string    `json:"description"`

	Images []AchievementImage `gorm:"foreignKey:AchievementID;constraint:OnDelete:CASCADE" json:"-"`

	TimeStamp
}
