package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/icekidtech/breakthrough-palace-mission/backend/config"
	"github.com/icekidtech/breakthrough-palace-mission/backend/database"
	"github.com/icekidtech/breakthrough-palace-mission/backend/middleware"
	"github.com/icekidtech/breakthrough-palace-mission/backend/routes"
	"github.com/joho/godotenv"
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

	// Setup server
	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, db, cfg)
	handler := middleware.CORSMiddleware(cfg)(mux)

	// Start server
	port := ":" + cfg.Port
	fmt.Printf("Server running on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
