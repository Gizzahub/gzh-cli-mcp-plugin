# .make/test.mk - Test targets

##@ Testing

.PHONY: test test-v test-coverage

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -timeout $(TEST_TIMEOUT) ./...
	@echo "✅ Tests passed"

test-v: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -timeout $(TEST_TIMEOUT) -coverprofile=$(COVERAGE_OUT) ./...
	$(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Coverage report: $(COVERAGE_HTML)"
	$(GO) tool cover -func=$(COVERAGE_OUT) | tail -1
