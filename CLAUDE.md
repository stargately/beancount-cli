# beancount-cli

## Docs

- [Testing guidelines](docs/testing.md)

## Code organization

- `cmd/` contains only Cobra command definitions and their direct handlers. Reusable logic (API clients, utilities, business logic) belongs in `internal/`.
