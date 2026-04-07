// handlers/product.go
// Product HTTP handlers

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

// ProductHandler handles product-related requests
type ProductHandler struct{}

// NewProductHandler creates new product handler
func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

// CreateProduct creates a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	companyID := middleware.GetCompanyID(c)
	userID := middleware.GetUserID(c)

	product := models.Product{
		CompanyID:        companyID,
		Code:             req.Code,
		SKU:              req.SKU,
		Barcode:          req.Barcode,
		Name:             req.Name,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		CategoryID:       req.CategoryID,
		Brand:            req.Brand,
		BaseUnitID:       req.BaseUnitID,
		Weight:           req.Weight,
		Dimensions:       req.Dimensions,
		TrackInventory:   req.TrackInventory,
		IsService:        req.IsService,
		IsActive:         req.IsActive,
		Metadata:         req.Metadata,
	}

	// Set created_by from user context
	if userID != "" {
		uid := parseUUID(userID)
		product.CreatedBy = &uid
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to create product: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewSuccessResponse(product, "Product created successfully"))
}

// GetProduct gets a product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid product ID"))
		return
	}

	companyID := middleware.GetCompanyID(c)

	var product models.Product
	if err := database.DB.
		Where("company_id = ?", companyID).
		Preload("Category").
		Preload("BaseUnit").
		Preload("Variants").
		First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Product not found"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(product, "Product retrieved successfully"))
}

// ListProducts lists all products with pagination
func (h *ProductHandler) ListProducts(c *gin.Context) {
	var req dto.ListProductsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}
	req.DefaultPagination()

	companyID := middleware.GetCompanyID(c)

	query := database.DB.Where("company_id = ?", companyID)

	// Apply filters
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}
	if req.TrackInventory != nil {
		query = query.Where("track_inventory = ?", *req.TrackInventory)
	}
	if req.IsService != nil {
		query = query.Where("is_service = ?", *req.IsService)
	}
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}
	if req.Query != "" {
		query = query.Where("name ILIKE ? OR code ILIKE ? OR sku ILIKE ?",
			"%"+req.Query+"%", "%"+req.Query+"%", "%"+req.Query+"%")
	}

	// Count total
	var total int64
	query.Model(&models.Product{}).Count(&total)

	// Get products
	var products []models.Product
	err := query.
		Preload("Category").
		Preload("BaseUnit").
		Order(req.SortBy + " " + req.SortOrder).
		Limit(req.PageSize).
		Offset(req.Offset()).
		Find(&products).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to fetch products"))
		return
	}

	pagination := dto.PaginationResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponseWithPagination(products, "Products retrieved successfully", pagination))
}

// UpdateProduct updates a product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid product ID"))
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	companyID := middleware.GetCompanyID(c)
	userID := middleware.GetUserID(c)

	var product models.Product
	if err := database.DB.Where("company_id = ?", companyID).First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Product not found"))
		return
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.SKU != "" {
		updates["sku"] = req.SKU
	}
	if req.Barcode != "" {
		updates["barcode"] = req.Barcode
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.ShortDescription != "" {
		updates["short_description"] = req.ShortDescription
	}
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.Brand != "" {
		updates["brand"] = req.Brand
	}
	if req.BaseUnitID != nil {
		updates["base_unit_id"] = *req.BaseUnitID
	}
	if req.Weight > 0 {
		updates["weight"] = req.Weight
	}
	if req.Dimensions != nil {
		updates["dimensions"] = req.Dimensions
	}
	if req.TrackInventory != nil {
		updates["track_inventory"] = *req.TrackInventory
	}
	if req.IsService != nil {
		updates["is_service"] = *req.IsService
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.Metadata != nil {
		updates["metadata"] = req.Metadata
	}

	// Set updated_by
	if userID != "" {
		uid := parseUUID(userID)
		updates["updated_by"] = uid
	}

	if err := database.DB.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to update product"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(product, "Product updated successfully"))
}

// DeleteProduct soft deletes a product
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid product ID"))
		return
	}

	companyID := middleware.GetCompanyID(c)

	var product models.Product
	if err := database.DB.Where("company_id = ?", companyID).First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Product not found"))
		return
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to delete product"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(nil, "Product deleted successfully"))
}

// Helper function
func parseUUID(s string) [16]byte {
	var uuid [16]byte
	// Simple parsing, implement proper UUID parsing if needed
	copy(uuid[:], []byte(s))
	return uuid
}
