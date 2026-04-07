// models/master_data.go
// Master Data Models - Products, Categories, Units, Attributes

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Unit represents units of measurement
type Unit struct {
	ID               uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID        uint64         `gorm:"not null;index;column:company_id" json:"company_id"`
	UUID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();column:uuid" json:"uuid"`
	Code             string         `gorm:"size:50;not null;uniqueIndex:idx_unit_code_company;column:code" json:"code"`
	Name             string         `gorm:"size:100;not null;column:name" json:"name"`
	Symbol           string         `gorm:"size:20;column:symbol" json:"symbol,omitempty"`
	UnitType         string         `gorm:"size:50;default:piece;column:unit_type" json:"unit_type"`
	ConversionFactor float64        `gorm:"type:decimal(18,6);default:1;column:conversion_factor" json:"conversion_factor"`
	BaseUnitID       *uint64        `gorm:"column:base_unit_id" json:"base_unit_id,omitempty"`
	IsBaseUnit       bool           `gorm:"default:false;column:is_base_unit" json:"is_base_unit"`
	Description      string         `gorm:"type:text;column:description" json:"description,omitempty"`
	IsActive         bool           `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedAt        time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Company  Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	BaseUnit *Unit   `gorm:"foreignKey:BaseUnitID" json:"base_unit,omitempty"`
}

func (Unit) TableName() string {
	return "master_data.units"
}

// Category represents product categories (hierarchical)
type Category struct {
	ID          uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID   uint64         `gorm:"not null;index;column:company_id" json:"company_id"`
	UUID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();column:uuid" json:"uuid"`
	Code        string         `gorm:"size:100;not null;uniqueIndex:idx_category_code_company;column:code" json:"code"`
	Name        string         `gorm:"size:200;not null;column:name" json:"name"`
	Description string         `gorm:"type:text;column:description" json:"description,omitempty"`
	ParentID    *uint64        `gorm:"index;column:parent_id" json:"parent_id,omitempty"`
	Level       int            `gorm:"default:0;column:level" json:"level"`
	Path        string         `gorm:"type:ltree;index:,type:gin;column:path" json:"path,omitempty"`
	SortOrder   int            `gorm:"default:0;column:sort_order" json:"sort_order"`
	ImageURL    string         `gorm:"size:500;column:image_url" json:"image_url,omitempty"`
	IsActive    bool           `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedAt   time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Company  Company     `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Parent   *Category   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Products []Product   `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}

func (Category) TableName() string {
	return "master_data.categories"
}

// Attribute represents product attributes for variants
type Attribute struct {
	ID                 uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID          uint64         `gorm:"not null;index;column:company_id" json:"company_id"`
	UUID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();column:uuid" json:"uuid"`
	Code               string         `gorm:"size:100;not null;uniqueIndex:idx_attr_code_company;column:code" json:"code"`
	Name               string         `gorm:"size:200;not null;column:name" json:"name"`
	Description        string         `gorm:"type:text;column:description" json:"description,omitempty"`
	DataType           string         `gorm:"size:50;default:string;column:data_type" json:"data_type"`
	IsRequired         bool           `gorm:"default:false;column:is_required" json:"is_required"`
	IsVariantAttribute bool           `gorm:"default:false;column:is_variant_attribute" json:"is_variant_attribute"`
	IsFilterable       bool           `gorm:"default:false;column:is_filterable" json:"is_filterable"`
	Options            map[string]interface{} `gorm:"type:jsonb;column:options" json:"options,omitempty"` // For select type: ["Red", "Blue", "Green"]
	SortOrder          int            `gorm:"default:0;column:sort_order" json:"sort_order"`
	IsActive           bool           `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedAt          time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (Attribute) TableName() string {
	return "master_data.attributes"
}

// Product represents the main product entity
type Product struct {
	ID               uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID        uint64         `gorm:"not null;index:idx_products_company;column:company_id" json:"company_id"`
	UUID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;column:uuid" json:"uuid"`
	Code             string         `gorm:"size:100;not null;uniqueIndex:idx_products_code_company;column:code" json:"code"`
	SKU              string         `gorm:"size:100;column:sku" json:"sku,omitempty"`
	Barcode          string         `gorm:"size:100;column:barcode" json:"barcode,omitempty"`
	Name             string         `gorm:"size:250;not null;column:name" json:"name"`
	Description      string         `gorm:"type:text;column:description" json:"description,omitempty"`
	ShortDescription string         `gorm:"size:500;column:short_description" json:"short_description,omitempty"`
	CategoryID       *uint64        `gorm:"index:idx_products_category;column:category_id" json:"category_id,omitempty"`
	Brand            string         `gorm:"size:100;column:brand" json:"brand,omitempty"`
	BaseUnitID       *uint64        `gorm:"column:base_unit_id" json:"base_unit_id,omitempty"`
	Weight           float64        `gorm:"type:decimal(18,4);column:weight" json:"weight,omitempty"`
	Dimensions       map[string]interface{} `gorm:"type:jsonb;column:dimensions" json:"dimensions,omitempty"` // {"length": 10, "width": 5, "height": 3}
	TrackInventory   bool           `gorm:"default:true;column:track_inventory" json:"track_inventory"`
	IsService        bool           `gorm:"default:false;column:is_service" json:"is_service"`
	IsActive         bool           `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedBy        *uuid.UUID     `gorm:"type:uuid;column:created_by" json:"created_by,omitempty"`
	CreatedAt        time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedBy        *uuid.UUID     `gorm:"type:uuid;column:updated_by" json:"updated_by,omitempty"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`
	Metadata         map[string]interface{} `gorm:"type:jsonb;default:'{}';column:metadata" json:"metadata,omitempty"`

	// Relationships
	Company   Company          `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Category  *Category        `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	BaseUnit  *Unit            `gorm:"foreignKey:BaseUnitID" json:"base_unit,omitempty"`
	Variants  []ProductVariant `gorm:"foreignKey:ProductID" json:"variants,omitempty"`
	Images    []ProductImage   `gorm:"foreignKey:ProductID" json:"images,omitempty"`
}

func (Product) TableName() string {
	return "master_data.products"
}

// ProductVariant represents SKU-level variants
type ProductVariant struct {
	ID          uint64         `gorm:"primaryKey;column:id" json:"id"`
	ProductID   uint64         `gorm:"not null;index:idx_variants_product;column:product_id" json:"product_id"`
	UUID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;column:uuid" json:"uuid"`
	SKU         string         `gorm:"size:100;not null;uniqueIndex:idx_variants_sku_product;column:sku" json:"sku"`
	Barcode     string         `gorm:"size:100;index:idx_variants_barcode;column:barcode" json:"barcode,omitempty"`
	Name        string         `gorm:"size:250;column:name" json:"name,omitempty"`
	Options     map[string]interface{} `gorm:"type:jsonb;column:options" json:"options,omitempty"` // {"color": "Red", "size": "Large"}
	Weight      float64        `gorm:"type:decimal(18,4);column:weight" json:"weight,omitempty"`
	Dimensions  map[string]interface{} `gorm:"type:jsonb;column:dimensions" json:"dimensions,omitempty"`
	IsDefault   bool           `gorm:"default:false;column:is_default" json:"is_default"`
	IsActive    bool           `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedAt   time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Product   Product          `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Images    []ProductImage   `gorm:"foreignKey:VariantID" json:"images,omitempty"`
	Stocks    []Stock          `gorm:"foreignKey:VariantID" json:"stocks,omitempty"`
	PriceTiers []PriceTier     `gorm:"foreignKey:VariantID" json:"price_tiers,omitempty"`
}

func (ProductVariant) TableName() string {
	return "master_data.product_variants"
}

// ProductImage stores references to Supabase Storage
type ProductImage struct {
	ID           uint64    `gorm:"primaryKey;column:id" json:"id"`
	ProductID    uint64    `gorm:"not null;index;column:product_id" json:"product_id"`
	VariantID    *uint64   `gorm:"index;column:variant_id" json:"variant_id,omitempty"`
	StoragePath  string    `gorm:"size:500;not null;column:storage_path" json:"storage_path"`
	FileName     string    `gorm:"size:250;column:file_name" json:"file_name,omitempty"`
	FileSize     int       `gorm:"column:file_size" json:"file_size,omitempty"`
	MimeType     string    `gorm:"size:50;column:mime_type" json:"mime_type,omitempty"`
	IsPrimary    bool      `gorm:"default:false;column:is_primary" json:"is_primary"`
	SortOrder    int       `gorm:"default:0;column:sort_order" json:"sort_order"`
	CreatedAt    time.Time `gorm:"not null;default:now();column:created_at" json:"created_at"`

	// Relationships
	Product  Product           `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Variant  *ProductVariant   `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
}

func (ProductImage) TableName() string {
	return "master_data.product_images"
}
