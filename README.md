# Purchase Cart Service

## Documentazione Swagger

Per generare la documentazione Swagger:

```sh
go install github.com/swaggo/swag/cmd/swag@latest
swag init --generalInfo internal/api/http/router.go --output docs
```

La documentazione sarà disponibile su [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) dopo aver avviato il servizio.
