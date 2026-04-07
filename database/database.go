// database/database.go
// Database connection and initialization

package database

import (
	"fmt"
	"os"

	"github.com/ikurniawann/wmsmicroservice/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// Connect initializes database connection
func Connect(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	return db, nil
}

// ConnectFromEnv connects using environment variables
func ConnectFromEnv() (*gorm.DB, error) {
	cfg := Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		Database: getEnv("DB_NAME", "neowms"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	return Connect(cfg)
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// System
		&models.Company{},
		&models.CompanyUser{},
		&models.Location{},
		&models.AuditLog{},

		// Master Data
		&models.Unit{},
		&models.Category{},
		&models.Attribute{},
		&models.Product{},
		&models.ProductVariant{},
		&models.ProductImage{},

		// Inventory
		&models.Stock{},
		&models.StorageBin{},
		&models.StockBin{},
		&models.Movement{},
		&models.MovementDetail{},
		&models.Reservation{},

		// Sales
		&models.CustomerGroup{},
		&models.Customer{},
		&models.CustomerAddress{},
		&models.Order{},
		&models.OrderItem{},
		&models.PriceTier{},
		&models.DiscountEvent{},

		// Purchase
		&models.Supplier{},
		&models.SupplierPrice{},
		&models.PurchaseOrder{},
		&models.PurchaseOrderItem{},
		&models.Receipt{},
		&models.ReceiptItem{},
	)
}

// Close closes database connection
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
