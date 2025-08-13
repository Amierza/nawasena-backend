package entity

import (
	"github.com/google/uuid"
)

type Partner struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name  string    `gorm:"type:varchar(150);not null" json:"name"`
	Image string    `json:"image"`

	TimeStamp
}
