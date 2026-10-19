package routes

import (
	"net/http"

	"github.com/icekidtech/breakthrough-palace-mission/backend/config"
	"github.com/icekidtech/breakthrough-palace-mission/backend/handlers"
	"github.com/icekidtech/breakthrough-palace-mission/backend/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(mux *http.ServeMux, db *gorm.DB, cfg *config.Config) {
	// Public routes (no auth required)
	mux.HandleFunc("GET /api/health", handlers.Health)
	mux.HandleFunc("POST /api/auth/signup", handlers.Signup(db))
	mux.HandleFunc("POST /api/auth/login", handlers.Login(db, cfg))
	mux.HandleFunc("POST /api/auth/logout", handlers.Logout)

	// Protected routes (auth required)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /api/auth/me", handlers.Me(db))

	// Biodata routes (priority - onboarding gate)
	protectedMux.HandleFunc("GET /api/biodata", handlers.GetBiodata(db))
	protectedMux.HandleFunc("PUT /api/biodata", handlers.UpdateBiodata(db))
	protectedMux.HandleFunc("GET /api/biodata/completion-status", handlers.CheckBiodataCompletion(db))
	protectedMux.HandleFunc("GET /api/biodata/stats", handlers.ListBiodataStats(db))
	protectedMux.HandleFunc("GET /api/members/{id}/biodata", handlers.GetMemberBiodata(db))

	// Family routes
	protectedMux.HandleFunc("GET /api/family", handlers.GetFamily(db))
	protectedMux.HandleFunc("GET /api/family/stats", handlers.GetFamilyStats(db))

	// Member routes
	protectedMux.HandleFunc("GET /api/members", handlers.ListMembers(db))
	protectedMux.HandleFunc("GET /api/members/", handlers.GetMember(db))
	protectedMux.HandleFunc("POST /api/members", handlers.AddMember(db))
	protectedMux.HandleFunc("PUT /api/members/", handlers.UpdateMember(db))
	protectedMux.HandleFunc("DELETE /api/members/", handlers.DeleteMember(db))

	// Meeting routes
	protectedMux.HandleFunc("GET /api/meetings", handlers.ListMeetings(db))
	protectedMux.HandleFunc("GET /api/meetings/", handlers.GetMeeting(db))
	protectedMux.HandleFunc("POST /api/meetings", handlers.CreateMeeting(db))
	protectedMux.HandleFunc("PUT /api/meetings/", handlers.UpdateMeeting(db))
	protectedMux.HandleFunc("DELETE /api/meetings/", handlers.DeleteMeeting(db))

	// Attendance routes
	protectedMux.HandleFunc("GET /api/attendance", handlers.ListAttendance(db))
	protectedMux.HandleFunc("POST /api/attendance", handlers.RecordAttendance(db))
	protectedMux.HandleFunc("GET /api/attendance/stats", handlers.GetAttendanceStats(db))
	protectedMux.HandleFunc("PUT /api/attendance/", handlers.UpdateAttendance(db))

	// Minutes routes
	protectedMux.HandleFunc("GET /api/minutes", handlers.ListMinutes(db))
	protectedMux.HandleFunc("GET /api/minutes/", handlers.GetMinutes(db))
	protectedMux.HandleFunc("POST /api/minutes", handlers.PublishMinutes(db))
	protectedMux.HandleFunc("DELETE /api/minutes/", handlers.DeleteMinutes(db))

	// Contribution routes
	protectedMux.HandleFunc("GET /api/contributions", handlers.ListContributions(db))
	protectedMux.HandleFunc("POST /api/contributions", handlers.RecordContribution(db))
	protectedMux.HandleFunc("GET /api/contributions/stats", handlers.GetContributionStats(db))
	protectedMux.HandleFunc("PUT /api/contributions/", handlers.UpdateContribution(db))

	// Admin routes
	protectedMux.HandleFunc("POST /api/admin/families", handlers.CreateFamily(db))
	protectedMux.HandleFunc("GET /api/admin/families", handlers.ListFamilies(db))
	protectedMux.HandleFunc("PUT /api/family", handlers.UpdateFamily(db))

	// Apply auth middleware to protected routes
	mux.Handle("/api/", middleware.AuthMiddleware(cfg)(protectedMux))
}
