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

type RecordContributionRequest struct {
	MemberID string  `json:"member_id"`
	Type     string  `json:"type"` // family_levy, monthly_welfare, building_fund
	Amount   float64 `json:"amount"`
	Date     time.Time `json:"date"`
	Notes    string  `json:"notes"`
}

type UpdateContributionRequest struct {
	Status      string `json:"status"` // confirmed, pending, rejected
	ConfirmedBy string `json:"confirmed_by"`
	Notes       string `json:"notes"`
}

type ContributionStats struct {
	TotalThisMonth      float64 `json:"total_this_month"`
	TotalThisYear       float64 `json:"total_this_year"`
	AnnualTarget        float64 `json:"annual_target"`
	ProgressPercentage  float64 `json:"progress_percentage"`
	ConfirmedCount      int64   `json:"confirmed_count"`
	PendingCount        int64   `json:"pending_count"`
	ContributorsCount   int64   `json:"contributors_count"`
	RejectedCount       int64   `json:"rejected_count"`
}

func ListContributions(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		memberID := middleware.ExtractMemberID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "contributions:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var contributions []models.Contribution
		query := db.Where("family_id = ?", familyID)

		// Members can only view their own contributions
		if role == "member" {
			query = query.Where("member_id = ?", memberID)
		}

		if err := query.Find(&contributions).Error; err != nil {
			http.Error(w, `{"error":"failed to fetch contributions"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contributions)
	}
}

func RecordContribution(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := middleware.ExtractFamilyID(r)
		memberID := middleware.ExtractMemberID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "contributions:record") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var req RecordContributionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if req.MemberID == "" || req.Type == "" || req.Amount == 0 {
			http.Error(w, `{"error":"member_id, type, and amount are required"}`, http.StatusBadRequest)
			return
		}

		// Members can only record for themselves
		if role == "member" && req.MemberID != memberID {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		// If recorded by member, status is pending confirmation
		status := "pending"
		if role != "member" {
			status = "confirmed"
		}

		contribution := &models.Contribution{
			FamilyID: utils.ParseUUID(familyID),
			MemberID: utils.ParseUUID(req.MemberID),
			Type:     req.Type,
			Amount:   req.Amount,
			Date:     req.Date,
			Status:   status,
			Notes:    req.Notes,
			Currency: "NGN",
		}

		// If admin records it, they confirm it
		if role != "member" {
			contribution.ConfirmedBy = utils.ParseUUID(memberID)
		}

		if err := db.Create(contribution).Error; err != nil {
			http.Error(w, `{"error":"failed to record contribution"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(contribution)
	}
}

func GetContributionStats(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "contributions:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		stats := ContributionStats{
			AnnualTarget: 500000, // Default annual target, should come from family settings
		}

		// Get current month and year
		now := time.Now()
		monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)

		// Total this month (confirmed only)
		db.Model(&models.Contribution{}).
			Where("family_id = ? AND status = ? AND date >= ?", familyID, "confirmed", monthStart).
			Select("COALESCE(SUM(amount), 0)").
			Row().Scan(&stats.TotalThisMonth)

		// Total this year (confirmed only)
		db.Model(&models.Contribution{}).
			Where("family_id = ? AND status = ? AND date >= ?", familyID, "confirmed", yearStart).
			Select("COALESCE(SUM(amount), 0)").
			Row().Scan(&stats.TotalThisYear)

		// Calculate progress percentage
		if stats.AnnualTarget > 0 {
			stats.ProgressPercentage = (stats.TotalThisYear / stats.AnnualTarget) * 100
		}

		// Count by status
		db.Model(&models.Contribution{}).
			Where("family_id = ? AND status = ?", familyID, "confirmed").
			Count(&stats.ConfirmedCount)

		db.Model(&models.Contribution{}).
			Where("family_id = ? AND status = ?", familyID, "pending").
			Count(&stats.PendingCount)

		db.Model(&models.Contribution{}).
			Where("family_id = ? AND status = ?", familyID, "rejected").
			Count(&stats.RejectedCount)

		// Count unique contributors (this year)
		db.Model(&models.Contribution{}).
			Where("family_id = ? AND status = ? AND date >= ?", familyID, "confirmed", yearStart).
			Select("COUNT(DISTINCT member_id)").
			Row().Scan(&stats.ContributorsCount)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(stats)
	}
}

func UpdateContribution(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := middleware.ExtractFamilyID(r)
		memberID := middleware.ExtractMemberID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "contributions:confirm") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		// Extract ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		contributionID := parts[len(parts)-1]

		var req UpdateContributionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if req.Status == "" {
			http.Error(w, `{"error":"status is required"}`, http.StatusBadRequest)
			return
		}

		// Set confirmed_by if confirming
		updates := map[string]interface{}{
			"status": req.Status,
			"notes":  req.Notes,
		}

		if req.Status == "confirmed" {
			updates["confirmed_by"] = memberID
		}

		if err := db.Model(&models.Contribution{}).
			Where("id = ? AND family_id = ?", contributionID, familyID).
			Updates(updates).Error; err != nil {
			http.Error(w, `{"error":"failed to update contribution"}`, http.StatusInternalServerError)
			return
		}

		var contribution models.Contribution
		db.First(&contribution, "id = ?", contributionID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contribution)
	}
}
