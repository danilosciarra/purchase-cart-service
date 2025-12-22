# Purchase Cart Service

REST service in **Go/Gin** to create orders from a list of products. Focused on simplicity and separation of concerns.

## Overview

The service exposes HTTP APIs to:
- create an order from a set of items
- return order ID, total price, total VAT, and per-line details
- browse a preloaded product catalog (not editable via API)

## Quickstart

- Local run:
```bash
go run ./main.go
```
- Docker:
```bash
docker build -t purchase-cart-service .
docker run -p 8080:8080 purchase-cart-service
```

## APIs (base path: /api/v1)

- Health Check
  - `GET /api/v1/health` → `200 OK` + `ok`
- Create Order
  - `POST /api/v1/orders`
  - Request (decimal amounts):
    ```json
    {
      "country_code":"it",
      "items": [
        { "product_id": "A123", "quantity": 2, }
      ]
    }
    ```
  - Response (22% VAT):
    ```json
    {
      "order_id": "uuid",
      "total_price": 24.40,
      "total_vat": 4.40,
      "items": [
        { "product_id": "A123","name":"product name", "quantity": 2, "unit_price": 10.00, "vat": 4.40 }
      ]
    }
    ```

- Products
  - Product catalog is preloaded at startup: products cannot be created or modified via API.
  - `GET /api/v1/products` → product list
  - `GET /api/v1/products/:id` → product details by ID
  // Note: `POST /api/v1/products` is not available in this project.

## Architecture

Layered architecture, organized in subfolders to separate responsibilities.

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

## Amount conventions

- Decimal amounts.
- `total_price` = net total + total VAT.
- `items[].vat` = VAT of the line total.
- Example VAT rate: 22%.

## API documentation

API Swagger is generated using **Swaggo**.

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
