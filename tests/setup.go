// tests/setup.go
// Test setup and utilities

package tests

import (
	"os"
	"testing"

	"github.com/ikurniawann/wmsmicroservice/database"
	"github.com/ikurniawann/wmsmicroservice/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestDB holds the test database connection
var TestDB *gorm.DB

// SetupTestDB initializes an in-memory SQLite database for testing
func SetupTestDB() (*gorm.DB, error) {
	// Use SQLite in-memory for fast tests
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate all models
	err = db.AutoMigrate(
		// Public schema
		&models.Company{},
		&models.CompanyUser{},
		&models.Location{},
		&models.CompanySetting{},
		&models.NumberingFormat{},
		&models.AuditLog{},
		// Master data
		&models.ProductCategory{},
		&models.UnitOfMeasure{},
		&models.Product{},
		&models.ProductVariant{},
		&models.BomHeader{},
		&models.BomLine{},
		&models.CustomerGroup{},
		&models.Customer{},
		&models.CustomerAddress{},
		&models.Supplier{},
		// Inventory
		&models.LocationType{},
		&models.WarehouseZone{},
		&models.BinLocation{},
		&models.Stock{},
		&models.StockLevel{},
		&models.Movement{},
		&models.MovementLine{},
		&models.StockAdjustment{},
		&models.StockAdjustmentLine{},
		&models.StockTransfer{},
		&models.StockTransferLine{},
		// Sales
		&models.PriceList{},
		&models.PriceListLine{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderStatusHistory{},
		&models.Invoice{},
		&models.InvoiceItem{},
		&models.Payment{},
		&models.PaymentAllocation{},
		&models.SalesReturn{},
		&models.SalesReturnItem{},
		// Purchase
		&models.PurchaseOrder{},
		&models.PurchaseOrderItem{},
		&models.PurchaseInvoice{},
		&models.PurchaseInvoiceItem{},
		&models.PurchaseReturn{},
		&models.PurchaseReturnItem{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// TeardownTestDB closes the test database
func TeardownTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// CreateTestCompany creates a test company
func CreateTestCompany(db *gorm.DB) (*models.Company, error) {
	company := &models.Company{
		Code:        "TEST-CO",
		Name:        "Test Company",
		Email:       "test@example.com",
		Phone:       "1234567890",
		Address:     "123 Test St",
		City:        "Test City",
		Province:    "Test Province",
		PostalCode:  "12345",
		Country:     "ID",
		TaxID:       "123456789",
		Currency:    "IDR",
		Language:    "id",
		DateFormat:  "DD/MM/YYYY",
		TimeZone:    "Asia/Jakarta",
		IsActive:    true,
		FiscalYear:  2024,
	}

	if err := db.Create(company).Error; err != nil {
		return nil, err
	}

	return company, nil
}

// CreateTestProduct creates a test product
func CreateTestProduct(db *gorm.DB, companyID uint64) (*models.Product, error) {
	product := &models.Product{
		CompanyID:      companyID,
		Code:           "PROD-001",
		SKU:            "SKU-001",
		Name:           "Test Product",
		Description:    "A test product",
		TrackInventory: true,
		IsActive:       true,
	}

	if err := db.Create(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

// CreateTestCustomer creates a test customer
func CreateTestCustomer(db *gorm.DB, companyID uint64) (*models.Customer, error) {
	customer := &models.Customer{
		CompanyID: companyID,
		Code:      "CUST-001",
		Name:      "Test Customer",
		Email:     "customer@test.com",
		Phone:     "987654321",
		IsActive:  true,
	}

	if err := db.Create(customer).Error; err != nil {
		return nil, err
	}

	return customer, nil
}

// SkipIfNoEnv skips test if environment variable is not set
func SkipIfNoEnv(t *testing.T, envVar string) {
	if os.Getenv(envVar) == "" {
		t.Skipf("Skipping test: %s not set", envVar)
	}
}
