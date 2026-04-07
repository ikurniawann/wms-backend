// middleware/company.go
// Company context middleware for multi-tenancy

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CompanyContextMiddleware ensures company isolation
func CompanyContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID := GetCompanyID(c)
		if companyID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Company context required"})
			c.Abort()
			return
		}

		c.Set("current_company_id", companyID)
		c.Next()
	}
}

// CompanyScope adds company filter to database queries
func CompanyScope(db *gorm.DB, c *gin.Context) *gorm.DB {
	companyID := GetCompanyID(c)
	if companyID > 0 {
		return db.Where("company_id = ?", companyID)
	}
	return db
}

// GetCurrentCompanyID returns company ID from context
func GetCurrentCompanyID(c *gin.Context) uint64 {
	companyID, exists := c.Get("current_company_id")
	if !exists {
		return GetCompanyID(c)
	}
	
	switch v := companyID.(type) {
	case uint64:
		return v
	case float64:
		return uint64(v)
	case int:
		return uint64(v)
	default:
		return 0
	}
}

// IsCompanyMember checks if user belongs to company
func IsCompanyMember(c *gin.Context) bool {
	return GetCurrentCompanyID(c) > 0
}
