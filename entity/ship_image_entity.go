package entity

import "github.com/google/uuid"

type ShipImage struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name string    `gorm:"type:varchar(150);not null" json:"name"`

	ShipID *uuid.UUID `gorm:"type:uuid" json:"ship_id,omitempty"`
	Ship   Ship       `gorm:"foreignKey:ShipID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ship,omitempty"`

	TimeStamp
}
