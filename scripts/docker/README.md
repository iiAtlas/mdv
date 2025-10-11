# Docker Build System for Linux GUI

This directory contains Docker configurations for building mdv Linux GUI binaries on non-Linux platforms (macOS, Windows) and in CI/CD environments.

## Quick Start

### Build Both Architectures (Recommended)

```bash
./scripts/docker/build-linux.sh
```

This builds both AMD64 and ARM64 binaries using QEMU emulation for ARM64.

### Build Specific Architecture

```bash
# AMD64 only (faster)
./scripts/docker/build-linux.sh amd64

# ARM64 only (slower due to QEMU)
./scripts/docker/build-linux.sh arm64
```

Binaries will be output to `cmd/mdv-gui/build/bin/mdv-gui_linux_*/`

## Available Dockerfiles

### 1. `Dockerfile.linux-amd64` (Native AMD64 - Simple)

**Build:**
```bash
docker buildx build --platform linux/amd64 -f scripts/docker/Dockerfile.linux-amd64 -t mdv-linux-amd64 .
```

**Run:**
```bash
docker run --rm -v "$PWD:/workspace" mdv-linux-amd64
```

**Output:** `cmd/mdv-gui/build/bin/mdv-gui_linux_amd64/mdv-gui`

---

### 2. `Dockerfile.linux-arm64` (Native ARM64 via QEMU)

**Build:**
```bash
docker buildx build --platform linux/arm64 -f scripts/docker/Dockerfile.linux-arm64 -t mdv-linux-arm64 .
```

**Run:**
```bash
docker run --rm -v "$PWD:/workspace" mdv-linux-arm64
```

**Output:** `cmd/mdv-gui/build/bin/mdv-gui_linux_arm64/mdv-gui`

**Note:** ARM64 builds are 2-4x slower due to QEMU emulation.

---

### 3. `Dockerfile.linux-crosscompile` (Cross-Compilation - Fastest)

**Build Image:**
```bash
docker build -f scripts/docker/Dockerfile.linux-crosscompile -t mdv-linux-crosscompile .
```

**Build AMD64:**
```bash
docker run --rm -v "$PWD:/workspace" \
  -e WAILS_PLATFORMS=linux/amd64 \
  -e CC=gcc \
  mdv-linux-crosscompile
```

**Build ARM64:**
```bash
docker run --rm -v "$PWD:/workspace" \
  -e WAILS_PLATFORMS=linux/arm64 \
  -e CC=aarch64-linux-gnu-gcc \
  -e PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig \
  -e PKG_CONFIG_LIBDIR=/usr/lib/aarch64-linux-gnu/pkgconfig:/usr/share/pkgconfig \
  mdv-linux-crosscompile
```

**Pros:** Faster ARM64 builds (native AMD64 execution)
**Cons:** More complex setup

## Prerequisites

### macOS (Docker Desktop)

Docker Desktop includes all necessary tools:
- Docker Buildx
- QEMU emulation for multi-platform builds

**Install:**
```bash
brew install --cask docker
```

No additional setup required.

### Linux (Native Docker)

Install Docker Buildx and QEMU:

```bash
# Install Docker Buildx plugin
sudo apt-get install docker-buildx-plugin

# Install QEMU for multi-platform support
docker run --privileged --rm tonistiigi/binfmt --install all
```

### Windows (Docker Desktop)

Same as macOS - Docker Desktop includes all tools.

## Build Strategy Comparison

| Strategy | Speed (AMD64) | Speed (ARM64) | Complexity | Best For |
|----------|---------------|---------------|------------|----------|
| Native (amd64/arm64) | 2 min | 6 min | Simple | Local dev |
| Cross-Compile | 2 min | 2 min | Complex | Speed-critical |
| GitHub Actions | 2 min | 2 min | Simple | CI/CD |

**Recommendation:** Use `build-linux.sh` for local development (QEMU-based native builds)

## Troubleshooting

### Build fails with "permission denied"

Make sure the script is executable:
```bash
chmod +x scripts/docker/build-linux.sh
```

### Docker daemon not running

Start Docker Desktop or the Docker daemon:
```bash
# macOS
open -a Docker

# Linux
sudo systemctl start docker
```

### Out of disk space

Clean up old Docker images:
```bash
docker system prune -a
```

### ARM64 build is very slow

**Expected:** ARM64 builds via QEMU are 2-4x slower than native AMD64

**Solutions:**
1. Use cross-compilation Dockerfile (faster but more complex)
2. Accept slower build time (QEMU emulation overhead)
3. Use native ARM64 machine

### Error: "Package libgtk-3-dev:arm64 has no installation candidate"

**Cause:** APT sources not configured for ARM64 packages

**Solution:** Use the updated Dockerfiles which include proper APT source configuration for ports.ubuntu.com

## Verifying Builds

### Check Output

```bash
ls -lh cmd/mdv-gui/build/bin/
```

Expected:
```
mdv-gui_linux_amd64/mdv-gui    # AMD64 binary
mdv-gui_linux_arm64/mdv-gui    # ARM64 binary
```

### Test Binaries

```bash
# Test AMD64
docker run --rm -v "$PWD:/workspace" ubuntu:22.04 \
  /workspace/cmd/mdv-gui/build/bin/mdv-gui_linux_amd64/mdv-gui --version

# Test ARM64 (requires QEMU)
docker run --rm --platform linux/arm64 -v "$PWD:/workspace" ubuntu:22.04 \
  /workspace/cmd/mdv-gui/build/bin/mdv-gui_linux_arm64/mdv-gui --version
```

## Optimized Dockerfiles (New!)

### 4. `Dockerfile.linux-amd64-optimized` (BuildKit Cache + Multi-Stage)

**Enhanced version with:**
- BuildKit cache mounts for 30-60s faster rebuilds
- Pinned Wails version for reproducibility
- Go module layer caching
- Optional multi-stage runtime image

**Build:**
```bash
docker buildx build --platform linux/amd64 \
  -f scripts/docker/Dockerfile.linux-amd64-optimized \
  -t mdv-linux-amd64-optimized --load .
```

**Run:**
```bash
docker run --rm -v "$PWD:/workspace" mdv-linux-amd64-optimized
```

**Benefits:** Subsequent builds are 50% faster due to cached Go modules.

---

### 5. `Dockerfile.tui` (TUI-Only Container - ~10MB)

**Minimal container for running mdv in TUI mode:**
- Tiny image size (~10MB)
- Multi-platform support (amd64, arm64, arm/v7)
- Static binary (no dependencies)
- Perfect for servers and CI/CD

**Build:**
```bash
docker buildx build --platform linux/amd64,linux/arm64 \
  -f scripts/docker/Dockerfile.tui -t mdv-tui:latest --load .
```

**Run:**
```bash
# View your markdown files
docker run --rm -it -v "$PWD:/data" mdv-tui /data/README.md

# View included examples
docker run --rm -it mdv-tui /examples/demo.md
```

---

## Performance Improvements

The project now includes a `.dockerignore` file that reduces build context size:
- **Before:** ~2.9GB (includes .cache, dist, build artifacts)
- **After:** ~50MB (only source code and dependencies)
- **Impact:** Faster builds and smaller image transfer

## Related Documentation

- **`../../DOCKER_ASSESSMENT.md`** - Comprehensive Docker audit and recommendations
- **`../../DOCKER_IMPLEMENTATION_GUIDE.md`** - Step-by-step implementation guide
- **`cross-plat.md`** - Comprehensive cross-platform build guide
- **`.github/workflows/release.yml`** - CI/CD configuration

---

**Last Updated:** 2025-10-10
