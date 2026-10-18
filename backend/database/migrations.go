package database

import (
	"github.com/icekidtech/breakthrough-palace-mission/backend/database/migrations"
	"gorm.io/gorm"
)

// RunMigrations executes all data migrations after schema is created
func RunMigrations(db *gorm.DB) error {
	// Seed roles and permissions
	if err := migrations.SeedRolesAndPermissions(db); err != nil {
		return err
	}

	return nil
}
