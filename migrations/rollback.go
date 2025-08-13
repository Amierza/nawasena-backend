package migrations

import (
	"github.com/Amierza/nawasena-backend/entity"
	"gorm.io/gorm"
)

func Rollback(db *gorm.DB) error {
	tables := []interface{}{
		&entity.Member{},
		&entity.Position{},

		&entity.CompetitionImage{},
		&entity.Competition{},

		&entity.ShipImage{},
		&entity.Ship{},

		&entity.AchievementImage{},
		&entity.Achievement{},

		&entity.Admin{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return err
		}
	}

	return nil
}
