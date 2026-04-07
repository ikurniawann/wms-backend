// models/inventory.go
// Inventory Models - Stocks, Movements, Storage

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Stock represents stock levels per variant per location
type Stock struct {
	ID                uint64     `gorm:"primaryKey;column:id" json:"id"`
	CompanyID         uint64     `gorm:"not null;index:idx_stocks_company;column:company_id" json:"company_id"`
	VariantID         uint64     `gorm:"not null;uniqueIndex:idx_stocks_variant_location;column:variant_id" json:"variant_id"`
	LocationID        uint64     `gorm:"not null;uniqueIndex:idx_stocks_variant_location;index:idx_stocks_location;column:location_id" json:"location_id"`
	QuantityAvailable float64    `gorm:"type:decimal(18,4);default:0;not null;column:quantity_available" json:"quantity_available"`
	QuantityReserved  float64    `gorm:"type:decimal(18,4);default:0;not null;column:quantity_reserved" json:"quantity_reserved"`
	QuantityOnHand    float64    `gorm:"type:decimal(18,4);column:quantity_on_hand" json:"quantity_on_hand"` // Generated
	QuantityIncoming  float64    `gorm:"type:decimal(18,4);default:0;column:quantity_incoming" json:"quantity_incoming"`
	QuantityOutgoing  float64    `gorm:"type:decimal(18,4);default:0;column:quantity_outgoing" json:"quantity_outgoing"`
	ReorderPoint      float64    `gorm:"type:decimal(18,4);default:0;column:reorder_point" json:"reorder_point"`
	ReorderQuantity   float64    `gorm:"type:decimal(18,4);column:reorder_quantity" json:"reorder_quantity,omitempty"`
	MaxStock          float64    `gorm:"type:decimal(18,4);column:max_stock" json:"max_stock,omitempty"`
	AvgCost           float64    `gorm:"type:decimal(18,4);default:0;column:avg_cost" json:"avg_cost"`
	LastCountedAt     *time.Time `gorm:"column:last_counted_at" json:"last_counted_at,omitempty"`
	LastMovementAt    *time.Time `gorm:"column:last_movement_at" json:"last_movement_at,omitempty"`
	CreatedAt         time.Time  `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// Relationships
	Company    Company      `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Variant    ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
	Location   Location     `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	StockBins  []StockBin   `gorm:"foreignKey:StockID" json:"stock_bins,omitempty"`
	Reservations []Reservation `gorm:"foreignKey:StockID" json:"reservations,omitempty"`
}

func (Stock) TableName() string {
	return "inventory.stocks"
}

// StorageBin represents rack/bin locations within a warehouse
type StorageBin struct {
	ID             uint64         `gorm:"primaryKey;column:id" json:"id"`
	CompanyID      uint64         `gorm:"not null;index;column:company_id" json:"company_id"`
	LocationID     uint64         `gorm:"not null;index;column:location_id" json:"location_id"`
	Code           string         `gorm:"size:50;not null;uniqueIndex:idx_bin_code_location;column:code" json:"code"`
	Name           string         `gorm:"size:100;column:name" json:"name,omitempty"`
	Zone           string         `gorm:"size:50;column:zone" json:"zone,omitempty"` // receiving, picking, storage, shipping, quarantine
	Aisle          string         `gorm:"size:50;column:aisle" json:"aisle,omitempty"`
	Rack           string         `gorm:"size:50;column:rack" json:"rack,omitempty"`
	Shelf          string         `gorm:"size:50;column:shelf" json:"shelf,omitempty"`
	BinType        string         `gorm:"size:50;default:standard;column:bin_type" json:"bin_type"`
	Capacity       float64        `gorm:"type:decimal(18,4);column:capacity" json:"capacity,omitempty"`
	CapacityUnitID *uint64        `gorm:"column:capacity_unit_id" json:"capacity_unit_id,omitempty"`
	IsActive       bool           `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedAt      time.Time      `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`

	// Relationships
	Company       Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Location      Location       `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	CapacityUnit  *Unit          `gorm:"foreignKey:CapacityUnitID" json:"capacity_unit,omitempty"`
	StockBins     []StockBin     `gorm:"foreignKey:BinID" json:"stock_bins,omitempty"`
}

func (StorageBin) TableName() string {
	return "inventory.storage_bins"
}

// StockBin represents stock stored in specific bins
type StockBin struct {
	ID          uint64     `gorm:"primaryKey;column:id" json:"id"`
	StockID     uint64     `gorm:"not null;index;column:stock_id" json:"stock_id"`
	BinID       uint64     `gorm:"not null;index;column:bin_id" json:"bin_id"`
	Quantity    float64    `gorm:"type:decimal(18,4);default:0;not null;column:quantity" json:"quantity"`
	LotNumber   string     `gorm:"size:100;column:lot_number" json:"lot_number,omitempty"`
	ExpiryDate  *time.Time `gorm:"type:date;column:expiry_date" json:"expiry_date,omitempty"`
	ReceivedAt  *time.Time `gorm:"column:received_at" json:"received_at,omitempty"`
	CreatedAt   time.Time  `gorm:"not null;default:now();column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// Relationships
	Stock Stock       `gorm:"foreignKey:StockID" json:"stock,omitempty"`
	Bin   StorageBin  `gorm:"foreignKey:BinID" json:"bin,omitempty"`
}

func (StockBin) TableName() string {
	return "inventory.stock_bins"
}

// Movement represents stock movements (in, out, transfer, adjustment)
type Movement struct {
	ID              uint64     `gorm:"primaryKey;column:id" json:"id"`
	CompanyID       uint64     `gorm:"not null;index:idx_movements_company;column:company_id" json:"company_id"`
	UUID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex;column:uuid" json:"uuid"`
	MovementNumber  string     `gorm:"size:100;not null;uniqueIndex:idx_movement_number_company;column:movement_number" json:"movement_number"`
	MovementDate    time.Time  `gorm:"type:date;not null;default:CURRENT_DATE;index:idx_movements_date;column:movement_date" json:"movement_date"`
	MovementType    string     `gorm:"size:20;not null;column:movement_type" json:"movement_type"` // in, out, transfer, adjustment, count
	Reason          string     `gorm:"size:20;not null;column:reason" json:"reason"` // purchase, sale, return, transfer, adjustment, damage, expired, production, count
	ReferenceType   string     `gorm:"size:50;column:reference_type" json:"reference_type,omitempty"`
	ReferenceID     *uint64    `gorm:"column:reference_id" json:"reference_id,omitempty"`
	ReferenceNumber string     `gorm:"size:100;column:reference_number" json:"reference_number,omitempty"`
	LocationID      uint64     `gorm:"not null;column:location_id" json:"location_id"`
	FromLocationID  *uint64    `gorm:"column:from_location_id" json:"from_location_id,omitempty"`
	ToLocationID    *uint64    `gorm:"column:to_location_id" json:"to_location_id,omitempty"`
	TotalItems      int        `gorm:"default:0;column:total_items" json:"total_items"`
	TotalQuantity   float64    `gorm:"type:decimal(18,4);default:0;column:total_quantity" json:"total_quantity"`
	Notes           string     `gorm:"type:text;column:notes" json:"notes,omitempty"`
	CreatedBy       *uuid.UUID `gorm:"type:uuid;column:created_by" json:"created_by,omitempty"`
	CreatedAt       time.Time  `gorm:"not null;default:now();column:created_at" json:"created_at"`
	PostedAt        *time.Time `gorm:"column:posted_at" json:"posted_at,omitempty"`

	// Relationships
	Company        Company           `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Location       Location          `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	FromLocation   *Location         `gorm:"foreignKey:FromLocationID" json:"from_location,omitempty"`
	ToLocation     *Location         `gorm:"foreignKey:ToLocationID" json:"to_location,omitempty"`
	Details        []MovementDetail  `gorm:"foreignKey:MovementID" json:"details,omitempty"`
}

func (Movement) TableName() string {
	return "inventory.movements"
}

// MovementDetail represents line items in a movement
type MovementDetail struct {
	ID          uint64     `gorm:"primaryKey;column:id" json:"id"`
	MovementID  uint64     `gorm:"not null;index;column:movement_id" json:"movement_id"`
	VariantID   uint64     `gorm:"not null;column:variant_id" json:"variant_id"`
	FromBinID   *uint64    `gorm:"column:from_bin_id" json:"from_bin_id,omitempty"`
	ToBinID     *uint64    `gorm:"column:to_bin_id" json:"to_bin_id,omitempty"`
	LotNumber   string     `gorm:"size:100;column:lot_number" json:"lot_number,omitempty"`
	Quantity    float64    `gorm:"type:decimal(18,4);not null;column:quantity" json:"quantity"`
	UnitCost    float64    `gorm:"type:decimal(18,4);column:unit_cost" json:"unit_cost,omitempty"`
	TotalCost   float64    `gorm:"type:decimal(18,4);column:total_cost" json:"total_cost,omitempty"`
	ExpiryDate  *time.Time `gorm:"type:date;column:expiry_date" json:"expiry_date,omitempty"`
	Notes       string     `gorm:"type:text;column:notes" json:"notes,omitempty"`
	CreatedAt   time.Time  `gorm:"not null;default:now();column:created_at" json:"created_at"`

	// Relationships
	Movement Movement      `gorm:"foreignKey:MovementID" json:"movement,omitempty"`
	Variant  ProductVariant `gorm:"foreignKey:VariantID" json:"variant,omitempty"`
	FromBin  *StorageBin   `gorm:"foreignKey:FromBinID" json:"from_bin,omitempty"`
	ToBin    *StorageBin   `gorm:"foreignKey:ToBinID" json:"to_bin,omitempty"`
}

func (MovementDetail) TableName() string {
	return "inventory.movement_details"
}

// Reservation represents stock reservations for orders
type Reservation struct {
	ID               uint64     `gorm:"primaryKey;column:id" json:"id"`
	CompanyID        uint64     `gorm:"not null;index;column:company_id" json:"company_id"`
	StockID          uint64     `gorm:"not null;index;column:stock_id" json:"stock_id"`
	ReferenceType    string     `gorm:"size:50;not null;column:reference_type" json:"reference_type"`
	ReferenceID      uint64     `gorm:"not null;column:reference_id" json:"reference_id"`
	QuantityReserved float64    `gorm:"type:decimal(18,4);not null;column:quantity_reserved" json:"quantity_reserved"`
	QuantityReleased float64    `gorm:"type:decimal(18,4);default:0;column:quantity_released" json:"quantity_released"`
	ExpiryDate       *time.Time `gorm:"type:date;column:expiry_date" json:"expiry_date,omitempty"`
	IsActive         bool       `gorm:"default:true;not null;column:is_active" json:"is_active"`
	CreatedAt        time.Time  `gorm:"not null;default:now();column:created_at" json:"created_at"`
	ReleasedAt       *time.Time `gorm:"column:released_at" json:"released_at,omitempty"`

	// Relationships
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Stock   Stock   `gorm:"foreignKey:StockID" json:"stock,omitempty"`
}

func (Reservation) TableName() string {
	return "inventory.reservations"
}

// StockSummaryView represents the stock summary view
type StockSummaryView struct {
	ID                uint64  `json:"id"`
	CompanyID         uint64  `json:"company_id"`
	VariantID         uint64  `json:"variant_id"`
	SKU               string  `json:"sku"`
	VariantName       string  `json:"variant_name"`
	ProductID         uint64  `json:"product_id"`
	ProductName       string  `json:"product_name"`
	CategoryName      string  `json:"category_name"`
	LocationID        uint64  `json:"location_id"`
	LocationName      string  `json:"location_name"`
	QuantityAvailable float64 `json:"quantity_available"`
	QuantityReserved  float64 `json:"quantity_reserved"`
	QuantityOnHand    float64 `json:"quantity_on_hand"`
	QuantityIncoming  float64 `json:"quantity_incoming"`
	QuantityOutgoing  float64 `json:"quantity_outgoing"`
	ReorderPoint      float64 `json:"reorder_point"`
	AvgCost           float64 `json:"avg_cost"`
	StockValue        float64 `json:"stock_value"`
	IsLowStock        bool    `json:"is_low_stock"`
}

func (StockSummaryView) TableName() string {
	return "inventory.v_stock_summary"
}
