#!/bin/sh
set -e

printf "Version to release (e.g. 0.2.0): "
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
