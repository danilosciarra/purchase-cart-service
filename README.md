# Purchase Cart Service

REST service written in **Go** using **Gin** that implements a simplified **purchase cart / order service**.

The project was developed as a coding test and focuses on **clarity, separation of concerns, and clean architecture**, keeping the scope intentionally limited.

---

## Overview

The service exposes HTTP APIs to:
- create an order from a list of items
- return order ID, total price, total VAT, and per-item details
- browse a **preloaded, read-only product catalog**

---

## Quickstart

### Local
```bash
go run ./main.go
```

### Docker
```bash
docker build -t purchase-cart-service .
docker run -p 8080:8080 purchase-cart-service
```

---

## Swagger

Swagger documentation is generated using **Swaggo**.
After run service, swagger its present at: http://localhost:8080/swagger/index.html
---


## API (base path: `/api/v1`)

### Health check
```
GET /health
```

### Create order
```
POST /orders
```

Request:
```json
{
  "country_code": "it",
  "items": [
    { "product_id": "A123", "quantity": 2 }
  ]
}
```

Response (example):
```json
{
  "order_id": "uuid",
  "total_price": 24.40,
  "total_vat": 4.40,
  "items": [
    { "product_id": "A123", "name": "product name", "quantity": 2, "unit_price": 10.00, "vat": 4.40 }
  ]
}
```

### Products
- `GET /products` → list products
- `GET /products/:id` → product details

Products are loaded at startup and cannot be modified via API.

---

## Architecture

Layered architecture with clear separation of responsibilities:

Structure (simplified):
```
purchase-cart-service/
├─ main.go                     # Binary entrypoint: loads config, starts Server
├─ cmd/
│  └─ server/
│     └─ server.go             # HTTP Server: initializes router and registers handlers
├─ internal/
│  ├─ config/                  # Configuration loading/validation (ServiceName, WebApp, Database)
│  │  ├─ loader.go
│  │  └─ model.go
│  ├─ domain/                  # Domain model and business logic (Order, Item, VAT/total calculations)
│  │  ├─ order/
│  │  │  └─ service.go
│  │  ├─ product/              # Product model and logic (validation, transformations)
│  │  │  ├─ model.go
│  │  │  └─ service.go
│  │  └─ ...
│  ├─ api/
│  │  └─ http/                 # HTTP transport with Gin: router and handlers
│  │     ├─ router.go          # router definition and prefix registration (/ and /api/v1)
│  │     └─ handlers/
│  │        ├─ healthcheck.go  # health handler
│  │        ├─ order.go        # order handlers
│  │        └─ product.go      # product handlers (list, detail)
├─ repository/                 # runtime repositories (e.g., InMemory) for Order/VatRate
│  ├─ order_repository.go
│  ├─ vat_rate_repository.go
│  └─ product_repository.go    # product storage (preloaded at startup, read-only via API)
├─ docs                        # Generated Swagger (Swaggo)
├─ config.json                 # Runtime configuration (ServiceName, WebApp, Database)
└─ go.mod / go.sum
```

Layers and responsibilities:
- main.go: process bootstrap; loads config, builds Server, starts listening.
- cmd/server: HTTP Server component; depends on config, router, and domain services.
- internal/api/http (Gin): exposes APIs, validates input, maps DTO ⇄ domain, handles errors.
- internal/service / internal/domain: 
  - order: orchestrates use cases and contains pure logic (totals/VAT calculations, order creation).
  - product: product catalog management (runtime validations; no creation via API).
- internal/config: configuration models and loading.
- repository: storage interfaces/implementations (InMemory/DB).
  - order_repository: order persistence.
  - vat_rate_repository: VAT rate source.
  - product_repository: product persistence/lookup; the catalog is loaded at startup and does not support inserts via API.
- docs: generated Swagger files.

---

## Configuration

The service reads configuration from `config.json` (project root). Supported fields:
- `ServiceName`: service name.
- `WebApp.Hostname`: HTTP server bind address (e.g., `0.0.0.0`).
- `WebApp.Port`: HTTP server port (e.g., `8080`).
- `Database`: persistence configuration.
  - `Type`: storage type (e.g., `InMemory`).
  - `Host`, `Port`, `User`, `Password`, `Name`: DB parameters (used if `Type` is not `InMemory`).

Example:
```json
{
  "ServiceName": "purchase-cart",
  "WebApp": {
    "Hostname": "0.0.0.0",
    "Port": 8080
  },
  "Database": {
    "Type": "InMemory",
    "Host": "localhost",
    "Port": 5432,
    "User": "db",
    "Password": "password",
    "Name": "purchase_cart_db"
  }
}
```

Notes:
- With `Database.Type = "InMemory"` DB parameters can be ignored.

---
