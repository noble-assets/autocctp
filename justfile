default: format lint test

help:
    just --list

###############################################################################
###                          Formatting & Linting                           ###
###############################################################################

gofumpt_cmd := "mvdan.cc/gofumpt"
golangci_lint_cmd := "github.com/golangci/golangci-lint/cmd/golangci-lint"

format:
	@echo "🤖 Running formatter..."
	@go run {{gofumpt_cmd}} -l -w .
	@echo "✅ Completed formatting!"

lint:
    @echo "🔍 Running linter..."
    @go run {{golangci_lint_cmd}} run --timeout=10m
    @echo "✅ Completed linting!"


###############################################################################
###                                 Testing                                 ###
###############################################################################

test: test-unit

test-unit:
	@echo "🤖 Running unit tests..."
	@go test -cover -coverprofile=coverage.out -race -v
	@echo "✅ Completed unit tests!"