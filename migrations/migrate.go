package migrations

import (
	"github.com/Amierza/nawasena-backend/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.Admin{},

		&entity.Achievement{},
		&entity.AchievementImage{},

		&entity.Ship{},
		&entity.ShipImage{},

		&entity.Competition{},
		&entity.CompetitionImage{},

		&entity.NewsCategory{},
		&entity.News{},
		&entity.NewsImage{},

		&entity.Partner{},
		&entity.Flyer{},

		&entity.Position{},
		&entity.Member{},
	); err != nil {
		return err
	}

	return nil
}
