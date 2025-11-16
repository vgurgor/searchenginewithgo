# Architecture

Yüksek seviye bileşenler:
- API Layer (Gin)
- Service Layer (iş mantığı, orchestration)
- Repository Layer (PostgreSQL, pgx)
- Provider Integration (Adapters + Factory)
- Background Jobs (Ticker ve hafif Job Manager)
- Cache/RateLimit (Redis)

Veri akışı: Request -> Handler -> Service -> Repository -> DB. Provider senkronizasyonu job/servis ile tetiklenir.


