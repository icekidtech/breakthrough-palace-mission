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
	mux.HandleFunc("POST /api/auth/login", handlers.Login(db, cfg))
	mux.HandleFunc("POST /api/auth/logout", handlers.Logout)

	// Protected routes (auth required)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /api/auth/me", handlers.Me(db))
	protectedMux.HandleFunc("GET /api/family", handlers.GetFamily(db))
	protectedMux.HandleFunc("GET /api/family/stats", handlers.GetFamilyStats(db))

	// Apply auth middleware to protected routes
	mux.Handle("/api/", middleware.AuthMiddleware(cfg)(protectedMux))
}
