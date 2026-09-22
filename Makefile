# TR-12 Client and Host — Go
# Run `make` to build, `make doctor` to check prerequisites.

SHELL := /bin/bash
.DEFAULT_GOAL := build

# Binaries
BIN_DIR    := client/bin
SDK_BIN    := $(BIN_DIR)/cdd-sdk
ARD_BIN    := $(BIN_DIR)/ard

# Go settings
export GOPROXY := direct
export GONOSUMCHECK := *

# Cross-compilation targets
GOOS   ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# ─────────────────────────────────────────────
# Targets
# ─────────────────────────────────────────────

.PHONY: all build clean deps doctor setup help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

all: setup build ## Full setup + build

setup: submodules deps ## Init submodules and download Go modules

submodules: ## Initialize and update git submodules
	git submodule update --init --recursive

deps: ## Download Go module dependencies
	cd client && go mod download

build: $(SDK_BIN) $(ARD_BIN) ## Build SDK and ARD binaries

$(SDK_BIN): submodules
	cd client && go build -o bin/cdd-sdk ./cmd/cdd-sdk/

$(ARD_BIN): submodules
	cd client && go build -o bin/ard ./cmd/application_reference_design/

clean: ## Remove built binaries
	rm -rf $(BIN_DIR)

# Cross-compilation
build-linux-amd64: submodules ## Cross-compile for Linux x86_64
	cd client && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/cdd-sdk-linux-amd64 ./cmd/cdd-sdk/
	cd client && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ard-linux-amd64 ./cmd/application_reference_design/

build-linux-arm64: submodules ## Cross-compile for Linux ARM64
	cd client && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/cdd-sdk-linux-arm64 ./cmd/cdd-sdk/
	cd client && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/ard-linux-arm64 ./cmd/application_reference_design/

# ─────────────────────────────────────────────
# Doctor — check prerequisites
# ─────────────────────────────────────────────

doctor: ## Check that all build prerequisites are met
	@echo "Checking prerequisites..."
	@echo ""
	@status=0; \
	\
	printf "  go:          "; \
	if command -v go >/dev/null 2>&1; then \
		printf "\033[32m✓\033[0m %s\n" "$$(go version | awk '{print $$3}')"; \
	else \
		printf "\033[31m✗ not found\033[0m — install from https://go.dev/dl/\n"; \
		status=1; \
	fi; \
	\
	printf "  git:         "; \
	if command -v git >/dev/null 2>&1; then \
		printf "\033[32m✓\033[0m %s\n" "$$(git --version | awk '{print $$3}')"; \
	else \
		printf "\033[31m✗ not found\033[0m\n"; \
		status=1; \
	fi; \
	\
	printf "  make:        "; \
	if command -v make >/dev/null 2>&1; then \
		printf "\033[32m✓\033[0m %s\n" "$$(make --version | head -1)"; \
	else \
		printf "\033[31m✗ not found\033[0m\n"; \
		status=1; \
	fi; \
	\
	printf "  submodules:  "; \
	if [ -f models/TR-12-Models/README.md ]; then \
		printf "\033[32m✓\033[0m initialized\n"; \
	else \
		printf "\033[33m○ not initialized\033[0m — run 'make setup'\n"; \
		status=1; \
	fi; \
	\
	printf "  go modules:  "; \
	if [ -d "$$(go env GOPATH)/pkg/mod/github.com/gin-gonic" ] 2>/dev/null; then \
		printf "\033[32m✓\033[0m cached\n"; \
	else \
		printf "\033[33m○ not downloaded\033[0m — run 'make deps'\n"; \
	fi; \
	\
	printf "  network:     "; \
	if curl -sfo /dev/null --max-time 3 https://github.com 2>/dev/null; then \
		printf "\033[32m✓\033[0m github.com reachable\n"; \
	else \
		printf "\033[33m○ github.com unreachable\033[0m — modules must already be cached\n"; \
	fi; \
	\
	echo ""; \
	if [ $$status -eq 0 ]; then \
		echo "All prerequisites met. Run 'make build' to compile."; \
	else \
		echo "Some prerequisites missing. See above."; \
		exit 1; \
	fi
