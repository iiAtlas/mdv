.PHONY: all tui gui clean install install-tui install-gui run-tui run-gui help

# Default target
all: tui gui

# Build TUI version
tui:
	@echo "Building TUI version..."
	@cd cmd/mdv && go build -o ../../mdv
	@echo "✓ Built: ./mdv"

# Build GUI version
gui:
	@echo "Building GUI version..."
	@cd cmd/mdv-gui && wails build
	@echo "✓ Built: ./cmd/mdv-gui/build/bin/mdv-gui"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f mdv
	@rm -rf cmd/mdv-gui/build
	@rm -f cmd/mdv-gui/frontend/wailsjs
	@echo "✓ Cleaned"

# Install TUI to $GOPATH/bin
install-tui:
	@echo "Installing mdv (TUI) to $(shell go env GOPATH)/bin..."
	@cd cmd/mdv && go install
	@echo "✓ Installed: $(shell go env GOPATH)/bin/mdv"

# Install GUI to $GOPATH/bin
install-gui:
	@echo "Installing mdv-gui to $(shell go env GOPATH)/bin..."
	@cd cmd/mdv-gui && wails build
	@cp cmd/mdv-gui/build/bin/mdv-gui $(shell go env GOPATH)/bin/
	@echo "✓ Installed: $(shell go env GOPATH)/bin/mdv-gui"

# Install both
install: install-tui install-gui

# Run TUI in development
run-tui:
	@go run ./cmd/mdv test.md

# Run GUI in development mode
run-gui:
	@cd cmd/mdv-gui && wails dev

# Show help
help:
	@echo "MDV - Markdown Viewer"
	@echo ""
	@echo "Build targets:"
	@echo "  make tui          Build TUI version to ./mdv"
	@echo "  make gui          Build GUI version with Wails"
	@echo "  make all          Build both versions (default)"
	@echo "  make clean        Remove built binaries"
	@echo ""
	@echo "Install targets:"
	@echo "  make install-tui  Install TUI to \$$GOPATH/bin/mdv"
	@echo "  make install-gui  Install GUI to \$$GOPATH/bin/mdv-gui"
	@echo "  make install      Install both"
	@echo ""
	@echo "Run targets:"
	@echo "  make run-tui      Quick test of TUI (uses test.md)"
	@echo "  make run-gui      Run GUI in dev mode"
	@echo ""
	@echo "Other:"
	@echo "  make help         Show this help"
