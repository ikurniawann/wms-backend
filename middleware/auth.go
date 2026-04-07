// middleware/auth.go
// Authentication middleware using Supabase JWT

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims from Supabase
type Claims struct {
	Sub        string                 `json:"sub"`
	Email      string                 `json:"email"`
	CompanyID  uint64                 `json:"company_id"`
	Role       string                 `json:"role"`
	Metadata   map[string]interface{} `json:"user_metadata"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT token from Supabase
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Parse token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.Sub)
		c.Set("email", claims.Email)
		c.Set("company_id", claims.CompanyID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// OptionalAuth allows optional authentication
func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(*Claims); ok {
				c.Set("user_id", claims.Sub)
				c.Set("email", claims.Email)
				c.Set("company_id", claims.CompanyID)
				c.Set("role", claims.Role)
			}
		}

		c.Next()
	}
}

// GetUserID returns user ID from context
func GetUserID(c *gin.Context) string {
	userID, _ := c.Get("user_id")
	if id, ok := userID.(string); ok {
		return id
	}
	return ""
}

// GetCompanyID returns company ID from context
func GetCompanyID(c *gin.Context) uint64 {
	companyID, exists := c.Get("company_id")
	if !exists {
		return 0
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

// GetUserEmail returns user email from context
func GetUserEmail(c *gin.Context) string {
	email, _ := c.Get("email")
	if e, ok := email.(string); ok {
		return e
	}
	return ""
}

// GetUserRole returns user role from context
func GetUserRole(c *gin.Context) string {
	role, _ := c.Get("role")
	if r, ok := role.(string); ok {
		return r
	}
	return ""
}

// RequireRole middleware checks user role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetUserRole(c)
		
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}
		
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}
