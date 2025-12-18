# Purchase Cart Service

## Descrizione

Microservizio Go per la gestione del carrello acquisti e degli ordini. Espone API REST documentate tramite Swagger e utilizza Gin come router HTTP.

---

## Architettura

- **Linguaggio:** Go
- **Framework HTTP:** Gin
- **Documentazione API:** Swagger (swaggo)
- **Struttura a livelli:** separazione tra handler HTTP, dominio (business logic) e infrastruttura.

### Struttura delle cartelle

```
purchase-cart-service/
├── cmd/                  # Entrypoint dell'applicazione (main.go)
├── internal/
│   └── api/
│       └── http/
│           ├── handlers/ # Handler HTTP (controller)
│           └── router.go # Configurazione router e Swagger
│   └── domain/
│       └── order/        # Logica di dominio ordini
├── docs/                 # Documentazione Swagger generata
├── go.mod
├── go.sum
├── README.md
```

---

## Dipendenze principali

- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [swaggo/swag](https://github.com/swaggo/swag)
- [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)
- [swaggo/files](https://github.com/swaggo/files)

---

## Avvio rapido

### 1. Clona il repository

```sh
git clone https://github.com/danilosciarra/purchase-cart-service.git
cd purchase-cart-service
```

### 2. Installa le dipendenze

```sh
go mod tidy
```

### 3. Genera la documentazione Swagger

```sh
swag init --generalInfo internal/api/http/router.go --output docs
```

### 4. Avvia il servizio

```sh
go run ./cmd/server/main.go
```

### 5. Accedi alla documentazione Swagger

Visita [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## Testing

Per eseguire i test:

```sh
go test ./...
```

---

## Note architetturali

- **Separation of Concerns:**  
  - Gli handler HTTP sono separati dalla logica di dominio.
  - La logica di dominio non dipende dall'infrastruttura.
- **Swagger:**  
  - I commenti sopra i metodi handler e le strutture dati generano la documentazione OpenAPI.
- **Estendibilità:**  
  - La struttura a moduli facilita l'aggiunta di nuove funzionalità (es. nuovi handler, servizi di dominio).
- **Error Handling:**  
  - Gli errori sono gestiti tramite risposte JSON standardizzate.

---

## Considerazioni progettuali

* Lo storage è **in-memory** per semplicità, ma astratto tramite interfaccia per una futura persistenza (es. database).
* L’API è pensata per essere facilmente versionabile (`/v1/orders`).

---