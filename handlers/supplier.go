// handlers/supplier.go
// Supplier handlers

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ikurniawann/wmsmicroservice/database"
	"github.com/ikurniawann/wmsmicroservice/dto"
	"github.com/ikurniawann/wmsmicroservice/middleware"
	"github.com/ikurniawann/wmsmicroservice/models"
)

// SupplierHandler handles supplier requests
type SupplierHandler struct{}

// NewSupplierHandler creates new supplier handler
func NewSupplierHandler() *SupplierHandler {
	return &SupplierHandler{}
}

// CreateSupplier creates a new supplier
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req struct {
		Code             string  `json:"code" binding:"required,max=100"`
		Name             string  `json:"name" binding:"required,max=250"`
		Email            string  `json:"email" binding:"omitempty,email,max=150"`
		Phone            string  `json:"phone" binding:"omitempty,max=50"`
		Mobile           string  `json:"mobile" binding:"omitempty,max=50"`
		TaxID            string  `json:"tax_id" binding:"omitempty,max=50"`
		Address          string  `json:"address"`
		City             string  `json:"city" binding:"omitempty,max=100"`
		Province         string  `json:"province" binding:"omitempty,max=100"`
		PostalCode       string  `json:"postal_code" binding:"omitempty,max=20"`
		Website          string  `json:"website" binding:"omitempty,url,max=150"`
		PaymentTermsDays int     `json:"payment_terms_days" binding:"omitempty,min=0"`
		CreditLimit      float64 `json:"credit_limit" binding:"omitempty,gte=0"`
		IsActive         bool    `json:"is_active"`
		Notes            string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	companyID := middleware.GetCompanyID(c)

	supplier := models.Supplier{
		CompanyID:        companyID,
		Code:             req.Code,
		Name:             req.Name,
		Email:            req.Email,
		Phone:            req.Phone,
		Mobile:           req.Mobile,
		TaxID:            req.TaxID,
		Address:          req.Address,
		City:             req.City,
		Province:         req.Province,
		PostalCode:       req.PostalCode,
		Website:          req.Website,
		PaymentTermsDays: req.PaymentTermsDays,
		CreditLimit:      req.CreditLimit,
		IsActive:         req.IsActive,
		Notes:            req.Notes,
	}

	if err := database.DB.Create(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to create supplier: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewSuccessResponse(supplier, "Supplier created successfully"))
}

// GetSupplier gets a supplier by ID
func (h *SupplierHandler) GetSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid supplier ID"))
		return
	}

	companyID := middleware.GetCompanyID(c)

	var supplier models.Supplier
	if err := database.DB.
		Where("company_id = ?", companyID).
		First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Supplier not found"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(supplier, "Supplier retrieved successfully"))
}

// ListSuppliers lists all suppliers
func (h *SupplierHandler) ListSuppliers(c *gin.Context) {
	var req dto.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}
	req.DefaultPagination()

	companyID := middleware.GetCompanyID(c)

	query := database.DB.Where("company_id = ?", companyID)

	// Apply filters
	if req.Query != "" {
		query = query.Where("name ILIKE ? OR code ILIKE ? OR email ILIKE ?",
			"%"+req.Query+"%", "%"+req.Query+"%", "%"+req.Query+"%")
	}

	// Count total
	var total int64
	query.Model(&models.Supplier{}).Count(&total)

	// Get suppliers
	var suppliers []models.Supplier
	err := query.
		Order(req.SortBy + " " + req.SortOrder).
		Limit(req.PageSize).
		Offset(req.Offset()).
		Find(&suppliers).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to fetch suppliers"))
		return
	}

	pagination := dto.PaginationResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponseWithPagination(suppliers, "Suppliers retrieved successfully", pagination))
}

// UpdateSupplier updates a supplier
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid supplier ID"))
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	companyID := middleware.GetCompanyID(c)

	var supplier models.Supplier
	if err := database.DB.Where("company_id = ?", companyID).First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Supplier not found"))
		return
	}

	if err := database.DB.Model(&supplier).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to update supplier"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(supplier, "Supplier updated successfully"))
}

// DeleteSupplier soft deletes a supplier
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid supplier ID"))
		return
	}

	companyID := middleware.GetCompanyID(c)

	var supplier models.Supplier
	if err := database.DB.Where("company_id = ?", companyID).First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Supplier not found"))
		return
	}

	if err := database.DB.Delete(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to delete supplier"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(nil, "Supplier deleted successfully"))
}
