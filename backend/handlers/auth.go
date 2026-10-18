package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/icekidtech/breakthrough-palace-mission/backend/auth"
	"github.com/icekidtech/breakthrough-palace-mission/backend/config"
	"github.com/icekidtech/breakthrough-palace-mission/backend/middleware"
	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
	"github.com/icekidtech/breakthrough-palace-mission/backend/utils"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Member  struct {
		ID       interface{} `json:"id"`
		Name     string      `json:"name"`
		Email    string      `json:"email"`
		Role     string      `json:"role"`
		FamilyID interface{} `json:"family_id"`
	} `json:"member"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

func Login(db *gorm.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		// Find member by email
		var member models.Member
		if err := db.Where("email = ?", req.Email).First(&member).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
				return
			}
			http.Error(w, `{"error":"database error"}`, http.StatusInternalServerError)
			return
		}

		// Verify password
		if !utils.VerifyPassword(member.Password, req.Password) {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		// Generate JWT (convert UUID to string)
		token, err := auth.GenerateJWT(member.ID.String(), member.FamilyID.String(), member.Name, member.Role, cfg.JWTSecret)
		if err != nil {
			http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
			return
		}

		// Set cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			HttpOnly: true,
			Secure:   cfg.CookieSecure,
			SameSite: http.SameSiteLaxMode,
			Domain:   cfg.CookieDomain,
			Path:     "/",
			MaxAge:   86400 * 7, // 7 days
		})

		// Return response
		resp := LoginResponse{
			Message: "login successful",
		}
		resp.Member.ID = member.ID
		resp.Member.Name = member.Name
		resp.Member.Email = member.Email
		resp.Member.Role = member.Role
		resp.Member.FamilyID = member.FamilyID

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1,
	})

	resp := LogoutResponse{
		Message: "logout successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func Me(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memberID := middleware.ExtractMemberID(r)
		if memberID == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		var member models.Member
		if err := db.Preload("Family").First(&member, "id = ?", memberID).Error; err != nil {
			http.Error(w, `{"error":"member not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(member)
	}
}
