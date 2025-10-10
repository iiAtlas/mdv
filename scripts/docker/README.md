# Docker Build Environment for Linux

This directory contains Docker configurations for building the mdv Linux GUI without requiring a Linux VM.

## Quick Start

```bash
# Build Linux GUI for both amd64 and arm64 using Docker (from repository root)
./scripts/docker/build-linux.sh

# Or build for a specific architecture:
./scripts/docker/build-linux.sh amd64    # x86_64 only
./scripts/docker/build-linux.sh arm64    # ARM64 only
```

This script will:
1. Build the Docker image with all Linux dependencies (including ARM64 cross-compiler)
2. Run the build inside the container
3. Output binaries to `cmd/mdv-gui/build/bin/`

## Manual Usage

If you prefer to run Docker commands manually:

```bash
# Build the Docker image (x86_64/amd64)
docker build --platform linux/amd64 -f scripts/docker/Dockerfile.linux-build -t mdv-linux-builder .

# Build Linux GUI (amd64 native)
docker run --platform linux/amd64 --rm -v "$PWD:/workspace" \
  -e WAILS_PLATFORMS=linux/amd64 \
  -e CC=gcc \
  mdv-linux-builder

# Build for ARM64 (cross-compile)
docker run --platform linux/amd64 --rm -v "$PWD:/workspace" \
  -e WAILS_PLATFORMS=linux/arm64 \
  -e CC=aarch64-linux-gnu-gcc \
  -e PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig \
  -e PKG_CONFIG_LIBDIR=/usr/lib/aarch64-linux-gnu/pkgconfig:/usr/share/pkgconfig \
  mdv-linux-builder

# Run a shell inside the container for debugging
docker run --platform linux/amd64 --rm -it -v "$PWD:/workspace" \
  mdv-linux-builder bash
```

**Note**: The `--platform linux/amd64` flag ensures the container runs in x86_64 mode, which works on both Intel and Apple Silicon Macs (via Rosetta 2 emulation on Apple Silicon).

## Requirements

- Docker Desktop (macOS/Windows) or Docker Engine (Linux)
- ~2GB disk space for the image
- ~5-10 minutes for first build (image creation)
- **Note**: On Apple Silicon Macs, the container runs via Rosetta 2 emulation (x86_64)

## What's Inside

The Docker image includes:
- Ubuntu 22.04 base (amd64)
- Go 1.24
- Wails CLI
- GTK3 development libraries
- WebKit2GTK development libraries
- Build essentials (gcc, make, etc.)
- ARM64 cross-compilation toolchain (gcc-aarch64-linux-gnu)

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

## Advanced: Using with CI/CD

This same Dockerfile can be used in CI/CD pipelines that don't have native Linux runners:

```yaml
# GitHub Actions example
- name: Build Linux GUI via Docker
  run: |
    docker build -f scripts/docker/Dockerfile.linux-build -t mdv-linux-builder .
    docker run --rm -v "$PWD:/workspace" mdv-linux-builder
```
