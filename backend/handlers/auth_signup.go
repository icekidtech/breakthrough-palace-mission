package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/icekidtech/breakthrough-palace-mission/backend/config"
	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
	"github.com/icekidtech/breakthrough-palace-mission/backend/services"
	"github.com/icekidtech/breakthrough-palace-mission/backend/utils"
	"gorm.io/gorm"
)

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	FamilyID string `json:"family_id"`
}

type SignupResponse struct {
	Message string `json:"message"`
	Member  struct {
		ID       interface{} `json:"id"`
		Name     string      `json:"name"`
		Email    string      `json:"email"`
		Phone    string      `json:"phone"`
		FamilyID interface{} `json:"family_id"`
		Status   string      `json:"status"`
	} `json:"member"`
}

func Signup(db *gorm.DB, cfg *config.Config, emailService *services.EmailService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		// Validate required fields
		if req.Name == "" || req.Email == "" || req.Phone == "" || req.Password == "" || req.FamilyID == "" {
			http.Error(w, `{"error":"all fields are required"}`, http.StatusBadRequest)
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

		// Create member as active (no pending approval)
		member := &models.Member{
			Name:     req.Name,
			Email:    req.Email,
			Phone:    req.Phone,
			Password: hashedPassword,
			FamilyID: utils.ParseUUID(req.FamilyID),
			Role:     "member", // Default role
			Status:   "active", // Active immediately
		}

		if err := db.Create(member).Error; err != nil {
			http.Error(w, `{"error":"failed to create account"}`, http.StatusInternalServerError)
			return
		}

		// Generate verification token
		verificationToken := utils.GenerateToken()

		// Create verification token record
		emailVerificationToken := &models.VerificationToken{
			MemberID:  member.ID,
			Token:     verificationToken,
			Type:      "email_verification",
			ExpiresAt: utils.TimeNowAddHours(24), // 24 hour expiry
		}

		if err := db.Create(emailVerificationToken).Error; err != nil {
			fmt.Printf("Failed to create verification token: %v\n", err)
		}

		// Send verification email
		verificationLink := fmt.Sprintf("%s/auth/verify-email?token=%s", cfg.FrontendURL, verificationToken)
		if err := emailService.SendVerificationEmail(member.Email, member.Name, verificationLink); err != nil {
			// Log error but continue (account was created successfully)
			fmt.Printf("Failed to send verification email: %v\n", err)
		}

		// Send welcome email
		dashboardURL := fmt.Sprintf("%s/dashboard", cfg.FrontendURL)
		if err := emailService.SendWelcomeEmail(member.Email, member.Name, dashboardURL); err != nil {
			// Log error but continue
			fmt.Printf("Failed to send welcome email: %v\n", err)
		}

		resp := SignupResponse{
			Message: "signup successful",
		}
		resp.Member.ID = member.ID
		resp.Member.Name = member.Name
		resp.Member.Email = member.Email
		resp.Member.Phone = member.Phone
		resp.Member.FamilyID = member.FamilyID
		resp.Member.Status = member.Status

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
