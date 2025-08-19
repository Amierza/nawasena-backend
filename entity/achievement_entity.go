package entity

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Achievement struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(150);not null" json:"name"`
	Year        int            `gorm:"not null" json:"year"`
	Description string         `json:"description"`
	Location    string         `json:"location"`
	Rank        string         `json:"rank"`
	Competition string         `json:"competition"`
	Team        pq.StringArray `gorm:"type:text[]" json:"team"`
	Impact      string         `json:"impact"`
	VideoURL    string         `json:"video_url"`
	Featured    bool           `gorm:"default:false" json:"featured"`
	Tags        pq.StringArray `gorm:"type:text[]" json:"tags"`

	Images []AchievementImage `gorm:"foreignKey:AchievementID;constraint:OnDelete:CASCADE" json:"-"`

	AchievementCategoryID *uuid.UUID          `gorm:"type:uuid" json:"achievement_category_id,omitempty"`
	AchievementCategory   AchievementCategory `gorm:"foreignKey:AchievementCategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"achievement_category,omitempty"`

	TimeStamp
}
