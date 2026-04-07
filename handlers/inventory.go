// handlers/inventory.go
// Inventory handlers

package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ikurniawann/wmsmicroservice/database"
	"github.com/ikurniawann/wmsmicroservice/dto"
	"github.com/ikurniawann/wmsmicroservice/middleware"
	"github.com/ikurniawann/wmsmicroservice/models"
)

// InventoryHandler handles inventory requests
type InventoryHandler struct{}

// NewInventoryHandler creates new inventory handler
func NewInventoryHandler() *InventoryHandler {
	return &InventoryHandler{}
}

// GetStock gets stock levels for a variant
func (h *InventoryHandler) GetStock(c *gin.Context) {
	variantID, err := strconv.ParseUint(c.Param("variant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid variant ID"))
		return
	}

	companyID := middleware.GetCompanyID(c)

	var stock models.Stock
	if err := database.DB.
		Where("company_id = ? AND variant_id = ?", companyID, variantID).
		Preload("Variant").
		Preload("Location").
		First(&stock).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Stock not found"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(stock, "Stock retrieved successfully"))
}

// ListStocks lists all stock levels
func (h *InventoryHandler) ListStocks(c *gin.Context) {
	var req dto.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}
	req.DefaultPagination()

	companyID := middleware.GetCompanyID(c)

	query := database.DB.Where("company_id = ?", companyID)

	// Apply filters
	if req.LocationID > 0 {
		query = query.Where("location_id = ?", req.LocationID)
	}

	// Count total
	var total int64
	query.Model(&models.Stock{}).Count(&total)

	// Get stocks with joins
	var stocks []models.Stock
	err := query.
		Preload("Variant").
		Preload("Variant.Product").
		Preload("Location").
		Order(req.SortBy + " " + req.SortOrder).
		Limit(req.PageSize).
		Offset(req.Offset()).
		Find(&stocks).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to fetch stocks"))
		return
	}

	pagination := dto.PaginationResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponseWithPagination(stocks, "Stocks retrieved successfully", pagination))
}

// AdjustStock adjusts stock quantity
func (h *InventoryHandler) AdjustStock(c *gin.Context) {
	variantID, err := strconv.ParseUint(c.Param("variant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid variant ID"))
		return
	}

	var req struct {
		LocationID uint64  `json:"location_id" binding:"required"`
		Quantity   float64 `json:"quantity" binding:"required"`
		Reason     string  `json:"reason" binding:"required"`
		Notes      string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	companyID := middleware.GetCompanyID(c)
	userID := middleware.GetUserID(c)

	// Find or create stock record
	var stock models.Stock
	result := database.DB.Where(
		"company_id = ? AND variant_id = ? AND location_id = ?",
		companyID, variantID, req.LocationID,
	).First(&stock)

	if result.Error != nil {
		// Create new stock record
		stock = models.Stock{
			CompanyID:         companyID,
			VariantID:           variantID,
			LocationID:          req.LocationID,
			QuantityAvailable: req.Quantity,
		}
		if err := database.DB.Create(&stock).Error; err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to create stock"))
			return
		}
	} else {
		// Update existing stock
		if err := database.DB.Model(&stock).Update("quantity_available", req.Quantity).Error; err != nil {
			c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to update stock"))
			return
		}
	}

	// Create movement record for audit
	movement := models.Movement{
		CompanyID:      companyID,
		MovementNumber: "ADJ-" + time.Now().Format("20060102-150405"),
		MovementDate:   time.Now(),
		MovementType:   "adjustment",
		Reason:         req.Reason,
		LocationID:     req.LocationID,
		TotalQuantity:  req.Quantity,
		Notes:          req.Notes,
	}

	if userID != "" {
		uid := parseUUID(userID)
		movement.CreatedBy = &uid
	}

	if err := database.DB.Create(&movement).Error; err != nil {
		// Log error but don't fail the request
		// movement creation is for audit only
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(stock, "Stock adjusted successfully"))
}

// GetStockSummary gets stock summary from view
func (h *InventoryHandler) GetStockSummary(c *gin.Context) {
	companyID := middleware.GetCompanyID(c)

	var summaries []models.StockSummaryView
	if err := database.DB.
		Where("company_id = ?", companyID).
		Find(&summaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to fetch stock summary"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(summaries, "Stock summary retrieved"))
}
