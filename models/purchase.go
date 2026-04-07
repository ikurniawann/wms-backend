// models/purchase.go
// Purchase Models - Suppliers, Purchase Orders, Receipts

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Supplier represents vendors/suppliers
type Supplier struct {
	ID              uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID       uint64         `gorm:"not null;index;column:company_id" json:"company_id"`
	UUID            uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;column:uuid" json:"uuid"`
	Code            string         `gorm:"size:100;not null;uniqueIndex:idx_supplier_code_company;column:code" json:"code"`
	Name            string         `gorm:"size:250;not null;column:name" json:"name"`
	Email           string         `gorm:"size:150;column:email" json:"email,omitempty"`
	Phone           string         `gorm:"size:50;column:phone" json:"phone,omitempty"`
	Mobile          string         `gorm:"size:50;column:mobile" json:"mobile,omitempty"`
	TaxID           string         `gorm:"size:50;column:tax_id" json:"tax_id,omitempty"`
	Address         string         `gorm:"type:text;column:address" json:"address,omitempty"`
	City            string         `gorm:"size:100;column:city" json:"city,omitempty"`
	Province        string         `gorm:"size:100;column:province" json:"province,omitempty"`
	PostalCode      string         `gorm:"size:20;column:postal_code" json:"postal_code,omitempty"`
	Country         string         `gorm:"size:100;default:Indonesia;column:country" json:"country"`
	Website         string         `gorm:"size:150;column:website" json:"website,omitempty"`
	PaymentTermsDays int           `gorm:"default:0;column:payment_terms_days" json:"payment_terms_days"`
	CreditLimit     float64        `gorm:"type:decimal(18,4);default:0;column:credit_limit" json:"credit_limit"`
	CurrentBalance  float64        `gorm:"type:decimal(18,4);default:0;column:current_balance" json:"current_balance"`
	IsActive        bool           `gorm:"default:true;not null;column:is_active" json:"is_active"`
	Notes           string         `gorm:"type:text;column:notes" json:"notes,omitempty"`
	CreatedAt       time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Company       Company          `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Prices        []SupplierPrice  `gorm:"foreignKey:SupplierID" json:"prices,omitempty"`
	PurchaseOrders []PurchaseOrder `gorm:"foreignKey:SupplierID" json:"purchase_orders,omitempty"`
	Receipts      []Receipt        `gorm:"foreignKey:SupplierID" json:"receipts,omitempty"`
}

func (Supplier) TableName() string {
	return "purchase.suppliers"
}

// SupplierPrice represents supplier-specific product pricing
type SupplierPrice struct {
	ID             uint64     `gorm:"primaryKey;column:id" json:"id"`
	CompanyID      uint64     `gorm:"not null;index;column:company_id" json:"company_id"`
	SupplierID     uint64     `gorm:"not null;column:supplier_id" json:"supplier_id"`
	VariantID      uint64     `gorm:"not null;column:variant_id" json:"variant_id"`
	UnitID         *uint64    `gorm:"column:unit_id" json:"unit_id,omitempty"`
	SupplierSKU    string     `gorm:"size:100;column:supplier_sku" json:"supplier_sku,omitempty"`
	MinOrderQty    float64    `gorm:"type:decimal(18,4);default:1;column:min_order_qty" json:"min_order_qty"`
	Price          float64    `gorm:"type:decimal(18,4);not null;column:price" json:"price"`
	LeadTimeDays   int        `gorm:"default:0;column:lead_time_days" json:"lead_time_days"`
	EffectiveFrom  *time.Time `gorm:"type:date;default:CURRENT_DATE;column:effective_from" json:"effective_from"`
	EffectiveUntil *time.Time `gorm:"type:date;column:effective_until" json:"effective_until,omitempty"`
	IsPreferred    bool       `gorm:"default:false;column:is_preferred" json:"is_preferred"`
	IsActive       bool       `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedAt      time.Time  `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// Relationships
	Company  Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Supplier Supplier       `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Variant  ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
	Unit     *Unit          `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
}

func (SupplierPrice) TableName() string {
	return "purchase.supplier_prices"
}

// PurchaseOrder represents purchase orders
type PurchaseOrder struct {
	ID                   uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID            uint64         `gorm:"not null;index;column:company_id" json:"company_id"`
	UUID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;column:uuid" json:"uuid"`
	PONumber             string         `gorm:"size:100;not null;uniqueIndex:idx_po_number_company;column:po_number" json:"po_number"`
	PODate               time.Time      `gorm:"type:date;not null;default:CURRENT_DATE;column:po_date" json:"po_date"`
	ExpectedDeliveryDate *time.Time     `gorm:"type:date;column:expected_delivery_date" json:"expected_delivery_date,omitempty"`
	LocationID           *uint64        `gorm:"index;column:location_id" json:"location_id,omitempty"`
	SupplierID           *uint64        `gorm:"index;column:supplier_id" json:"supplier_id,omitempty"`
	SupplierName         string         `gorm:"size:250;column:supplier_name" json:"supplier_name,omitempty"`
	Subtotal             float64        `gorm:"type:decimal(18,4);default:0;not null;column:subtotal" json:"subtotal"`
	DiscountAmount       float64        `gorm:"type:decimal(18,4);default:0;column:discount_amount" json:"discount_amount"`
	TaxPercent           float64        `gorm:"type:decimal(5,2);default:0;column:tax_percent" json:"tax_percent"`
	TaxAmount            float64        `gorm:"type:decimal(18,4);default:0;column:tax_amount" json:"tax_amount"`
	ShippingCost         float64        `gorm:"type:decimal(18,4);default:0;column:shipping_cost" json:"shipping_cost"`
	TotalAmount          float64        `gorm:"type:decimal(18,4);default:0;not null;column:total_amount" json:"total_amount"`
	Status               string         `gorm:"size:20;default:draft;column:status" json:"status"`
	Notes                string         `gorm:"type:text;column:notes" json:"notes,omitempty"`
	InternalNotes        string         `gorm:"type:text;column:internal_notes" json:"internal_notes,omitempty"`
	CreatedBy            *uuid.UUID     `gorm:"type:uuid;column:created_by" json:"created_by,omitempty"`
	CreatedAt            time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at,omitempty"`

	// Relationships
	Company      Company           `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Location     *Location         `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Supplier     *Supplier         `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Items        []PurchaseOrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Receipts     []Receipt         `gorm:"foreignKey:POID" json:"receipts,omitempty"`
}

func (PurchaseOrder) TableName() string {
	return "purchase.orders"
}

// PurchaseOrderItem represents line items in a purchase order
type PurchaseOrderItem struct {
	ID                uint64     `gorm:"primaryKey;column:id" json:"id"`
	OrderID           uint64     `gorm:"not null;index;column:order_id" json:"order_id"`
	VariantID         *uint64    `gorm:"column:variant_id" json:"variant_id,omitempty"`
	VariantSKU        string     `gorm:"size:100;column:variant_sku" json:"variant_sku,omitempty"`
	VariantName       string     `gorm:"size:250;column:variant_name" json:"variant_name,omitempty"`
	QuantityOrdered   float64    `gorm:"type:decimal(18,4);not null;column:quantity_ordered" json:"quantity_ordered"`
	QuantityReceived  float64    `gorm:"type:decimal(18,4);default:0;column:quantity_received" json:"quantity_received"`
	QuantityCancelled float64    `gorm:"type:decimal(18,4);default:0;column:quantity_cancelled" json:"quantity_cancelled"`
	UnitID            *uint64    `gorm:"column:unit_id" json:"unit_id,omitempty"`
	UnitPrice         float64    `gorm:"type:decimal(18,4);not null;column:unit_price" json:"unit_price"`
	DiscountPercent   float64    `gorm:"type:decimal(5,2);default:0;column:discount_percent" json:"discount_percent"`
	TotalPrice        float64    `gorm:"type:decimal(18,4);not null;column:total_price" json:"total_price"`
	Notes             string     `gorm:"type:text;column:notes" json:"notes,omitempty"`
	CreatedAt         time.Time  `gorm:"not null;default:now();column:created_at" json:"created_at"`

	// Relationships
	Order    PurchaseOrder   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Variant  *ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
	Unit     *Unit           `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	Receipts []ReceiptItem   `gorm:"foreignKey:POItemID" json:"receipts,omitempty"`
}

func (PurchaseOrderItem) TableName() string {
	return "purchase.order_items"
}

// Receipt represents goods receipts from suppliers
type Receipt struct {
	ID            uint64     `gorm:"primaryKey;column:id" json:"id"`
	CompanyID     uint64     `gorm:"not null;index;column:company_id" json:"company_id"`
	UUID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;column:uuid" json:"uuid"`
	ReceiptNumber string     `gorm:"size:100;not null;uniqueIndex:idx_receipt_number_company;column:receipt_number" json:"receipt_number"`
	ReceiptDate   time.Time  `gorm:"type:date;not null;default:CURRENT_DATE;column:receipt_date" json:"receipt_date"`
	POID          *uint64    `gorm:"index;column:po_id" json:"po_id,omitempty"`
	SupplierID    *uint64    `gorm:"index;column:supplier_id" json:"supplier_id,omitempty"`
	LocationID    *uint64    `gorm:"index;column:location_id" json:"location_id,omitempty"`
	TotalItems    int        `gorm:"default:0;column:total_items" json:"total_items"`
	TotalQuantity float64    `gorm:"type:decimal(18,4);default:0;column:total_quantity" json:"total_quantity"`
	Notes         string     `gorm:"type:text;column:notes" json:"notes,omitempty"`
	CreatedBy     *uuid.UUID `gorm:"type:uuid;column:created_by" json:"created_by,omitempty"`
	CreatedAt     time.Time  `gorm:"not null;default:now();column:created_at" json:"created_at"`

	// Relationships
	Company  Company       `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	PO       *PurchaseOrder `gorm:"foreignKey:POID" json:"po,omitempty"`
	Supplier *Supplier      `gorm:"foreignKey:SupplierID" json:"supplier,omitempty"`
	Location *Location      `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Items    []ReceiptItem  `gorm:"foreignKey:ReceiptID" json:"items,omitempty"`
}

func (Receipt) TableName() string {
	return "purchase.receipts"
}

// ReceiptItem represents line items in a receipt
type ReceiptItem struct {
	ID               uint64     `gorm:"primaryKey;column:id" json:"id"`
	ReceiptID        uint64     `gorm:"not null;index;column:receipt_id" json:"receipt_id"`
	POItemID         *uint64    `gorm:"index;column:po_item_id" json:"po_item_id,omitempty"`
	VariantID        *uint64    `gorm:"column:variant_id" json:"variant_id,omitempty"`
	QuantityReceived float64    `gorm:"type:decimal(18,4);not null;column:quantity_received" json:"quantity_received"`
	UnitID           *uint64    `gorm:"column:unit_id" json:"unit_id,omitempty"`
	UnitCost         float64    `gorm:"type:decimal(18,4);not null;column:unit_cost" json:"unit_cost"`
	BinID            *uint64    `gorm:"index;column:bin_id" json:"bin_id,omitempty"`
	LotNumber        string     `gorm:"size:100;column:lot_number" json:"lot_number,omitempty"`
	ExpiryDate       *time.Time `gorm:"type:date;column:expiry_date" json:"expiry_date,omitempty"`
	Notes            string     `gorm:"type:text;column:notes" json:"notes,omitempty"`
	CreatedAt        time.Time  `gorm:"not null;default:now();column:created_at" json:"created_at"`

	// Relationships
	Receipt  Receipt         `gorm:"foreignKey:ReceiptID" json:"receipt,omitempty"`
	POItem   *PurchaseOrderItem `gorm:"foreignKey:POItemID" json:"po_item,omitempty"`
	Variant  *ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
	Unit     *Unit            `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	Bin      *StorageBin      `gorm:"foreignKey:BinID" json:"bin,omitempty"`
}

func (ReceiptItem) TableName() string {
	return "purchase.receipt_items"
}
