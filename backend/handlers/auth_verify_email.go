package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
	"gorm.io/gorm"
)

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type VerifyEmailResponse struct {
	Message string `json:"message"`
}

func VerifyEmail(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req VerifyEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if req.Token == "" {
			http.Error(w, `{"error":"token required"}`, http.StatusBadRequest)
			return
		}

		// Find verification token
		var token models.VerificationToken
		if err := db.Where("token = ? AND type = ? AND expires_at > ?", req.Token, "email_verification", time.Now()).First(&token).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}
			http.Error(w, `{"error":"failed to verify token"}`, http.StatusInternalServerError)
			return
		}

		// Update member email verified status
		if err := db.Model(&models.Member{}).Where("id = ?", token.MemberID).Update("email_verified", true).Error; err != nil {
			http.Error(w, `{"error":"failed to verify email"}`, http.StatusInternalServerError)
			return
		}

		// Delete used token
		db.Delete(&token)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VerifyEmailResponse{Message: "email verified successfully"})
	}
}
