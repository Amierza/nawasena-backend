package migrations

import (
	"github.com/Amierza/nawasena-backend/entity"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	err := SeedFromJSON[entity.Admin](db, "./migrations/json/admins.json", entity.Admin{}, "Email")
	if err != nil {
		return err
	}

	err = SeedFromJSON[entity.Position](db, "./migrations/json/positions.json", entity.Position{}, "Name")
	if err != nil {
		return err
	}

	return nil
}
