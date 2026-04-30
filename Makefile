BINARY := beancount-cli
BUILD_DIR := dist

.PHONY: build clean codegen update-schema lint format test release help

## help: show this help message
help:
	@echo "Usage: make <target>"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /' | column -t -s ':'

## build: compile the binary into ./dist/beancount-cli
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) .

## codegen: regenerate the GraphQL client from graphql/schema.graphql
codegen:
	go generate ./...

## update-schema: download the SDL schema from the live API and regenerate the GraphQL client
update-schema:
	@echo "Downloading schema from https://beancount.io/api-gateway/schema.graphql..."
	curl -sf "https://beancount.io/api-gateway/schema.graphql" > graphql/schema.graphql
	$(MAKE) codegen
	@echo "Schema updated and client regenerated."

## test: run all tests
test:
	go test ./...

## format: format Go source files with gofmt
format:
	gofmt -w .

## lint: run go vet
lint:
	go vet ./...

## release: interactively tag a new version and publish via goreleaser
release:
	@sh scripts/release.sh

## clean: remove build artifacts
clean:
	rm -rf $(BUILD_DIR)
