package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleAssignment struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	MemberID  uuid.UUID      `gorm:"type:uuid" json:"member_id"`
	RoleID    uuid.UUID      `gorm:"type:uuid" json:"role_id"`
	FamilyID  uuid.UUID      `gorm:"type:uuid;index" json:"family_id"` // Scope: which family (if family role)
	StartDate time.Time      `json:"start_date"`
	EndDate   *time.Time     `json:"end_date"` // NULL = no end date
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// Relations
	Member Family `json:"member,omitempty" gorm:"foreignKey:MemberID;references:ID"`
	Role   Role   `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID"`
	Family Family `json:"family,omitempty" gorm:"foreignKey:FamilyID;references:ID"`
}

func (RoleAssignment) TableName() string {
	return "role_assignments"
}

func (ra *RoleAssignment) BeforeCreate(tx *gorm.DB) error {
	if ra.ID == uuid.Nil {
		ra.ID = uuid.New()
	}
	if ra.StartDate.IsZero() {
		ra.StartDate = time.Now()
	}
	return nil
}
