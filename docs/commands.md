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
| `account open`             | Add an account open directive to a ledger            |
| `account close`            | Add an account close directive to a ledger           |
| `transaction add`          | Add a transaction to a ledger                        |
| `balance add`              | Add a balance assertion directive to a ledger        |
| `note add`                 | Add a note directive to a ledger                     |
| `event add`                | Add an event directive to a ledger                   |

## account open

Add an account open directive to a ledger. Optionally restrict the account to one or more currencies by repeating `--currency`.

```sh
beancount-cli account open \
  --ledger user/mybook \
  --account Expenses:Food:Coffee \
  --date 2024-01-01 \
  --currency USD
```

### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--ledger` | yes | — | Ledger full name (find with `ledger list`) |
| `--account` | yes | — | Account name to open (e.g. `Expenses:Food`) |
| `--date` | no | today | Open date, `YYYY-MM-DD` |
| `--currency` | no | — | Allowed currency (repeatable, omit to allow any) |

### Examples

**Open an account with no currency restriction:**

```sh
beancount-cli account open \
  --ledger user/mybook \
  --account Assets:Savings \
  --date 2024-01-01
```

**Open an account restricted to multiple currencies:**

```sh
beancount-cli account open \
  --ledger user/mybook \
  --account Assets:Brokerage \
  --date 2024-01-01 \
  --currency USD \
  --currency EUR
```

## account close

Add an account close directive to a ledger. After this date the account must carry a zero balance.

```sh
beancount-cli account close \
  --ledger user/mybook \
  --account Expenses:Food:Coffee \
  --date 2024-12-31
```

### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--ledger` | yes | — | Ledger full name (find with `ledger list`) |
| `--account` | yes | — | Account name to close |
| `--date` | no | today | Close date, `YYYY-MM-DD` |

### Examples

**Close an account on a specific date:**

```sh
beancount-cli account close \
  --ledger user/mybook \
  --account Liabilities:OldCreditCard \
  --date 2024-06-30
```

**Close an account today (date omitted):**

```sh
beancount-cli account close \
  --ledger user/mybook \
  --account Assets:OldSavings
```

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

## balance add

Add a balance assertion directive to a ledger. This asserts that an account holds exactly the specified amount on a given date.

```sh
beancount-cli balance add \
  --ledger user/mybook \
  --account Assets:Checking \
  --amount 1000.00 \
  --currency USD \
  --date 2024-01-01
```

### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--ledger` | yes | — | Ledger full name (find with `ledger list`) |
| `--account` | yes | — | Account name to assert (e.g. `Assets:Checking`) |
| `--amount` | yes | — | Expected balance amount (e.g. `1000.00`) |
| `--currency` | yes | — | Currency code (e.g. `USD`) |
| `--date` | no | today | Assertion date, `YYYY-MM-DD` |

### Examples

**Assert a checking account balance:**

```sh
beancount-cli balance add \
  --ledger user/mybook \
  --account Assets:Checking \
  --amount 4250.00 \
  --currency USD \
  --date 2024-03-31
```

**Assert a credit card balance (negative amount):**

```sh
beancount-cli balance add \
  --ledger user/mybook \
  --account Liabilities:CreditCard \
  --amount -340.50 \
  --currency USD
```

## note add

Add a note directive to a ledger. Notes attach a free-text comment to an account on a specific date.

```sh
beancount-cli note add \
  --ledger user/mybook \
  --account Assets:Checking \
  --content "Switched to paperless statements" \
  --date 2024-01-01
```

### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--ledger` | yes | — | Ledger full name (find with `ledger list`) |
| `--account` | yes | — | Account to attach the note to |
| `--content` | yes | — | Note text |
| `--date` | no | today | Note date, `YYYY-MM-DD` |

### Examples

**Note an account event:**

```sh
beancount-cli note add \
  --ledger user/mybook \
  --account Assets:Savings \
  --content "Opened high-yield savings account" \
  --date 2024-06-01
```

**Quick note without a date (defaults to today):**

```sh
beancount-cli note add \
  --ledger user/mybook \
  --account Expenses:Medical \
  --content "Insurance reimbursement pending"
```

## event add

Add an event directive to a ledger. Events record a named value that changes over time (e.g. current location, employer, status).

```sh
beancount-cli event add \
  --ledger user/mybook \
  --type "location" \
  --description "New York" \
  --date 2024-01-01
```

### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--ledger` | yes | — | Ledger full name (find with `ledger list`) |
| `--type` | yes | — | Event type (e.g. `location`, `employer`) |
| `--description` | yes | — | Value for this event |
| `--date` | no | today | Event date, `YYYY-MM-DD` |

### Examples

**Record a change of location:**

```sh
beancount-cli event add \
  --ledger user/mybook \
  --type "location" \
  --description "San Francisco" \
  --date 2024-03-01
```

**Record an event today (date omitted):**

```sh
beancount-cli event add \
  --ledger user/mybook \
  --type "employer" \
  --description "Acme Corp"
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

## bean-check

Validate a local `.beancount` file using `bean-check`.

```sh
beancount-cli bean-check path/to/ledger.beancount
```

## bean-query

Run a BQL query against a local `.beancount` file using `bean-query`.

```sh
beancount-cli bean-query path/to/ledger.beancount "SELECT account, sum(position) GROUP BY account"
```
