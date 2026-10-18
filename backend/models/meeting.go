package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Meeting struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	FamilyID  uuid.UUID `gorm:"type:uuid" json:"family_id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"` // general_meeting, committee_meeting, prayer_gathering
	Date      time.Time `json:"date"`
	StartTime string    `json:"start_time"`                  // HH:MM format
	EndTime   string    `json:"end_time"`                    // HH:MM format
	Location  string    `json:"location"`                    // Physical address or "Online · Google Meet"
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"` // Member ID
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Family    Family       `json:"family,omitempty" gorm:"foreignKey:FamilyID;references:ID"`
	Creator   Member       `json:"creator,omitempty" gorm:"foreignKey:CreatedBy;references:ID"`
	Attendees []Attendance `json:"attendees,omitempty"`
	Minutes   []Minute     `json:"minutes,omitempty"`
}

func (Meeting) TableName() string {
	return "meetings"
}

func (m *Meeting) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
