package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	Name        string         `gorm:"uniqueIndex" json:"name"` // "members:view", "meetings:schedule"
	Description string         `json:"description"`
	Category    string         `json:"category"` // "members", "meetings", "attendance", "contributions", etc
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Relations
	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

func (Permission) TableName() string {
	return "permissions"
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
