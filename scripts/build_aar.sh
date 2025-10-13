#!/usr/bin/env bash
set -euo pipefail

# ===== Config =====
PKG="github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/mobile"  # gomobile bind package
OUT_DIR="dist/android"
AAR_NAME="hy2core"
MOD_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Версию берём из тега, либо из ENV, либо fallback
VERSION="${VERSION:-$(git -C "$MOD_ROOT" describe --tags --abbrev=0 2>/dev/null || echo 0.1.0)}"
BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
COMMIT_HASH="$(git -C "$MOD_ROOT" rev-parse --short HEAD 2>/dev/null || echo dev)"

# ===== Checks =====
need() { command -v "$1" >/dev/null 2>&1 || { echo "❌ Требуется $1"; exit 1; }; }
need go
need gomobile

if [[ -z "${ANDROID_HOME:-}" && -z "${ANDROID_SDK_ROOT:-}" ]]; then
  echo "❌ ANDROID_HOME или ANDROID_SDK_ROOT не установлены"; exit 1;
fi

# ===== Init =====
mkdir -p "$OUT_DIR"
pushd "$MOD_ROOT" >/dev/null

# Один раз на машине:
gomobile version >/dev/null 2>&1 || gomobile init

# ===== Build =====
# Совет: задайте package name через -javapkg, чтобы классы попадали в tech.bereznev.hy2
JAVAPKG="tech.bereznev.hy2"

echo "▶ gomobile bind → AAR (${AAR_NAME}-${VERSION}.aar)"
GOFLAGS="-ldflags=-X 'mobile.sdkVersion=${VERSION}' -X 'mobile.engineID=skeleton' -X 'mobile.buildTime=${BUILD_TIME}' -X 'mobile.commitHash=${COMMIT_HASH}'"

gomobile bind -target=android -androidapi=21 -v \
  -javapkg "${JAVAPKG}" \
  -ldflags="${GOFLAGS}" \
  -o "${OUT_DIR}/${AAR_NAME}-${VERSION}.aar" \
  "${PKG}"

# Стабильное имя «последней» сборки
cp -f "${OUT_DIR}/${AAR_NAME}-${VERSION}.aar" "${OUT_DIR}/${AAR_NAME}.aar"

echo "✅ Готово: ${OUT_DIR}/${AAR_NAME}-${VERSION}.aar"
popd >/dev/null