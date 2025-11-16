# GitHub Actions Workflows - Değişiklik Geçmişi

## v1.1.0 - 2025-11-16

### 🔄 Action Güncellemeleri (Deprecated Versiyonlar Kaldırıldı)

#### Güncellenen Action'lar:

1. **actions/upload-artifact**: v3 → v4
   - Dosyalar: `ci.yml`, `scheduled-tasks.yml`, `release.yml`
   - Sebep: v3 deprecated, GitHub'ın yeni artifact sistemi

2. **actions/download-artifact**: v3 → v4
   - Dosyalar: `release.yml`
   - Sebep: v3 deprecated

3. **codecov/codecov-action**: v3 → v4
   - Dosyalar: `ci.yml`
   - Değişiklik: Token parametresi eklendi (v4 gereksinimi)
   - Gerekli secret: `CODECOV_TOKEN`

4. **github/codeql-action/upload-sarif**: v2 → v4
   - Dosyalar: `ci.yml`
   - Sebep: v2 ve v3 deprecated (v3, December 2026'da kaldırılacak)
   - Eklenen: `security-events: write` permission

5. **appleboy/ssh-action**: v1.0.0 → v1.0.3
   - Dosyalar: `deploy.yml`
   - Sebep: Bug fixes ve iyileştirmeler

6. **golangci/golangci-lint-action**: v3 → v6
   - Dosyalar: `ci.yml`, `backend/.golangci.yml`
   - Sebep: v3 eski, v6 en son stable versiyon
   - Değişiklikler:
     - `golint` → `revive` (golint deprecated)
     - `gomnd` → `mnd` (yeni isim)
     - `maligned` → kaldırıldı (deprecated)
     - `govet.check-shadowing` → `govet.enable: [shadow]`
     - `nolintlint.allow-leading-space` → kaldırıldı
     - `run.skip-dirs/skip-files` → `issues.exclude-dirs/exclude-files`

7. **actions/dependency-review-action**: v3 → v4
   - Dosyalar: `pr-check.yml`
   - Sebep: v4 daha fazla özellik ve iyileştirme

8. **actions/create-release**: v1 → KALDIRILDI
   - Dosyalar: `release.yml`
   - Alternatif: `softprops/action-gh-release@v1`
   - Sebep: v1 deprecated ve artık desteklenmiyor

9. **actions/upload-release-asset**: v1 → KALDIRILDI
   - Dosyalar: `release.yml`
   - Alternatif: `softprops/action-gh-release@v1` (asset upload dahil)
   - Sebep: v1 deprecated

### 🔧 Release Workflow Yeniden Yapılandırıldı

`release.yml` dosyası tamamen yeniden yazıldı:
- ✅ Modern `softprops/action-gh-release` action'ı kullanılıyor
- ✅ Otomatik changelog oluşturma (git commits'ten)
- ✅ Multi-platform binary build korundu
- ✅ Asset upload otomatik
- ✅ GitHub Release Notes otomatik oluşturma eklendi

### 📝 Yeni Gereksinimler

#### GitHub Repository Secrets:

**Opsiyonel (Codecov kullanacaksanız):**
```
CODECOV_TOKEN - Codecov API token (codecov.io'dan alınır)
```

#### Codecov Token Nasıl Alınır:
1. https://codecov.io adresine gidin
2. GitHub ile giriş yapın
3. Repository'nizi ekleyin
4. Settings → Repository Upload Token
5. Token'ı kopyalayıp GitHub Secrets'a ekleyin

**Not:** Codecov kullanmayacaksanız, `ci.yml` dosyasından coverage upload adımını kaldırabilirsiniz.

### 🐛 Düzeltilen Hatalar

- ✅ "deprecated version of actions/upload-artifact: v3" hatası
- ✅ "deprecated version of actions/create-release: v1" hatası
- ✅ "deprecated version of actions/upload-release-asset: v1" hatası
- ✅ "CodeQL Action v3 will be deprecated" uyarısı
- ✅ "Resource not accessible by integration" permission hatası
- ✅ golangci-lint konfigürasyon validation hataları
- ✅ "additional properties not allowed" hataları (.golangci.yml)
- ✅ Deprecated linter kullanımları (golint, gomnd, maligned)
- ✅ Go code formatting hataları (54 dosya gofmt ile formatlandı)
- ✅ Import sıralaması ve kod stili iyileştirildi
- ✅ Go version uyumsuzluğu (1.24.0 → 1.22, Dockerfile ile uyumlu)
- ✅ Docker build "go.mod requires go >= 1.24.0" hatası
- ✅ go vet "time imported and not used" hatası (`redis_limiter_test.go`)
- ✅ go vet "redis redeclared in this block" hatası (`test_helpers.go`)
- ✅ Import çakışması (redis client vs redis testcontainer)
- ✅ golangci-lint "no go files to analyze" hatası
- ✅ backend-lint job'unda eksik dependency download
- ✅ Go version uyumsuzluğu (CI: 1.21 vs go.mod: 1.22)
- ✅ Test panic: nil pointer dereference (`content_sync_service_test.go`)
- ✅ Test'te nil logger kullanımı (ScoreCalculatorService)
- ✅ Docker build: grpc-gateway v2.27.3 Go 1.24.0 gerektirme hatası
- ✅ Dependency version uyumsuzluğu
- ✅ Release workflow'unda asset upload sorunları

### 📦 Go Module ve Version Güncellemeleri

**Go Version:**
- `go 1.22` (Dockerfile ve CI ile uyumlu)
- GitHub Actions `GO_VERSION` 1.21 → 1.22 güncellendi
- `backend-lint` job'una `go mod download` eklendi

**Eklenen Test Dependencies:**
- `testcontainers/testcontainers-go` v0.40.0
- `testcontainers/testcontainers-go/modules/postgres` v0.40.0
- `testcontainers/testcontainers-go/modules/redis` v0.40.0
- `stretchr/testify` v1.11.1

**Güncellenen Dependencies:**
- `redis/go-redis/v9` v9.5.1 → v9.7.3
- `grpc-ecosystem/grpc-gateway/v2` v2.27.3 → v2.22.0 (Go 1.22 uyumluluğu için downgrade)

**Not:** `go mod tidy` çalıştırıldı ve tüm eksik dependencies eklendi. Go version 1.22'de sabitlendi (Dockerfile uyumu için). grpc-gateway v2.27.3, Go 1.24.0 gerektirdiği için v2.22.0'a downgrade edildi.

### 🚀 İyileştirmeler

1. **Daha Hızlı Artifact İşleme**: v4 artifact sistemi daha hızlı ve güvenilir
2. **Daha İyi Release Yönetimi**: `softprops/action-gh-release` daha esnek ve özellik zengin
3. **Otomatik Changelog**: Git commit'lerinden otomatik changelog oluşturma
4. **GitHub Release Notes**: Otomatik release notes oluşturma aktif
5. **Security Scanning İyileştirildi**: CodeQL v4 + doğru permission'lar
6. **Code Formatting**: 54 Go dosyası `gofmt -s -w` ile formatlandı
7. **Module Management**: `go mod tidy` ile eksik dependencies eklendi

### ⚠️ Breaking Changes

**YOK** - Tüm değişiklikler backward compatible

### 🔐 Permission Güncellemeleri

**Security Scanning Job:**
```yaml
permissions:
  contents: read
  security-events: write
```

Bu permission'lar olmadan CodeQL SARIF upload başarısız olur.

### 📚 Referanslar

- [GitHub Blog: Artifact Actions v3 Deprecation](https://github.blog/changelog/2024-04-16-deprecation-notice-v3-of-the-artifact-actions/)
- [GitHub Blog: CodeQL Action v3 Deprecation](https://github.blog/changelog/2025-10-28-upcoming-deprecation-of-codeql-action-v3/)
- [actions/upload-artifact v4 Migration Guide](https://github.com/actions/upload-artifact/blob/main/docs/MIGRATION.md)
- [Codecov Action v4 Documentation](https://github.com/codecov/codecov-action)
- [softprops/action-gh-release](https://github.com/softprops/action-gh-release)
- [GitHub Permissions Documentation](https://docs.github.com/en/actions/using-jobs/assigning-permissions-to-jobs)

---

## v1.0.0 - 2025-11-16

### 🎉 İlk Sürüm

- ✅ CI Pipeline (test, lint, build)
- ✅ Deployment Pipeline
- ✅ PR Check Automation
- ✅ Release Management
- ✅ Scheduled Tasks
- ✅ Dependabot Configuration
- ✅ Auto-labeling
- ✅ CODEOWNERS

