package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid" json:"id"`
	Name        string         `gorm:"uniqueIndex" json:"name"` // "secretary", "family_leader", etc
	Description string         `json:"description"`
	Scope       string         `json:"scope"` // "family", "platform", "department"
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Relations
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
	Assignments []RoleAssignment `json:"assignments,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
