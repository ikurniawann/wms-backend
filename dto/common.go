// dto/common.go
// Common DTOs for request/response

package dto

import (
	"time"
)

// PaginationRequest for list endpoints
type PaginationRequest struct {
	Page     int    `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
	SortBy   string `form:"sort_by" json:"sort_by"`
	SortOrder string `form:"sort_order" json:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// DefaultPagination returns default values
func (p *PaginationRequest) DefaultPagination() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 20
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.SortOrder == "" {
		p.SortOrder = "desc"
	}
}

// Offset returns SQL offset
func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// PaginationResponse for paginated results
type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// CalculateTotalPages calculates total pages
func (p *PaginationResponse) CalculateTotalPages() {
	if p.PageSize > 0 {
		p.TotalPages = int(p.Total) / p.PageSize
		if p.Total%int64(p.PageSize) > 0 {
			p.TotalPages++
		}
	}
}

// APIResponse standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta contains metadata for response
type Meta struct {
	*PaginationResponse
	Timestamp time.Time `json:"timestamp"`
}

// NewSuccessResponse creates success response
func NewSuccessResponse(data interface{}, message string) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewSuccessResponseWithPagination creates success response with pagination
func NewSuccessResponseWithPagination(data interface{}, message string, pagination PaginationResponse) APIResponse {
	pagination.CalculateTotalPages()
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &Meta{
			PaginationResponse: &pagination,
			Timestamp:          time.Now(),
		},
	}
}

// NewErrorResponse creates error response
func NewErrorResponse(message string) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
	}
}

// IDRequest for endpoints requiring only ID
type IDRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

// UUIDRequest for endpoints requiring UUID
type UUIDRequest struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

// CodeRequest for lookup by code
type CodeRequest struct {
	Code string `uri:"code" binding:"required"`
}

// SoftDeleteRequest for soft delete
type SoftDeleteRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

// BatchIDsRequest for batch operations
type BatchIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"required,min=1"`
}

// SearchRequest for search endpoints
type SearchRequest struct {
	PaginationRequest
	Query      string `form:"q" json:"q"`
	CategoryID uint64 `form:"category_id" json:"category_id"`
	LocationID uint64 `form:"location_id" json:"location_id"`
	Status     string `form:"status" json:"status"`
	StartDate  string `form:"start_date" json:"start_date"`
	EndDate    string `form:"end_date" json:"end_date"`
}
