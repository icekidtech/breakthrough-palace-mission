package migrations

import (
	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
	"gorm.io/gorm"
)

// CreateBiodata creates the biodata table with all necessary fields
func CreateBiodata(db *gorm.DB) error {
	if !db.Migrator().HasTable(&models.Biodata{}) {
		if err := db.Migrator().CreateTable(&models.Biodata{}); err != nil {
			return err
		}
	}
	return nil
}
