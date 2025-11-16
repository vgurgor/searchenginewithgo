# API Documentation

Bu dosya, Public ve Admin API uçlarını, istek/yanıt formatlarını ve hata yapısını özetler. Swagger UI: /swagger

## Authentication
- Public API: auth gerektirmez
- Admin API: `X-API-Key` header zorunlu

## Response Format
```json
{ "success": true, "data": {}, "error": { "code": "", "message": "" } }
```

## Public API
- GET `/api/v1/contents/search`
- GET `/api/v1/contents/:id`
- GET `/api/v1/contents/stats`

Örnek:
```bash
curl "http://localhost:8080/api/v1/contents/search?q=test&page=1&page_size=20"
```

## Admin API (X-API-Key)
- POST `/api/v1/admin/sync`
- GET `/api/v1/admin/sync/history`
- POST `/api/v1/admin/scores/recalculate`
- GET `/api/v1/admin/providers`
- POST `/api/v1/admin/providers/health-check`
- DELETE `/api/v1/admin/contents/:id`
- GET `/api/v1/admin/metrics/dashboard`
- GET `/api/v1/admin/jobs/:jobId`


