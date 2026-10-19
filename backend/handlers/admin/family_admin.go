package admin

import (
	"encoding/json"
	"net/http"

	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
	"github.com/icekidtech/breakthrough-palace-mission/backend/rbac"
	"gorm.io/gorm"
)

type CreateFamilyRequest struct {
	Name    string `json:"name"`
	Founder string `json:"founder"`
	Country string `json:"country"`
	Branch  string `json:"branch"`
}

func CreateFamily(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		// Check permission
		role := r.Header.Get("X-Role")
		if !rbac.IsSuperAdmin(role) {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		var req CreateFamilyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		family := &models.Family{
			Name:    req.Name,
			Founder: req.Founder,
			Country: req.Country,
			Branch:  req.Branch,
		}

		if err := db.Create(family).Error; err != nil {
			http.Error(w, `{"error":"failed to create family"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(family)
	}
}

// ListFamiliesAdmin returns all families for admin panel (super admin only)
func ListFamiliesAdmin(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check permission
		role := r.Header.Get("X-Role")
		if !rbac.IsSuperAdmin(role) {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

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

func UpdateFamily(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		familyID := r.PathValue("id")
		if familyID == "" {
			http.Error(w, `{"error":"family id required"}`, http.StatusBadRequest)
			return
		}

		role := r.Header.Get("X-Role")

		// Super admin or family leader/secretary can edit
		if !rbac.IsSuperAdmin(role) {
			if role != "family_leader" && role != "family_assistant" && role != "secretary" {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}
		}

		var req CreateFamilyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if err := db.Model(&models.Family{}).Where("id = ?", familyID).Updates(req).Error; err != nil {
			http.Error(w, `{"error":"failed to update family"}`, http.StatusInternalServerError)
			return
		}

		var family models.Family
		db.First(&family, "id = ?", familyID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(family)
	}
}

// DeleteFamily deletes a family (super admin only)
func DeleteFamily(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		// Check permission
		role := r.Header.Get("X-Role")
		if !rbac.IsSuperAdmin(role) {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}

		familyID := r.PathValue("id")
		if familyID == "" {
			http.Error(w, `{"error":"family id required"}`, http.StatusBadRequest)
			return
		}

		// Check family exists
		var family models.Family
		if err := db.First(&family, "id = ?", familyID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, `{"error":"family not found"}`, http.StatusNotFound)
				return
			}
			http.Error(w, `{"error":"failed to fetch family"}`, http.StatusInternalServerError)
			return
		}

		// Delete family
		if err := db.Delete(&family).Error; err != nil {
			http.Error(w, `{"error":"failed to delete family"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "family deleted successfully"})
	}
}
