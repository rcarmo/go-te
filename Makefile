SHELL := /bin/sh

.DEFAULT_GOAL := help

# Generic Makefile template.
# Projects can extend/override these targets.

IMAGE ?= $(notdir $(CURDIR))
TAG ?= latest
FULL_IMAGE := $(IMAGE):$(TAG)
REGISTRY ?= ghcr.io
GHCR_OWNER ?= $(shell whoami)
GHCR_IMAGE := $(REGISTRY)/$(GHCR_OWNER)/$(IMAGE):$(TAG)

.PHONY: help
help: ## Show targets
	@grep -E '^[a-zA-Z0-9_.-]+:.*?##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "%-18s %s\\n", $$1, $$2}'

# =============================================================================
# Docker
# =============================================================================

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(FULL_IMAGE) .

.PHONY: dual-tag
dual-tag: docker-build ## Tag image as ghcr.io/<user>/<image>:<tag>
	docker tag $(FULL_IMAGE) $(GHCR_IMAGE)

.PHONY: tag-ghcr
tag-ghcr: dual-tag ## Convenience alias for dual-tag

# =============================================================================
# Dependency install
# =============================================================================

.PHONY: install
install: ## Install project dependencies
	@$(MAKE) deps

.PHONY: install-dev
install-dev: ## Install dev dependencies
	@$(MAKE) deps

.PHONY: deps
deps: ## Install Go toolchain dependencies
	@command -v go >/dev/null 2>&1 || { echo "Go is required"; exit 1; }
	go mod download
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@if ! command -v gosec >/dev/null 2>&1; then \
		echo "Installing gosec..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	fi

# =============================================================================
# Quality
# =============================================================================

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run --timeout=5m

.PHONY: security
security: ## Run gosec
	gosec ./...

.PHONY: test
test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

.PHONY: coverage
coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...

.PHONY: check
check: ## Run standard validation pipeline
	@$(MAKE) vet && $(MAKE) lint

# =============================================================================
# Release
# =============================================================================

.PHONY: bump-patch
bump-patch: ## Bump patch version in VERSION and create git tag
	@if [ ! -f VERSION ]; then echo "VERSION file not found"; exit 1; fi
	@OLD=$$(cat VERSION); \
	MAJOR=$$(echo $$OLD | cut -d. -f1); \
	MINOR=$$(echo $$OLD | cut -d. -f2); \
	PATCH=$$(echo $$OLD | cut -d. -f3); \
	NEW="$$MAJOR.$$MINOR.$$((PATCH + 1))"; \
	echo $$NEW > VERSION; \
	git add VERSION; \
	git commit -m "Bump version to $$NEW"; \
	git tag "v$$NEW"; \
	echo "Bumped version: $$OLD -> $$NEW (tagged v$$NEW)"

.PHONY: push
push: ## Push commits and current tag to origin
	@TAG=$$(git describe --tags --exact-match 2>/dev/null); \
	git push origin main; \
	if [ -n "$$TAG" ]; then \
		echo "Pushing tag $$TAG..."; \
		git push origin "$$TAG"; \
	else \
		echo "No tag on current commit"; \
	fi

# =============================================================================
# Cleanup
# =============================================================================

.PHONY: clean
clean: ## Remove local build/test artifacts
	@rm -f coverage.out
	@go clean -testcache
