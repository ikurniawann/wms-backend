// main.go - Updated with Router
// Application entry point

package main

import (
	"flag"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ikurniawann/wmsmicroservice/database"
	"github.com/ikurniawann/wmsmicroservice/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// Parse flags
	var migrateFlag bool
	flag.BoolVar(&migrateFlag, "migrate", false, "Run database migrations")
	flag.Parse()

	// Connect to database
	db, err := database.ConnectFromEnv()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	log.Println("✅ Database connected successfully!")

	// Run migration if flag is set
	if migrateFlag {
		log.Println("🔄 Running migrations...")
		if err := database.Migrate(db); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Println("✅ Migrations completed!")
		return
	}

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	if os.Getenv("ENV") == "development" {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📚 API Documentation: http://localhost:%s/health", port)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
