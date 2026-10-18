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

type RecordAttendanceRequest struct {
	MeetingID string `json:"meeting_id"`
	MemberID  string `json:"member_id"`
	Status    string `json:"status"` // present, absent, late, excused
}

type UpdateAttendanceRequest struct {
	Status string `json:"status"`
}

type AttendanceStats struct {
	TotalRecords      int64   `json:"total_records"`
	PresentCount      int64   `json:"present_count"`
	AbsentCount       int64   `json:"absent_count"`
	LateCount         int64   `json:"late_count"`
	AttendanceRate    float64 `json:"attendance_rate"`
	AveragePerMeeting int64   `json:"average_per_meeting"`
}

func ListAttendance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		memberID := middleware.ExtractMemberID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "attendance:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var attendance []models.Attendance
		query := db.Where("family_id = ?", familyID)

		// Members can only view their own attendance
		if role == "member" {
			query = query.Where("member_id = ?", memberID)
		}

		if err := query.Find(&attendance).Error; err != nil {
			http.Error(w, `{"error":"failed to fetch attendance"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(attendance)
	}
}

func RecordAttendance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := middleware.ExtractFamilyID(r)
		memberID := middleware.ExtractMemberID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "attendance:record") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var req RecordAttendanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if req.MeetingID == "" || req.MemberID == "" || req.Status == "" {
			http.Error(w, `{"error":"meeting_id, member_id, and status are required"}`, http.StatusBadRequest)
			return
		}

		// Members can only record for themselves
		if role == "member" && req.MemberID != memberID {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		attendance := &models.Attendance{
			FamilyID:   utils.ParseUUID(familyID),
			MeetingID:  utils.ParseUUID(req.MeetingID),
			MemberID:   utils.ParseUUID(req.MemberID),
			Status:     req.Status,
			RecordedBy: utils.ParseUUID(memberID),
			RecordedAt: time.Now(),
		}

		if err := db.Create(attendance).Error; err != nil {
			http.Error(w, `{"error":"failed to record attendance"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(attendance)
	}
}

func GetAttendanceStats(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "attendance:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		stats := AttendanceStats{}

		// Total records
		db.Model(&models.Attendance{}).Where("family_id = ?", familyID).Count(&stats.TotalRecords)

		// Count by status
		db.Model(&models.Attendance{}).Where("family_id = ? AND status = ?", familyID, "present").Count(&stats.PresentCount)
		db.Model(&models.Attendance{}).Where("family_id = ? AND status = ?", familyID, "absent").Count(&stats.AbsentCount)
		db.Model(&models.Attendance{}).Where("family_id = ? AND status = ?", familyID, "late").Count(&stats.LateCount)

		// Calculate attendance rate
		if stats.TotalRecords > 0 {
			stats.AttendanceRate = float64(stats.PresentCount) / float64(stats.TotalRecords) * 100
		}

		// Average per meeting
		var meetingCount int64
		db.Model(&models.Meeting{}).Where("family_id = ?", familyID).Count(&meetingCount)
		if meetingCount > 0 {
			stats.AveragePerMeeting = stats.TotalRecords / meetingCount
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(stats)
	}
}

func UpdateAttendance(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "attendance:edit") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		// Extract ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		attendanceID := parts[len(parts)-1]

		var req UpdateAttendanceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if err := db.Model(&models.Attendance{}).Where("id = ? AND family_id = ?", attendanceID, familyID).Update("status", req.Status).Error; err != nil {
			http.Error(w, `{"error":"failed to update attendance"}`, http.StatusInternalServerError)
			return
		}

		var attendance models.Attendance
		db.First(&attendance, "id = ?", attendanceID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(attendance)
	}
}
