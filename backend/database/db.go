package database

import (
	"log"

	"github.com/icekidtech/breakthrough-palace-mission/backend/config"
	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
	return db
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Family{},
		&models.Member{},
		&models.Role{},
		&models.Permission{},
		&models.RoleAssignment{},
		&models.Meeting{},
		&models.Attendance{},
		&models.Contribution{},
		&models.Minute{},
		&models.Biodata{},
	)
}
