// routes/routes.go
// API route definitions

package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ikurniawann/wmsmicroservice/handlers"
	"github.com/ikurniawann/wmsmicroservice/middleware"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine) {
	// Get JWT secret from env
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // Default for development
	}

	// Public routes (no auth required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"version": "1.0.0",
			"service": "wms-backend",
		})
	})

	// API v1 group
	api := r.Group("/api/v1")

	// Apply CORS
	api.Use(middleware.DefaultCORS())

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Use Supabase Auth",
				"docs": "https://supabase.com/docs/guides/auth",
			})
		})
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	protected.Use(middleware.CompanyContextMiddleware())

	// Product routes
	productHandler := handlers.NewProductHandler()
	products := protected.Group("/products")
	{
		products.GET("", productHandler.ListProducts)
		products.POST("", productHandler.CreateProduct)
		products.GET("/:id", productHandler.GetProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}

	// Inventory routes
	inventoryHandler := handlers.NewInventoryHandler()
	inventory := protected.Group("/inventory")
	{
		inventory.GET("/stocks", inventoryHandler.ListStocks)
		inventory.GET("/stocks/:variant_id", inventoryHandler.GetStock)
		inventory.POST("/stocks/:variant_id/adjust", inventoryHandler.AdjustStock)
		inventory.GET("/summary", inventoryHandler.GetStockSummary)
	}

	// Customer routes
	customerHandler := handlers.NewCustomerHandler()
	customers := protected.Group("/customers")
	{
		customers.GET("", customerHandler.ListCustomers)
		customers.POST("", customerHandler.CreateCustomer)
		customers.GET("/:id", customerHandler.GetCustomer)
		customers.PUT("/:id", customerHandler.UpdateCustomer)
		customers.DELETE("/:id", customerHandler.DeleteCustomer)
	}

	// Supplier routes
	supplierHandler := handlers.NewSupplierHandler()
	suppliers := protected.Group("/suppliers")
	{
		suppliers.GET("", supplierHandler.ListSuppliers)
		suppliers.POST("", supplierHandler.CreateSupplier)
		suppliers.GET("/:id", supplierHandler.GetSupplier)
		suppliers.PUT("/:id", supplierHandler.UpdateSupplier)
		suppliers.DELETE("/:id", supplierHandler.DeleteSupplier)
	}

	// Sales Order routes
	orderHandler := handlers.NewOrderHandler()
	sales := protected.Group("/sales")
	{
		sales.GET("/orders", orderHandler.ListSalesOrders)
		sales.POST("/orders", orderHandler.CreateSalesOrder)
		sales.GET("/orders/:id", orderHandler.GetSalesOrder)
		sales.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus)
	}

	// Purchase Order routes
	purchaseHandler := handlers.NewPurchaseHandler()
	purchases := protected.Group("/purchase")
	{
		purchases.GET("/orders", purchaseHandler.ListPurchaseOrders)
		purchases.POST("/orders", purchaseHandler.CreatePurchaseOrder)
		purchases.GET("/orders/:id", purchaseHandler.GetPurchaseOrder)
		purchases.PATCH("/orders/:id/status", purchaseHandler.UpdatePOStatus)
	}

	// Reports routes (placeholder)
	reports := protected.Group("/reports")
	{
		reports.GET("/daily-sales", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Daily sales report endpoint"})
		})
		reports.GET("/product-performance", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Product performance report endpoint"})
		})
	}
}
