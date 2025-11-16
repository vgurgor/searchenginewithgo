# 🚀 CI/CD Pipeline Kurulum ve Kullanım Kılavuzu

## ⚠️ KRİTİK: Go Version Sorunu

### 🔴 Problem
`backend/go.mod` dosyasında Go version **sürekli 1.24.0'a değişiyor**. Bu şu sorunlara sebep oluyor:

- ❌ Docker build başarısız oluyor
- ❌ CI/CD pipeline başarısız oluyor
- ❌ Production deployment imkansız oluyor

**Sebep:** Go 1.24.0 henüz yayınlanmadı ama bazı komutlar (`go mod tidy`, `go get`) otomatik olarak yükseltiyor.

### ✅ ÇÖZÜM: 3 Katmanlı Koruma

#### 1️⃣ Git Hook (Local Koruma)
Commit öncesi otomatik kontrol:

```bash
# Proje root'da
make install-hooks
```

Bu hook her commit öncesi `backend/go.mod`'u kontrol eder ve Go 1.24.0 varsa commit'i engeller.

#### 2️⃣ CI Pipeline Kontrolü (Remote Koruma)
GitHub Actions'da ilk iş olarak Go version kontrol edilir:

```yaml
jobs:
  check-go-version:
    name: Verify Go Version
    # go.mod'da 1.24.0 varsa pipeline başarısız olur
```

#### 3️⃣ go.mod Dosyası Uyarısı
`backend/go.mod` başında şu uyarı mevcut:

```go
// IMPORTANT: Keep go version at 1.22 for Dockerfile compatibility
// DO NOT change to 1.24.0 - it's not released yet and will break Docker build
// See backend/GO_VERSION.md for details
```

### 🛠️ Düzeltme Komutları

#### Manuel Kontrol
```bash
make check-go-version
```

#### Otomatik Düzeltme
```bash
make fix-go-version
```

#### Backend'de Direkt
```bash
cd backend
make check-go-version
make fix-go-version
```

### 📋 Version Uyumu Tablosu

**MUTLAKA TÜM YERLERDEKİ VERSION 1.22 OLMALI:**

| Bileşen | Dosya | Satır | Değer |
|---------|-------|-------|-------|
| go.mod | `backend/go.mod` | 6 | `go 1.22` |
| Dockerfile | `backend/Dockerfile` | 2 | `golang:1.22-alpine` |
| CI Env | `.github/workflows/ci.yml` | 10 | `GO_VERSION: '1.22'` |
| Release | `.github/workflows/release.yml` | 30 | `go-version: '1.22'` |
| Scheduled | `.github/workflows/scheduled-tasks.yml` | 21, 109 | `go-version: '1.22'` |

---

## 📦 Oluşturulan CI/CD Dosyaları

### Workflow Dosyaları (5 adet)

1. **`.github/workflows/ci.yml`** - Ana CI Pipeline
   - Backend/Frontend testler
   - Linting
   - Docker build
   - Security scanning
   
2. **`.github/workflows/deploy.yml`** - Production Deployment
   - Docker image push
   - Server deployment
   - Health checks
   
3. **`.github/workflows/pr-check.yml`** - PR Otomasyonu
   - PR başlık kontrolü
   - Auto-labeling
   - Security scanning
   
4. **`.github/workflows/release.yml`** - Release Yönetimi
   - Multi-platform binary build
   - Changelog oluşturma
   - GitHub Release
   
5. **`.github/workflows/scheduled-tasks.yml`** - Günlük Bakım
   - Dependency checks
   - Security audit
   - Image cleanup

### Yapılandırma Dosyaları

6. **`.github/dependabot.yml`** - Otomatik dependency güncellemeleri
7. **`.github/labeler.yml`** - PR auto-labeling kuralları
8. **`.github/CODEOWNERS`** - Code review otomasyonu
9. **`.github/PULL_REQUEST_TEMPLATE.md`** - PR şablonu
10. **`backend/.golangci.yml`** - Go linting yapılandırması

### Dokümantasyon

11. **`.github/workflows/README.md`** - Pipeline dokümantasyonu
12. **`.github/workflows/CHANGELOG.md`** - Değişiklik geçmişi
13. **`backend/.golangci-migration.md`** - Linter migration kılavuzu
14. **`backend/GO_VERSION.md`** - Go version politikası
15. **`.github/hooks/README.md`** - Git hooks dokümantasyonu
16. **`.github/CI_CD_SETUP.md`** - Bu dosya

### Git Hooks

17. **`.github/hooks/pre-commit`** - Go version kontrolü

---

## 🔐 Gerekli GitHub Secrets

Repository Settings → Secrets and variables → Actions:

### Deployment (Production için gerekli)
```bash
DEPLOY_HOST        # Sunucu adresi
DEPLOY_USER        # SSH kullanıcı
DEPLOY_SSH_KEY     # SSH private key
DEPLOY_PORT        # SSH port (varsayılan: 22)
DEPLOY_URL         # Health check URL
```

### İsteğe Bağlı
```bash
SLACK_WEBHOOK      # Slack bildirimleri
CODECOV_TOKEN      # Code coverage tracking
```

---

## 🚀 Kullanım

### İlk Kurulum

```bash
# 1. Hooks'u yükle
make install-hooks

# 2. Go version'ı kontrol et
make check-go-version

# 3. Servisleri başlat
make up
```

### Geliştirme Akışı

```bash
# 1. Feature branch oluştur
git checkout -b feature/yeni-ozellik

# 2. Değişiklik yap

# 3. Go version kontrolü (otomatik pre-commit hook ile)
git commit -m "feat: yeni özellik"

# 4. Push
git push origin feature/yeni-ozellik

# 5. PR aç - Otomatik checks çalışır
```

### Commit Öncesi Checklist

- [ ] `make check-go-version` başarılı
- [ ] `make test` başarılı
- [ ] `make build` başarılı
- [ ] Linter hataları düzeltildi

---

## 🐛 Sık Karşılaşılan Sorunlar

### ❌ Docker Build Hatası

**Hata:**
```
go: go.mod requires go >= 1.24.0 (running go 1.22.12)
ERROR: process "/bin/sh -c go mod download" did not complete successfully
```

**Çözüm:**
```bash
make fix-go-version
git add backend/go.mod backend/go.sum
git commit --amend --no-edit
```

### ❌ CI Pipeline "check-go-version" Failed

**Sebep:** `backend/go.mod` dosyasında Go 1.24.0 var.

**Çözüm:**
```bash
# Local'de düzelt
make fix-go-version

# Commit et
git add backend/go.mod
git commit -m "fix: correct Go version to 1.22"
git push
```

### ❌ golangci-lint "no go files to analyze"

**Sebep:** Dependencies yüklenmemiş.

**Çözüm:** CI'da otomatik olarak `go mod download` çalışıyor. Local'de:
```bash
cd backend
go mod download
golangci-lint run
```

---

## 📊 CI/CD Status

Pipeline'ların durumunu görmek için:

**GitHub Repository → Actions Tab**

### Pipeline Başarı Kriterleri

✅ **check-go-version** - Go 1.22 kontrolü
✅ **backend-test** - Tüm testler geçiyor
✅ **backend-lint** - Linter başarılı
✅ **frontend-test** - Build başarılı
✅ **docker-build** - Image oluşuyor
✅ **security-scan** - Güvenlik taraması temiz

---

## 🔄 Güncelleme Geçmişi

| Tarih | Değişiklik | Durum |
|-------|-----------|-------|
| 2025-11-16 | İlk oluşturma | ✅ |
| 2025-11-16 | Go version koruma eklendi | ✅ |
| 2025-11-16 | Tüm deprecated actions güncellendi | ✅ |
| 2025-11-16 | golangci-lint v6 migration | ✅ |
| 2025-11-16 | Test'ler düzeltildi | ✅ |

---

## 💡 Öneriler

1. **Her zaman hooks kullanın**: `make install-hooks`
2. **Commit öncesi kontrol**: `make check-go-version`
3. **Test edin**: `make test`
4. **Build test**: `cd backend && go build ./cmd/api`
5. **Docker test**: `make build`

---

## 📞 Yardım

Sorun yaşarsanız:

1. **Go Version:** `backend/GO_VERSION.md`
2. **Pipeline:** `.github/workflows/README.md`
3. **Changelog:** `.github/workflows/CHANGELOG.md`
4. **Linter Migration:** `backend/.golangci-migration.md`

## 🎯 Sonuç

CI/CD pipeline tamamen hazır! Ancak **Go 1.22 versiyonunu korumak kritik önemde**.

Hooks kurulu olduğunda, yanlış version ile commit yapamazsınız. 🔒

