# Format, build and test the simapp
default: format tidy build test

# Run the help command to see all the available commands
help:
    just --list

###############################################################################
###                                 Golang                                  ###
###############################################################################

# Clean up the go.mod on all modules
tidy:
	@go mod tidy
	@go work sync
	@cd simapp && go mod tidy

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
###                                  Build                                  ###
###############################################################################

# Build the simapp binary
build:
	@echo "🤖 Building simd..."
	@cd simapp && just build
	@echo "✅ Completed build!"

###############################################################################
###                                 Testing                                 ###
###############################################################################

# Run all the tests
test: test-unit

test-unit:
	@echo "🤖 Running unit tests..."
	@go test -cover -coverprofile=coverage.out -race -v
	@echo "✅ Completed unit tests!"

