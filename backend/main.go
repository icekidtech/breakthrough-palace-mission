package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/icekidtech/breakthrough-palace-mission/backend/config"
	"github.com/icekidtech/breakthrough-palace-mission/backend/database"
	"github.com/icekidtech/breakthrough-palace-mission/backend/handlers"
	"github.com/icekidtech/breakthrough-palace-mission/backend/middleware"
)

func main() {
	// Load .env file
	godotenv.Load()

	// Load config
	cfg := config.Load()
	log.Printf("Starting server in %s mode\n", cfg.Env)

	// Initialize database
	db := database.Init(cfg)
	log.Println("Running migrations...")
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Seed database (only if empty)
	if err := database.Seed(db); err != nil {
		log.Fatalf("Seed failed: %v", err)
	}

	// Setup routes
	mux := http.NewServeMux()

	// Auth routes (no auth required)
	mux.HandleFunc("POST /api/auth/login", handlers.Login(db, cfg))
	mux.HandleFunc("POST /api/auth/logout", handlers.Logout)

	// Protected routes (auth required)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /api/auth/me", handlers.Me(db))
	protectedMux.HandleFunc("GET /api/family", handlers.GetFamily(db))
	protectedMux.HandleFunc("GET /api/family/stats", handlers.GetFamilyStats(db))

	// Apply auth middleware to protected routes
	mux.Handle("/api/", middleware.AuthMiddleware(cfg)(protectedMux))

	// Wrap all routes with CORS
	handler := middleware.CORSMiddleware(cfg)(mux)

	// Start server
	port := ":" + cfg.Port
	fmt.Printf("Server running on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
