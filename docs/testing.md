# Testing Guidelines

## Philosophy

Test the **data transformation pipeline**, not mock infrastructure.

The valuable seam is: CLI flags → pure transformation function → exact struct sent to the GraphQL API. Avoid mocking the GraphQL client or HTTP transport unless there is no alternative. Pure function tests are faster, simpler, and catch the real bug surface (nil vs empty slice, pointer semantics, default values, field ordering).

## How to make a command testable

Extract the core logic out of the `RunE` handler into a pure function that accepts plain values and returns a typed struct — no I/O, no `utils.NewAuthedClient()` call:

```go
func buildXxxInput(flag, date string, items []string) (generated.XxxInput, error)
```

The `RunE` handler becomes a thin wiring layer: authenticate → build input → call API.

See `cmd/transaction_add.go` for a worked example (`buildTransactionInput`).

## What to test

Cover every field mapping from CLI input to the output struct:

| Scenario | What to assert |
|----------|---------------|
| All fields provided | Each field equals the expected value |
| Optional string empty | Pointer field is `nil` (not `ptr("")`) |
| Optional string non-empty | Pointer field is non-nil with correct value |
| Slice flag not provided | Field is `[]string{}`, not `nil` (avoids GraphQL serialization bugs) |
| Date flag empty | Field defaults to `time.Now().Format("2006-01-02")` |
| Invalid input format | Error contains expected substring |
| No required items | Error with descriptive message |
| Multiple items | Correct count and order preserved |

## Test structure

Use table-driven tests with a `check` function for struct assertions and a `wantErr` string for error cases:

```go
tests := []struct {
    name    string
    // ... inputs ...
    check   func(t *testing.T, got SomeStruct)
    wantErr string
}{
    {
        name:    "optional field empty becomes nil",
        // ...
        check: func(t *testing.T, got SomeStruct) {
            t.Helper()
            if got.Field != nil {
                t.Errorf("Field: got %v, want nil", got.Field)
            }
        },
    },
    {
        name:    "invalid input returns error",
        // ...
        wantErr: "expected substring",
    },
}
```

Check errors with `strings.Contains(err.Error(), tc.wantErr)` — avoids brittle exact-match coupling.

## Run tests

```
make test
```
