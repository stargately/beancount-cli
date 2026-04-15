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

## update-schema: re-introspect the live API and regenerate the GraphQL client
update-schema:
	@command -v python3 >/dev/null 2>&1 || { echo "python3 required for schema download"; exit 1; }
	@echo "Introspecting $(BEANCOUNT_API_URL)$(if $(BEANCOUNT_API_URL),,https://beancount.io/api-gateway/)..."
	curl -s -X POST "$${BEANCOUNT_API_URL:-https://beancount.io/api-gateway/}" \
	  -H "Content-Type: application/json" \
	  -d '{"query":"fragment FullType on __Type { kind name description fields(includeDeprecated: true) { name description args { ...InputValue } type { ...TypeRef } isDeprecated deprecationReason } inputFields { ...InputValue } interfaces { ...TypeRef } enumValues(includeDeprecated: true) { name description isDeprecated deprecationReason } possibleTypes { ...TypeRef } } fragment InputValue on __InputValue { name description type { ...TypeRef } defaultValue } fragment TypeRef on __Type { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name } } } } } } } } query IntrospectionQuery { __schema { queryType { name } mutationType { name } subscriptionType { name } types { ...FullType } directives { name description locations args { ...InputValue } } } }"}' \
	  | python3 scripts/introspection-to-sdl.py > graphql/schema.graphql
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
