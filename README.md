## Arama Motoru Servisi (Monorepo)

Ana hedef: Farklı sağlayıcılardan (JSON/XML) gelen içerikleri toplayıp tek bir standart modele dönüştürmek, puanlamak ve arama/sıralama kriterleriyle API üzerinden sunmak. Basit bir React dashboard ile sonuçları listelemek.

### Monorepo Yapısı
```
search_engine/
  backend/           # Go (Gin) API, DB, Redis, Docker + Migrations + Swagger
  frontend/          # Vite + React + Tailwind UI (dashboard)
  docker-compose.yml # API + DB + Redis + Frontend (dev)
  provider1.json     # JSON provider mock verisi
  provider2.xml      # XML provider mock verisi
```

### Özellikler
- **Gelişmiş Arama**: Anahtar kelimeye göre arama, içerik türü (video/metin) filtreleme
- **Akıllı Sıralama**: Popülerlik ve alakalılık skorlarına göre sıralama
- **Sayfalama**: Sayfa bazlı navigasyon (page/page_size)
- **Full-Text Search**: PostgreSQL FTS ile gelişmiş arama, fuzzy matching
- **Provider Entegrasyonu**: JSON/XML provider'lardan standart formata çevirme
- **Puanlama Algoritması**: İçerik türü ağırlıklandırma, güncellik ve etkileşim puanı
- **Caching**: Redis ile çok katmanlı cache sistemi, invalidation
- **Rate Limiting**: Redis tabanlı istek limiti yönetimi
- **Background Jobs**: Periyodik senkronizasyon ve skor yeniden hesaplama
- **Monitoring**: Detaylı health checks ve sistem metrikleri
- **Admin Dashboard**: Yönetim arayüzü ile sistem kontrolü
- **Modern UI**: React + TypeScript, search suggestions, skeleton loading

### Puanlama Formülü
Final Skor = (Temel Puan × İçerik Türü Katsayısı) + Güncellik Puanı + Etkileşim Puanı
- Temel Puan:
  - Video: `views / 1000 + (likes / 100)`
  - Metin: `reading_time + (reactions / 50)`
- İçerik Türü Katsayısı:
  - Video: `1.5`
  - Metin: `1.0`
- Güncellik Puanı:
  - 1 hafta: `+5`, 1 ay: `+3`, 3 ay: `+1`, daha eski: `+0`
- Etkileşim Puanı:
  - Video: `(likes / views) * 10`
  - Metin: `(reactions / reading_time) * 5`

---

## Hızlı Başlangıç (Geliştirme)
Önkoşullar: Docker ve Docker Compose

**⚠️ ÖNEMLİ: Go Version Kontrolü**
```bash
# Git hooks'u yükleyin (önerilir - Go version hatalarını önler)
make install-hooks

# Veya manuel kontrol
make check-go-version

# Eğer Go version yanlışsa düzeltin
make fix-go-version
```

1) Ortam değişkenlerini oluşturun:
```
cp backend/env.example backend/.env
```

2) Servisleri başlatın:
```
make up
```

3) (İlk kurulum) Veritabanı migration:
```
make migrate
```

4) Erişim adresleri:
- API: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/swagger/index.html`
- Frontend (Vite dev): `http://localhost:5173`

Yardımcı komutlar:
```
make logs             # Tüm servis logları
make api-logs         # API logları
make fe-logs          # Frontend logları
make down             # Servisleri durdur
make seed             # Örnek verilerle doldur
make check-go-version # Go version kontrolü
make fix-go-version   # Go version'ı 1.22'ye düzelt
make install-hooks    # Git hooks'u yükle
```

### Test Çalıştırma
Backend test'leri:
```bash
# Unit test'ler
make test-unit

# Integration test'ler (Docker gerekli)
make test-integration

# Tüm test'ler
make test

# Test coverage raporu
make test-coverage

# Test environment kontrolü
make test-env-up    # Test DB + Redis başlat
make test-env-down  # Test environment durdur
```

ENV ve servis konfigürasyonu: `backend/env.example`

---

## Örnek API Kullanımı
Arama:
```bash
curl "http://localhost:8080/api/v1/contents/search?q=python&page=1&page_size=20&type=video"
```
Detay:
```bash
curl "http://localhost:8080/api/v1/contents/123"
```
İstatistikler:
```bash
curl "http://localhost:8080/api/v1/contents/stats"
```

Admin (X-API-Key):
```bash
curl -H "X-API-Key: your-secret-key" -X POST "http://localhost:8080/api/v1/admin/sync"
```

Swagger UI üzerinden uçları keşfedebilirsiniz: `http://localhost:8080/swagger/index.html`

---

## Provider Mock’ları
- Yerel dosyalar docker-compose ile API konteynerine mount edilir:
  - `provider1.json` → JSON provider
  - `provider2.xml` → XML provider
- Önemli ENV’ler (bkz. `backend/env.example`):
  - `PROVIDER1_FILE_PATH=/app/data/providers/provider1.json`
  - `PROVIDER2_FILE_PATH=/app/data/providers/provider2.xml`
  - `PROVIDERS_FILE_ONLY=true` (dosya tabanlı mocku zorlar)
- Mock endpoint referansları:
  - `GET /mock/provider1/...`
  - `GET /mock/provider2/...`

---

## API Uçları (Özet)
- **Public Endpoints**
  - `GET /health` - Sistem durumu kontrolü
  - `GET /api/v1/contents/search` - İçerik arama (full-text search)
  - `GET /api/v1/contents/:id` - İçerik detayları (cached)
  - `GET /api/v1/contents/stats` - İstatistikler

- **Admin Endpoints** (header: `X-API-Key`)
  - `POST /api/v1/admin/sync` - Manuel senkronizasyon
  - `GET /api/v1/admin/sync/history` - Senkronizasyon geçmişi
  - `POST /api/v1/admin/scores/recalculate` - Skor yeniden hesaplama
  - `GET /api/v1/admin/providers` - Provider istatistikleri
  - `POST /api/v1/admin/providers/health-check` - Provider sağlık kontrolü
  - `DELETE /api/v1/admin/contents/:id` - İçerik soft delete
  - `GET /api/v1/admin/metrics/dashboard` - Dashboard metrikleri
  - `GET /api/v1/admin/jobs/:jobId` - Job durumu takibi
  - `GET /api/v1/admin/metrics/system` - Sistem metrikleri

- **Mock Endpoints** (test için)
  - `GET /mock/provider1/contents` - JSON provider mock
  - `GET /mock/provider2/feed` - XML provider mock

Yanıt formatı (örnek):
```json
{ "success": true, "data": {}, "error": { "code": "", "message": "" } }
```

---

## Teknoloji Tercihleri
- Backend: Go (Gin), PostgreSQL (pgx), Redis
- Frontend: React + Vite + Tailwind
- Docker Compose ile çoklu servis geliştirme ortamı

Kısa Mimari Notlar: Handler → Service → Repository → DB akışı, provider adaptörleri (JSON/XML) ve factory, Redis cache + rate limit, arkaplan işler (senkronizasyon ve skor yeniden hesaplama). Ayrıntı: `backend/ARCHITECTURE.md`

API şemaları ve örnekler: Swagger UI veya `backend/docs/openapi.yaml`

---

## Üretim Notu (Özet)
- Backend prod compose: `backend/docker-compose.production.yml` (API + Postgres + Redis + Nginx)
- Örnek image alma/derleme, güvenlik ve CI notları: `backend/DEPLOYMENT.md`

---

## Dağıtım (Detay)
- Image oluşturma:
```bash
docker build -t content-api:latest backend
```
- Production Compose:
```bash
docker compose -f backend/docker-compose.production.yml up -d
```
- Güvenlik:
  - `ADMIN_API_KEY` zorunlu
  - Nginx üzerinden SSL/TLS reverse proxy önerilir (`backend/deployment/nginx.conf`)
- CI/CD:
  - **GitHub Actions**: Otomatik test, lint, build ve deployment pipeline'ları
  - **Workflows**: `.github/workflows/` klasöründe tanımlı 5 workflow
    - `ci.yml` - Otomatik test ve build
    - `deploy.yml` - Production deployment
    - `pr-check.yml` - Pull request otomasyonu
    - `release.yml` - Version release yönetimi
    - `scheduled-tasks.yml` - Günlük maintenance işleri
  - **Dependabot**: Otomatik dependency güncellemeleri
  - Detaylı bilgi: `.github/workflows/README.md`

---

## SOLID ve Mimari İlkeler
- **Single Responsibility**: Handler/Service/Repository ayrı sorumluluklar; her paket odaklıdır.
- **Open/Closed**: Provider entegrasyonları `providers.factory` ile genişlemeye açık, var olan kodu değiştirmeye kapalıdır.
- **Liskov Substitution**: Provider arayüzleri aynı sözleşmeyi uygular; JSON/XML sağlayıcılar birbirlerinin yerine geçebilir.
- **Interface Segregation**: İnce arayüzler (ör. `ContentRepository`, `Provider`) tüketen katmanlara yalnızca gerekeni sunar.
- **Dependency Inversion**: Üst katmanlar soyutlamalara (interface) bağımlıdır; somut bağımlılıklar `config/factory` ile enjekte edilir.

## İleri Düzey Özellikler
- **Full-Text Search**: PostgreSQL tsvector ile gelişmiş arama, trigram similarity ile fuzzy matching
- **Caching Strategy**: Redis ile çok katmanlı cache (search results, content details, invalidation)
- **Error Handling**: Structured error responses, HTTP status code mapping, detailed logging
- **Monitoring**: Comprehensive health checks, system metrics, performance monitoring
- **Testing**: Integration tests with Docker containers, comprehensive test coverage
- **Frontend UX**: Search suggestions, recent searches, skeleton loading states, responsive design

---

## Sorun Giderme
- Port çakışmaları için: `.env` içindeki `API_PORT` değiştirilebilir.
- Frontend API tabanı: `VITE_API_BASE_URL` docker-compose ile `http://localhost:8080/api/v1`
- Migration/container sırası için: `make up` sonrası `make migrate` komutunu çalıştırın.

---

## Lisans
Bu depo teknik değerlendirme amacıyla hazırlanmıştır.

