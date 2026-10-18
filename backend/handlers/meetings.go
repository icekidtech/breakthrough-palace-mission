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

type CreateMeetingRequest struct {
	Title     string    `json:"title"`
	Type      string    `json:"type"` // general_meeting, committee_meeting, prayer_gathering
	Date      time.Time `json:"date"`
	StartTime string    `json:"start_time"` // HH:MM
	EndTime   string    `json:"end_time"`   // HH:MM
	Location  string    `json:"location"`
}

type UpdateMeetingRequest struct {
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Date      time.Time `json:"date"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Location  string    `json:"location"`
}

func ListMeetings(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "meetings:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var meetings []models.Meeting
		if err := db.Where("family_id = ?", familyID).Find(&meetings).Error; err != nil {
			http.Error(w, `{"error":"failed to fetch meetings"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(meetings)
	}
}

func GetMeeting(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "meetings:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		// Extract ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		meetingID := parts[len(parts)-1]

		var meeting models.Meeting
		if err := db.Where("id = ? AND family_id = ?", meetingID, familyID).First(&meeting).Error; err != nil {
			http.Error(w, `{"error":"meeting not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(meeting)
	}
}

func CreateMeeting(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := middleware.ExtractFamilyID(r)
		memberID := middleware.ExtractMemberID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "meetings:schedule") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var req CreateMeetingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if req.Title == "" || req.Date.IsZero() || req.Location == "" {
			http.Error(w, `{"error":"title, date, and location are required"}`, http.StatusBadRequest)
			return
		}

		meeting := &models.Meeting{
			Title:     req.Title,
			Type:      req.Type,
			Date:      req.Date,
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			Location:  req.Location,
			FamilyID:  utils.ParseUUID(familyID),
			CreatedBy: utils.ParseUUID(memberID),
		}

		if err := db.Create(meeting).Error; err != nil {
			http.Error(w, `{"error":"failed to create meeting"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(meeting)
	}
}

func UpdateMeeting(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "meetings:edit") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		// Extract ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		meetingID := parts[len(parts)-1]

		var req UpdateMeetingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if err := db.Model(&models.Meeting{}).Where("id = ? AND family_id = ?", meetingID, familyID).Updates(req).Error; err != nil {
			http.Error(w, `{"error":"failed to update meeting"}`, http.StatusInternalServerError)
			return
		}

		var meeting models.Meeting
		db.First(&meeting, "id = ?", meetingID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(meeting)
	}
}

func DeleteMeeting(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "meetings:delete") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		// Extract ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		meetingID := parts[len(parts)-1]

		if err := db.Where("id = ? AND family_id = ?", meetingID, familyID).Delete(&models.Meeting{}).Error; err != nil {
			http.Error(w, `{"error":"failed to delete meeting"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "meeting deleted"})
	}
}
