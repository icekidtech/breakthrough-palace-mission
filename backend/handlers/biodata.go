package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/icekidtech/breakthrough-palace-mission/backend/middleware"
	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
	"github.com/icekidtech/breakthrough-palace-mission/backend/rbac"
	"github.com/icekidtech/breakthrough-palace-mission/backend/utils"
	"gorm.io/gorm"
)

type BiodataUpdateRequest struct {
	FullName                 string     `json:"full_name"`
	DateOfBirth              *time.Time `json:"date_of_birth,omitempty"`
	Gender                   string     `json:"gender"`
	MaritalStatus            string     `json:"marital_status"`
	Nationality              string     `json:"nationality"`
	PhoneNumber              string     `json:"phone_number"`
	EmailAddress             string     `json:"email_address"`
	HomeAddress              string     `json:"home_address"`
	EmergencyContact         string     `json:"emergency_contact"`
	EmergencyPhone           string     `json:"emergency_phone"`
	DateJoinedChurch         *time.Time `json:"date_joined_church,omitempty"`
	PreviousChurch           string     `json:"previous_church,omitempty"`
	Department               string     `json:"department"`
	FamilyName               string     `json:"family_name"`
	Occupation               string     `json:"occupation"`
	CompanyBusinessName      string     `json:"company_business_name,omitempty"`
	EducationalQualification string     `json:"educational_qualification"`
	Birthday                 *time.Time `json:"birthday,omitempty"`
	WeddingAnniversary       *time.Time `json:"wedding_anniversary,omitempty"`
	SpecialOccasions         string     `json:"special_occasions,omitempty"`
	WelfareSupportNeeds      string     `json:"welfare_support_needs,omitempty"`
	Photograph               *string    `json:"photograph,omitempty"`
	SocialMediaHandles       string     `json:"social_media_handles,omitempty"`
	SpecialInterests         string     `json:"special_interests,omitempty"`
	ContactConsent           *bool      `json:"contact_consent,omitempty"`
}

type BiodataResponse struct {
	Biodata    models.Biodata `json:"biodata"`
	IsComplete bool           `json:"is_complete"`
}

// initializeBiodataFromMember creates a new biodata record pre-filled with member signup data
func initializeBiodataFromMember(db *gorm.DB, memberID string) (*models.Biodata, error) {
	// Fetch member to get signup data
	var member models.Member
	if err := db.Where("id = ?", memberID).First(&member).Error; err != nil {
		return nil, err
	}

	biodata := &models.Biodata{
		MemberID:     utils.ParseUUID(memberID),
		FullName:     member.Name,
		PhoneNumber:  member.Phone,
		EmailAddress: member.Email,
		IsComplete:   false, // Will be updated by BeforeSave hook
	}

	if err := db.Create(biodata).Error; err != nil {
		return nil, err
	}

	return biodata, nil
}

// GetBiodata retrieves the biodata for the current user
func GetBiodata(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memberID := middleware.ExtractMemberID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "biodata:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var biodata models.Biodata
		result := db.Where("member_id = ?", memberID).First(&biodata)

		// If biodata doesn't exist, create one with pre-filled member data
		if result.Error == gorm.ErrRecordNotFound {
			newBiodata, err := initializeBiodataFromMember(db, memberID)
			if err != nil {
				http.Error(w, `{"error":"failed to create biodata"}`, http.StatusInternalServerError)
				return
			}
			biodata = *newBiodata
		} else if result.Error != nil {
			http.Error(w, `{"error":"failed to fetch biodata"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BiodataResponse{
			Biodata:    biodata,
			IsComplete: biodata.IsComplete,
		})
	}
}

// UpdateBiodata updates the biodata for the current user
func UpdateBiodata(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		memberID := middleware.ExtractMemberID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "biodata:edit") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var req BiodataUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		// Get or create biodata
		var biodata models.Biodata
		result := db.Where("member_id = ?", memberID).First(&biodata)

		if result.Error == gorm.ErrRecordNotFound {
			biodata = models.Biodata{
				MemberID: utils.ParseUUID(memberID),
			}
		} else if result.Error != nil {
			http.Error(w, `{"error":"failed to fetch biodata"}`, http.StatusInternalServerError)
			return
		}

		// Update biodata fields
		biodata.FullName = req.FullName
		biodata.DateOfBirth = req.DateOfBirth
		biodata.Gender = req.Gender
		biodata.MaritalStatus = req.MaritalStatus
		biodata.Nationality = req.Nationality
		biodata.PhoneNumber = req.PhoneNumber
		biodata.EmailAddress = req.EmailAddress
		biodata.HomeAddress = req.HomeAddress
		biodata.EmergencyContact = req.EmergencyContact
		biodata.EmergencyPhone = req.EmergencyPhone
		biodata.DateJoinedChurch = req.DateJoinedChurch
		biodata.PreviousChurch = req.PreviousChurch
		biodata.Department = req.Department
		biodata.FamilyName = req.FamilyName
		biodata.Occupation = req.Occupation
		biodata.CompanyBusinessName = req.CompanyBusinessName
		biodata.EducationalQualification = req.EducationalQualification
		biodata.Birthday = req.Birthday
		biodata.WeddingAnniversary = req.WeddingAnniversary
		biodata.SpecialOccasions = req.SpecialOccasions
		biodata.WelfareSupportNeeds = req.WelfareSupportNeeds
		biodata.Photograph = req.Photograph
		biodata.SocialMediaHandles = req.SocialMediaHandles
		biodata.SpecialInterests = req.SpecialInterests
		biodata.ContactConsent = req.ContactConsent

		// Save biodata (BeforeSave hook will update IsComplete)
		if err := db.Save(&biodata).Error; err != nil {
			http.Error(w, `{"error":"failed to update biodata"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BiodataResponse{
			Biodata:    biodata,
			IsComplete: biodata.IsComplete,
		})
	}
}

// CheckBiodataCompletion checks if the user's biodata is complete
func CheckBiodataCompletion(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memberID := middleware.ExtractMemberID(r)

		var biodata models.Biodata
		if err := db.Where("member_id = ?", memberID).First(&biodata).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"is_complete": false,
					"message":     "biodata not yet started",
				})
				return
			}
			http.Error(w, `{"error":"failed to fetch biodata"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"is_complete": biodata.IsComplete,
		})
	}
}

// GetMemberBiodata retrieves biodata for a specific member (admin only)
func GetMemberBiodata(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "biodata:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		// Extract member ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		targetMemberID := parts[len(parts)-1]

		var biodata models.Biodata
		if err := db.Where("member_id = ?", targetMemberID).First(&biodata).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, `{"error":"biodata not found"}`, http.StatusNotFound)
				return
			}
			http.Error(w, `{"error":"failed to fetch biodata"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BiodataResponse{
			Biodata:    biodata,
			IsComplete: biodata.IsComplete,
		})
	}
}

// ListBiodataStats returns stats about biodata completion in a family (admin only)
func ListBiodataStats(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "biodata:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var completedCount int64
		var totalCount int64

		// Count total members in family
		db.Model(&models.Member{}).Where("family_id = ?", familyID).Count(&totalCount)

		// Count members with complete biodata
		db.Model(&models.Biodata{}).
			Joins("JOIN members ON biodatas.member_id = members.id").
			Where("members.family_id = ? AND biodatas.is_complete = ?", familyID, true).
			Count(&completedCount)

		var percentage float64
		if totalCount > 0 {
			percentage = (float64(completedCount) / float64(totalCount)) * 100
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total_members":      totalCount,
			"completed_biodata":  completedCount,
			"pending_biodata":    totalCount - completedCount,
			"completion_percent": percentage,
		})
	}
}
