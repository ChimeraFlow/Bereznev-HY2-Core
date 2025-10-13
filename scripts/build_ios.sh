#!/usr/bin/env bash
set -euo pipefail

# ===== Config =====
PKG="github.com/ChimeraFlow/Bereznev-HY2-Core/core-go/mobile"
OUT_DIR="dist/ios"
FW_NAME="Hy2Core"
MOD_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

VERSION="${VERSION:-$(git -C "$MOD_ROOT" describe --tags --abbrev=0 2>/dev/null || echo 0.1.0)}"
BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
COMMIT_HASH="$(git -C "$MOD_ROOT" rev-parse --short HEAD 2>/dev/null || echo dev)"

need() { command -v "$1" >/dev/null 2>&1 || { echo "❌ Требуется $1"; exit 1; }; }
need go
need gomobile
need xcodebuild

mkdir -p "$OUT_DIR"
pushd "$MOD_ROOT" >/dev/null

gomobile version >/dev/null 2>&1 || gomobile init

# iOS: создаём универсальную XCFramework
GOFLAGS="-ldflags=-X 'mobile.sdkVersion=${VERSION}' -X 'mobile.engineID=skeleton' -X 'mobile.buildTime=${BUILD_TIME}' -X 'mobile.commitHash=${COMMIT_HASH}'"

echo "▶ gomobile bind → XCFramework (${FW_NAME}-${VERSION}.xcframework)"
gomobile bind \
  -target=ios \
  -ldflags="${GOFLAGS}" \
  -o "${OUT_DIR}/${FW_NAME}-${VERSION}.xcframework" \
  "${PKG}"

cp -Rf "${OUT_DIR}/${FW_NAME}-${VERSION}.xcframework" "${OUT_DIR}/${FW_NAME}.xcframework"
echo "✅ Готово: ${OUT_DIR}/${FW_NAME}-${VERSION}.xcframework"
popd >/dev/null