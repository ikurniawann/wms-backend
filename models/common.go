// models/common.go
// Common types and utilities

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// CompanyScopedModel adds company_id to base model
type CompanyScopedModel struct {
	BaseModel
	CompanyID uint64 `gorm:"not null;index" json:"company_id"`
}

// AuditModel adds audit fields
type AuditModel struct {
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`
}

// Pagination represents pagination parameters
type Pagination struct {
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	SortBy   string `json:"sort_by" query:"sort_by"`
	SortOrder string `json:"sort_order" query:"sort_order"`
}

// DefaultPagination returns default pagination values
func DefaultPagination() Pagination {
	return Pagination{
		Page:      1,
		PageSize:  20,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}

// GetOffset returns the offset for SQL query
func (p Pagination) GetOffset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetLimit()
}

// GetLimit returns the limit for SQL query
func (p Pagination) GetLimit() int {
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	return p.PageSize
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta represents metadata for paginated responses
type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(err string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   err,
	}
}

// NewPaginatedResponse creates a paginated response
func NewPaginatedResponse(data interface{}, pagination Pagination, total int64) APIResponse {
	totalPages := int(total) / pagination.GetLimit()
	if int(total)%pagination.GetLimit() > 0 {
		totalPages++
	}

	return APIResponse{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
