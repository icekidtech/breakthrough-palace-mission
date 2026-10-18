package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Contribution struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`
	FamilyID    uuid.UUID `gorm:"type:uuid" json:"family_id"`
	MemberID    uuid.UUID `gorm:"type:uuid" json:"member_id"`
	Type        string    `json:"type"` // family_levy, monthly_welfare, building_fund
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"` // NGN, USD
	Date        time.Time `json:"date"`
	Status      string    `json:"status"`                        // confirmed, pending, rejected
	ConfirmedBy uuid.UUID `gorm:"type:uuid" json:"confirmed_by"` // Member ID
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Family    Family `json:"family,omitempty" gorm:"foreignKey:FamilyID;references:ID"`
	Member    Member `json:"member,omitempty" gorm:"foreignKey:MemberID;references:ID"`
	Confirmer Member `json:"confirmer,omitempty" gorm:"foreignKey:ConfirmedBy;references:ID"`
}

func (Contribution) TableName() string {
	return "contributions"
}

func (c *Contribution) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.Currency == "" {
		c.Currency = "NGN"
	}
	return nil
}
