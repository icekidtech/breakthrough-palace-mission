package migrations

import (
	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
	"gorm.io/gorm"
)

func CreateVerificationTokensTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.VerificationToken{})
}

func AddEmailVerifiedToMembers(db *gorm.DB) error {
	return db.Migrator().AddColumn(&models.Member{}, "email_verified")
}
