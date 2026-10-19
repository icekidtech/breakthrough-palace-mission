package routes

import (
	"net/http"

	"github.com/icekidtech/breakthrough-palace-mission/backend/config"
	"github.com/icekidtech/breakthrough-palace-mission/backend/handlers"
	"github.com/icekidtech/breakthrough-palace-mission/backend/handlers/admin"
	"github.com/icekidtech/breakthrough-palace-mission/backend/middleware"
	"gorm.io/gorm"
)

func RegisterRoutes(mux *http.ServeMux, db *gorm.DB, cfg *config.Config) {
	authMiddleware := middleware.AuthMiddleware(cfg)

	// Public routes (no auth required)
	mux.HandleFunc("GET /api/health", handlers.Health)
	mux.HandleFunc("GET /api/families/", handlers.ListFamilies(db))
	mux.HandleFunc("POST /api/auth/signup", handlers.Signup(db))
	mux.HandleFunc("POST /api/auth/login", handlers.Login(db, cfg))
	mux.HandleFunc("POST /api/auth/logout", handlers.Logout)

	// Auth me (protected)
	mux.Handle("GET /api/auth/me", authMiddleware(http.HandlerFunc(handlers.Me(db))))

	// Biodata routes (priority - onboarding gate)
	mux.Handle("GET /api/biodata", authMiddleware(http.HandlerFunc(handlers.GetBiodata(db))))
	mux.Handle("PUT /api/biodata", authMiddleware(http.HandlerFunc(handlers.UpdateBiodata(db))))
	mux.Handle("GET /api/biodata/completion-status", authMiddleware(http.HandlerFunc(handlers.CheckBiodataCompletion(db))))
	mux.Handle("GET /api/biodata/stats", authMiddleware(http.HandlerFunc(handlers.ListBiodataStats(db))))
	mux.Handle("GET /api/members/{id}/biodata", authMiddleware(http.HandlerFunc(handlers.GetMemberBiodata(db))))

	// Family routes (protected)
	mux.Handle("GET /api/family", authMiddleware(http.HandlerFunc(handlers.GetFamily(db))))
	mux.Handle("GET /api/family/stats", authMiddleware(http.HandlerFunc(handlers.GetFamilyStats(db))))

	// Member routes (protected)
	mux.Handle("GET /api/members", authMiddleware(http.HandlerFunc(handlers.ListMembers(db))))
	mux.Handle("GET /api/members/", authMiddleware(http.HandlerFunc(handlers.GetMember(db))))
	mux.Handle("POST /api/members", authMiddleware(http.HandlerFunc(handlers.AddMember(db))))
	mux.Handle("PUT /api/members/", authMiddleware(http.HandlerFunc(handlers.UpdateMember(db))))
	mux.Handle("DELETE /api/members/", authMiddleware(http.HandlerFunc(handlers.DeleteMember(db))))

	// Meeting routes (protected)
	mux.Handle("GET /api/meetings", authMiddleware(http.HandlerFunc(handlers.ListMeetings(db))))
	mux.Handle("GET /api/meetings/", authMiddleware(http.HandlerFunc(handlers.GetMeeting(db))))
	mux.Handle("POST /api/meetings", authMiddleware(http.HandlerFunc(handlers.CreateMeeting(db))))
	mux.Handle("PUT /api/meetings/", authMiddleware(http.HandlerFunc(handlers.UpdateMeeting(db))))
	mux.Handle("DELETE /api/meetings/", authMiddleware(http.HandlerFunc(handlers.DeleteMeeting(db))))

	// Attendance routes (protected)
	mux.Handle("GET /api/attendance", authMiddleware(http.HandlerFunc(handlers.ListAttendance(db))))
	mux.Handle("POST /api/attendance", authMiddleware(http.HandlerFunc(handlers.RecordAttendance(db))))
	mux.Handle("GET /api/attendance/stats", authMiddleware(http.HandlerFunc(handlers.GetAttendanceStats(db))))
	mux.Handle("PUT /api/attendance/", authMiddleware(http.HandlerFunc(handlers.UpdateAttendance(db))))

	// Minutes routes (protected)
	mux.Handle("GET /api/minutes", authMiddleware(http.HandlerFunc(handlers.ListMinutes(db))))
	mux.Handle("GET /api/minutes/", authMiddleware(http.HandlerFunc(handlers.GetMinutes(db))))
	mux.Handle("POST /api/minutes", authMiddleware(http.HandlerFunc(handlers.PublishMinutes(db))))
	mux.Handle("DELETE /api/minutes/", authMiddleware(http.HandlerFunc(handlers.DeleteMinutes(db))))

	// Contribution routes (protected)
	mux.Handle("GET /api/contributions", authMiddleware(http.HandlerFunc(handlers.ListContributions(db))))
	mux.Handle("POST /api/contributions", authMiddleware(http.HandlerFunc(handlers.RecordContribution(db))))
	mux.Handle("GET /api/contributions/stats", authMiddleware(http.HandlerFunc(handlers.GetContributionStats(db))))
	mux.Handle("PUT /api/contributions/", authMiddleware(http.HandlerFunc(handlers.UpdateContribution(db))))

	// Admin routes (protected)
	mux.Handle("GET /api/admin/families", authMiddleware(http.HandlerFunc(admin.ListFamiliesAdmin(db))))
	mux.Handle("POST /api/admin/families", authMiddleware(http.HandlerFunc(admin.CreateFamily(db))))
	mux.Handle("PUT /api/admin/families/{id}", authMiddleware(http.HandlerFunc(admin.UpdateFamily(db))))
	mux.Handle("DELETE /api/admin/families/{id}", authMiddleware(http.HandlerFunc(admin.DeleteFamily(db))))
}
