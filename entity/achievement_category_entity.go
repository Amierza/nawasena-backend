package entity

import (
	"github.com/google/uuid"
)

type AchievementCategory struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(150);not null" json:"name"`

	Achievement []Achievement `gorm:"foreignKey:AchievementCategoryID;constraint:OnDelete:CASCADE" json:"-"`

	TimeStamp
}
