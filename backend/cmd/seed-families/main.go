package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/icekidtech/breakthrough-palace-mission/backend/config"
	"github.com/icekidtech/breakthrough-palace-mission/backend/database"
	"github.com/icekidtech/breakthrough-palace-mission/backend/models"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect to database
	db := database.Init(cfg)

	// Run migrations first
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Default families with details
	defaultFamilies := []models.Family{
		{ID: uuid.New(), Name: "Abundance", Founder: "TBD", Country: "Nigeria", Branch: "Main"},
		{ID: uuid.New(), Name: "Achievers", Founder: "TBD", Country: "Nigeria", Branch: "Main"},
		{ID: uuid.New(), Name: "Breakthrough", Founder: "TBD", Country: "Nigeria", Branch: "Main"},
		{ID: uuid.New(), Name: "Dominion", Founder: "TBD", Country: "Nigeria", Branch: "Main"},
		{ID: uuid.New(), Name: "Eagles", Founder: "TBD", Country: "Nigeria", Branch: "Main"},
		{ID: uuid.New(), Name: "Glorious", Founder: "TBD", Country: "Nigeria", Branch: "Main"},
		{ID: uuid.New(), Name: "Omega", Founder: "TBD", Country: "Nigeria", Branch: "Main"},
		{ID: uuid.New(), Name: "Peculiar", Founder: "TBD", Country: "Nigeria", Branch: "Main"},
		{ID: uuid.New(), Name: "Restoration", Founder: "TBD", Country: "Nigeria", Branch: "Main"},
	}

	fmt.Println("=== Breakthrough Palace Mission - Family Seeder ===\n")
	fmt.Println("Seeding families...\n")

	var familiesCreated int
	var familiesSkipped int

	for _, family := range defaultFamilies {
		// Check if family already exists
		var existing models.Family
		if err := db.Where("name = ?", family.Name).First(&existing).Error; err == nil {
			fmt.Printf("⊘ Skipped: %s (already exists)\n", family.Name)
			familiesSkipped++
			continue
		}

		if err := db.Create(&family).Error; err != nil {
			fmt.Printf("❌ Failed to create family '%s': %v\n", family.Name, err)
			continue
		}

		fmt.Printf("✓ Created: %s\n", family.Name)
		familiesCreated++
	}

	// Display summary
	var allFamilies []models.Family
	if err := db.Find(&allFamilies).Error; err != nil {
		log.Fatalf("Failed to fetch families: %v", err)
	}

	fmt.Printf("\n=== Seeding Summary ===\n")
	fmt.Printf("Created: %d\n", familiesCreated)
	fmt.Printf("Skipped: %d\n", familiesSkipped)
	fmt.Printf("Total in database: %d\n\n", len(allFamilies))

	fmt.Println("All families:")
	for i, family := range allFamilies {
		fmt.Printf("%d. %s | Founder: %s | Country: %s | Branch: %s\n",
			i+1, family.Name, family.Founder, family.Country, family.Branch)
	}
}
