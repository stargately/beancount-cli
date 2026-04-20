# Repo Migration Guide

## Migrating to a new GitHub repository

The GitHub repository identity is centralized in `internal/config/config.go`:

```go
RepoOwner = "stargately"
RepoName  = "beancount-cli"
RepoSlug  = RepoOwner + "/" + RepoName
```

Changing `RepoName` (and `RepoOwner` if needed) propagates automatically to all Go code: the upgrade checker, the update notifier URL, and the upgrade command.

The following files are **not** covered by those constants and must be updated manually at migration time:

### `scripts/install.sh`
Already parameterized — update the two variables near the top:
```sh
OWNER="stargately"
REPO="beancount-cli"
```

### `.goreleaser.yaml`
Update the project name and Homebrew tap fields:
```yaml
project_name: beancount-cli

brews:
  - repository:
      owner: stargately
      name: homebrew-beancount-cli
```

### `README.md`
Update all installation URLs that reference the raw GitHub path, e.g.:
```
https://raw.githubusercontent.com/stargately/beancount-cli/main/scripts/install.sh
```
