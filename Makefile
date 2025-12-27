# Makefile for mcp-plugin (MCP Plugin Management CLI)
# ==============================================================================

.DEFAULT_GOAL := help

# Include modular makefiles
include .make/vars.mk
include .make/build.mk
include .make/test.mk
include .make/quality.mk

# ==============================================================================
# Help Target
# ==============================================================================

.PHONY: help
help: ## Display this help message
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "  mcp-plugin - MCP Plugin Management CLI"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST) .make/*.mk 2>/dev/null || true
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "  Quick Start: make build && make test"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""

# ==============================================================================
# Cleanup
# ==============================================================================

.PHONY: clean
clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) coverage.out coverage.html
	@echo "✅ Cleaned"

# ==============================================================================
# Project Info
# ==============================================================================

.PHONY: version
version: ## Display version information
	@echo "Version:    $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $$(go version | awk '{print $$3}')"
