.PHONY: proto-format proto-lint proto-gen license test-unit
all: proto-all test-unit

###############################################################################
###                                 Tooling                                 ###
###############################################################################

FILES := $(shell find $(shell go list -f '{{.Dir}}' ./...) -name "*.go" -a -not -name "*.pulsar.go" | sed "s|$(shell pwd)/||g")
license:
	@go-license --config .github/license.yml $(FILES)

###############################################################################
###                                Protobuf                                 ###
###############################################################################

BUF_VERSION=1.42
BUILDER_VERSION=0.15.1

proto-all: proto-format proto-lint proto-gen

proto-format:
	@echo "🤖 Running protobuf formatter..."
	@cd proto
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) format --diff --write
	@echo "✅ Completed protobuf formatting!"

proto-gen:
	@echo "🤖 Generating code from protobuf..."
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		ghcr.io/cosmos/proto-builder:$(BUILDER_VERSION) sh ./proto/generate.sh
	@echo "✅ Completed code generation!"

proto-lint:
	@echo "🤖 Running protobuf linter..."
	@cd proto
	@docker run --rm --volume "$(PWD)":/workspace --workdir /workspace \
		bufbuild/buf:$(BUF_VERSION) lint
	@echo "✅ Completed protobuf linting!"

###############################################################################
###                                 Testing                                 ###
###############################################################################

test-unit:
	@echo "🤖 Running unit tests..."
	@go test -cover -coverprofile=coverage.out -race -v .
	@echo "✅ Completed unit tests!"
