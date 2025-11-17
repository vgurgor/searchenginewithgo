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

## 🚀 Kurulum ve Çalıştırma

### Gereksinimler

**Zorunlu:**
- **Docker** (v20.10+) ve **Docker Compose** (v2.0+)
- **Git**

**Opsiyonel (Geliştirme için):**
- **Go** 1.24+ (local development için)
- **Node.js** 20+ (frontend geliştirme için)
- **Make** (komut kısayolları için)

### 1️⃣ Projeyi İndirin

```bash
git clone <repository-url>
cd search_engine
```

### 2️⃣ Ortam Değişkenlerini Ayarlayın

```bash
# Backend .env dosyasını oluşturun
cp backend/env.example backend/.env

# Varsayılan ayarlar çoğu durumda yeterlidir
# Gerekirse düzenleyin:
# - API_PORT=8080
# - ADMIN_API_KEY=your-secret-key
# - DATABASE_URL, REDIS_URL, vb.
```

### 3️⃣ Servisleri Başlatın

```bash
# Docker Compose ile tüm servisleri başlat
make up

# Alternatif (Make yoksa):
docker compose up -d
```

Bu komut şunları başlatır:
- ✅ PostgreSQL (port 5432)
- ✅ Redis (port 6379)
- ✅ Backend API (port 8080)
- ✅ Frontend (port 5173)

### 4️⃣ Veritabanı Migration (İlk Kurulumda)

```bash
# Migration'ları çalıştır
make migrate

# Alternatif:
docker compose run --rm migrate up

# (Opsiyonel) Örnek verilerle doldur
make seed
```

### 5️⃣ Erişim ve Kullanım

- 🌐 **API**: http://localhost:8080
- 📚 **Swagger UI**: http://localhost:8080/swagger/index.html
- 🎨 **Frontend Dashboard**: http://localhost:5173
- 💚 **Health Check**: http://localhost:8080/health

### 🛠️ Yaygın Komutlar

```bash
# Servisleri görüntüle
make logs              # Tüm servis logları
make api-logs          # Sadece API logları
make fe-logs           # Sadece Frontend logları

# Servisleri yönet
make down              # Servisleri durdur
make restart           # Yeniden başlat
make ps                # Çalışan container'ları listele

# Temizlik
make clean             # Build artifacts temizle
make clean-all         # Tüm volumes dahil temizle

# Veritabanı
make migrate           # Migration'ları çalıştır
make migrate-down      # Migration'ları geri al
make seed              # Örnek veri yükle

# Geliştirme
make check-go-version  # Go version kontrol
make install-hooks     # Git hooks yükle
```

---

## 🧪 Test Çalıştırma

### Backend Testleri

```bash
# Unit testler (hızlı)
cd backend
make test-unit

# Integration testler (Docker gerekli)
make test-integration

# Tüm testler
make test

# Test coverage raporu
make test-coverage
# Sonuç: coverage.html dosyası oluşur
```

### Frontend Testleri

```bash
cd frontend
npm test
```

### Docker Test Environment

```bash
# Test için ayrı DB + Redis başlat
cd backend
make test-env-up

# Testleri çalıştır
make test-integration

# Test environment'ı durdur
make test-env-down
```

---

## 💻 Local Development (Docker Olmadan)

### Backend

```bash
cd backend

# Dependencies
go mod download

# PostgreSQL ve Redis gerekli (Docker ile):
docker compose up -d postgres redis

# Migration
make migrate

# API'yi local çalıştır
go run ./cmd/api

# Veya build edip çalıştır
make build
./bin/api
```

### Frontend

```bash
cd frontend

# Dependencies
npm install

# Development server
npm run dev

# Build
npm run build
```

### Swagger Docs Güncelleme

```bash
cd backend

# Swagger docs'u yeniden oluştur
make swagger

# swag CLI gerekli (yüklü değilse):
go install github.com/swaggo/swag/cmd/swag@latest
```

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

## 🔧 Sorun Giderme

### Port Çakışması

```bash
# .env dosyasında portları değiştirin
API_PORT=8081           # Backend
POSTGRES_PORT=5433      # PostgreSQL
REDIS_PORT=6380         # Redis

# Frontend için (frontend/.env)
VITE_API_PORT=8081
```

### Docker Container'lar Başlamıyor

```bash
# Mevcut container'ları temizle
make down
docker system prune -f

# Yeniden başlat
make up
```

### Migration Hataları

```bash
# Migration durumunu kontrol et
docker compose run --rm migrate version

# Migration'ları sıfırla (DİKKAT: Veri kaybı!)
make migrate-down
make migrate

# Manuel migration (PostgreSQL container'ı içinde)
docker compose exec postgres psql -U postgres -d searchdb
```

### API Erişim Sorunları

```bash
# Health check
curl http://localhost:8080/health

# Container loglarını kontrol et
make api-logs

# Database bağlantısını test et
docker compose exec postgres psql -U postgres -d searchdb -c "SELECT 1;"
```

### Go Version Sorunları

```bash
# Go version kontrol
cd backend
make check-go-version

# go.mod'da Go 1.24 olmalı
grep "^go " go.mod
# Çıktı: go 1.24.0

# Eğer farklıysa dependencies'i yenileyin
go mod tidy
```

### Frontend Build/Dev Sorunları

```bash
cd frontend

# node_modules temizle ve yeniden yükle
rm -rf node_modules package-lock.json
npm install

# Cache temizle
npm cache clean --force

# Dev server'ı yeniden başlat
npm run dev
```

### Test Container'ları Temizleme

```bash
# Eski test container'larını temizle
docker ps -a | grep testcontainers | awk '{print $1}' | xargs docker rm -f

# Test environment'ı temizle
cd backend
make test-env-down
docker volume prune -f
```

### Veritabanı Bağlantı Sorunları

```bash
# PostgreSQL container'ının çalıştığını kontrol et
docker compose ps postgres

# PostgreSQL loglarını kontrol et
docker compose logs postgres

# Elle bağlantı test et
docker compose exec postgres psql -U postgres -d searchdb -c "\dt"
```

---

## Lisans
Bu depo teknik değerlendirme amacıyla hazırlanmıştır.

