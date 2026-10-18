package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Biodata struct {
	ID uuid.UUID `gorm:"primaryKey;type:uuid" json:"id"`

	// Foreign key to Member
	MemberID uuid.UUID `gorm:"type:uuid;uniqueIndex" json:"member_id"`

	// Personal Information
	FullName           string    `json:"full_name"`
	DateOfBirth        *time.Time `json:"date_of_birth,omitempty"`
	Gender             string    `json:"gender"` // Male, Female, Other
	MaritalStatus      string    `json:"marital_status"` // Single, Married, Divorced, Widowed
	Nationality        string    `json:"nationality"`
	PhoneNumber        string    `json:"phone_number"`
	EmailAddress       string    `json:"email_address"`
	HomeAddress        string    `json:"home_address"`
	EmergencyContact   string    `json:"emergency_contact"`
	EmergencyPhone     string    `json:"emergency_phone"`

	// Church Information
	DateJoinedChurch   *time.Time `json:"date_joined_church,omitempty"`
	PreviousChurch     string    `json:"previous_church,omitempty"`
	Department         string    `json:"department"` // Department/unit in church
	FamilyName         string    `json:"family_name"` // Church family name

	// Skills & Profession
	Occupation            string `json:"occupation"`
	CompanyBusinessName   string `json:"company_business_name,omitempty"`
	EducationalQualification string `json:"educational_qualification"`

	// Welfare & Support
	Birthday           *time.Time `json:"birthday,omitempty"`
	WeddingAnniversary *time.Time `json:"wedding_anniversary,omitempty"`
	SpecialOccasions   string    `json:"special_occasions,omitempty"` // JSON array or comma-separated
	WelfareSupportNeeds string   `json:"welfare_support_needs,omitempty"` // Areas where member needs support

	// Optional Information
	Photograph        *string `json:"photograph,omitempty"` // URL to photo
	SocialMediaHandles string  `json:"social_media_handles,omitempty"` // JSON object
	SpecialInterests  string  `json:"special_interests,omitempty"` // Hobbies/interests
	ContactConsent    *bool   `json:"contact_consent,omitempty"` // Consent to be contacted

	// Metadata
	IsComplete bool      `json:"is_complete"` // Flag to indicate if required fields are filled
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  sql.NullTime `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Member *Member `json:"member,omitempty" gorm:"foreignKey:MemberID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Biodata) TableName() string {
	return "biodatas"
}

func (b *Biodata) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	b.UpdateCompletionStatus()
	return nil
}

func (b *Biodata) BeforeSave(tx *gorm.DB) error {
	b.UpdateCompletionStatus()
	return nil
}

// UpdateCompletionStatus checks if all required fields are filled
func (b *Biodata) UpdateCompletionStatus() {
	// Required fields: full_name, phone_number, email_address, home_address, 
	// emergency_contact, emergency_phone, date_of_birth, gender, marital_status,
	// date_joined_church, department, family_name, occupation, educational_qualification
	requiredFields := []bool{
		b.FullName != "",
		b.PhoneNumber != "",
		b.EmailAddress != "",
		b.HomeAddress != "",
		b.EmergencyContact != "",
		b.EmergencyPhone != "",
		b.DateOfBirth != nil,
		b.Gender != "",
		b.MaritalStatus != "",
		b.DateJoinedChurch != nil,
		b.Department != "",
		b.FamilyName != "",
		b.Occupation != "",
		b.EducationalQualification != "",
	}

	// Check if all required fields are filled
	b.IsComplete = true
	for _, field := range requiredFields {
		if !field {
			b.IsComplete = false
			break
		}
	}
}
