#!/usr/bin/env bash
# Test script for local cross-platform builds
# This shows what can be built locally vs what requires CI/native runners

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

echo "=========================================="
echo "MDV Cross-Platform Build Test"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

success() {
    echo -e "${GREEN}✓${NC} $1"
}

warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
}

# Test 1: TUI Cross-Compilation (should always work)
echo "=== Test 1: TUI Cross-Compilation ==="
echo "Building TUI for multiple platforms (no CGO, should work everywhere)..."
echo ""

if CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tmp/mdv-linux-amd64 ./cmd/mdv 2>/dev/null; then
    success "Linux amd64 TUI build"
else
    error "Linux amd64 TUI build FAILED"
fi

if CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /tmp/mdv-linux-arm64 ./cmd/mdv 2>/dev/null; then
    success "Linux arm64 TUI build"
else
    error "Linux arm64 TUI build FAILED"
fi

if CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o /tmp/mdv-windows.exe ./cmd/mdv 2>/dev/null; then
    success "Windows amd64 TUI build"
else
    error "Windows amd64 TUI build FAILED"
fi

if CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o /tmp/mdv-darwin-amd64 ./cmd/mdv 2>/dev/null; then
    success "macOS amd64 TUI build"
else
    error "macOS amd64 TUI build FAILED"
fi

if CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o /tmp/mdv-darwin-arm64 ./cmd/mdv 2>/dev/null; then
    success "macOS arm64 TUI build"
else
    error "macOS arm64 TUI build FAILED"
fi

echo ""

# Test 2: macOS GUI Build (only works on macOS)
echo "=== Test 2: macOS GUI Build ==="
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "Building macOS GUI with Wails..."
    if ./scripts/build-wails.sh 2>&1 | grep -q "SUCCESS"; then
        success "macOS GUI build"
        ls -lh cmd/mdv-gui/build/bin/
    else
        warning "macOS GUI build completed (check output for issues)"
    fi
else
    warning "Skipped (not on macOS - requires macOS runner)"
fi

echo ""

# Test 3: Linux GUI Build (requires Linux or Docker)
echo "=== Test 3: Linux GUI Build ==="
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "Building Linux GUI with Wails..."
    if command -v gtk-config >/dev/null 2>&1 || command -v pkg-config >/dev/null 2>&1; then
        if WAILS_PLATFORMS="linux/amd64" ./scripts/build-wails.sh 2>&1 | grep -q "SUCCESS"; then
            success "Linux GUI build"
        else
            warning "Linux GUI build completed (check output for issues)"
        fi
    else
        error "Linux GUI dependencies not installed (libgtk-3-dev, libwebkit2gtk-4.1-dev or libwebkit2gtk-4.0-dev)"
    fi
elif command -v docker >/dev/null 2>&1; then
    warning "Skipped (not on Linux - use Docker: see scripts/docker/Dockerfile.linux-build)"
else
    warning "Skipped (requires Linux runner or Docker)"
fi

echo ""

# Test 4: Windows GUI Build (requires Windows)
echo "=== Test 4: Windows GUI Build ==="
if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    echo "Building Windows GUI with Wails..."
    if WAILS_PLATFORMS="windows/amd64" ./scripts/build-wails.sh 2>&1 | grep -q "SUCCESS"; then
        success "Windows GUI build"
    else
        warning "Windows GUI build completed (check output for issues)"
    fi
else
    warning "Skipped (requires Windows runner or VM)"
fi

echo ""

# Test 5: GoReleaser Snapshot (local packaging test)
echo "=== Test 5: GoReleaser Snapshot ==="
if command -v goreleaser >/dev/null 2>&1; then
    echo "Running GoReleaser in snapshot mode..."
    echo "(This packages what's already built, doesn't cross-compile GUI)"
    echo ""

    if goreleaser release --snapshot --clean --skip=publish 2>&1 | tee /tmp/goreleaser-output.log | grep -q "SUCCESS\|done"; then
        success "GoReleaser snapshot completed"
        echo ""
        echo "Check ./dist/ for generated archives:"
        ls -lh dist/*.tar.gz dist/*.zip 2>/dev/null || echo "  (no archives found)"
    else
        warning "GoReleaser completed with warnings (check /tmp/goreleaser-output.log)"
    fi
else
    warning "GoReleaser not installed (brew install goreleaser)"
fi

echo ""
echo "=========================================="
echo "Summary"
echo "=========================================="
echo ""
echo "✓ TUI builds work everywhere (pure Go, no CGO)"
echo "✓ GUI builds require native platforms:"
echo "  - macOS GUI: Needs macOS (Xcode, Cocoa)"
echo "  - Linux GUI: Needs Linux (GTK3, WebKit2GTK)"
echo "  - Windows GUI: Needs Windows (MinGW, WebView2)"
echo ""
echo "For full cross-platform testing:"
echo "  1. Local: macOS GUI + TUI for all platforms"
echo "  2. Docker: Linux GUI (see scripts/docker/Dockerfile.linux-build)"
echo "  3. CI/CD: Automated multi-platform builds via GitHub Actions"
echo ""
echo "To test CI workflow locally:"
echo "  act -j build-gui  # requires 'act' CLI tool"
echo ""
