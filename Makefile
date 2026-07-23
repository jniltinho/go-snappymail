## UPX
UPX_VERSION := 5.1.1
UPX_ARCHIVE := upx-$(UPX_VERSION)-amd64_linux.tar.xz
UPX_DIR     := upx-$(UPX_VERSION)-amd64_linux
UPX_BIN     := /usr/local/bin/upx
UPX_URL     := https://github.com/upx/upx/releases/download/v$(UPX_VERSION)/$(UPX_ARCHIVE)

## Application
APP         := go-snappymail
MODULE      := go-snappymail
CMD_PKG     := $(MODULE)/cmd
DIST_DIR    := dist
BIN         := $(DIST_DIR)/$(APP)
VERSION     := $(shell (git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//') || echo dev)
BUILD_TIME  := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GOOS        ?= $(shell go env GOOS)
GOARCH      ?= $(shell go env GOARCH)

LDFLAGS := -trimpath -ldflags "-s -w \
	-X $(CMD_PKG).Version=$(VERSION) \
	-X $(CMD_PKG).BuildDate=$(BUILD_TIME) \
	-X $(CMD_PKG).GitCommit=$(GIT_COMMIT)"

.PHONY: all build build-prod release run migrate clean tidy deps frontend frontend-dev dev \
        install-upx test test-integration test-short check-git new-skin validate-skins help

## Default: build binary into dist/
all: build

## Build embedded SPA (placeholder until frontend/ Vue project exists)
frontend:
	@if [ -f frontend/package.json ]; then \
		echo "Building Vue frontend..."; \
		cd frontend && npm run build; \
	else \
		echo "Using checked-in web/dist/ (no frontend/ yet)"; \
	fi

frontend-dev dev:
	@echo "Start backend: make run  (or: go run . serve --debug)"
	@if [ -f frontend/package.json ]; then cd frontend && npm run dev; \
	else echo "frontend/ not initialized yet"; fi

## Development build → dist/go-snappymail
build: frontend
	@mkdir -p $(DIST_DIR)
	@echo "Building $(APP) $(VERSION) ($(GOOS)/$(GOARCH)) → $(BIN)"
	CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BIN) $(LDFLAGS) .

## Production build with UPX compression
build-prod: build
	@command -v upx >/dev/null 2>&1 || { echo "UPX not found. Run: make install-upx"; exit 1; }
	@echo "Compressing $(BIN) with UPX..."
	upx --best --lzma $(BIN)
	@ls -lh $(BIN)

## Cross-compile release binaries into dist/
release: frontend
	@mkdir -p $(DIST_DIR)
	@for target in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64; do \
		os=$${target%/*}; arch=$${target#*/}; \
		out="$(DIST_DIR)/$(APP)_$(VERSION)_$${os}_$${arch}"; \
		echo "Building $$out ..."; \
		CGO_ENABLED=1 GOOS=$$os GOARCH=$$arch go build -o "$$out" $(LDFLAGS) . || exit 1; \
		upx --best --lzma "$$out" 2>/dev/null || true; \
	done
	@echo "Release artifacts:"
	@ls -lh $(DIST_DIR)/$(APP)_$(VERSION)_*

run:
	@test -x $(BIN) || $(MAKE) build
	./$(BIN) serve

migrate:
	@test -x $(BIN) || $(MAKE) build
	./$(BIN) migrate

## Unit tests (race + coverage)
test:
	go test -v -race -coverprofile=coverage.out ./...

## Integration tests (IMAP lab — optional)
test-integration:
	IMAP_TEST_HOST=mailserver IMAP_TEST_USER=user@test.local IMAP_TEST_PASS='Password1@' \
	IMAP_TEST_INSECURE=1 IMAP_TEST_TLS_SERVER_NAME=mail.test.local \
	go test -v -tags=integration ./internal/handler/...

test-short:
	go test -v -short ./...

## Reject base/ and dist/ binaries from git
check-git:
	@bash scripts/check-git-clean.sh

## Scaffold a new UI skin (see docs/skins.md)
new-skin:
	@test -n "$(ID)" || { echo "Usage: make new-skin ID=mybrand [REGISTER=1]"; exit 1; }
	@REGISTER=$(REGISTER) bash scripts/new-skin.sh "$(ID)" $(if $(filter 1,$(REGISTER)),--register,)

## Verify Go catalog ↔ TS manifest ↔ CSS imports
validate-skins:
	@bash scripts/validate-skins.sh

clean:
	@echo "Removing binaries from $(DIST_DIR)/"
	rm -rf $(DIST_DIR)/$(APP) $(DIST_DIR)/$(APP)_*

tidy:
	go mod tidy

deps:
	go mod download

install-upx:
	@echo "Installing UPX $(UPX_VERSION)..."
	curl -ksSL "$(UPX_URL)" -o "$(UPX_ARCHIVE)"
	tar -xf "$(UPX_ARCHIVE)"
	chmod +x "$(UPX_DIR)/upx"
	sudo mv "$(UPX_DIR)/upx" "$(UPX_BIN)"
	rm -rf "$(UPX_DIR)" "$(UPX_ARCHIVE)"
	@$(UPX_BIN) --version

help:
	@echo "Targets:"
	@echo "  all          Build binary into dist/ (default)"
	@echo "  build        Build dist/$(APP)"
	@echo "  build-prod   Build + UPX compress dist/$(APP)"
	@echo "  release      Cross-compile + UPX → dist/$(APP)_VERSION_GOOS_GOARCH"
	@echo "  run          Run dist/$(APP) serve"
	@echo "  migrate      Run dist/$(APP) migrate"
	@echo "  test         Run unit tests with race and coverage"
	@echo "  test-integration  Run IMAP login integration test (Docker lab)"
	@echo "  check-git    Fail if base/ or dist/ binaries are tracked"
	@echo "  new-skin     Scaffold skin — make new-skin ID=mybrand [REGISTER=1]"
	@echo "  validate-skins  Check Go catalog, TS manifest, CSS imports are in sync"
	@echo "  frontend     Build Vue SPA to web/dist/ (when frontend/ exists)"
	@echo "  clean        Remove dist/$(APP)* binaries"
	@echo "  tidy         go mod tidy"
	@echo "  deps         go mod download"
	@echo "  install-upx  Install UPX $(UPX_VERSION)"
	@echo ""
	@echo "Variables: VERSION=$(VERSION) GOOS=$(GOOS) GOARCH=$(GOARCH)"
