// models/sales.go
// Sales Models - Customers, Orders, Pricing

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CustomerGroup represents customer segmentation
type CustomerGroup struct {
	ID                uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID         uint64         `gorm:"not null;index;column:company_id" json:"company_id"`
	UUID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();column:uuid" json:"uuid"`
	Code              string         `gorm:"size:100;not null;uniqueIndex:idx_cust_group_code_company;column:code" json:"code"`
	Name              string         `gorm:"size:200;not null;column:name" json:"name"`
	Description       string         `gorm:"type:text;column:description" json:"description,omitempty"`
	DiscountType      string         `gorm:"size:20;default:percentage;column:discount_type" json:"discount_type"`
	DiscountValue     float64        `gorm:"type:decimal(18,4);default:0;column:discount_value" json:"discount_value"`
	MinPurchaseAmount float64        `gorm:"type:decimal(18,4);column:min_purchase_amount" json:"min_purchase_amount,omitempty"`
	CreditLimit       float64        `gorm:"type:decimal(18,4);column:credit_limit" json:"credit_limit,omitempty"`
	PaymentTermDays   int            `gorm:"default:0;column:payment_term_days" json:"payment_term_days"`
	PriceTier         int            `gorm:"default:1;column:price_tier" json:"price_tier"`
	IsActive          bool           `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedAt         time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Company   Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Customers []Customer `gorm:"foreignKey:CustomerGroupID" json:"customers,omitempty"`
}

func (CustomerGroup) TableName() string {
	return "sales.customer_groups"
}

// Customer represents customers/clients
type Customer struct {
	ID               uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID        uint64         `gorm:"not null;index;column:company_id" json:"company_id"`
	UUID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;column:uuid" json:"uuid"`
	Code             string         `gorm:"size:100;not null;uniqueIndex:idx_customer_code_company;column:code" json:"code"`
	CustomerGroupID  *uint64        `gorm:"index;column:customer_group_id" json:"customer_group_id,omitempty"`
	Name             string         `gorm:"size:250;not null;column:name" json:"name"`
	Email            string         `gorm:"size:150;column:email" json:"email,omitempty"`
	Phone            string         `gorm:"size:50;column:phone" json:"phone,omitempty"`
	Mobile           string         `gorm:"size:50;column:mobile" json:"mobile,omitempty"`
	TaxID            string         `gorm:"size:50;column:tax_id" json:"tax_id,omitempty"`
	BillingAddress   string         `gorm:"type:text;column:billing_address" json:"billing_address,omitempty"`
	BillingCity      string         `gorm:"size:100;column:billing_city" json:"billing_city,omitempty"`
	BillingProvince  string         `gorm:"size:100;column:billing_province" json:"billing_province,omitempty"`
	BillingPostalCode string        `gorm:"size:20;column:billing_postal_code" json:"billing_postal_code,omitempty"`
	BillingCountry   string         `gorm:"size:100;default:Indonesia;column:billing_country" json:"billing_country"`
	CreditLimit      float64        `gorm:"type:decimal(18,4);default:0;column:credit_limit" json:"credit_limit"`
	CurrentBalance   float64        `gorm:"type:decimal(18,4);default:0;column:current_balance" json:"current_balance"`
	Points           int            `gorm:"default:0;column:points" json:"points"`
	PriceTier        int            `gorm:"default:1;column:price_tier" json:"price_tier"`
	IsWholesale      bool           `gorm:"default:false;column:is_wholesale" json:"is_wholesale"`
	IsActive         bool           `gorm:"default:true;not null;column:is_active" json:"is_active"`
	BirthDate        *time.Time     `gorm:"type:date;column:birth_date" json:"birth_date,omitempty"`
	Gender           string         `gorm:"size:20;column:gender" json:"gender,omitempty"`
	Notes            string         `gorm:"type:text;column:notes" json:"notes,omitempty"`
	CreatedAt        time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Company     Company             `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	CustGroup   *CustomerGroup      `gorm:"foreignKey:CustomerGroupID" json:"customer_group,omitempty"`
	Addresses   []CustomerAddress   `gorm:"foreignKey:CustomerID" json:"addresses,omitempty"`
	Orders      []Order             `gorm:"foreignKey:CustomerID" json:"orders,omitempty"`
}

func (Customer) TableName() string {
	return "sales.customers"
}

// CustomerAddress represents shipping addresses for customers
type CustomerAddress struct {
	ID           uint64     `gorm:"primaryKey;column:id" json:"id"`
	CustomerID   uint64     `gorm:"not null;index;column:customer_id" json:"customer_id"`
	Label        string     `gorm:"size:100;not null;default:Default;column:label" json:"label"`
	RecipientName string    `gorm:"size:250;column:recipient_name" json:"recipient_name,omitempty"`
	Phone        string     `gorm:"size:50;column:phone" json:"phone,omitempty"`
	Address      string     `gorm:"type:text;not null;column:address" json:"address"`
	City         string     `gorm:"size:100;column:city" json:"city,omitempty"`
	Province     string     `gorm:"size:100;column:province" json:"province,omitempty"`
	PostalCode   string     `gorm:"size:20;column:postal_code" json:"postal_code,omitempty"`
	Country      string     `gorm:"size:100;default:Indonesia;column:country" json:"country"`
	IsDefault    bool       `gorm:"default:false;column:is_default" json:"is_default"`
	Latitude     *float64   `gorm:"type:decimal(10,8);column:latitude" json:"latitude,omitempty"`
	Longitude    *float64   `gorm:"type:decimal(11,8);column:longitude" json:"longitude,omitempty"`
	CreatedAt    time.Time  `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// Relationships
	Customer Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (CustomerAddress) TableName() string {
	return "sales.customer_addresses"
}

// Order represents sales orders
type Order struct {
	ID                uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID         uint64         `gorm:"not null;index:idx_orders_company;column:company_id" json:"company_id"`
	UUID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;column:uuid" json:"uuid"`
	OrderNumber       string         `gorm:"size:100;not null;uniqueIndex:idx_order_number_company;column:order_number" json:"order_number"`
	OrderDate         time.Time      `gorm:"type:date;not null;default:CURRENT_DATE;index:idx_orders_date;column:order_date" json:"order_date"`
	DeliveryDate      *time.Time     `gorm:"type:date;column:delivery_date" json:"delivery_date,omitempty"`
	LocationID        *uint64        `gorm:"index;column:location_id" json:"location_id,omitempty"`
	CustomerID        *uint64        `gorm:"index:idx_orders_customer;column:customer_id" json:"customer_id,omitempty"`
	CustomerName      string         `gorm:"size:250;column:customer_name" json:"customer_name,omitempty"`
	CustomerAddressID *uint64        `gorm:"column:customer_address_id" json:"customer_address_id,omitempty"`
	SalesPersonID     *uuid.UUID     `gorm:"type:uuid;column:sales_person_id" json:"sales_person_id,omitempty"`
	Subtotal          float64        `gorm:"type:decimal(18,4);default:0;not null;column:subtotal" json:"subtotal"`
	DiscountAmount    float64        `gorm:"type:decimal(18,4);default:0;column:discount_amount" json:"discount_amount"`
	DiscountPercent   float64        `gorm:"type:decimal(5,2);default:0;column:discount_percent" json:"discount_percent"`
	TaxPercent        float64        `gorm:"type:decimal(5,2);default:0;column:tax_percent" json:"tax_percent"`
	TaxAmount         float64        `gorm:"type:decimal(18,4);default:0;column:tax_amount" json:"tax_amount"`
	ShippingCost      float64        `gorm:"type:decimal(18,4);default:0;column:shipping_cost" json:"shipping_cost"`
	TotalAmount       float64        `gorm:"type:decimal(18,4);default:0;not null;column:total_amount" json:"total_amount"`
	PaymentStatus     string         `gorm:"size:20;default:unpaid;column:payment_status" json:"payment_status"`
	PaidAmount        float64        `gorm:"type:decimal(18,4);default:0;column:paid_amount" json:"paid_amount"`
	Status            string         `gorm:"size:20;default:draft;index:idx_orders_status;column:status" json:"status"`
	Notes             string         `gorm:"type:text;column:notes" json:"notes,omitempty"`
	InternalNotes     string         `gorm:"type:text;column:internal_notes" json:"internal_notes,omitempty"`
	Source            string         `gorm:"size:50;default:manual;column:source" json:"source"`
	SourceReference   string         `gorm:"size:100;column:source_reference" json:"source_reference,omitempty"`
	ConfirmedAt       *time.Time     `gorm:"column:confirmed_at" json:"confirmed_at,omitempty"`
	ShippedAt         *time.Time     `gorm:"column:shipped_at" json:"shipped_at,omitempty"`
	DeliveredAt       *time.Time     `gorm:"column:delivered_at" json:"delivered_at,omitempty"`
	CancelledAt       *time.Time     `gorm:"column:cancelled_at" json:"cancelled_at,omitempty"`
	CreatedBy         *uuid.UUID     `gorm:"type:uuid;column:created_by" json:"created_by,omitempty"`
	CreatedAt         time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedBy         *uuid.UUID     `gorm:"type:uuid;column:updated_by" json:"updated_by,omitempty"`
	UpdatedAt         time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Company         Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Location        *Location       `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Customer        *Customer       `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	CustAddress     *CustomerAddress `gorm:"foreignKey:CustomerAddressID" json:"customer_address,omitempty"`
	Items           []OrderItem     `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func (Order) TableName() string {
	return "sales.orders"
}

// OrderItem represents line items in a sales order
type OrderItem struct {
	ID                uint64     `gorm:"primaryKey;column:id" json:"id"`
	OrderID           uint64     `gorm:"not null;index;column:order_id" json:"order_id"`
	VariantID         *uint64    `gorm:"column:variant_id" json:"variant_id,omitempty"`
	VariantSKU        string     `gorm:"size:100;column:variant_sku" json:"variant_sku,omitempty"`
	VariantName       string     `gorm:"size:250;column:variant_name" json:"variant_name,omitempty"`
	QuantityOrdered   float64    `gorm:"type:decimal(18,4);not null;column:quantity_ordered" json:"quantity_ordered"`
	QuantityDelivered float64    `gorm:"type:decimal(18,4);default:0;column:quantity_delivered" json:"quantity_delivered"`
	QuantityReturned  float64    `gorm:"type:decimal(18,4);default:0;column:quantity_returned" json:"quantity_returned"`
	UnitID            *uint64    `gorm:"column:unit_id" json:"unit_id,omitempty"`
	UnitPrice         float64    `gorm:"type:decimal(18,4);not null;column:unit_price" json:"unit_price"`
	DiscountPercent   float64    `gorm:"type:decimal(5,2);default:0;column:discount_percent" json:"discount_percent"`
	DiscountAmount    float64    `gorm:"type:decimal(18,4);default:0;column:discount_amount" json:"discount_amount"`
	TaxPercent        float64    `gorm:"type:decimal(5,2);default:0;column:tax_percent" json:"tax_percent"`
	TaxAmount         float64    `gorm:"type:decimal(18,4);default:0;column:tax_amount" json:"tax_amount"`
	TotalPrice        float64    `gorm:"type:decimal(18,4);not null;column:total_price" json:"total_price"`
	Cost              float64    `gorm:"type:decimal(18,4);column:cost" json:"cost,omitempty"`
	FromLocationID    *uint64    `gorm:"column:from_location_id" json:"from_location_id,omitempty"`
	PickedAt          *time.Time `gorm:"column:picked_at" json:"picked_at,omitempty"`
	PickedBy          *uuid.UUID `gorm:"type:uuid;column:picked_by" json:"picked_by,omitempty"`
	Notes             string     `gorm:"type:text;column:notes" json:"notes,omitempty"`
	CreatedAt         time.Time  `gorm:"not null;default:now();column:created_at" json:"created_at"`

	// Relationships
	Order    Order           `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Variant  *ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
	Unit     *Unit           `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	Location *Location       `gorm:"foreignKey:FromLocationID" json:"location,omitempty"`
}

func (OrderItem) TableName() string {
	return "sales.order_items"
}

// PriceTier represents multi-tier pricing
type PriceTier struct {
	ID               uint64     `gorm:"primaryKey;column:id" json:"id"`
	CompanyID        uint64     `gorm:"not null;index;column:company_id" json:"company_id"`
	VariantID        uint64     `gorm:"not null;column:variant_id" json:"variant_id"`
	CustomerGroupID  *uint64    `gorm:"column:customer_group_id" json:"customer_group_id,omitempty"`
	UnitID           *uint64    `gorm:"column:unit_id" json:"unit_id,omitempty"`
	TierLevel        int        `gorm:"not null;default:1;column:tier_level" json:"tier_level"`
	MinQuantity      float64    `gorm:"type:decimal(18,4);default:1;column:min_quantity" json:"min_quantity"`
	Price            float64    `gorm:"type:decimal(18,4);not null;column:price" json:"price"`
	EffectiveFrom    *time.Time `gorm:"type:date;default:CURRENT_DATE;column:effective_from" json:"effective_from"`
	EffectiveUntil   *time.Time `gorm:"type:date;column:effective_until" json:"effective_until,omitempty"`
	IsActive         bool       `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedAt        time.Time  `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// Relationships
	Company       Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Variant       ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
	CustGroup     *CustomerGroup `gorm:"foreignKey:CustomerGroupID" json:"customer_group,omitempty"`
	Unit          *Unit          `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
}

func (PriceTier) TableName() string {
	return "sales.price_tiers"
}

// DiscountEvent represents promotional discounts
type DiscountEvent struct {
	ID                uint64                 `gorm:"primaryKey;column:id" json:"id"`
	CompanyID         uint64                 `gorm:"not null;index;column:company_id" json:"company_id"`
	UUID              uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();column:uuid" json:"uuid"`
	Code              string                 `gorm:"size:100;not null;uniqueIndex:idx_discount_code_company;column:code" json:"code"`
	Name              string                 `gorm:"size:200;not null;column:name" json:"name"`
	Description       string                 `gorm:"type:text;column:description" json:"description,omitempty"`
	DiscountType      string                 `gorm:"size:20;not null;column:discount_type" json:"discount_type"`
	DiscountValue     float64                `gorm:"type:decimal(18,4);not null;column:discount_value" json:"discount_value"`
	MinPurchaseAmount float64                `gorm:"type:decimal(18,4);default:0;column:min_purchase_amount" json:"min_purchase_amount"`
	MaxDiscountAmount float64                `gorm:"type:decimal(18,4);column:max_discount_amount" json:"max_discount_amount,omitempty"`
	StartAt           time.Time              `gorm:"not null;column:start_at" json:"start_at"`
	EndAt             time.Time              `gorm:"not null;column:end_at" json:"end_at"`
	UsageLimit        *int                   `gorm:"column:usage_limit" json:"usage_limit,omitempty"`
	UsageCount        int                    `gorm:"default:0;column:usage_count" json:"usage_count"`
	AppliesTo         string                 `gorm:"size:50;default:all;column:applies_to" json:"applies_to"`
	Conditions        map[string]interface{} `gorm:"type:jsonb;column:conditions" json:"conditions,omitempty"`
	IsActive          bool                   `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedBy         *uuid.UUID             `gorm:"type:uuid;column:created_by" json:"created_by,omitempty"`
	CreatedAt         time.Time              `gorm:"not null;default:now();column:created_at" json:"created_at"`

	// Relationships
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

func (DiscountEvent) TableName() string {
	return "sales.discount_events"
}

// DailySalesView represents the daily sales summary view
type DailySalesView struct {
	CompanyID      uint64    `json:"company_id"`
	SalesDate      time.Time `json:"sales_date"`
	OrderCount     int64     `json:"order_count"`
	TotalRevenue   float64   `json:"total_revenue"`
	TotalSubtotal  float64   `json:"total_subtotal"`
	TotalDiscounts float64   `json:"total_discounts"`
	TotalTax       float64   `json:"total_tax"`
	TotalShipping  float64   `json:"total_shipping"`
}

func (DailySalesView) TableName() string {
	return "sales.v_daily_sales"
}

// ProductPerformanceView represents product sales performance
type ProductPerformanceView struct {
	VariantID   uint64  `json:"variant_id"`
	SKU         string  `json:"sku"`
	ProductName string  `json:"product_name"`
	CategoryName string `json:"category_name,omitempty"`
	TotalQtySold float64 `json:"total_qty_sold"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalCost    float64 `json:"total_cost"`
	GrossProfit  float64 `json:"gross_profit"`
}

func (ProductPerformanceView) TableName() string {
	return "sales.v_product_performance"
}
