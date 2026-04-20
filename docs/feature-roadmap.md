# beancount-cli Feature Roadmap

## Overview

This document outlines planned features to extend `beancount-cli` so users can interact with their full [beancount.io](https://beancount.io) ledger from the terminal without needing the dashboard.

**Implementation stack:** Go + Cobra + genqlient (GraphQL)  
**Pattern:** add operations to `graphql/operations.graphql` → `go generate ./...` → implement `cmd/*.go`

---

## Current State

| Area | Commands |
|---|---|
| Auth | `login`, `logout`, `whoami`, `upgrade` |
| Ledger | `ledger create/list/delete` |
| Entries | `account open/close`, `transaction add`, `balance add`, `note add`, `event add` |
| Local tools | `bean-check <file>`, `bean-query <file>` |

**Key gaps:** no file management, no read access, no reports, no remote BQL, missing entry types.

---

## Tier 1 — Quick Wins

### 1. Missing Entry Directives

Four entry types are supported by the API but have no CLI commands. Each follows the same pattern as the existing `balance add` command.

| Command | GraphQL Mutation | Required Flags |
|---|---|---|
| `price add` | `addEntryPrice` | `--ledger`, `--currency`, `--amount <num,cur>`, `--date` |
| `commodity add` | `addEntryCommodity` | `--ledger`, `--currency`, `--date` |
| `budget add` | `addEntryBudget` | `--ledger`, `--account`, `--amount <num,cur>`, `--interval DAILY\|WEEKLY\|MONTHLY\|QUARTERLY\|YEARLY` |
| `document add` | `addEntryDocument` | `--ledger`, `--account`, `--filename`, `--date`, `--tag`, `--link` |

**New files:** `cmd/price.go`, `cmd/price_add.go`, `cmd/commodity.go`, `cmd/commodity_add.go`, `cmd/budget.go`, `cmd/budget_add.go`, `cmd/document.go`, `cmd/document_add.go`

**Example usage:**
```bash
beancount-cli price add --ledger user/mybook --currency BTC --amount 60000,USD
beancount-cli commodity add --ledger user/mybook --currency AAPL
beancount-cli budget add --ledger user/mybook --account Expenses:Food --amount 500,USD --interval MONTHLY
beancount-cli document add --ledger user/mybook --account Assets:Bank --filename receipts/2024-01-01.pdf
```

---

### 2. Ledger Metadata Management

Extends the `ledger` command with three new subcommands.

| Command | GraphQL Mutation | Notes |
|---|---|---|
| `ledger update <fullname>` | `updateLedger` | flags: `--name`, `--description`, `--private` |
| `ledger star <fullname>` | `starLedger` | prints confirmation |
| `ledger unstar <fullname>` | `unstarLedger` | prints confirmation |

**New files:** `cmd/ledger_update.go`, `cmd/ledger_star.go`, `cmd/ledger_unstar.go`

**Example usage:**
```bash
beancount-cli ledger update user/mybook --description "Personal finances 2024"
beancount-cli ledger star user/mybook
```

---

## Tier 2 — Core Read & Write Access

### 3. File Management (`file` subcommand)

Users can view and edit their beancount source files directly from the terminal.

| Command | GraphQL Op | Notes |
|---|---|---|
| `file list <ledger> [--path /]` | `getLedgerDirContent` | table: name, type, size, sha |
| `file view <ledger> <path>` | `getLedgerFile` | prints content to stdout (pipeable) |
| `file create <ledger> <path> --message <msg>` | `createLedgerFile` | reads content from stdin |
| `file update <ledger> <path> --message <msg>` | `getLedgerFile` + `updateLedgerFile` | auto-fetches current SHA; reads new content from stdin |
| `file delete <ledger> <path> --message <msg>` | `getLedgerFile` + `deleteLedgerFile` | auto-fetches SHA; asks for confirmation |
| `file rename <ledger> <old> <new> --message <msg>` | `renameLedgerFile` | |

> **Implementation note:** `updateLedgerFile` and `deleteLedgerFile` require the current file's SHA. These commands first call `getLedgerFile` to fetch it automatically, so the user doesn't need to provide it manually.

**New files:** `cmd/file.go`, `cmd/file_list.go`, `cmd/file_view.go`, `cmd/file_create.go`, `cmd/file_update.go`, `cmd/file_delete.go`, `cmd/file_rename.go`

**Example usage:**
```bash
beancount-cli file list user/mybook
beancount-cli file view user/mybook main.bean
cat new-entries.bean | beancount-cli file create user/mybook 2024/january.bean --message "add january"
beancount-cli file view user/mybook main.bean | sed 's/foo/bar/' | beancount-cli file update user/mybook main.bean --message "fix typo"
```

---

### 4. Transaction List (`transaction list`)

Browse and search ledger transactions without opening the dashboard.

**Command:** `transaction list <ledger> [--account] [--from] [--to] [--payee] [--tag] [--limit 20] [--format text|json]`

- Calls `getLedgerJournal` with a `JournalQueryInput`
- Output: table with date, flag, payee, narration, accounts, amount
- `--format json` for scripting/piping

**New file:** `cmd/transaction_list.go`

**Example usage:**
```bash
beancount-cli transaction list user/mybook --account Assets:Checking --limit 10
beancount-cli transaction list user/mybook --from 2024-01-01 --to 2024-03-31
beancount-cli transaction list user/mybook --payee Starbucks --format json
```

---

### 5. Journal Entry Edit & Delete (`journal` subcommand)

Fix mistakes without editing files manually.

| Command | GraphQL Mutation | Notes |
|---|---|---|
| `journal delete <ledger> <entry-hash>` | `deleteLedgerEntrySourceSlice` | fetches `sha256sum` automatically via `getLedgerEntryContext`; asks for confirmation |
| `journal edit <ledger> <entry-hash>` | `updateLedgerEntrySourceSlice` | fetches current source slice → opens in `$EDITOR` → submits on save |

**New files:** `cmd/journal.go`, `cmd/journal_delete.go`, `cmd/journal_edit.go`

**Example usage:**
```bash
# Get the hash from transaction list
beancount-cli journal delete user/mybook abc123def456
beancount-cli journal edit user/mybook abc123def456  # opens $EDITOR
```

---

### 6. Financial Reports (`report` subcommand)

Surface the most-used financial views.

| Command | GraphQL Query | Notes |
|---|---|---|
| `report balance-sheet <ledger> [--time] [--conversion USD]` | `getLedgerBalanceSheet` | tree: Assets / Liabilities / Equity |
| `report income <ledger> [--time] [--interval monthly]` | `getLedgerIncomeStatement` | Income vs Expenses |
| `report trial-balance <ledger> [--time]` | `getLedgerTrialBalance` | all accounts + balances |

**New files:** `cmd/report.go`, `cmd/report_balance_sheet.go`, `cmd/report_income.go`, `cmd/report_trial_balance.go`

**Example usage:**
```bash
beancount-cli report balance-sheet user/mybook
beancount-cli report income user/mybook --time 2024
beancount-cli report trial-balance user/mybook --time 2024-Q1
```

---

### 7. Remote BQL Query (`query run`)

Execute BQL against a remote ledger (as opposed to `bean-query` which operates on local files).

**Command:** `query run <ledger> <bql>` or `query run <ledger>` (interactive mode)

- Calls `queryShell(ledgerId, query)`
- Output: formatted table, `--format csv` for export

**New file:** `cmd/query_run.go` (extends existing `query.go` parent)

**Example usage:**
```bash
beancount-cli query run user/mybook "SELECT account, sum(position) GROUP BY account"
beancount-cli query run user/mybook "SELECT date, narration, position WHERE account ~ 'Expenses'"
beancount-cli query run user/mybook  # interactive REPL
```

---

## Tier 3 — Power User & Team Features

### 8. Ledger Attributes

Quick lookups useful for scripting and shell autocomplete.

| Command | GraphQL Query |
|---|---|
| `ledger accounts <fullname>` | `getLedgerAccounts` |
| `ledger tags <fullname>` | `getLedgerTags` |
| `ledger payees <fullname>` | `getLedgerPayees` |
| `ledger currencies <fullname>` | `getLedgerCurrencies` |

**Example usage:**
```bash
beancount-cli ledger accounts user/mybook | grep Expenses
beancount-cli ledger payees user/mybook
```

---

### 9. SSH Key Management (`key` subcommand)

| Command | GraphQL Op | Notes |
|---|---|---|
| `key list` | `listPublicKeys` | table: id, title, fingerprint |
| `key add --title <name> [--key <pubkey>]` | `createPublicKey` | reads from `--key` flag, file, or stdin |
| `key delete <key-id>` | `deletePublicKey` | asks for confirmation |

**New files:** `cmd/key.go`, `cmd/key_list.go`, `cmd/key_add.go`, `cmd/key_delete.go`

**Example usage:**
```bash
beancount-cli key list
beancount-cli key add --title "MacBook" --key ~/.ssh/id_ed25519.pub
beancount-cli key delete key_abc123
```

---

### 10. Collaborator Management (`collaborator` subcommand)

| Command | GraphQL Op | Notes |
|---|---|---|
| `collaborator list <ledger>` | `listLedgerCollaborators` | table: username, permission |
| `collaborator add <ledger> <username> [--permission read]` | `addOrUpdateLedgerCollaborator` | default: read |
| `collaborator remove <ledger> <username>` | `deleteLedgerCollaborator` | asks for confirmation |

**New files:** `cmd/collaborator.go`, `cmd/collaborator_list.go`, `cmd/collaborator_add.go`, `cmd/collaborator_remove.go`

**Example usage:**
```bash
beancount-cli collaborator list user/mybook
beancount-cli collaborator add user/mybook alice --permission write
beancount-cli collaborator remove user/mybook alice
```

---

## Recommended Implementation Order

1. `price add`, `commodity add`, `budget add`, `document add` — copy-paste from existing entry commands
2. `ledger update`, `ledger star/unstar` — trivial, same pattern as `ledger create`
3. `file list`, `file view` — read-only file access
4. `transaction list` — most-used read path
5. `file create/update/delete/rename` — file writes (after read path is validated)
6. `report balance-sheet`, `report income`, `report trial-balance`
7. `query run` — remote BQL
8. `journal delete`, `journal edit`
9. `ledger accounts/tags/payees/currencies`
10. `key list/add/delete`, `collaborator list/add/remove`

---

## GraphQL Codegen Workflow

For each new feature, add operations to `graphql/operations.graphql`, then regenerate:

```bash
go generate ./...   # regenerates generated/genqlient.go
make build          # compile binary
make lint           # check for issues
```
