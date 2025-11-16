# Deployment Guide

## Production Build
```bash
docker build -t content-api:latest .
```

## Compose (production)
`docker-compose.production.yml` ile API, Postgres, Redis ve Nginx ayağa kaldırılır.

## CI/CD
`.github/workflows/ci.yml` ile test, lint ve image publish adımları tanımlıdır.

## Güvenlik
- `ADMIN_API_KEY` zorunlu
- SSL/TLS reverse proxy (Nginx) önerilir


