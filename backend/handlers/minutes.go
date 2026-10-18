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

type CreateMinutesRequest struct {
	MeetingID string `json:"meeting_id"`
	Title     string `json:"title"`
	Content   string `json:"content"` // Markdown or HTML
	PDFUrl    string `json:"pdf_url"`
}

type UpdateMinutesRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	PDFUrl  string `json:"pdf_url"`
}

func ListMinutes(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "minutes:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var minutes []models.Minute
		if err := db.Where("family_id = ?", familyID).Find(&minutes).Error; err != nil {
			http.Error(w, `{"error":"failed to fetch minutes"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(minutes)
	}
}

func GetMinutes(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "minutes:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		// Extract ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		minutesID := parts[len(parts)-1]

		var minute models.Minute
		if err := db.Where("id = ? AND family_id = ?", minutesID, familyID).First(&minute).Error; err != nil {
			http.Error(w, `{"error":"minutes not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(minute)
	}
}

func PublishMinutes(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := middleware.ExtractFamilyID(r)
		memberID := middleware.ExtractMemberID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "minutes:publish") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var req CreateMinutesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if req.MeetingID == "" || req.Title == "" || req.Content == "" {
			http.Error(w, `{"error":"meeting_id, title, and content are required"}`, http.StatusBadRequest)
			return
		}

		minute := &models.Minute{
			FamilyID:    utils.ParseUUID(familyID),
			MeetingID:   utils.ParseUUID(req.MeetingID),
			Title:       req.Title,
			Content:     req.Content,
			PDFUrl:      req.PDFUrl,
			PublishedBy: utils.ParseUUID(memberID),
			PublishedAt: time.Now(),
		}

		if err := db.Create(minute).Error; err != nil {
			http.Error(w, `{"error":"failed to publish minutes"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(minute)
	}
}

func DeleteMinutes(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "minutes:delete") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		// Extract ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		minutesID := parts[len(parts)-1]

		if err := db.Where("id = ? AND family_id = ?", minutesID, familyID).Delete(&models.Minute{}).Error; err != nil {
			http.Error(w, `{"error":"failed to delete minutes"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "minutes deleted"})
	}
}
