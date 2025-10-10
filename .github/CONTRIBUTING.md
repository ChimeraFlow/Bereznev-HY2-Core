# Contributing to Bereznev HY2 Core

Thanks for wanting to help!  
This document explains how to work with the repository and submit contributions properly.

---

## 🔗 Quick Links
- 🛡️ Security reports → see [`SECURITY.md`](../SECURITY.md)
- 🧭 Architecture → [`docs/ARCHITECTURE.md`](../docs/ARCHITECTURE.md)
- ⚙️ GitHub flow → [`docs/GITHUB_SETTINGS.md`](../docs/GITHUB_SETTINGS.md)

---

## 🧰 Development Setup
Requirements:
- JDK **17+**
- Go **1.22+**
- Android SDK + NDK
- Kotlin / Gradle  
Optional:
- [`gomobile`](https://pkg.go.dev/golang.org/x/mobile/cmd/gomobile) for AAR builds (see README)

---

## 🌿 Branching Model
- `main` – stable branch, deploys to **prod**
- `develop` – integration branch, deploys to **staging**
- `feature/*` – new work → PR to `develop`
- `fix/*` – bugfix → PR to `develop`
- `hotfix/*` – critical fix → PR to `main`, then cherry-pick to `develop`

```bash
git checkout -b feature/my-change
# work...
git commit -m "feat(core): add obfs mode"
git push -u origin feature/my-change
# open PR → develop
```