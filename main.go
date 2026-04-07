package main

import (
	"flag"
	"log"
	"os"

	"github.com/ikurniawann/wmsmicroservice/database"
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

	// TODO: Initialize Gin router and start server
	// router := gin.Default()
	// ... setup routes ...
	// router.Run(":" + os.Getenv("PORT"))

	log.Println("🚀 Server would start on port:", os.Getenv("PORT"))
	log.Println("💡 Use -migrate flag to run migrations")
}
