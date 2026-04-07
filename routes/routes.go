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
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 group
	api := r.Group("/api/v1")

	// Apply CORS
	api.Use(middleware.DefaultCORS())

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Use Supabase Auth"})
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

	// Sales Order routes
	orderHandler := handlers.NewOrderHandler()
	sales := protected.Group("/sales")
	{
		sales.GET("/orders", orderHandler.ListSalesOrders)
		sales.POST("/orders", orderHandler.CreateSalesOrder)
		sales.GET("/orders/:id", orderHandler.GetSalesOrder)
		sales.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus)
	}

	// TODO: Add more routes
	// - Purchase orders
	// - Customers
	// - Suppliers
	// - Reports
	// - Movements
}
