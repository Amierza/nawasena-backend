package entity

import "github.com/google/uuid"

type AchievementImage struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(150);not null" json:"name"`

	AchievementID *uuid.UUID  `gorm:"type:uuid" json:"achievement_id,omitempty"`
	Achievement   Achievement `gorm:"foreignKey:AchievementID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"achievement,omitempty"`

	TimeStamp
}
