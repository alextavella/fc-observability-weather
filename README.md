# Fullcycle | Observability

## Getting started

```bash
go run cmd/api/main.go
```

## Docker

```bash
docker compose up -d --build
```

## Tests

```bash
curl -X POST http://localhost:8080/zipcode \
-H "Content-Type: application/json" \
-d '{"cep": "09861160"}'
```
