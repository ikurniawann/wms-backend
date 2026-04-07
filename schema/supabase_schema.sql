-- Supabase WMS Schema - Best Practice & Scalable
-- Compatible with: PostgreSQL 15+ / Supabase
-- Date: 2026-04-07
-- Architecture: Multi-schema with RLS

-- ============================================
-- ENABLE EXTENSIONS
-- ============================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================
-- CREATE SCHEMAS
-- ============================================
CREATE SCHEMA IF NOT EXISTS system;
CREATE SCHEMA IF NOT EXISTS master_data;
CREATE SCHEMA IF NOT EXISTS inventory;
CREATE SCHEMA IF NOT EXISTS sales;
CREATE SCHEMA IF NOT EXISTS purchase;
CREATE SCHEMA IF NOT EXISTS accounting;

-- Set search path (optional, for convenience)
-- ALTER DATABASE your_db SET search_path TO public, system, master_data, inventory, sales, purchase, accounting;

-- ============================================
-- PUBLIC SCHEMA - SYSTEM TABLES
-- ============================================

-- Companies (Multi-tenant root)
CREATE TABLE public.companies (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    legal_name VARCHAR(200),
    tax_id VARCHAR(50),
    email VARCHAR(150),
    phone VARCHAR(50),
    website VARCHAR(150),
    logo_url VARCHAR(500),
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) DEFAULT 'Indonesia',
    timezone VARCHAR(50) DEFAULT 'Asia/Jakarta',
    currency VARCHAR(3) DEFAULT 'IDR',
    subscription_tier VARCHAR(50) DEFAULT 'basic',
    subscription_expires_at TIMESTAMPTZ,
    settings JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Company Users (link to Supabase Auth)
CREATE TABLE public.company_users (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    email VARCHAR(150) NOT NULL,
    full_name VARCHAR(200),
    avatar_url VARCHAR(500),
    role VARCHAR(50) DEFAULT 'member', -- owner, admin, manager, member
    department VARCHAR(100),
    phone VARCHAR(50),
    is_active BOOLEAN DEFAULT true NOT NULL,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, user_id)
);

-- Locations / Warehouses / Stores
CREATE TABLE public.locations (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    location_type VARCHAR(50) DEFAULT 'warehouse', -- warehouse, store, virtual, supplier
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(20),
    phone VARCHAR(50),
    email VARCHAR(150),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    is_main_location BOOLEAN DEFAULT false,
    is_sales_location BOOLEAN DEFAULT false, -- Can sell from here
    is_purchase_location BOOLEAN DEFAULT true, -- Can receive PO here
    settings JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(company_id, code)
);

-- Location Types Enum
CREATE TYPE public.location_type AS ENUM ('warehouse', 'store', 'virtual', 'supplier', 'customer');

-- Add constraint (optional, if want strict enum)
-- ALTER TABLE public.locations ADD CONSTRAINT chk_location_type 
--     CHECK (location_type IN ('warehouse', 'store', 'virtual', 'supplier', 'customer'));

-- Audit Logs (system-wide)
CREATE TABLE public.audit_logs (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    user_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL, -- CREATE, UPDATE, DELETE, LOGIN, etc
    table_name VARCHAR(100) NOT NULL,
    record_id VARCHAR(100),
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Create index on audit_logs
CREATE INDEX idx_audit_logs_company ON public.audit_logs(company_id);
CREATE INDEX idx_audit_logs_created ON public.audit_logs(created_at);
CREATE INDEX idx_audit_logs_action ON public.audit_logs(action);

-- ============================================
-- MASTER_DATA SCHEMA
-- ============================================

-- Units of Measurement
CREATE TABLE master_data.units (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(20),
    unit_type VARCHAR(50) DEFAULT 'piece', -- weight, volume, length, piece, pack
    conversion_factor DECIMAL(18, 6) DEFAULT 1, -- to base unit
    base_unit_id BIGINT REFERENCES master_data.units(id) ON DELETE SET NULL,
    is_base_unit BOOLEAN DEFAULT false,
    description TEXT,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(company_id, code)
);

-- Product Categories (Hierarchical)
CREATE TABLE master_data.categories (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    parent_id BIGINT REFERENCES master_data.categories(id) ON DELETE SET NULL,
    level INTEGER DEFAULT 0,
    path LTREE, -- For hierarchical queries (needs ltree extension)
    sort_order INTEGER DEFAULT 0,
    image_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(company_id, code)
);

-- Enable ltree for hierarchical categories
CREATE EXTENSION IF NOT EXISTS "ltree";
CREATE INDEX idx_categories_path ON master_data.categories USING GIST(path);

-- Product Attributes (for variants)
CREATE TABLE master_data.attributes (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    data_type VARCHAR(50) DEFAULT 'string', -- string, number, boolean, date, select
    is_required BOOLEAN DEFAULT false,
    is_variant_attribute BOOLEAN DEFAULT false, -- Used for product variants
    is_filterable BOOLEAN DEFAULT false, -- Show in filters
    options JSONB, -- For select type: ["Red", "Blue", "Green"]
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(company_id, code)
);

-- Products
CREATE TABLE master_data.products (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    code VARCHAR(100) NOT NULL,
    sku VARCHAR(100),
    barcode VARCHAR(100),
    name VARCHAR(250) NOT NULL,
    description TEXT,
    short_description VARCHAR(500),
    category_id BIGINT REFERENCES master_data.categories(id) ON DELETE SET NULL,
    brand VARCHAR(100),
    base_unit_id BIGINT REFERENCES master_data.units(id),
    weight DECIMAL(18, 4),
    dimensions JSONB, -- {"length": 10, "width": 5, "height": 3}
    track_inventory BOOLEAN DEFAULT true,
    is_service BOOLEAN DEFAULT false, -- Non-stock item
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_by UUID REFERENCES auth.users(id),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    metadata JSONB DEFAULT '{}',
    UNIQUE(company_id, code)
);

CREATE INDEX idx_products_company ON master_data.products(company_id);
CREATE INDEX idx_products_category ON master_data.products(category_id);
CREATE INDEX idx_products_name ON master_data.products USING gin(to_tsvector('english', name));

-- Product Variants (SKU level)
CREATE TABLE master_data.product_variants (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT REFERENCES master_data.products(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    sku VARCHAR(100) NOT NULL,
    barcode VARCHAR(100),
    name VARCHAR(250), -- e.g., "T-Shirt Red Large"
    options JSONB, -- {"color": "Red", "size": "Large"}
    weight DECIMAL(18, 4),
    dimensions JSONB,
    is_default BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(product_id, sku)
);

CREATE INDEX idx_variants_product ON master_data.product_variants(product_id);
CREATE INDEX idx_variants_sku ON master_data.product_variants(sku);
CREATE INDEX idx_variants_barcode ON master_data.product_variants(barcode);

-- Product Images (stored in Supabase Storage, reference here)
CREATE TABLE master_data.product_images (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT REFERENCES master_data.products(id) ON DELETE CASCADE,
    variant_id BIGINT REFERENCES master_data.product_variants(id) ON DELETE CASCADE,
    storage_path VARCHAR(500) NOT NULL, -- bucket/folder/filename
    file_name VARCHAR(250),
    file_size INTEGER,
    mime_type VARCHAR(50),
    is_primary BOOLEAN DEFAULT false,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- ============================================
-- INVENTORY SCHEMA
-- ============================================

-- Stock (per variant per location)
CREATE TABLE inventory.stocks (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    variant_id BIGINT REFERENCES master_data.product_variants(id) ON DELETE CASCADE,
    location_id BIGINT REFERENCES public.locations(id) ON DELETE CASCADE,
    quantity_available DECIMAL(18, 4) DEFAULT 0 NOT NULL,
    quantity_reserved DECIMAL(18, 4) DEFAULT 0 NOT NULL,
    quantity_on_hand DECIMAL(18, 4) GENERATED ALWAYS AS (quantity_available + quantity_reserved) STORED,
    quantity_incoming DECIMAL(18, 4) DEFAULT 0, -- From PO
    quantity_outgoing DECIMAL(18, 4) DEFAULT 0, -- From SO
    reorder_point DECIMAL(18, 4) DEFAULT 0,
    reorder_quantity DECIMAL(18, 4),
    max_stock DECIMAL(18, 4),
    avg_cost DECIMAL(18, 4) DEFAULT 0,
    last_counted_at TIMESTAMPTZ,
    last_movement_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(variant_id, location_id)
);

CREATE INDEX idx_stocks_company ON inventory.stocks(company_id);
CREATE INDEX idx_stocks_location ON inventory.stocks(location_id);
CREATE INDEX idx_stocks_low_stock ON inventory.stocks(quantity_available) WHERE quantity_available <= reorder_point;

-- Storage Bins / Rack locations
CREATE TABLE inventory.storage_bins (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    location_id BIGINT REFERENCES public.locations(id) ON DELETE CASCADE,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100),
    zone VARCHAR(50), -- receiving, picking, storage, shipping, quarantine
    aisle VARCHAR(50),
    rack VARCHAR(50),
    shelf VARCHAR(50),
    bin_type VARCHAR(50) DEFAULT 'standard', -- standard, bulk, cold, hazardous
    capacity DECIMAL(18, 4),
    capacity_unit_id BIGINT REFERENCES master_data.units(id),
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(location_id, code)
);

-- Stock by Bin (detailed location)
CREATE TABLE inventory.stock_bins (
    id BIGSERIAL PRIMARY KEY,
    stock_id BIGINT REFERENCES inventory.stocks(id) ON DELETE CASCADE,
    bin_id BIGINT REFERENCES inventory.storage_bins(id) ON DELETE CASCADE,
    quantity DECIMAL(18, 4) DEFAULT 0 NOT NULL,
    lot_number VARCHAR(100),
    expiry_date DATE,
    received_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(stock_id, bin_id, lot_number)
);

-- Stock Movements (transaction log)
CREATE TYPE inventory.movement_type AS ENUM ('in', 'out', 'transfer', 'adjustment', 'count');
CREATE TYPE inventory.movement_reason AS ENUM ('purchase', 'sale', 'return', 'transfer', 'adjustment', 'damage', 'expired', 'production', 'count');

CREATE TABLE inventory.movements (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    movement_number VARCHAR(100) NOT NULL,
    movement_date DATE NOT NULL DEFAULT CURRENT_DATE,
    movement_type inventory.movement_type NOT NULL,
    reason inventory.movement_reason NOT NULL,
    
    -- References
    reference_type VARCHAR(50), -- PO, SO, Transfer, Adjustment
    reference_id BIGINT,
    reference_number VARCHAR(100),
    
    -- Location
    location_id BIGINT REFERENCES public.locations(id) ON DELETE CASCADE,
    from_location_id BIGINT REFERENCES public.locations(id) ON DELETE SET NULL,
    to_location_id BIGINT REFERENCES public.locations(id) ON DELETE SET NULL,
    
    -- Totals
    total_items INTEGER DEFAULT 0,
    total_quantity DECIMAL(18, 4) DEFAULT 0,
    notes TEXT,
    
    -- Audit
    created_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    posted_at TIMESTAMPTZ, -- When stock was actually updated
    
    UNIQUE(company_id, movement_number)
);

CREATE INDEX idx_movements_company ON inventory.movements(company_id);
CREATE INDEX idx_movements_date ON inventory.movements(movement_date);
CREATE INDEX idx_movements_type ON inventory.movements(movement_type);

-- Movement Details
CREATE TABLE inventory.movement_details (
    id BIGSERIAL PRIMARY KEY,
    movement_id BIGINT REFERENCES inventory.movements(id) ON DELETE CASCADE,
    variant_id BIGINT REFERENCES master_data.product_variants(id) ON DELETE CASCADE,
    from_bin_id BIGINT REFERENCES inventory.storage_bins(id) ON DELETE SET NULL,
    to_bin_id BIGINT REFERENCES inventory.storage_bins(id) ON DELETE SET NULL,
    lot_number VARCHAR(100),
    quantity DECIMAL(18, 4) NOT NULL,
    unit_cost DECIMAL(18, 4),
    total_cost DECIMAL(18, 4),
    expiry_date DATE,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Stock Reservations (for sales orders)
CREATE TABLE inventory.reservations (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    stock_id BIGINT REFERENCES inventory.stocks(id) ON DELETE CASCADE,
    reference_type VARCHAR(50) NOT NULL, -- sales_order, production_order
    reference_id BIGINT NOT NULL,
    quantity_reserved DECIMAL(18, 4) NOT NULL,
    quantity_released DECIMAL(18, 4) DEFAULT 0,
    expiry_date DATE, -- Reservation expires
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    released_at TIMESTAMPTZ
);

-- ============================================
-- SALES SCHEMA
-- ============================================

-- Customer Groups
CREATE TABLE sales.customer_groups (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    discount_type VARCHAR(20) DEFAULT 'percentage', -- percentage, fixed
    discount_value DECIMAL(18, 4) DEFAULT 0,
    min_purchase_amount DECIMAL(18, 4),
    credit_limit DECIMAL(18, 4),
    payment_term_days INTEGER DEFAULT 0,
    price_tier INTEGER DEFAULT 1, -- 1-10 for multi-tier pricing
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(company_id, code)
);

-- Customers
CREATE TABLE sales.customers (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    code VARCHAR(100) NOT NULL,
    customer_group_id BIGINT REFERENCES sales.customer_groups(id) ON DELETE SET NULL,
    
    -- Basic Info
    name VARCHAR(250) NOT NULL,
    email VARCHAR(150),
    phone VARCHAR(50),
    mobile VARCHAR(50),
    tax_id VARCHAR(50),
    
    -- Address
    billing_address TEXT,
    billing_city VARCHAR(100),
    billing_province VARCHAR(100),
    billing_postal_code VARCHAR(20),
    billing_country VARCHAR(100) DEFAULT 'Indonesia',
    
    -- Credit
    credit_limit DECIMAL(18, 4) DEFAULT 0,
    current_balance DECIMAL(18, 4) DEFAULT 0,
    points INTEGER DEFAULT 0,
    
    -- Settings
    price_tier INTEGER DEFAULT 1,
    is_wholesale BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true NOT NULL,
    
    -- Metadata
    birth_date DATE,
    gender VARCHAR(20),
    notes TEXT,
    
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(company_id, code)
);

-- Customer Shipping Addresses
CREATE TABLE sales.customer_addresses (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT REFERENCES sales.customers(id) ON DELETE CASCADE,
    label VARCHAR(100) NOT NULL DEFAULT 'Default', -- Home, Office, Warehouse
    recipient_name VARCHAR(250),
    phone VARCHAR(50),
    address TEXT NOT NULL,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) DEFAULT 'Indonesia',
    is_default BOOLEAN DEFAULT false,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Sales Orders
CREATE TYPE sales.order_status AS ENUM ('draft', 'confirmed', 'processing', 'picked', 'packed', 'shipped', 'delivered', 'cancelled', 'returned');
CREATE TYPE sales.payment_status AS ENUM ('unpaid', 'partial', 'paid', 'refunded', 'overpaid');

CREATE TABLE sales.orders (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    
    -- Order Info
    order_number VARCHAR(100) NOT NULL,
    order_date DATE NOT NULL DEFAULT CURRENT_DATE,
    delivery_date DATE,
    
    -- Location
    location_id BIGINT REFERENCES public.locations(id) ON DELETE SET NULL,
    
    -- Customer
    customer_id BIGINT REFERENCES sales.customers(id) ON DELETE SET NULL,
    customer_name VARCHAR(250), -- Denormalized
    customer_address_id BIGINT REFERENCES sales.customer_addresses(id) ON DELETE SET NULL,
    
    -- Sales Person
    sales_person_id UUID REFERENCES auth.users(id) ON DELETE SET NULL,
    
    -- Pricing
    subtotal DECIMAL(18, 4) DEFAULT 0 NOT NULL,
    discount_amount DECIMAL(18, 4) DEFAULT 0,
    discount_percent DECIMAL(5, 2) DEFAULT 0,
    tax_percent DECIMAL(5, 2) DEFAULT 0,
    tax_amount DECIMAL(18, 4) DEFAULT 0,
    shipping_cost DECIMAL(18, 4) DEFAULT 0,
    total_amount DECIMAL(18, 4) DEFAULT 0 NOT NULL,
    
    -- Payment
    payment_status sales.payment_status DEFAULT 'unpaid',
    paid_amount DECIMAL(18, 4) DEFAULT 0,
    
    -- Status
    status sales.order_status DEFAULT 'draft',
    notes TEXT,
    internal_notes TEXT,
    
    -- Source
    source VARCHAR(50) DEFAULT 'manual', -- manual, pos, ecommerce, api
    source_reference VARCHAR(100),
    
    -- Timestamps
    confirmed_at TIMESTAMPTZ,
    shipped_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    
    created_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_by UUID REFERENCES auth.users(id),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    UNIQUE(company_id, order_number)
);

CREATE INDEX idx_orders_company ON sales.orders(company_id);
CREATE INDEX idx_orders_customer ON sales.orders(customer_id);
CREATE INDEX idx_orders_status ON sales.orders(status);
CREATE INDEX idx_orders_date ON sales.orders(order_date);

-- Sales Order Items
CREATE TABLE sales.order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES sales.orders(id) ON DELETE CASCADE,
    
    -- Product
    variant_id BIGINT REFERENCES master_data.product_variants(id) ON DELETE SET NULL,
    variant_sku VARCHAR(100), -- Denormalized
    variant_name VARCHAR(250), -- Denormalized
    
    -- Quantity
    quantity_ordered DECIMAL(18, 4) NOT NULL,
    quantity_delivered DECIMAL(18, 4) DEFAULT 0,
    quantity_returned DECIMAL(18, 4) DEFAULT 0,
    unit_id BIGINT REFERENCES master_data.units(id),
    
    -- Pricing
    unit_price DECIMAL(18, 4) NOT NULL,
    discount_percent DECIMAL(5, 2) DEFAULT 0,
    discount_amount DECIMAL(18, 4) DEFAULT 0,
    tax_percent DECIMAL(5, 2) DEFAULT 0,
    tax_amount DECIMAL(18, 4) DEFAULT 0,
    total_price DECIMAL(18, 4) NOT NULL,
    cost DECIMAL(18, 4), -- For profit calculation
    
    -- Fulfillment
    from_location_id BIGINT REFERENCES public.locations(id) ON DELETE SET NULL,
    picked_at TIMESTAMPTZ,
    picked_by UUID REFERENCES auth.users(id),
    
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Price Tiers (Multi-tier pricing)
CREATE TABLE sales.price_tiers (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    variant_id BIGINT REFERENCES master_data.product_variants(id) ON DELETE CASCADE,
    customer_group_id BIGINT REFERENCES sales.customer_groups(id) ON DELETE SET NULL,
    unit_id BIGINT REFERENCES master_data.units(id),
    
    tier_level INTEGER NOT NULL DEFAULT 1, -- 1-10
    min_quantity DECIMAL(18, 4) DEFAULT 1,
    price DECIMAL(18, 4) NOT NULL,
    
    effective_from DATE DEFAULT CURRENT_DATE,
    effective_until DATE,
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(variant_id, customer_group_id, tier_level, unit_id)
);

-- Discount Events
CREATE TABLE sales.discount_events (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    
    discount_type VARCHAR(20) NOT NULL, -- percentage, fixed, bogo
    discount_value DECIMAL(18, 4) NOT NULL,
    min_purchase_amount DECIMAL(18, 4) DEFAULT 0,
    max_discount_amount DECIMAL(18, 4),
    
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    
    usage_limit INTEGER,
    usage_count INTEGER DEFAULT 0,
    
    applies_to VARCHAR(50) DEFAULT 'all', -- all, products, categories, customers
    conditions JSONB, -- Flexible conditions
    
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    
    UNIQUE(company_id, code)
);

-- ============================================
-- PURCHASE SCHEMA
-- ============================================

-- Suppliers
CREATE TABLE purchase.suppliers (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(250) NOT NULL,
    
    email VARCHAR(150),
    phone VARCHAR(50),
    mobile VARCHAR(50),
    tax_id VARCHAR(50),
    
    address TEXT,
    city VARCHAR(100),
    province VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) DEFAULT 'Indonesia',
    website VARCHAR(150),
    
    payment_terms_days INTEGER DEFAULT 0,
    credit_limit DECIMAL(18, 4) DEFAULT 0,
    current_balance DECIMAL(18, 4) DEFAULT 0,
    
    is_active BOOLEAN DEFAULT true NOT NULL,
    notes TEXT,
    
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(company_id, code)
);

-- Supplier Product Pricing
CREATE TABLE purchase.supplier_prices (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    supplier_id BIGINT REFERENCES purchase.suppliers(id) ON DELETE CASCADE,
    variant_id BIGINT REFERENCES master_data.product_variants(id) ON DELETE CASCADE,
    unit_id BIGINT REFERENCES master_data.units(id),
    
    supplier_sku VARCHAR(100),
    min_order_qty DECIMAL(18, 4) DEFAULT 1,
    price DECIMAL(18, 4) NOT NULL,
    lead_time_days INTEGER DEFAULT 0,
    
    effective_from DATE DEFAULT CURRENT_DATE,
    effective_until DATE,
    is_preferred BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true NOT NULL,
    
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(supplier_id, variant_id, unit_id)
);

-- Purchase Orders
CREATE TYPE purchase.po_status AS ENUM ('draft', 'sent', 'partial', 'received', 'closed', 'cancelled');

CREATE TABLE purchase.orders (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    
    po_number VARCHAR(100) NOT NULL,
    po_date DATE NOT NULL DEFAULT CURRENT_DATE,
    expected_delivery_date DATE,
    
    location_id BIGINT REFERENCES public.locations(id) ON DELETE SET NULL,
    supplier_id BIGINT REFERENCES purchase.suppliers(id) ON DELETE SET NULL,
    supplier_name VARCHAR(250), -- Denormalized
    
    subtotal DECIMAL(18, 4) DEFAULT 0 NOT NULL,
    discount_amount DECIMAL(18, 4) DEFAULT 0,
    tax_percent DECIMAL(5, 2) DEFAULT 0,
    tax_amount DECIMAL(18, 4) DEFAULT 0,
    shipping_cost DECIMAL(18, 4) DEFAULT 0,
    total_amount DECIMAL(18, 4) DEFAULT 0 NOT NULL,
    
    status purchase.po_status DEFAULT 'draft',
    
    notes TEXT,
    internal_notes TEXT,
    
    created_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    UNIQUE(company_id, po_number)
);

-- Purchase Order Items
CREATE TABLE purchase.order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES purchase.orders(id) ON DELETE CASCADE,
    
    variant_id BIGINT REFERENCES master_data.product_variants(id) ON DELETE SET NULL,
    variant_sku VARCHAR(100), -- Denormalized
    variant_name VARCHAR(250), -- Denormalized
    
    quantity_ordered DECIMAL(18, 4) NOT NULL,
    quantity_received DECIMAL(18, 4) DEFAULT 0,
    quantity_cancelled DECIMAL(18, 4) DEFAULT 0,
    
    unit_id BIGINT REFERENCES master_data.units(id),
    unit_price DECIMAL(18, 4) NOT NULL,
    discount_percent DECIMAL(5, 2) DEFAULT 0,
    total_price DECIMAL(18, 4) NOT NULL,
    
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- Goods Receipts
CREATE TABLE purchase.receipts (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES public.companies(id) ON DELETE CASCADE,
    uuid UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    
    receipt_number VARCHAR(100) NOT NULL,
    receipt_date DATE NOT NULL DEFAULT CURRENT_DATE,
    
    po_id BIGINT REFERENCES purchase.orders(id) ON DELETE SET NULL,
    supplier_id BIGINT REFERENCES purchase.suppliers(id) ON DELETE SET NULL,
    location_id BIGINT REFERENCES public.locations(id) ON DELETE SET NULL,
    
    total_items INTEGER DEFAULT 0,
    total_quantity DECIMAL(18, 4) DEFAULT 0,
    
    notes TEXT,
    
    created_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    
    UNIQUE(company_id, receipt_number)
);

-- Receipt Items
CREATE TABLE purchase.receipt_items (
    id BIGSERIAL PRIMARY KEY,
    receipt_id BIGINT REFERENCES purchase.receipts(id) ON DELETE CASCADE,
    po_item_id BIGINT REFERENCES purchase.order_items(id) ON DELETE SET NULL,
    
    variant_id BIGINT REFERENCES master_data.product_variants(id) ON DELETE SET NULL,
    quantity_received DECIMAL(18, 4) NOT NULL,
    unit_id BIGINT REFERENCES master_data.units(id),
    unit_cost DECIMAL(18, 4) NOT NULL,
    
    bin_id BIGINT REFERENCES inventory.storage_bins(id) ON DELETE SET NULL,
    lot_number VARCHAR(100),
    expiry_date DATE,
    
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- ============================================
-- FUNCTIONS & TRIGGERS
-- ============================================

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply to all tables with updated_at
CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON public.companies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_company_users_updated_at BEFORE UPDATE ON public.company_users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON master_data.products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_variants_updated_at BEFORE UPDATE ON master_data.product_variants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_stocks_updated_at BEFORE UPDATE ON inventory.stocks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON sales.orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
CREATE TRIGGER update_po_updated_at BEFORE UPDATE ON purchase.orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================

-- Enable RLS on all tables
ALTER TABLE public.companies ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.company_users ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.locations ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.audit_logs ENABLE ROW LEVEL SECURITY;

ALTER TABLE master_data.units ENABLE ROW LEVEL SECURITY;
ALTER TABLE master_data.categories ENABLE ROW LEVEL SECURITY;
ALTER TABLE master_data.attributes ENABLE ROW LEVEL SECURITY;
ALTER TABLE master_data.products ENABLE ROW LEVEL SECURITY;
ALTER TABLE master_data.product_variants ENABLE ROW LEVEL SECURITY;

ALTER TABLE inventory.stocks ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory.storage_bins ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory.stock_bins ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory.movements ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory.movement_details ENABLE ROW LEVEL SECURITY;
ALTER TABLE inventory.reservations ENABLE ROW LEVEL SECURITY;

ALTER TABLE sales.customer_groups ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales.customers ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales.customer_addresses ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales.orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales.order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales.price_tiers ENABLE ROW LEVEL SECURITY;
ALTER TABLE sales.discount_events ENABLE ROW LEVEL SECURITY;

ALTER TABLE purchase.suppliers ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase.supplier_prices ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase.orders ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase.order_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase.receipts ENABLE ROW LEVEL SECURITY;
ALTER TABLE purchase.receipt_items ENABLE ROW LEVEL SECURITY;

-- Create helper function to get current user's company
CREATE OR REPLACE FUNCTION public.current_user_company()
RETURNS BIGINT AS $$
DECLARE
    company_id BIGINT;
BEGIN
    -- Get company_id from JWT claims
    SELECT (auth.jwt() ->> 'company_id')::BIGINT INTO company_id;
    RETURN company_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Create policies

-- Companies: Users can only see their own company
CREATE POLICY companies_isolation ON public.companies
    FOR ALL USING (id = public.current_user_company());

-- Company Users: Users can see users in their company
CREATE POLICY company_users_isolation ON public.company_users
    FOR ALL USING (company_id = public.current_user_company());

-- Locations
CREATE POLICY locations_isolation ON public.locations
    FOR ALL USING (company_id = public.current_user_company());

-- Audit Logs
CREATE POLICY audit_logs_isolation ON public.audit_logs
    FOR ALL USING (company_id = public.current_user_company());

-- Master Data
CREATE POLICY units_isolation ON master_data.units
    FOR ALL USING (company_id = public.current_user_company());
    
CREATE POLICY categories_isolation ON master_data.categories
    FOR ALL USING (company_id = public.current_user_company());
    
CREATE POLICY attributes_isolation ON master_data.attributes
    FOR ALL USING (company_id = public.current_user_company());
    
CREATE POLICY products_isolation ON master_data.products
    FOR ALL USING (company_id = public.current_user_company());
    
CREATE POLICY variants_isolation ON master_data.product_variants
    FOR ALL USING (
        EXISTS (
            SELECT 1 FROM master_data.products p 
            WHERE p.id = product_variants.product_id 
            AND p.company_id = public.current_user_company()
        )
    );

-- Inventory
CREATE POLICY stocks_isolation ON inventory.stocks
    FOR ALL USING (company_id = public.current_user_company());
    
CREATE POLICY storage_bins_isolation ON inventory.storage_bins
    FOR ALL USING (company_id = public.current_user_company());
    
CREATE POLICY movements_isolation ON inventory.movements
    FOR ALL USING (company_id = public.current_user_company());

-- Sales
CREATE POLICY customers_isolation ON sales.customers
    FOR ALL USING (company_id = public.current_user_company());
    
CREATE POLICY orders_isolation ON sales.orders
    FOR ALL USING (company_id = public.current_user_company());

-- Purchase
CREATE POLICY suppliers_isolation ON purchase.suppliers
    FOR ALL USING (company_id = public.current_user_company());
    
CREATE POLICY po_isolation ON purchase.orders
    FOR ALL USING (company_id = public.current_user_company());

-- ============================================
-- VIEWS FOR REPORTING
-- ============================================

-- Stock Summary
CREATE VIEW inventory.v_stock_summary AS
SELECT 
    s.id,
    s.company_id,
    pv.id AS variant_id,
    pv.sku,
    pv.name AS variant_name,
    p.id AS product_id,
    p.name AS product_name,
    c.name AS category_name,
    l.id AS location_id,
    l.name AS location_name,
    s.quantity_available,
    s.quantity_reserved,
    s.quantity_on_hand,
    s.quantity_incoming,
    s.quantity_outgoing,
    s.reorder_point,
    s.avg_cost,
    (s.quantity_available * s.avg_cost) AS stock_value,
    CASE 
        WHEN s.quantity_available <= s.reorder_point AND s.reorder_point > 0 
        THEN true 
        ELSE false 
    END AS is_low_stock
FROM inventory.stocks s
JOIN master_data.product_variants pv ON s.variant_id = pv.id
JOIN master_data.products p ON pv.product_id = p.id
LEFT JOIN master_data.categories c ON p.category_id = c.id
JOIN public.locations l ON s.location_id = l.id
WHERE s.deleted_at IS NULL;

-- Daily Sales Summary
CREATE VIEW sales.v_daily_sales AS
SELECT 
    company_id,
    DATE(order_date) AS sales_date,
    COUNT(*) AS order_count,
    SUM(total_amount) AS total_revenue,
    SUM(subtotal) AS total_subtotal,
    SUM(discount_amount) AS total_discounts,
    SUM(tax_amount) AS total_tax,
    SUM(shipping_cost) AS total_shipping
FROM sales.orders
WHERE status NOT IN ('draft', 'cancelled')
GROUP BY company_id, DATE(order_date);

-- Product Sales Performance
CREATE VIEW sales.v_product_performance AS
SELECT 
    oi.variant_id,
    pv.sku,
    COALESCE(pv.name, p.name) AS product_name,
    c.name AS category_name,
    SUM(oi.quantity_ordered) AS total_qty_sold,
    SUM(oi.total_price) AS total_revenue,
    SUM(oi.quantity_ordered * oi.cost) AS total_cost,
    SUM(oi.total_price) - SUM(oi.quantity_ordered * oi.cost) AS gross_profit
FROM sales.order_items oi
JOIN sales.orders o ON oi.order_id = o.id
JOIN master_data.product_variants pv ON oi.variant_id = pv.id
JOIN master_data.products p ON pv.product_id = p.id
LEFT JOIN master_data.categories c ON p.category_id = c.id
WHERE o.status IN ('delivered', 'shipped')
GROUP BY oi.variant_id, pv.sku, COALESCE(pv.name, p.name), c.name;

-- ============================================
-- COMMENTS
-- ============================================
COMMENT ON TABLE public.companies IS 'Root table for multi-tenancy';
COMMENT ON TABLE public.company_users IS 'Link between Supabase Auth users and companies';
COMMENT ON TABLE inventory.stocks IS 'Stock levels per variant per location';
COMMENT ON TABLE sales.orders IS 'Sales orders with RLS by company';
COMMENT ON TABLE purchase.orders IS 'Purchase orders with RLS by company';

-- ============================================
-- END OF SCHEMA
-- ============================================
