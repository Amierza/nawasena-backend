package entity

import (
	"time"

	"github.com/google/uuid"
)

type News struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(150);not null" json:"name"`
	Description string    `json:"description"`
	PublishedAt time.Time `gorm:"type:timestamp;not null" json:"date"`
	Location    string    `json:"location"`
	URL         string    `json:"url"`
	Status      string    `gorm:"type:varchar(50)" json:"status"`
	Views       int       `gorm:"default:0" json:"views"`
	Featured    bool      `gorm:"default:false" json:"featured"`

	Images []NewsImage `gorm:"foreignKey:NewsID;constraint:OnDelete:CASCADE" json:"-"`

	NewsCategoryID *uuid.UUID   `gorm:"type:uuid" json:"news_category_id,omitempty"`
	NewsCategory   NewsCategory `gorm:"foreignKey:NewsCategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"news_category,omitempty"`

	TimeStamp
}
