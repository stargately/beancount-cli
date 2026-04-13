# beancount-cli

A Go command-line tool for [Beancount](https://beancount.io). It authenticates via a browser-based device authorization flow and communicates with the Beancount GraphQL backend.

## Commands

| Command   | Description                                      |
|-----------|--------------------------------------------------|
| `login`   | Authenticate via browser (device auth flow)      |
| `logout`  | Revoke the current session and clear credentials |
| `whoami`  | Display the currently authenticated user         |

## Installation

**Prerequisites:** Go 1.23+

```sh
# Build the binary
make build

# Binary is output to ./dist/beancount-cli
./dist/beancount-cli --help

# Optional: move to PATH
cp dist/beancount-cli /usr/local/bin/beancount-cli
```

## Environment Variables

| Variable                 | Default                              | Description                             |
|--------------------------|--------------------------------------|-----------------------------------------|
| `BEANCOUNT_API_URL`      | `https://beancount.io/api-gateway/`  | GraphQL API endpoint                    |
| `BEANCOUNT_DASHBOARD_URL`| `https://beancount.io`               | Dashboard URL (used for the auth flow)  |

## Connecting to a Local Dev Server

To test against your local development environment, start both the backend and dashboard, then point the CLI at them via environment variables.

**1. Start the backend** (GraphQL API on port 4104):

```sh
cd backend-cluster/backend-v2
yarn start
```

**2. Start the dashboard** (required for the browser-based login flow, port 5173):

```sh
cd beancount-dashboard
yarn dev
```

**3. Configure the CLI and log in:**

```sh
export BEANCOUNT_API_URL=http://localhost:4104/
export BEANCOUNT_DASHBOARD_URL=http://localhost:5173

./dist/beancount-cli login
```

The `login` command will open `http://localhost:5173/auth/login/device?session_id=...` in your browser. Approve the session there and the CLI will save a token to `~/.beancount/credentials.json`.

## Credentials

Credentials are stored locally at `~/.beancount/credentials.json` (permissions `0600`):

```json
{
  "token": "<jwt>",
  "expireAt": "2026-12-31T00:00:00Z"
}
```

Run `beancount-cli logout` to revoke the token and delete this file.

## Development

### Make Targets

| Target          | Description                                                      |
|-----------------|------------------------------------------------------------------|
| `make build`    | Compile binary to `./dist/beancount-cli`                         |
| `make lint`     | Run `go vet ./...`                                               |
| `make codegen`  | Regenerate `generated/genqlient.go` from GraphQL operations      |
| `make update-schema` | Re-introspect the live API, update schema, and run codegen  |
| `make clean`    | Remove `./dist/`                                                 |

### GraphQL Files

There are two distinct GraphQL files with different roles:

| File | Role | Edited by |
|------|------|-----------|
| `graphql/schema.graphql` | Full copy of the server's type system — every type, enum, query, and mutation the backend exposes. Used by the code generator as a reference to validate operations and derive Go types. | Never — overwritten by `make update-schema` |
| `graphql/operations.graphql` | The specific queries and mutations this CLI actually calls. A small, hand-authored subset of what the schema allows. | You, when adding new API calls |

At code-gen time (`make codegen`), genqlient validates every operation in `operations.graphql` against `schema.graphql` — referencing a field that doesn't exist on the server is a compile-time error. It then generates fully-typed Go functions into `generated/genqlient.go`.

```
schema.graphql        ← what the server CAN do  (generated, never edit)
operations.graphql    ← what the CLI WANTS to do (hand-authored)
        ↓
    make codegen
        ↓
generated/genqlient.go ← type-safe Go functions  (generated, never edit)
```

### Adding a GraphQL Operation

To call a new backend endpoint from the CLI:

**1. Check `graphql/schema.graphql`** to find the query or mutation name and its available fields. For example, to fetch ledgers:

```graphql
# In schema.graphql (for reference only — do not edit)
type Query {
  ledgerList: [Ledger!]!
}
type Ledger {
  id: ID!
  name: String!
  ...
}
```

**2. Add the operation to `graphql/operations.graphql`:**

```graphql
query GetLedgerList {
  ledgerList {
    id
    name
  }
}
```

**3. Run codegen** to generate the Go client function:

```sh
make codegen
```

This produces a `GetLedgerList(ctx, client)` function in `generated/genqlient.go` that you can call directly in your command code.

**4. Use it in a command** (e.g. `cmd/ledgers.go`):

```go
client := gqlclient.NewAuthenticatedClient(cfg.APIURL, creds.Token)
resp, err := generated.GetLedgerList(cmd.Context(), client)
```

### Adding a New Command

1. Create `cmd/<name>.go` implementing a `cobra.Command`
2. Register it in `cmd/root.go` with `rootCmd.AddCommand(...)`
3. Add any required GraphQL operations to `graphql/operations.graphql` (see above)
4. Run `make codegen` to regenerate `generated/genqlient.go`

### Updating the GraphQL Schema

When the backend schema changes, re-introspect and regenerate the client:

```sh
# Against production (default)
make update-schema

# Against local dev server
BEANCOUNT_API_URL=http://localhost:4104/ make update-schema
```

Requires Python 3 (used by `scripts/introspection-to-sdl.py` to convert introspection JSON to SDL).

### Project Structure

```
beancount-cli/
├── cmd/                  # Command implementations (login, logout, whoami)
├── internal/
│   ├── config/           # Environment variable config
│   ├── credentials/      # Token storage (~/.beancount/credentials.json)
│   └── gqlclient/        # Authenticated GraphQL HTTP client
├── graphql/
│   ├── schema.graphql    # Introspected server schema (auto-generated)
│   └── operations.graphql# CLI GraphQL operations
├── generated/
│   └── genqlient.go      # Auto-generated type-safe GraphQL client
├── scripts/
│   └── introspection-to-sdl.py
├── genqlient.yaml        # Codegen config
└── Makefile
```
