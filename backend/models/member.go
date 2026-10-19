package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Member struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	FamilyID      uuid.UUID `gorm:"type:uuid" json:"family_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Password      string    `json:"-"`      // Never send password to frontend
	Role          string    `json:"role"`   // coordinator, secretary, welfare_lead, member
	Status        string    `json:"status"` // active, pending_invite, inactive
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	JoinedAt      time.Time `json:"joined_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relations
	Family        Family         `json:"family,omitempty" gorm:"foreignKey:FamilyID;references:ID"`
	Attendances   []Attendance   `json:"attendances,omitempty"`
	Contributions []Contribution `json:"contributions,omitempty"`
}

func (Member) TableName() string {
	return "members"
}

func (m *Member) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if m.JoinedAt.IsZero() {
		m.JoinedAt = time.Now()
	}
	return nil
}
