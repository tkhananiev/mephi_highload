# go-microservice (High-load HW2)

## Run (Docker Compose)
```bash
docker compose up --build -d
```

Service:
- API: http://localhost:8080
- Metrics: http://localhost:8080/metrics
- MinIO UI: http://localhost:9001 (minioadmin/minioadmin)

## CRUD examples
Create:
```bash
curl -s -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Timur","email":"timur@example.com"}'
```

List:
```bash
curl -s http://localhost:8080/api/users
```

## MinIO integration demo
Upload a small audit snapshot to S3 (MinIO):
```bash
curl -s -X POST http://localhost:8080/api/integrations/audit/upload
```

## Load test (wrk)
```bash
wrk -t12 -c500 -d60s http://localhost:8080/api/users
```
