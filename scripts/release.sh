#!/bin/sh
set -e

LATEST_VERSION=""
if command -v gh >/dev/null 2>&1; then
  LATEST_VERSION="$(gh release list --limit 1 --json tagName --jq '.[0].tagName' 2>/dev/null || true)"
fi

if [ -n "$LATEST_VERSION" ]; then
  printf "Version to release (current: %s, e.g. 0.2.0): " "$LATEST_VERSION"
else
  printf "Version to release (e.g. 0.2.0): "
fi
read -r VERSION

# Strip a leading 'v' if the user included one
VERSION="${VERSION#v}"

if [ -z "$VERSION" ]; then
  echo "error: version cannot be empty" >&2
  exit 1
fi

TAG="v${VERSION}"

echo "Tagging ${TAG}"
git tag -a "${TAG}" -m "Release ${TAG}"

echo "Pushing ${TAG}"
git push origin "${TAG}"

echo "Running goreleaser"
GITHUB_TOKEN="$(gh auth token)" \
TAP_GITHUB_TOKEN="$(gh auth token)" \
goreleaser release --clean
