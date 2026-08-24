MIGRATE_VERSION ?= v4.19.0
CODEGEN_VERSION ?= v2.3.0
GOLANGCI_VERSION ?= v2.12.0

.PHONY: setup
setup: ## Install tools
	@echo ">> install migrate"
	go install -tags 'postgres' -ldflags="-X main.Version=$(MIGRATE_VERSION)" github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)
	@echo ">> install codegen"
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(CODEGEN_VERSION)
	@echo ">> install golangci"
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)

.PHONY: test
test: ## Run integration tests
	go test -race -count=1 -shuffle=on -timeout=12m ./...

.PHONY: lint
lint: ## Run golangci linters
	golangci-lint run --output.text.print-issued-lines=false

.PHONY: fmt
fmt: ## Run golangci formatters
	golangci-lint fmt

.PHONY: help
help: ## Show available targets
	@grep -E '^[a-z]+(-[a-z]+)*:.*?## .+$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS=":.*?## "} {printf "%-12s %s\n", $$1, $$2}'
