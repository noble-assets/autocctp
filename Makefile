.PHONY: proto-format proto-lint proto-gen license build

###############################################################################
###                                Protobuf                                 ###
###############################################################################

BUF_VERSION=1.50
BUILDER_VERSION=0.15.3

proto-all: proto-format proto-lint proto-gen

proto-format:
	@echo "🤖 Running protobuf formatter..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) format --diff --write
	@echo "✅ Completed protobuf formatting!"

proto-lint:
	@echo "🤖 Running protobuf linter..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) lint
	@echo "✅ Completed protobuf linting!"

proto-gen:
	@echo "🤖 Generating code from protobuf..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		ghcr.io/cosmos/proto-builder:$(BUILDER_VERSION) sh ./proto/generate.sh
	@echo "✅ Completed code generation!"


###############################################################################
###                                 Tooling                                 ###
###############################################################################

goimports_reviser=github.com/incu6us/goimports-reviser/v3
gofumpt_cmd=mvdan.cc/gofumpt
golangci_lint_cmd=github.com/golangci/golangci-lint/cmd/golangci-lint

FILES := $(shell find . -name "*.go" -not -path "./simapp/*" -not -name "*.pb.go" -not -name "*.pb.gw.go" -not -name "*.pulsar.go")
license:
	@echo "🤖 Adding license to files..."
	@go-license --config .github/license.yaml $(FILES)
	@echo "✅ Completed license added!"

PREFIXES="github.com/cosmos,cosmossdk.io,github.com/cometbft,github.com/grpc-ecosystem"
format:
	@echo "🤖 Running formatters..."
	@go run $(goimports_reviser) -company-prefixes $PREFIXES -excludes 'utils/tools.go' -rm-unused -set-alias ./...
	@go run $(gofumpt_cmd) -l -w .
	@echo "✅ Completed formatting!"

lint:
	@echo "🤖 Running linter..."
	@go run $(golangci_lint_cmd) run --timeout=10m
	@echo "✅ Completed linting!"


###############################################################################
###                                 Testing                                 ###
###############################################################################

test-unit:
	@echo "🤖 Running unit tests for types package..."
	@go test -v ./types/...
	@echo "✅ Completed unit tests!"
