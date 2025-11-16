## Search Engine - Backend Altyapı (Go + Gin)

### Gereksinimler
- Docker ve Docker Compose

### Hızlı Başlangıç
1) Ortam değişkenlerini hazırlayın:
   - `cp env.example .env`
2) Servisleri ayağa kaldırın:
   - `docker compose up --build`
3) (İlk defa) veritabanı migration çalıştırın:
   - `bash scripts/migrate.sh`

### Servisler
- API: `http://localhost:8080`
  - Health: `GET /health`
  - Swagger UI: `http://localhost:8080/swagger/index.html`
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`

### Provider Mock API'leri
- JSON Provider (Provider1): `GET /mock/provider1/contents?limit=20&offset=0`
- XML Provider (Provider2): `GET /mock/provider2/feed?page=1&size=20`
- ENV:
  - `PROVIDER1_BASE_URL=http://localhost:8080/mock/provider1`
  - `PROVIDER2_BASE_URL=http://localhost:8080/mock/provider2`
  - `PROVIDER_TIMEOUT=10s`
  - `RATE_LIMIT_ENABLED=true`

### Docker Compose
- Hot-reload: `api` servisi `air` ile çalışır, volume mapping aktiftir.
- Persistans: PostgreSQL için `db_data` volume tanımlıdır.

### Ortam Değişkenleri
- API_PORT, LOG_LEVEL, DATABASE_URL, REDIS_URL
- Scoring: `SCORE_RECALCULATION_ENABLED`, `SCORE_RECALCULATION_INTERVAL`, `SCORE_BATCH_SIZE`, `VIDEO_TYPE_MULTIPLIER`, `TEXT_TYPE_MULTIPLIER`, `FRESHNESS_1_WEEK`, `FRESHNESS_1_MONTH`, `FRESHNESS_3_MONTHS`
- Sync: `CONTENT_SYNC_ENABLED`, `CONTENT_SYNC_INTERVAL`, `CONTENT_SYNC_RETRY_COUNT`, `CONTENT_SYNC_RETRY_DELAY`, `METRICS_CHANGE_THRESHOLD_PERCENT`, `METRICS_CHANGE_THRESHOLD_ABS_VIEWS`, `METRICS_CHANGE_THRESHOLD_ABS_LIKES`, `METRICS_CHANGE_THRESHOLD_ABS_REACTIONS`, `ADMIN_API_KEY`
- Pagination: `DEFAULT_PAGE_SIZE`, `MAX_PAGE_SIZE`
- Örnekler için `env.example` dosyasına bakın.

### Migration
- Araç: `migrate/migrate` (Docker)
- Komutlar:
  - `docker compose run --rm migrate up`
  - `docker compose run --rm migrate down`
- Dosyalar: `scripts/migrations`

### Geliştirme
- Loglar zap ile yapılandırılmış ve console'a structured formatta yazılır.
- Global error middleware entegredir.

### Kabul Kriterleri
- `docker compose up` ile tüm servisler ayağa kalkar.
- `GET /health` 200 döner (DB/Redis durumu JSON içinde).
- DB bağlantısı ve migration çalışır.
- Swagger UI açılır (gin-swagger) ve boş döküman gösterir.
- Kurulum adımları bu dosyadadır.
- Loglar düzenli formatta console'a yazılır.


