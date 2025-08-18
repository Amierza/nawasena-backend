package entity

import "github.com/google/uuid"

type NewsImage struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(150);not null" json:"name"`

	NewsID *uuid.UUID `gorm:"type:uuid" json:"news_id,omitempty"`
	News   News       `gorm:"foreignKey:NewsID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"news_image,omitempty"`

	TimeStamp
}
