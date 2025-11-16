# CI/CD Pipeline Dokümantasyonu

Bu dizin, Search Engine projesi için GitHub Actions CI/CD pipeline'larını içerir.

## 🚀 Workflow'lar

### 1. CI Pipeline (`ci.yml`)

**Tetikleme:** Push ve Pull Request (main, develop branch'leri)

**İşlemler:**
- ✅ Backend Go testleri ve linting
- ✅ Frontend TypeScript/React testleri ve linting
- ✅ Docker image build testleri
- ✅ Security scanning (Trivy)
- ✅ Code coverage raporlama

**Servisler:**
- PostgreSQL 15
- Redis 7

### 2. Deployment Pipeline (`deploy.yml`)

**Tetikleme:** 
- Push to main branch
- Version tags (v*.*.*)
- Manuel tetikleme

**İşlemler:**
- 🐳 Docker image build ve GitHub Container Registry'ye push
- 🚢 Production/Staging sunucularına deployment
- 🏥 Health check
- 📢 Slack bildirimi

**Gerekli GitHub Secrets:**
```
DEPLOY_HOST        # Deployment sunucu adresi
DEPLOY_USER        # SSH kullanıcı adı
DEPLOY_SSH_KEY     # SSH private key
DEPLOY_PORT        # SSH portu (varsayılan: 22)
DEPLOY_URL         # Health check için URL
SLACK_WEBHOOK      # Slack webhook URL (opsiyonel)
```

### 3. Pull Request Checks (`pr-check.yml`)

**Tetikleme:** Pull Request açılışı/güncellenmesi

**İşlemler:**
- 📝 PR başlık formatı kontrolü (Conventional Commits)
- 🌿 Branch isim kontrolü
- 📊 Kod kalitesi kontrolleri
- 🔍 Hardcoded secret taraması (TruffleHog)
- 🏷️ Otomatik etiketleme
- 📏 PR boyut etiketleme

**Branch İsimlendirme:**
- `feature/*` - Yeni özellikler
- `bugfix/*` - Bug düzeltmeleri
- `hotfix/*` - Acil düzeltmeler
- `release/*` - Release hazırlıkları

**PR Başlık Formatı:**
```
<type>(<scope>): <description>

Örnekler:
- feat: add new search filter
- fix(backend): resolve database connection issue
- docs: update README
```

### 4. Scheduled Tasks (`scheduled-tasks.yml`)

**Tetikleme:** Her gün saat 02:00 UTC

**İşlemler:**
- 🔄 Dependency güncellemelerini kontrol
- 🔒 Security audit
- 🧹 Eski Docker image'larını temizleme
- 📈 Coverage raporu oluşturma

## 🛠️ Kurulum

### 1. GitHub Repository Secrets Ekleme

Repository Settings → Secrets and variables → Actions:

```bash
# Deployment secrets
DEPLOY_HOST=your-server.com
DEPLOY_USER=deploy-user
DEPLOY_SSH_KEY=<private-key-content>
DEPLOY_PORT=22
DEPLOY_URL=https://your-domain.com
SLACK_WEBHOOK=https://hooks.slack.com/services/...

# Code Coverage (Opsiyonel)
CODECOV_TOKEN=<your-codecov-token>
```

**Codecov Token Nasıl Alınır:**
1. https://codecov.io adresine gidin
2. GitHub ile giriş yapın
3. Repository'nizi ekleyin
4. Settings → Repository Upload Token
5. Token'ı kopyalayıp GitHub Secrets'a ekleyin

**Not:** Codecov kullanmayacaksanız, `ci.yml` dosyasından coverage upload adımını silebilirsiniz.

### 2. GitHub Container Registry Yapılandırması

Package settings'den public/private ayarlarını yapılandırın.

### 3. Branch Protection Rules

Repository Settings → Branches → Add rule:

**main** branch için:
- ✅ Require pull request reviews (1 reviewer)
- ✅ Require status checks to pass
  - backend-test
  - frontend-test
  - backend-lint
  - docker-build
- ✅ Require branches to be up to date
- ✅ Include administrators

## 📊 Status Badges

README.md dosyanıza ekleyebileceğiniz badge'ler:

```markdown
![CI](https://github.com/{owner}/{repo}/workflows/CI%20Pipeline/badge.svg)
![Deploy](https://github.com/{owner}/{repo}/workflows/Deploy%20to%20Production/badge.svg)
```

## 🔧 Yerel Test

Workflow'ları yerel olarak test etmek için [act](https://github.com/nektos/act) kullanabilirsiniz:

```bash
# CI pipeline'ı test et
act -j backend-test

# Tüm workflow'u çalıştır
act push
```

## 📈 Monitoring

### Coverage Reports
- Backend: Codecov'a otomatik upload
- Artifacts: 30 gün boyunca GitHub'da saklanır

### Build Artifacts
- Frontend build: 1 gün boyunca saklanır
- Coverage reports: 30 gün boyunca saklanır

## 🤝 Katkıda Bulunma

1. Feature branch oluşturun
2. Değişikliklerinizi commit edin
3. Branch'i push edin
4. Pull Request açın
5. CI checks'lerin geçmesini bekleyin
6. Code review sürecini tamamlayın

## 📞 Destek

Sorun yaşarsanız:
1. GitHub Actions logs'larını kontrol edin
2. Issue açın
3. Team ile iletişime geçin

## 🔄 Güncelleme Geçmişi

| Tarih | Değişiklik | Versiyon |
|-------|-----------|----------|
| 2025-11 | İlk oluşturma | v1.0.0 |

