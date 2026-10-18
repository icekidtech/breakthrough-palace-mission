package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/icekidtech/breakthrough-palace-mission/backend/middleware"
	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
	"github.com/icekidtech/breakthrough-palace-mission/backend/rbac"
	"github.com/icekidtech/breakthrough-palace-mission/backend/utils"
	"gorm.io/gorm"
)

type MemberUpdateRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func ListMembers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "members:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var members []models.Member
		if err := db.Where("family_id = ?", familyID).Find(&members).Error; err != nil {
			http.Error(w, `{"error":"failed to fetch members"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(members)
	}
}

func GetMember(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		familyID := middleware.ExtractFamilyID(r)
		memberID := middleware.ExtractMemberID(r)
		role := r.Header.Get("X-Role")

		// Extract ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		targetMemberID := parts[len(parts)-1]

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "members:view") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		// Members can only view their own profile
		if targetMemberID != memberID && role == "member" {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var member models.Member
		if err := db.Where("id = ? AND family_id = ?", targetMemberID, familyID).First(&member).Error; err != nil {
			http.Error(w, `{"error":"member not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(member)
	}
}

func AddMember(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "members:add") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.Email == "" || req.Phone == "" || req.Password == "" {
			http.Error(w, `{"error":"name, email, phone, and password are required"}`, http.StatusBadRequest)
			return
		}

		// Check if email already exists
		var existingMember models.Member
		if err := db.Where("email = ?", req.Email).First(&existingMember).Error; err == nil {
			http.Error(w, `{"error":"email already exists"}`, http.StatusConflict)
			return
		}

		// Hash password
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			http.Error(w, `{"error":"failed to process password"}`, http.StatusInternalServerError)
			return
		}

		member := &models.Member{
			Name:     req.Name,
			Email:    req.Email,
			Phone:    req.Phone,
			Password: hashedPassword,
			FamilyID: utils.ParseUUID(familyID),
			Role:     "member",
			Status:   "active", // Admin can add members directly as active
		}

		if err := db.Create(member).Error; err != nil {
			http.Error(w, `{"error":"failed to create member"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(member)
	}
}

func UpdateMember(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := middleware.ExtractFamilyID(r)
		memberID := middleware.ExtractMemberID(r)
		role := r.Header.Get("X-Role")

		// Extract ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		targetMemberID := parts[len(parts)-1]

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "members:edit") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		// Members can only edit their own profile
		if targetMemberID != memberID && role == "member" {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var req MemberUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if err := db.Model(&models.Member{}).Where("id = ? AND family_id = ?", targetMemberID, familyID).Updates(req).Error; err != nil {
			http.Error(w, `{"error":"failed to update member"}`, http.StatusInternalServerError)
			return
		}

		var member models.Member
		db.First(&member, "id = ?", targetMemberID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(member)
	}
}

func DeleteMember(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := middleware.ExtractFamilyID(r)
		role := r.Header.Get("X-Role")

		// Check permission
		if !rbac.IsSuperAdmin(role) && !hasPermission(db, role, "members:remove") {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		// Extract ID from URL path
		parts := strings.Split(r.URL.Path, "/")
		targetMemberID := parts[len(parts)-1]

		if err := db.Where("id = ? AND family_id = ?", targetMemberID, familyID).Delete(&models.Member{}).Error; err != nil {
			http.Error(w, `{"error":"failed to delete member"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "member deleted"})
	}
}

// Helper function to check if role has permission
func hasPermission(db *gorm.DB, role string, permissionName string) bool {
	var count int64
	db.Table("role_permissions").
		Joins("JOIN roles ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("roles.name = ? AND permissions.name = ?", role, permissionName).
		Count(&count)
	return count > 0
}
