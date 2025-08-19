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

	err = SeedFromJSON[entity.NewsCategory](db, "./migrations/json/news_categories.json", entity.NewsCategory{}, "Name")
	if err != nil {
		return err
	}

	err = SeedFromJSON[entity.AchievementCategory](db, "./migrations/json/achievement_categories.json", entity.AchievementCategory{}, "Name")
	if err != nil {
		return err
	}

	return nil
}
