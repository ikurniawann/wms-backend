# NEO WMS Backend

Go-based backend for NEO WMS using Gin + GORM + PostgreSQL (Supabase).

## 📁 Structure

```
backend/
├── main.go              # Application entry point
├── config/              # Configuration management
│   └── config.go
├── database/            # Database connection & migration
│   └── database.go
├── models/              # GORM models
│   ├── company.go       # System models (companies, users, locations)
│   ├── master_data.go   # Master data (products, categories, units)
│   ├── inventory.go     # Inventory (stocks, movements, bins)
│   ├── sales.go         # Sales (customers, orders, pricing)
│   ├── purchase.go      # Purchase (suppliers, POs, receipts)
│   └── common.go        # Common types & utilities
├── handlers/            # HTTP handlers/controllers
│   ├── auth.go
│   ├── product.go
│   ├── inventory.go
│   ├── sales.go
│   └── purchase.go
├── middleware/          # Gin middleware
│   ├── auth.go
│   ├── cors.go
│   └── company.go
├── services/            # Business logic
│   ├── product.go
│   ├── inventory.go
│   └── sales.go
├── repositories/        # Data access layer
│   └── ...
├── dto/                 # Data transfer objects
│   └── ...
├── go.mod
└── go.sum
```

## 🚀 Quick Start

### 1. Setup Environment

Create `.env` file:

```bash
# Database (Supabase)
DB_HOST=db.xxxxx.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=postgres
DB_SSLMODE=require

# JWT (Supabase)
JWT_SECRET=your-jwt-secret
SUPABASE_URL=https://xxxxx.supabase.co
SUPABASE_KEY=your-anon-key

# Server
PORT=8080
ENV=development
```

### 2. Install Dependencies

```bash
go mod init github.com/ikurniawann/wmsmicroservice
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/gin-gonic/gin
go get -u github.com/google/uuid
go get -u github.com/joho/godotenv
```

### 3. Run Migration

```bash
go run main.go migrate
```

### 4. Start Server

```bash
go run main.go
```

## 📚 API Documentation

See [API_DOCS.md](./API_DOCS.md) for full API documentation.

### Base URL

```
Development: http://localhost:8080/api/v1
Production:  https://your-domain.com/api/v1
```

### Authentication

All endpoints (except `/auth/*`) require JWT token in header:

```
Authorization: Bearer <token>
```

### Endpoints Overview

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/login` | Login with Supabase Auth |
| GET | `/products` | List products |
| POST | `/products` | Create product |
| GET | `/products/:id` | Get product detail |
| PUT | `/products/:id` | Update product |
| DELETE | `/products/:id` | Delete product |
| GET | `/stocks` | List stock levels |
| POST | `/stocks/adjust` | Adjust stock |
| GET | `/sales/orders` | List sales orders |
| POST | `/sales/orders` | Create sales order |
| GET | `/purchase/orders` | List purchase orders |
| POST | `/purchase/orders` | Create purchase order |

## 🗄️ Database Schema

Generated from Supabase schema with multi-tenant support:

- `public` - System tables (companies, users, locations)
- `master_data` - Product catalog
- `inventory` - Stock management
- `sales` - Sales orders & customers
- `purchase` - Purchase orders & suppliers

## 🔐 RLS (Row Level Security)

All queries automatically filter by `company_id` using JWT claims:

```go
// Middleware injects company_id from JWT
// GORM queries automatically apply RLS
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## 🏗️ Architecture

### Multi-Tenancy
- Schema-based isolation with `company_id`
- JWT-based authentication (Supabase Auth)
- RLS policies enforce data isolation

### Clean Architecture
```
HTTP Handler → Service → Repository → Database
     ↑              ↑           ↑
   DTO          Business    GORM
   Validation   Logic       Models
```

## 📦 Dependencies

- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL (via Supabase)
- **Auth**: Supabase Auth (JWT)
- **UUID**: [google/uuid](https://github.com/google/uuid)
- **Config**: [godotenv](https://github.com/joho/godotenv)

## 📝 Notes for Developers

### Model Conventions
- Table names: `schema.table_name` (lowercase, plural)
- JSON tags: `snake_case`
- GORM tags: `column:name`
- UUID fields: `type:uuid;default:gen_random_uuid()`

### Adding New Models
1. Create model in `models/` directory
2. Add to `AutoMigrate` in `database/database.go`
3. Create repository in `repositories/`
4. Create service in `services/`
5. Create handler in `handlers/`
6. Add routes in `main.go` or route file

### Database Migrations
For production, use migration files instead of AutoMigrate:

```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate create -ext sql -dir migrations -seq create_products_table
```

## 🔗 Related Repositories

- Frontend: [wmsmicroservicefe](https://github.com/ikurniawann/wmsmicroservicefe)
- Schema: `../supabase_schema.sql`

---

**Maintained by**: Arie Anggono
**Project**: NEO WMS Refactor (Go + PostgreSQL + Next.js)
