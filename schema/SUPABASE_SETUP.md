# Supabase WMS Schema - Setup Guide

## 📁 Files

| File | Description |
|------|-------------|
| `supabase_schema.sql` | Complete schema with RLS, ~70 tables |
| `postgresql_schema.sql` | Basic schema (alternative) |
| `postgresql_schema_part2.sql` | Extended tables |

## 🚀 Quick Setup

### 1. Create Supabase Project
```bash
# Via dashboard or CLI
supabase projects create neo-wms
```

### 2. Run Schema

#### Option A: Via SQL Editor (Dashboard)
1. Go to Supabase Dashboard → SQL Editor
2. Open `supabase_schema.sql`
3. Copy all content
4. Paste and Run

#### Option B: Via Supabase CLI
```bash
# Link to project
supabase link --project-ref your-project-ref

# Push schema
supabase db push

# Or execute directly
psql "postgresql://postgres:[password]@db.[ref].supabase.co:5432/postgres" -f supabase_schema.sql
```

#### Option C: Via Migration (Recommended for production)
```bash
# Create migration
supabase migration new initial_schema

# Copy content to migration file
cp supabase_schema.sql supabase/migrations/000001_initial_schema.sql

# Apply migration
supabase db push
```

### 3. Verify Setup
```sql
-- Check schemas
SELECT schema_name FROM information_schema.schemata 
WHERE schema_name IN ('public', 'master_data', 'inventory', 'sales', 'purchase', 'accounting');

-- Check tables
SELECT table_schema, table_name 
FROM information_schema.tables 
WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
ORDER BY table_schema, table_name;

-- Check RLS
SELECT schemaname, tablename, rowsecurity 
FROM pg_tables 
WHERE rowsecurity = true;
```

## 🔐 RLS (Row Level Security) Setup

### Test RLS
```sql
-- Set company context (simulates JWT claim)
SET LOCAL app.current_company_id = '1';

-- Query will only return rows for company_id = 1
SELECT * FROM master_data.products;

-- Reset
RESET app.current_company_id;
```

### Add Company ID to JWT
In your application middleware:
```typescript
// When creating Supabase client
const supabase = createClient(url, key, {
  auth: {
    persistSession: true,
    autoRefreshToken: true,
  },
  global: {
    headers: {
      'X-Company-ID': companyId, // Custom header
    },
  },
});

// Or use custom claim
const { data: { user } } = await supabase.auth.getUser();
const companyId = user?.user_metadata?.company_id;
```

## 📊 Schema Overview

### Public Schema (System)
```
public.companies          - Multi-tenant root
public.company_users      - Auth user ↔ Company link
public.locations          - Warehouses, stores, virtual
public.audit_logs         - System audit trail
```

### Domain Schemas
```
master_data.*    - Products, categories, units, attributes
inventory.*      - Stocks, movements, storage bins, reservations
sales.*          - Customers, orders, pricing, discounts
purchase.*       - Suppliers, POs, receipts
accounting.*      - Reserved for future (accounts, journals)
```

## 🔄 Realtime Subscriptions

### Enable Realtime for Tables
```sql
-- Enable realtime for stock updates
alter publication supabase_realtime add table inventory.stocks;
alter publication supabase_realtime add table inventory.movements;
alter publication supabase_realtime add table sales.orders;
```

### Subscribe from Frontend
```typescript
// React/Vue/Angular
const subscription = supabase
  .from('stocks:company_id=eq.1')
  .on('*', (payload) => {
    console.log('Stock changed!', payload);
  })
  .subscribe();

// Unsubscribe
supabase.removeSubscription(subscription);
```

## 📦 Storage Setup

### Create Buckets
```sql
-- Via SQL (Storage schema)
INSERT INTO storage.buckets (id, name, public)
VALUES 
  ('product-images', 'product-images', true),
  ('documents', 'documents', false),
  ('exports', 'exports', false);
```

### Storage Policies
```sql
-- Product images: Company-scoped access
CREATE POLICY "Company images" ON storage.objects
  FOR ALL 
  USING (
    bucket_id = 'product-images' 
    AND (storage.foldername(name))[1] = auth.jwt() ->> 'company_id'
  );
```

## 🧪 Test Data

### Insert Sample Data
```sql
-- Company
INSERT INTO public.companies (code, name, subscription_tier)
VALUES ('DEMO', 'Demo Company', 'basic')
RETURNING id;

-- Assume company_id = 1

-- Location
INSERT INTO public.locations (company_id, code, name, location_type)
VALUES (1, 'WH01', 'Main Warehouse', 'warehouse');

-- Unit
INSERT INTO master_data.units (company_id, code, name, is_base_unit)
VALUES (1, 'PCS', 'Pieces', true);

-- Category
INSERT INTO master_data.categories (company_id, code, name)
VALUES (1, 'ELEC', 'Electronics');

-- Product
INSERT INTO master_data.products (
  company_id, code, name, category_id, base_unit_id, track_inventory
) VALUES (1, 'PROD001', 'Laptop', 1, 1, true)
RETURNING id;

-- Variant
INSERT INTO master_data.product_variants (product_id, sku, name, is_default)
VALUES (1, 'LAPTOP-001', 'Standard Laptop', true);
```

## 🔗 Integration with Go Backend

### Connection String
```bash
# From Supabase Dashboard → Settings → Database
DATABASE_URL="postgresql://postgres:[password]@db.[ref].supabase.co:5432/postgres"
```

### GORM Models (Example)
```go
package models

import (
    "gorm.io/gorm"
    "github.com/google/uuid"
)

type Company struct {
    ID        uint64         `gorm:"primaryKey"`
    UUID      uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid()"`
    Code      string         `gorm:"uniqueIndex"`
    Name      string
    IsActive  bool
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Product struct {
    ID          uint64    `gorm:"primaryKey"`
    CompanyID   uint64    `gorm:"index"`
    UUID        uuid.UUID `gorm:"type:uuid"`
    Code        string    `gorm:"uniqueIndex:idx_products_code_company"`
    Name        string
    CategoryID  *uint64
    Category    Category  `gorm:"foreignKey:CategoryID"`
    IsActive    bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}
```

### Middleware for RLS
```go
// Set company context for RLS
func CompanyContextMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        companyID := c.GetHeader("X-Company-ID")
        if companyID != "" {
            // Set in database session
            db.Exec("SET LOCAL app.current_company_id = ?", companyID)
        }
        c.Next()
    }
}
```

## 📝 Naming Conventions

### Tables
- Schema-prefixed: `master_data.products`, `inventory.stocks`
- Singular or plural? **Plural**: `products`, `orders`
- Snake_case: `product_variants`, `order_items`

### Columns
- Primary key: `id` (BIGSERIAL)
- UUID: `uuid` (for external reference)
- Foreign key: `{table}_id` (e.g., `company_id`, `product_id`)
- Timestamps: `created_at`, `updated_at`, `deleted_at`
- Audit: `created_by`, `updated_by` (UUID)

### Enums
- PostgreSQL native enums for type safety
- Or use CHECK constraints for flexibility

## 🛡️ Security Best Practices

1. **Always enable RLS** - Never disable on production
2. **Use service role carefully** - Bypasses RLS
3. **Validate company_id** - In application layer too
4. **Audit sensitive operations** - Log to audit_logs
5. **Use prepared statements** - Prevent SQL injection

## 📈 Performance Optimization

### Indexes Created
- All foreign keys
- UUID columns
- Status fields
- Date ranges
- Full-text search (GIN indexes)

### Query Optimization
```sql
-- Use views for complex queries
SELECT * FROM inventory.v_stock_summary 
WHERE company_id = 1 AND is_low_stock = true;

-- Use materialized views for reporting
REFRESH MATERIALIZED VIEW sales.v_daily_sales;
```

## 🔧 Troubleshooting

### RLS Not Working?
```sql
-- Check if RLS is enabled
SELECT relname, relrowsecurity 
FROM pg_class 
WHERE relrowsecurity = true;

-- Check policies
SELECT schemaname, tablename, policyname, permissive, roles, cmd, qual
FROM pg_policies;
```

### Permission Denied?
```sql
-- Grant access to authenticated users
GRANT USAGE ON SCHEMA master_data TO authenticated;
GRANT ALL ON ALL TABLES IN SCHEMA master_data TO authenticated;
```

## 📚 Resources

- [Supabase Docs](https://supabase.com/docs)
- [PostgreSQL RLS](https://www.postgresql.org/docs/current/ddl-rowsecurity.html)
- [GORM with PostgreSQL](https://gorm.io/docs/connecting_to_the_database.html)

---

**Ready to deploy?** Copy `supabase_schema.sql` to your Supabase SQL Editor and run! 🚀
