package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/icekidtech/breakthrough-palace-mission/backend/middleware"
	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
	"gorm.io/gorm"
)

// ListFamilies returns all families (public endpoint - no auth required, used for signup and landing page)
func ListFamilies(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var families []models.Family
		if err := db.Find(&families).Error; err != nil {
			http.Error(w, `{"error":"failed to fetch families"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(families)
	}
}

func GetFamily(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		if familyID == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		var family models.Family
		if err := db.Where("id = ?", familyID).First(&family).Error; err != nil {
			http.Error(w, `{"error":"family not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(family)
	}
}

func GetFamilyStats(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		if familyID == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		stats := map[string]interface{}{}

		// Total members
		var totalMembers int64
		db.Model(&models.Member{}).Where("family_id = ?", familyID).Count(&totalMembers)
		stats["total_members"] = totalMembers

		// Active members
		var activeMembers int64
		db.Model(&models.Member{}).Where("family_id = ? AND status = ?", familyID, "active").Count(&activeMembers)
		stats["active_members"] = activeMembers

		// Pending invites
		var pendingMembers int64
		db.Model(&models.Member{}).Where("family_id = ? AND status = ?", familyID, "pending_invite").Count(&pendingMembers)
		stats["pending_invites"] = pendingMembers

		// Total meetings
		var totalMeetings int64
		db.Model(&models.Meeting{}).Where("family_id = ?", familyID).Count(&totalMeetings)
		stats["total_meetings"] = totalMeetings

		// Total contributions
		var totalContributions float64
		db.Model(&models.Contribution{}).Where("family_id = ? AND status = ?", familyID, "confirmed").Select("COALESCE(SUM(amount), 0)").Row().Scan(&totalContributions)
		stats["total_contributions"] = totalContributions

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(stats)
	}
}
