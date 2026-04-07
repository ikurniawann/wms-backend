// handlers/customer.go
// Customer handlers

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

// CustomerHandler handles customer requests
type CustomerHandler struct{}

// NewCustomerHandler creates new customer handler
func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{}
}

// CreateCustomer creates a new customer
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req struct {
		Code            string  `json:"code" binding:"required,max=100"`
		CustomerGroupID *uint64 `json:"customer_group_id"`
		Name            string  `json:"name" binding:"required,max=250"`
		Email           string  `json:"email" binding:"omitempty,email,max=150"`
		Phone           string  `json:"phone" binding:"omitempty,max=50"`
		Mobile          string  `json:"mobile" binding:"omitempty,max=50"`
		TaxID           string  `json:"tax_id" binding:"omitempty,max=50"`
		BillingAddress  string  `json:"billing_address"`
		BillingCity     string  `json:"billing_city" binding:"omitempty,max=100"`
		BillingProvince string  `json:"billing_province" binding:"omitempty,max=100"`
		BillingPostal   string  `json:"billing_postal_code" binding:"omitempty,max=20"`
		CreditLimit     float64 `json:"credit_limit" binding:"omitempty,gte=0"`
		PriceTier       int     `json:"price_tier" binding:"omitempty,min=1,max=10"`
		IsWholesale     bool    `json:"is_wholesale"`
		IsActive        bool    `json:"is_active"`
		BirthDate       string  `json:"birth_date" binding:"omitempty,datetime=2006-01-02"`
		Gender          string  `json:"gender" binding:"omitempty,oneof=male female other"`
		Notes           string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	companyID := middleware.GetCompanyID(c)

	customer := models.Customer{
		CompanyID:           companyID,
		Code:                req.Code,
		CustomerGroupID:     req.CustomerGroupID,
		Name:                req.Name,
		Email:               req.Email,
		Phone:               req.Phone,
		Mobile:              req.Mobile,
		TaxID:               req.TaxID,
		BillingAddress:      req.BillingAddress,
		BillingCity:         req.BillingCity,
		BillingProvince:     req.BillingProvince,
		BillingPostalCode:   req.BillingPostal,
		CreditLimit:         req.CreditLimit,
		PriceTier:           req.PriceTier,
		IsWholesale:         req.IsWholesale,
		IsActive:            req.IsActive,
		Gender:              req.Gender,
		Notes:               req.Notes,
	}

	// Parse birth date
	if req.BirthDate != "" {
		if d, err := parseDate(req.BirthDate); err == nil {
			customer.BirthDate = &d
		}
	}

	if err := database.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to create customer: "+err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.NewSuccessResponse(customer, "Customer created successfully"))
}

// GetCustomer gets a customer by ID
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid customer ID"))
		return
	}

	companyID := middleware.GetCompanyID(c)

	var customer models.Customer
	if err := database.DB.
		Where("company_id = ?", companyID).
		Preload("CustGroup").
		Preload("Addresses").
		First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Customer not found"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(customer, "Customer retrieved successfully"))
}

// ListCustomers lists all customers
func (h *CustomerHandler) ListCustomers(c *gin.Context) {
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
	query.Model(&models.Customer{}).Count(&total)

	// Get customers
	var customers []models.Customer
	err := query.
		Preload("CustGroup").
		Order(req.SortBy + " " + req.SortOrder).
		Limit(req.PageSize).
		Offset(req.Offset()).
		Find(&customers).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to fetch customers"))
		return
	}

	pagination := dto.PaginationResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponseWithPagination(customers, "Customers retrieved successfully", pagination))
}

// UpdateCustomer updates a customer
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid customer ID"))
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse(err.Error()))
		return
	}

	companyID := middleware.GetCompanyID(c)

	var customer models.Customer
	if err := database.DB.Where("company_id = ?", companyID).First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Customer not found"))
		return
	}

	if err := database.DB.Model(&customer).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to update customer"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(customer, "Customer updated successfully"))
}

// DeleteCustomer soft deletes a customer
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("Invalid customer ID"))
		return
	}

	companyID := middleware.GetCompanyID(c)

	var customer models.Customer
	if err := database.DB.Where("company_id = ?", companyID).First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.NewErrorResponse("Customer not found"))
		return
	}

	if err := database.DB.Delete(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("Failed to delete customer"))
		return
	}

	c.JSON(http.StatusOK, dto.NewSuccessResponse(nil, "Customer deleted successfully"))
}
