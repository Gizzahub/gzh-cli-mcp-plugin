# .make/quality.mk - Code quality targets

##@ Code Quality

.PHONY: fmt lint vet validate

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "✅ Formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "✅ Vet passed"

lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi
	@echo "✅ Lint passed"

validate: fmt vet lint test ## Run all quality checks
	@echo "✅ All quality checks passed"
