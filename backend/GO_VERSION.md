# ⚠️ ÖNEMLİ: Go Version Politikası

## 🔒 Go 1.22 Kullanılmalı

Bu proje **Go 1.22** ile geliştirilmiştir ve tüm sistem bu versiyona göre yapılandırılmıştır.

### ❌ Go 1.24.0 Kullanmayın!

**Sebep:**
1. Go 1.24.0 henüz resmi olarak yayınlanmadı (gelecek versiyon)
2. Dockerfile `golang:1.22-alpine` kullanıyor
3. GitHub Actions `GO_VERSION: '1.22'` ayarlı
4. Bazı dependencies Go 1.24+ gerektiriyor ve build başarısız oluyor

### 🐛 go.mod Go Version'ı Neden Değişiyor?

`go mod tidy` veya `go get` komutları bazen Go version'ı otomatik olarak güncelleyebilir. Bu sorunu önlemek için:

**1. go.mod başına yorum eklendi:**
```go
// IMPORTANT: Keep go version at 1.22 for Dockerfile compatibility
// DO NOT change to 1.24.0 - it's not released yet and will break Docker build
module search_engine

go 1.22
```

**2. Dockerfile'da GOTOOLCHAIN=auto erkene alındı:**
```dockerfile
FROM golang:1.22-alpine AS base
WORKDIR /app
ENV GOTOOLCHAIN=auto  # ✅ En başta!
RUN apk add --no-cache git ca-certificates build-base
COPY go.mod go.sum ./
RUN go mod download  # GOTOOLCHAIN=auto sayesinde Go 1.23+ modüller çalışır
```

**3. Otomatik düzeltme komutu:**
```bash
make fix-go-version  # go mod edit -go=1.22 kullanır
```

### ✅ Doğru Kurulum

```bash
# Backend dizininde
cd backend

# go.mod'u kontrol edin
grep "^go " go.mod
# Çıktı: go 1.22

# Eğer 1.24.0 ise düzeltin:
# 1. go.mod'u açın
# 2. "go 1.24.0" → "go 1.22" değiştirin
# 3. Kaydedin

# Dependencies'i yükleyin
go mod download
go mod verify

# Build test edin
go build -o /dev/null ./cmd/api
```

### 🔧 go mod tidy Kullanırken

`go mod tidy` çalıştırırken Go version değişirse:

```bash
# 1. go mod tidy çalıştırın
go mod tidy

# 2. go.mod'u kontrol edin
grep "^go " go.mod

# 3. Eğer 1.24.0 ise düzeltin
sed -i '' 's/^go 1\.24\.0/go 1.22/' go.mod

# veya manuel olarak düzeltin

# 4. Tekrar tidy çalıştırın
go mod tidy
```

### 📊 Version Uyumu Tablosu

| Bileşen | Go Version | Dosya |
|---------|-----------|-------|
| go.mod | **1.22** | `backend/go.mod` |
| Dockerfile | **golang:1.22-alpine** | `backend/Dockerfile` |
| CI/CD | **GO_VERSION: '1.22'** | `.github/workflows/ci.yml` |
| Release | **go-version: '1.22'** | `.github/workflows/release.yml` |
| Scheduled | **go-version: '1.22'** | `.github/workflows/scheduled-tasks.yml` |

**HEPSI 1.22 OLMALI!**

### 🚨 Docker Build Hatası

Eğer bu hatayı alırsanız:
```
go: go.mod requires go >= 1.24.0 (running go 1.22.12; GOTOOLCHAIN=local)
ERROR: process "/bin/sh -c go mod download" did not complete successfully
```

**Çözüm:** `backend/go.mod` dosyasında Go version'ı 1.22'ye değiştirin!

### 💡 Go 1.24'e Geçmek İçin

Go 1.24.0 resmi olarak yayınlandığında (muhtemelen 2025 Q1-Q2):

1. Dockerfile güncelleyin:
   ```dockerfile
   FROM golang:1.24-alpine AS base
   ```

2. GitHub Actions güncelleyin:
   ```yaml
   env:
     GO_VERSION: '1.24'
   ```

3. go.mod güncelleyin:
   ```
   go 1.24
   ```

4. Dependencies'i güncelleyin:
   ```bash
   go get -u ./...
   go mod tidy
   ```

**Ancak şu an için Go 1.22 kullanın!** 🔒

### 📞 Yardım

Sorun yaşarsanız:
1. `backend/go.mod` dosyasında `go 1.22` olduğundan emin olun
2. `go mod tidy` çalıştırın
3. `go build ./cmd/api` ile test edin
4. Docker build deneyin: `docker build -t test ./backend`

## ✅ Özet

**GO VERSION = 1.22**

Her zaman, her yerde, her dosyada!

