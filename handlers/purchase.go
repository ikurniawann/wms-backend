// handlers/purchase.go
// Purchase Order handlers

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

// PurchaseHandler handles purchase order requests
type PurchaseHandler struct{}

// NewPurchaseHandler creates new purchase handler
func NewPurchaseHandler() *PurchaseHandler {
	return &PurchaseHandler{}
}

// CreatePurchaseOrder creates a new PO
func (h *PurchaseHandler) CreatePurchaseOrder(c *gin.Context) {
	var req dto.CreatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	companyID := middleware.GetCompanyID(c)
	userID := middleware.GetUserID(c)

	// Parse dates
	poDate := time.Now()
	if req.PODate != "" {
		if d, err := time.Parse("2006-01-02", req.PODate); err == nil {
			poDate = d
		}
	}

	var expectedDate *time.Time
	if req.ExpectedDeliveryDate != "" {
		if d, err := time.Parse("2006-01-02", req.ExpectedDeliveryDate); err == nil {
			expectedDate = &d
		}
	}

	// Calculate totals
	var subtotal, totalDiscount, totalTax, totalAmount float64
	var poItems []models.PurchaseOrderItem

	for _, itemReq := range req.Items {
		itemSubtotal := itemReq.Quantity * itemReq.UnitPrice
		itemDiscount := itemSubtotal * itemReq.DiscountPercent / 100
		taxableAmount := itemSubtotal - itemDiscount
		itemTax := taxableAmount * req.TaxPercent / 100

		item := models.PurchaseOrderItem{
			VariantID:         itemReq.VariantID,
			QuantityOrdered:   itemReq.Quantity,
			UnitID:            itemReq.UnitID,
			UnitPrice:         itemReq.UnitPrice,
			DiscountPercent:   itemReq.DiscountPercent,
			TotalPrice:        taxableAmount + itemTax,
			Notes:             itemReq.Notes,
		}
		poItems = append(poItems, item)

		subtotal += itemSubtotal
		totalDiscount += itemDiscount
		totalTax += itemTax
	}

	totalAmount = subtotal - totalDiscount + totalTax + req.ShippingCost

	// Generate PO number if not provided
	poNumber := req.PONumber
	if poNumber == "" {
		poNumber = generatePONumber(companyID, "PO")
	}

	po := models.PurchaseOrder{
		CompanyID:            companyID,
		PONumber:             poNumber,
		PODate:               poDate,
		ExpectedDeliveryDate: expectedDate,
		LocationID:           req.LocationID,
		SupplierID:           req.SupplierID,
		Subtotal:             subtotal,
		DiscountAmount:       totalDiscount + req.DiscountAmount,
		TaxAmount:            totalTax,
		TaxPercent:           req.TaxPercent,
		ShippingCost:         req.ShippingCost,
		TotalAmount:          totalAmount,
		Status:               "draft",
		Notes:                req.Notes,
		InternalNotes:        req.InternalNotes,
		Items:                poItems,
	}

	// Set created_by
	if userID != "" {
		uid := parseUUID(userID)
		po.CreatedBy = &uid
	}

	if err := database.DB.Create(&po).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to create purchase order: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewSuccessResponse(po, "Purchase order created successfully"))
}

// GetPurchaseOrder gets a PO by ID
func (h *PurchaseHandler) GetPurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid PO ID"))
		return
	}

	companyID := middleware.GetCompanyID(c)

	var po models.PurchaseOrder
	if err := database.DB.
		Where("company_id = ?", companyID).
		Preload("Items").
		Preload("Supplier").
		Preload("Location").
		First(&po, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Purchase order not found"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(po, "Purchase order retrieved successfully"))
}

// ListPurchaseOrders lists all POs
func (h *PurchaseHandler) ListPurchaseOrders(c *gin.Context) {
	var req dto.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}
	req.DefaultPagination()

	companyID := middleware.GetCompanyID(c)

	query := database.DB.Where("company_id = ?", companyID)

	// Apply filters
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("po_date BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	// Count total
	var total int64
	query.Model(&models.PurchaseOrder{}).Count(&total)

	// Get POs
	var pos []models.PurchaseOrder
	err := query.
		Preload("Supplier").
		Order(req.SortBy + " " + req.SortOrder).
		Limit(req.PageSize).
		Offset(req.Offset()).
		Find(&pos).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to fetch purchase orders"))
		return
	}

	pagination := dto.PaginationResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponseWithPagination(pos, "Purchase orders retrieved successfully", pagination))
}

// UpdatePOStatus updates PO status
func (h *PurchaseHandler) UpdatePOStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid PO ID"))
		return
	}

	var req dto.UpdatePOStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	companyID := middleware.GetCompanyID(c)
	userID := middleware.GetUserID(c)

	var po models.PurchaseOrder
	if err := database.DB.Where("company_id = ?", companyID).First(&po, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Purchase order not found"))
		return
	}

	updates := map[string]interface{}{
		"status": req.Status,
	}

	if userID != "" {
		uid := parseUUID(userID)
		updates["updated_by"] = uid
	}

	if err := database.DB.Model(&po).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to update PO status"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(po, "Purchase order status updated"))
}

// generatePONumber generates unique PO number
func generatePONumber(companyID uint64, prefix string) string {
	now := time.Now()
	return prefix + "-" + now.Format("20060102") + "-" + strconv.FormatUint(companyID, 10)
}
