# Purchase Cart Service

Servizio REST in **Go/Gin** per creare ordini a partire da una lista di prodotti. Orientato a semplicità e separazione delle responsabilità.

## Panoramica

Il servizio espone API HTTP per:
- creare un ordine dato un insieme di items
- restituire ID ordine, prezzo totale, IVA totale e dettaglio per riga
- consultare un catalogo prodotti precaricato (non modificabile via API)

## Quickstart

- Avvio locale:
```bash
go run ./main.go
```
- Docker:
```bash
docker build -t purchase-cart-service .
docker run -p 8080:8080 purchase-cart-service
```

## API (base path: /api/v1)

- Health Check
  - `GET /api/v1/health` → `200 OK` + `ok`
- Create Order
  - `POST /api/v1/orders`
  - Request (decimali):
    ```json
    {
      "items": [
        { "product_id": "A123", "quantity": 2, "unit_price": 10.00 }
      ]
    }
    ```
  - Response (IVA 22%):
    ```json
    {
      "order_id": "uuid",
      "total_price": 24.40,
      "total_vat": 4.40,
      "items": [
        { "product_id": "A123", "quantity": 2, "unit_price": 10.00, "vat": 4.40 }
      ]
    }
    ```

- Products
  - Catalogo prodotti precaricato all’avvio: non è possibile inserire o modificare prodotti via API.
  - `GET /api/v1/products` → lista prodotti
  - `GET /api/v1/products/:id` → dettaglio prodotto per ID
  // Nota: l’endpoint `POST /api/v1/products` non è disponibile in questo progetto.

## Architettura

Architettura a layer, organizzata in sottocartelle per separare le responsabilità.

Struttura (semplificata):
```
purchase-cart-service/
├─ main.go                     # Entrypoint del binario: carica config, avvia Server
├─ cmd/
│  └─ server/
│     └─ server.go             # Server HTTP: inizializza router e registra handler
├─ internal/
│  ├─ config/                  # Lettura/validazione configurazione (ServiceName, WebApp, Database)
│  │  ├─ loader.go
│  │  └─ model.go
│  ├─ domain/                  # Modello di dominio e logiche (Order, Item, calcoli IVA/totali)
│  │  ├─ order/
│  │  │  └─ service.go
│  │  ├─ product/              # Modello e logiche di prodotto (validazione, trasformazioni)
│  │  │  ├─ model.go
│  │  │  └─ service.go
│  │  └─ ...
│  ├─ api/
│  │  └─ http/                 # Transport HTTP con Gin: router e handler
│  │     ├─ router.go          # definizione router e registrazione prefissi (/ e /api/v1)
│  │     └─ handlers/
│  │        ├─ healthcheck.go  # handler health
│  │        ├─ order.go        # handler ordini (lista, dettaglio, creazione)
│  │        └─ product.go      # handler prodotti (lista, dettaglio)
├─ repository/                 # repository runtime (es. InMemory) per Order/VatRate
│  ├─ order_repository.go
│  ├─ vat_rate_repository.go
│  └─ product_repository.go    # storage prodotti (precaricati all’avvio, sola lettura via API)
├─ docs/                       # Swagger generato (Swaggo)
├─ config.json                 # Configurazione runtime (ServiceName, WebApp, Database)
└─ go.mod / go.sum
```

Layer e responsabilità:
- main.go: bootstrap del processo; carica config, costruisce Server e ne avvia l’ascolto.
- cmd/server: componente Server HTTP; dipende da config, router e servizi di dominio.
- internal/api/http (Gin): espone le API, valida input, traduce DTO ⇄ dominio, gestisce errori.
- internal/service / internal/domain: 
  - order: coordinano i casi d’uso e contengono le logiche pure (calcoli totali/IVA, creazione ordini).
  - product: gestione del catalogo prodotti (CRUD lato runtime, validazioni).
- internal/config: modelli e caricamento della configurazione.
- repository: interfacce/implementazioni di storage (InMemory/DB).
  - order_repository: persistenza ordini.
  - vat_rate_repository: fonte aliquote IVA.
  - product_repository: persistenza/lookup prodotti; il catalogo è caricato all’avvio e non supporta inserimenti via API.
- docs: file Swagger generati.

## Convenzioni sugli importi

- Importi in formato decimale.
- `total_price` = totale netto + totale IVA.
- `items[].vat` = IVA del totale riga.
- Aliquota degli esempi: 22%.

## Documentazione API

Lo Swagger delle API è generato con **Swaggo**.

## Configurazione

Il servizio legge la configurazione da `config.json` (root del progetto). Campi supportati:
- `ServiceName`: nome del servizio.
- `WebApp.Hostname`: indirizzo di bind del server HTTP (es. `0.0.0.0`).
- `WebApp.Port`: porta del server HTTP (es. `8080`).
- `Database`: configurazione della persistenza.
  - `Type`: tipo di storage (es. `InMemory`).
  - `Host`, `Port`, `User`, `Password`, `Name`: parametri del DB (utilizzati se `Type` non è `InMemory`).

Esempio:
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

Note:
- Con `Database.Type = "InMemory"` i parametri del DB possono essere ignorati.
