// models/company.go
// Generated from Supabase Schema - NEO WMS
// Date: 2026-04-07

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Company represents the root tenant entity
type Company struct {
	ID                    uint64         `gorm:"primaryKey;column:id" json:"id"`
	UUID                  uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;column:uuid" json:"uuid"`
	Code                  string         `gorm:"uniqueIndex;size:50;not null;column:code" json:"code"`
	Name                  string         `gorm:"size:200;not null;column:name" json:"name"`
	LegalName             string         `gorm:"size:200;column:legal_name" json:"legal_name,omitempty"`
	TaxID                 string         `gorm:"size:50;column:tax_id" json:"tax_id,omitempty"`
	Email                 string         `gorm:"size:150;column:email" json:"email,omitempty"`
	Phone                 string         `gorm:"size:50;column:phone" json:"phone,omitempty"`
	Website               string         `gorm:"size:150;column:website" json:"website,omitempty"`
	LogoURL               string         `gorm:"size:500;column:logo_url" json:"logo_url,omitempty"`
	Address               string         `gorm:"type:text;column:address" json:"address,omitempty"`
	City                  string         `gorm:"size:100;column:city" json:"city,omitempty"`
	Province              string         `gorm:"size:100;column:province" json:"province,omitempty"`
	PostalCode            string         `gorm:"size:20;column:postal_code" json:"postal_code,omitempty"`
	Country               string         `gorm:"size:100;default:Indonesia;column:country" json:"country,omitempty"`
	Timezone              string         `gorm:"size:50;default:Asia/Jakarta;column:timezone" json:"timezone"`
	Currency              string         `gorm:"size:3;default:IDR;column:currency" json:"currency"`
	SubscriptionTier      string         `gorm:"size:50;default:basic;column:subscription_tier" json:"subscription_tier"`
	SubscriptionExpiresAt *time.Time     `gorm:"column:subscription_expires_at" json:"subscription_expires_at,omitempty"`
	Settings              map[string]interface{} `gorm:"type:jsonb;default:'{}';column:settings" json:"settings,omitempty"`
	IsActive              bool           `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedAt             time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Locations    []Location    `gorm:"foreignKey:CompanyID" json:"locations,omitempty"`
	CompanyUsers []CompanyUser `gorm:"foreignKey:CompanyID" json:"company_users,omitempty"`
}

func (Company) TableName() string {
	return "public.companies"
}

// CompanyUser links Supabase Auth users to companies
type CompanyUser struct {
	ID           uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID    uint64         `gorm:"not null;index;column:company_id" json:"company_id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;column:user_id" json:"user_id"`
	Email        string         `gorm:"size:150;not null;column:email" json:"email"`
	FullName     string         `gorm:"size:200;column:full_name" json:"full_name,omitempty"`
	AvatarURL    string         `gorm:"size:500;column:avatar_url" json:"avatar_url,omitempty"`
	Role         string         `gorm:"size:50;default:member;column:role" json:"role"`
	Department   string         `gorm:"size:100;column:department" json:"department,omitempty"`
	Phone        string         `gorm:"size:50;column:phone" json:"phone,omitempty"`
	IsActive     bool           `gorm:"default:true;not null;column:is_active" json:"is_active"`
	LastLoginAt  *time.Time     `gorm:"column:last_login_at" json:"last_login_at,omitempty"`
	CreatedAt    time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`

	// Relationships
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (CompanyUser) TableName() string {
	return "public.company_users"
}

// Location represents warehouses, stores, or virtual locations
type Location struct {
	ID                   uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID            uint64         `gorm:"not null;index;column:company_id" json:"company_id"`
	UUID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();column:uuid" json:"uuid"`
	Code                 string         `gorm:"size:50;not null;uniqueIndex:idx_location_code_company;column:code" json:"code"`
	Name                 string         `gorm:"size:200;not null;column:name" json:"name"`
	LocationType         string         `gorm:"size:50;default:warehouse;column:location_type" json:"location_type"`
	Address              string         `gorm:"type:text;column:address" json:"address,omitempty"`
	City                 string         `gorm:"size:100;column:city" json:"city,omitempty"`
	Province             string         `gorm:"size:100;column:province" json:"province,omitempty"`
	PostalCode           string         `gorm:"size:20;column:postal_code" json:"postal_code,omitempty"`
	Phone                string         `gorm:"size:50;column:phone" json:"phone,omitempty"`
	Email                string         `gorm:"size:150;column:email" json:"email,omitempty"`
	Latitude             *float64       `gorm:"type:decimal(10,8);column:latitude" json:"latitude,omitempty"`
	Longitude            *float64       `gorm:"type:decimal(11,8);column:longitude" json:"longitude,omitempty"`
	IsMainLocation       bool           `gorm:"default:false;column:is_main_location" json:"is_main_location"`
	IsSalesLocation      bool           `gorm:"default:false;column:is_sales_location" json:"is_sales_location"`
	IsPurchaseLocation   bool           `gorm:"default:true;column:is_purchase_location" json:"is_purchase_location"`
	Settings             map[string]interface{} `gorm:"type:jsonb;default:'{}';column:settings" json:"settings,omitempty"`
	IsActive             bool           `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedAt            time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Company     Company      `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	StorageBins []StorageBin `gorm:"foreignKey:LocationID" json:"storage_bins,omitempty"`
	Stocks      []Stock      `gorm:"foreignKey:LocationID" json:"stocks,omitempty"`
}

func (Location) TableName() string {
	return "public.locations"
}

// AuditLog tracks system changes
type AuditLog struct {
	ID         uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID  uint64         `gorm:"not null;index;column:company_id" json:"company_id"`
	UserID     *uuid.UUID     `gorm:"type:uuid;column:user_id" json:"user_id,omitempty"`
	Action     string         `gorm:"size:100;not null;column:action" json:"action"`
	TableName  string         `gorm:"size:100;not null;column:table_name" json:"table_name"`
	RecordID   string         `gorm:"size:100;column:record_id" json:"record_id,omitempty"`
	OldValues  map[string]interface{} `gorm:"type:jsonb;column:old_values" json:"old_values,omitempty"`
	NewValues  map[string]interface{} `gorm:"type:jsonb;column:new_values" json:"new_values,omitempty"`
	IPAddress  string         `gorm:"type:inet;column:ip_address" json:"ip_address,omitempty"`
	UserAgent  string         `gorm:"type:text;column:user_agent" json:"user_agent,omitempty"`
	CreatedAt  time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`

	// Relationships
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (AuditLog) TableName() string {
	return "public.audit_logs"
}
