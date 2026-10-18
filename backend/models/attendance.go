package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attendance struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	FamilyID   uuid.UUID `gorm:"type:uuid" json:"family_id"`
	MeetingID  uuid.UUID `gorm:"type:uuid" json:"meeting_id"`
	MemberID   uuid.UUID `gorm:"type:uuid" json:"member_id"`
	Status     string    `json:"status"`                       // present, absent, late
	RecordedBy uuid.UUID `gorm:"type:uuid" json:"recorded_by"` // Member ID who recorded
	RecordedAt time.Time `json:"recorded_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	Family   Family  `json:"family,omitempty" gorm:"foreignKey:FamilyID;references:ID"`
	Meeting  Meeting `json:"meeting,omitempty" gorm:"foreignKey:MeetingID;references:ID"`
	Member   Member  `json:"member,omitempty" gorm:"foreignKey:MemberID;references:ID"`
	Recorder Member  `json:"recorder,omitempty" gorm:"foreignKey:RecordedBy;references:ID"`
}

func (Attendance) TableName() string {
	return "attendance"
}

func (a *Attendance) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.RecordedAt.IsZero() {
		a.RecordedAt = time.Now()
	}
	return nil
}
