# Commands

| Command                    | Description                                          |
|----------------------------|------------------------------------------------------|
| `login`                    | Authenticate via browser (device auth flow)          |
| `logout`                   | Revoke the current session and clear credentials     |
| `whoami`                   | Display the currently authenticated user             |
| `check <file>`             | Validate a local `.beancount` file with `bean-check` |
| `query <file> [query]`     | Query a local `.beancount` file with `bean-query`    |
| `ledger list`              | List your ledgers                                    |
| `ledger create`            | Create a new ledger                                  |
| `ledger delete <fullname>` | Delete a ledger by full name                         |
| `transaction add`          | Add a transaction to a ledger                        |

## transaction add

Append a new transaction to a ledger. Each `--posting` flag represents one leg of the transaction in `account,amount,currency` format.

```sh
beancount-cli transaction add \
  --ledger user/mybook \
  --date 2024-01-15 \
  --payee "Starbucks" \
  --narration "Coffee" \
  --posting "Expenses:Food:Coffee,5.00,USD" \
  --posting "Assets:Checking,-5.00,USD"
```

### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--ledger` | yes | — | Ledger full name (find with `ledger list`) |
| `--posting` | yes (≥1) | — | `account,amount,currency` (repeatable) |
| `--date` | no | today | Transaction date, `YYYY-MM-DD` |
| `--flag` | no | `*` | `*` (cleared) or `!` (pending) |
| `--payee` | no | — | Payee name |
| `--narration` | no | — | Description / narration |
| `--tag` | no | — | Tag (repeatable) |
| `--link` | no | — | Link (repeatable) |

### Examples

**Basic expense with two postings:**

```sh
beancount-cli transaction add \
  --ledger user/mybook \
  --narration "Groceries" \
  --posting "Expenses:Groceries,120.50,USD" \
  --posting "Assets:Checking,-120.50,USD"
```

**Pending transaction (flag `!`) with a tag:**

```sh
beancount-cli transaction add \
  --ledger user/mybook \
  --date 2024-03-01 \
  --flag "!" \
  --payee "Landlord" \
  --narration "March rent" \
  --tag "rent" \
  --posting "Expenses:Rent,1500.00,USD" \
  --posting "Assets:Checking,-1500.00,USD"
```

## ledger create

Create a new Beancount ledger. The name must be in slug format (lowercase letters, digits, hyphens, underscores).

```sh
beancount-cli ledger create --name my-budget --description "Personal finances"
```

## ledger list

List all ledgers owned by the authenticated user.

```sh
beancount-cli ledger list
```

## ledger delete

Delete a ledger permanently by its full name. Find the full name with `ledger list`.

```sh
beancount-cli ledger delete user/my-budget
```

## check

Validate a local `.beancount` file using `bean-check`.

```sh
beancount-cli check path/to/ledger.beancount
```

## query

Run a BQL query against a local `.beancount` file using `bean-query`.

```sh
beancount-cli query path/to/ledger.beancount "SELECT account, sum(position) GROUP BY account"
```
