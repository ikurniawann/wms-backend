# ⚠️ DEPRECATED - Repository Archived

> **This repository is no longer maintained.**

## 🚀 Please Use

**Active Repository:** [wms-platform](https://github.com/ikurniawann/wms-platform)

```
https://github.com/ikurniawann/wms-platform
```

---

## Why Archived?

| Aspect | wms-backend (Old) | wms-platform (New) |
|--------|-------------------|-------------------|
| Architecture | Monolithic | **Modular Monolith** |
| Frontend | ❌ None | ✅ Next.js 14 |
| Docker | ❌ Single container | ✅ Full stack (5 services) |
| Scalability | Low | **High** |
| Code Quality | Medium | **Clean Architecture** |

---

## Migration Guide

### Old (wms-backend)
```bash
cd wms-backend
go run main.go
# Just backend, no frontend
```

### New (wms-platform)
```bash
cd wms-platform
docker-compose up -d
# Full stack: PostgreSQL + Redis + Backend + Worker + Frontend
```

---

## What Was Here?

This was the **Phase 3** implementation:
- Monolithic Go backend
- 7 handlers (product, order, inventory, customer, supplier, purchase, sales)
- 27 models
- Gin + GORM + PostgreSQL

**Evolution:** Monolithic → Modular Monolith → (Future: Microservices)

---

## Archive Date
📅 April 7, 2026

## Reference Only
This repo is kept for historical reference. All new development happens in [wms-platform](https://github.com/ikurniawann/wms-platform).
