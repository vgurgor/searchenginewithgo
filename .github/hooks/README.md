# Git Hooks

Bu dizin, projeye özel Git hook'larını içerir.

## 🪝 Mevcut Hook'lar

### pre-commit
Go version kontrolü yapar ve yanlışlıkla Go 1.24.0'ın commit edilmesini önler.

## 📥 Kurulum

Local bilgisayarınızda bu hook'ları kullanmak için:

```bash
# Proje root dizininde
cd /Users/volkan/blexi/developmentarea/search_engine

# Hook'ları git hooks dizinine kopyalayın
cp .github/hooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

# Test edin
.git/hooks/pre-commit
```

Veya tek komutla:

```bash
make install-hooks  # (root Makefile'a eklenebilir)
```

## ✅ pre-commit Hook - Ne Yapar?

Her commit öncesi otomatik olarak:
- ✅ `backend/go.mod` dosyasını kontrol eder
- ✅ Go version 1.24.0 ise commit'i engeller
- ✅ Düzeltme komutu önerir: `make fix-go-version`

## 🔧 Hook'u Bypass Etmek (Acil Durum)

**DİKKAT:** Sadece acil durumlarda kullanın!

```bash
git commit --no-verify -m "your message"
```

## 📝 Örnek Çıktı

### ✅ Başarılı Durum
```bash
$ git commit -m "feat: add new feature"
🔍 Checking Go version in go.mod...
✅ Go version is correct (1.22)
[main abc1234] feat: add new feature
```

### ❌ Hatalı Durum
```bash
$ git commit -m "feat: add new feature"
🔍 Checking Go version in go.mod...

❌ ERROR: go.mod has Go 1.24.0!

This will break Docker build because:
  - Go 1.24.0 is not released yet
  - Dockerfile uses golang:1.22-alpine
  - CI uses GO_VERSION: '1.22'

Please change backend/go.mod to:
  go 1.22

Quick fix:
  cd backend && make fix-go-version

See backend/GO_VERSION.md for details
```

## 🎯 Neden Gerekli?

Go 1.24.0 henüz yayınlanmadı ve:
- ❌ Docker build başarısız olur
- ❌ CI/CD pipeline başarısız olur
- ❌ Production deployment imkansız olur

Hook sayesinde bu sorun **commit aşamasında** engellenir!

## 🔄 Hook Güncelleme

Hook dosyası güncellendiğinde:

```bash
# Yeni versiyonu kopyalayın
cp .github/hooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

## 📚 Daha Fazla Bilgi

- [Git Hooks Documentation](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks)
- [backend/GO_VERSION.md](../../backend/GO_VERSION.md) - Go version politikası

