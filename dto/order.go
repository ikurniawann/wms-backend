// handlers/order.go
// Sales and Purchase Order handlers

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

// OrderHandler handles sales order requests
type OrderHandler struct{}

// NewOrderHandler creates new order handler
func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

// CreateSalesOrder creates a new sales order
func (h *OrderHandler) CreateSalesOrder(c *gin.Context) {
	var req dto.CreateSalesOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	companyID := middleware.GetCompanyID(c)
	userID := middleware.GetUserID(c)

	// Parse dates
	orderDate := time.Now()
	if req.OrderDate != "" {
		if d, err := time.Parse("2006-01-02", req.OrderDate); err == nil {
			orderDate = d
		}
	}

	var deliveryDate *time.Time
	if req.DeliveryDate != "" {
		if d, err := time.Parse("2006-01-02", req.DeliveryDate); err == nil {
			deliveryDate = &d
		}
	}

	// Calculate totals
	var subtotal, totalDiscount, totalTax, totalAmount float64
	var orderItems []models.OrderItem

	for _, itemReq := range req.Items {
		itemSubtotal := itemReq.Quantity * itemReq.UnitPrice
		itemDiscount := itemSubtotal * itemReq.DiscountPercent / 100
		taxableAmount := itemSubtotal - itemDiscount
		itemTax := taxableAmount * itemReq.TaxPercent / 100

		item := models.OrderItem{
			VariantID:       itemReq.VariantID,
			QuantityOrdered: itemReq.Quantity,
			UnitID:          itemReq.UnitID,
			UnitPrice:       itemReq.UnitPrice,
			DiscountPercent: itemReq.DiscountPercent,
			DiscountAmount:  itemDiscount,
			TaxPercent:      itemReq.TaxPercent,
			TaxAmount:       itemTax,
			TotalPrice:      taxableAmount + itemTax,
			Notes:           itemReq.Notes,
		}
		orderItems = append(orderItems, item)

		subtotal += itemSubtotal
		totalDiscount += itemDiscount
		totalTax += itemTax
	}

	totalAmount = subtotal - totalDiscount + totalTax + req.ShippingCost

	// Generate order number if not provided
	orderNumber := req.OrderNumber
	if orderNumber == "" {
		orderNumber = generateOrderNumber(companyID, "SO")
	}

	order := models.Order{
		CompanyID:         companyID,
		OrderNumber:       orderNumber,
		OrderDate:         orderDate,
		DeliveryDate:      deliveryDate,
		LocationID:        req.LocationID,
		CustomerID:        req.CustomerID,
		CustomerAddressID: req.CustomerAddressID,
		Subtotal:          subtotal,
		DiscountAmount:    totalDiscount + req.DiscountAmount,
		DiscountPercent:   req.DiscountPercent,
		TaxAmount:         totalTax,
		TaxPercent:        req.TaxPercent,
		ShippingCost:      req.ShippingCost,
		TotalAmount:       totalAmount,
		Status:            "draft",
		PaymentStatus:     "unpaid",
		Notes:             req.Notes,
		InternalNotes:     req.InternalNotes,
		Items:             orderItems,
	}

	// Set created_by
	if userID != "" {
		uid := parseUUID(userID)
		order.CreatedBy = &uid
	}

	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to create order: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewSuccessResponse(order, "Order created successfully"))
}

// GetSalesOrder gets an order by ID
func (h *OrderHandler) GetSalesOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid order ID"))
		return
	}

	companyID := middleware.GetCompanyID(c)

	var order models.Order
	if err := database.DB.
		Where("company_id = ?", companyID).
		Preload("Items").
		First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Order not found"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(order, "Order retrieved successfully"))
}

// ListSalesOrders lists orders with pagination
func (h *OrderHandler) ListSalesOrders(c *gin.Context) {
	var req dto.ListOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}
	req.DefaultPagination()

	companyID := middleware.GetCompanyID(c)

	query := database.DB.Where("company_id = ?", companyID)

	// Apply filters
	if req.CustomerID > 0 {
		query = query.Where("customer_id = ?", req.CustomerID)
	}
	if req.LocationID > 0 {
		query = query.Where("location_id = ?", req.LocationID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("order_date BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	// Count total
	var total int64
	query.Model(&models.Order{}).Count(&total)

	// Get orders
	var orders []models.Order
	err := query.
		Order(req.SortBy + " " + req.SortOrder).
		Limit(req.PageSize).
		Offset(req.Offset()).
		Find(&orders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to fetch orders"))
		return
	}

	pagination := dto.PaginationResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponseWithPagination(orders, "Orders retrieved successfully", pagination))
}

// UpdateOrderStatus updates order status
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid order ID"))
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	companyID := middleware.GetCompanyID(c)
	userID := middleware.GetUserID(c)

	var order models.Order
	if err := database.DB.Where("company_id = ?", companyID).First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Order not found"))
		return
	}

	updates := map[string]interface{}{
		"status": req.Status,
	}

	// Update timestamps based on status
	now := time.Now()
	switch req.Status {
	case "confirmed":
		updates["confirmed_at"] = now
	case "shipped":
		updates["shipped_at"] = now
	case "delivered":
		updates["delivered_at"] = now
	case "cancelled":
		updates["cancelled_at"] = now
	}

	if userID != "" {
		uid := parseUUID(userID)
		updates["updated_by"] = uid
	}

	if err := database.DB.Model(&order).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to update order status"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(order, "Order status updated successfully"))
}

// generateOrderNumber generates unique order number
func generateOrderNumber(companyID uint64, prefix string) string {
	// Simple implementation - use proper numbering in production
	now := time.Now()
	return prefix + "-" + now.Format("20060102") + "-" + strconv.FormatUint(companyID, 10)
}
