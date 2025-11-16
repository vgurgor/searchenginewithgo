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

### Docker Compose
- Hot-reload: `api` servisi `air` ile çalışır, volume mapping aktiftir.
- Persistans: PostgreSQL için `db_data` volume tanımlıdır.

### Ortam Değişkenleri
- API_PORT, LOG_LEVEL, DATABASE_URL, REDIS_URL
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


