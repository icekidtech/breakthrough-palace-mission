package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VerificationToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MemberID  uuid.UUID `gorm:"type:uuid;not null;index" json:"member_id"`
	Token     string    `gorm:"type:text;not null;uniqueIndex" json:"token"`
	Type      string    `gorm:"type:varchar(50);not null;index" json:"type"` // "email_verification" or "password_reset"
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Member *Member `gorm:"foreignKey:MemberID;constraint:OnDelete:CASCADE" json:"member,omitempty"`
}

func (vt *VerificationToken) BeforeCreate(tx *gorm.DB) error {
	if vt.ID == uuid.Nil {
		vt.ID = uuid.New()
	}
	return nil
}
