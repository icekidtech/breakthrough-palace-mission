package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Minute struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	FamilyID    uuid.UUID `gorm:"type:uuid" json:"family_id"`
	MeetingID   uuid.UUID `gorm:"type:uuid" json:"meeting_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"` // Markdown or HTML
	PDFUrl      string    `json:"pdf_url"`
	PublishedBy uuid.UUID `gorm:"type:uuid" json:"published_by"` // Member ID
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Family    Family  `json:"family,omitempty" gorm:"foreignKey:FamilyID;references:ID"`
	Meeting   Meeting `json:"meeting,omitempty" gorm:"foreignKey:MeetingID;references:ID"`
	Publisher Member  `json:"publisher,omitempty" gorm:"foreignKey:PublishedBy;references:ID"`
}

func (Minute) TableName() string {
	return "minutes"
}

func (m *Minute) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
