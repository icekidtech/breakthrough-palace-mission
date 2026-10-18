package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Family struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string    `json:"name"`
	Founder   string    `json:"founder"`
	Country   string    `json:"country"`
	Branch    string    `json:"branch"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Members       []Member       `json:"members,omitempty"`
	Meetings      []Meeting      `json:"meetings,omitempty"`
	Attendance    []Attendance   `json:"attendance,omitempty"`
	Contributions []Contribution `json:"contributions,omitempty"`
}

func (Family) TableName() string {
	return "families"
}

func (f *Family) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
