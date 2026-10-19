package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/icekidtech/breakthrough-palace-mission/backend/config"
	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
	"github.com/icekidtech/breakthrough-palace-mission/backend/services"
	"github.com/icekidtech/breakthrough-palace-mission/backend/utils"
	"gorm.io/gorm"
)

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type PasswordResetResponse struct {
	Message string `json:"message"`
}

// ForgotPassword generates a password reset token and sends email
func ForgotPassword(db *gorm.DB, cfg *config.Config, emailService *services.EmailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req ForgotPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if req.Email == "" {
			http.Error(w, `{"error":"email required"}`, http.StatusBadRequest)
			return
		}

		// Find member by email
		var member models.Member
		if err := db.Where("email = ?", req.Email).First(&member).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Don't reveal if email exists (security best practice)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"message": "if email exists, reset link has been sent"})
				return
			}
			http.Error(w, `{"error":"failed to process request"}`, http.StatusInternalServerError)
			return
		}

		// Generate reset token
		resetToken := utils.GenerateToken()

		// Create verification token record
		verificationToken := &models.VerificationToken{
			MemberID:  member.ID,
			Token:     resetToken,
			Type:      "password_reset",
			ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
		}

		if err := db.Create(verificationToken).Error; err != nil {
			http.Error(w, `{"error":"failed to generate reset link"}`, http.StatusInternalServerError)
			return
		}

		// Send email with reset link
		resetLink := fmt.Sprintf("%s/auth/reset-password?token=%s", cfg.FrontendURL, resetToken)
		if err := emailService.SendPasswordResetEmail(member.Email, member.Name, resetLink, 1); err != nil {
			// Log error but still return success (token was created)
			fmt.Printf("Failed to send password reset email: %v\n", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PasswordResetResponse{Message: "if email exists, reset link has been sent"})
	}
}

// ResetPassword resets password using a valid token
func ResetPassword(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req ResetPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		if req.Token == "" || req.Password == "" {
			http.Error(w, `{"error":"token and password required"}`, http.StatusBadRequest)
			return
		}

		if len(req.Password) < 6 {
			http.Error(w, `{"error":"password must be at least 6 characters"}`, http.StatusBadRequest)
			return
		}

		// Find verification token
		var token models.VerificationToken
		if err := db.Where("token = ? AND type = ? AND expires_at > ?", req.Token, "password_reset", time.Now()).First(&token).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}
			http.Error(w, `{"error":"failed to verify token"}`, http.StatusInternalServerError)
			return
		}

		// Hash new password
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			http.Error(w, `{"error":"failed to process password"}`, http.StatusInternalServerError)
			return
		}

		// Update member password
		if err := db.Model(&models.Member{}).Where("id = ?", token.MemberID).Update("password", hashedPassword).Error; err != nil {
			http.Error(w, `{"error":"failed to reset password"}`, http.StatusInternalServerError)
			return
		}

		// Delete used token
		db.Delete(&token)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PasswordResetResponse{Message: "password reset successfully"})
	}
}
