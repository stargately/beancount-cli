#!/bin/sh
# shellcheck shell=dash
# Install beancount-cli — godownloader-style
#
# Usage:
#   curl -sf https://raw.githubusercontent.com/stargately/beancount-cli/main/scripts/install.sh | sh
#   curl -sf https://raw.githubusercontent.com/stargately/beancount-cli/main/scripts/install.sh | VERSION=v0.1.0 sh
#   curl -sf https://raw.githubusercontent.com/stargately/beancount-cli/main/scripts/install.sh | BINDIR=$HOME/.local/bin sh

set -e
set -u

ORIGINAL_PATH="$PATH"

BINARY="beancount-cli"
OWNER="stargately"
REPO="beancount-cli"
BINDIR="${BINDIR:-${HOME}/.local/bin}"
VERSION="${VERSION:-}"

usage() {
  cat <<EOF
Install ${BINARY}

USAGE:
  curl -sf https://raw.githubusercontent.com/${OWNER}/${REPO}/main/scripts/install.sh | sh

ENVIRONMENT VARIABLES:
  VERSION   Release version to install (default: latest)
            e.g. VERSION=v0.1.0
  BINDIR    Directory to install the binary (default: ~/.local/bin)
            e.g. BINDIR=/usr/local/bin
EOF
  exit 2
}

main() {
  need_cmd uname
  need_cmd curl
  need_cmd tar

  OS="$(uname_os)"
  ARCH="$(uname_arch)"

  if [ -z "${VERSION}" ]; then
    VERSION="$(get_latest_version)"
    if [ -z "${VERSION}" ]; then
      echo "error: could not determine latest release version" >&2
      exit 1
    fi
  fi

  # Strip leading 'v' for archive name, but keep it in the tag
  TAG="${VERSION}"
  case "${VERSION}" in
    v*) ;;
    *)  TAG="v${VERSION}" ;;
  esac

  NAME="${BINARY}_${OS}_${ARCH}"
  TARBALL="${NAME}.tar.gz"
  CHECKSUM_FILE="checksums.txt"
  BASE_URL="https://github.com/${OWNER}/${REPO}/releases/download/${TAG}"
  TARBALL_URL="${BASE_URL}/${TARBALL}"
  CHECKSUM_URL="${BASE_URL}/${CHECKSUM_FILE}"

  echo "Installing ${BINARY} ${TAG} for ${OS}/${ARCH}"

  TMPDIR="$(mktemp -d)"
  trap 'rm -rf "${TMPDIR}"' EXIT

  echo "Downloading ${TARBALL_URL}"
  http_download "${TMPDIR}/${TARBALL}" "${TARBALL_URL}"

  echo "Downloading ${CHECKSUM_URL}"
  http_download "${TMPDIR}/${CHECKSUM_FILE}" "${CHECKSUM_URL}"

  verify_checksum "${TMPDIR}" "${TARBALL}" "${CHECKSUM_FILE}"

  echo "Extracting ${TARBALL}"
  tar -xzf "${TMPDIR}/${TARBALL}" -C "${TMPDIR}"

  install_binary "${TMPDIR}"

  if [ "${BINDIR}" = "${HOME}/.local/bin" ]; then
    ensure_user_local_bin_on_path
    warn_shell_path_missing_dir "${BINDIR}"
  fi

  echo "Installed ${BINARY} to ${BINDIR}/${BINARY}"
  "${BINDIR}/${BINARY}" --version 2>/dev/null || true
}

ensure_user_local_bin_on_path() {
  target="$HOME/.local/bin"
  mkdir -p "$target"
  export PATH="$target:$PATH"
  path_line='export PATH="$HOME/.local/bin:$PATH"'
  for rc in "$HOME/.bashrc" "$HOME/.zshrc"; do
    if [ -f "$rc" ] && ! grep -q ".local/bin" "$rc"; then
      echo "$path_line" >> "$rc"
    fi
  done
}

warn_shell_path_missing_dir() {
  dir="${1%/}"
  if [ -z "$dir" ]; then
    return 0
  fi
  case ":${ORIGINAL_PATH}:" in
    *":${dir}:"*) return 0 ;;
  esac
  echo ""
  echo "NOTE: ${dir} was not found in your PATH."
  echo "      It has been added to your current session and appended to ~/.bashrc / ~/.zshrc."
  echo "      Restart your terminal (or run: export PATH=\"${dir}:\$PATH\") to use ${BINARY}."
}

uname_os() {
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  case "${os}" in
    darwin)  echo "darwin" ;;
    linux)   echo "linux" ;;
    msys*|cygwin*|mingw*|windows*) echo "windows" ;;
    *)
      echo "error: unsupported OS: ${os}" >&2
      exit 1
      ;;
  esac
}

uname_arch() {
  arch="$(uname -m)"
  case "${arch}" in
    x86_64|amd64) echo "amd64" ;;
    aarch64|arm64) echo "arm64" ;;
    *)
      echo "error: unsupported architecture: ${arch}" >&2
      exit 1
      ;;
  esac
}

get_latest_version() {
  # GitHub API returns JSON; extract tag_name without needing jq
  url="https://api.github.com/repos/${OWNER}/${REPO}/releases/latest"
  version="$(http_get "${url}" \
    | grep '"tag_name"' \
    | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')"
  echo "${version}"
}

verify_checksum() {
  dir="$1"
  file="$2"
  checksum_file="$3"

  if ! need_cmd_optional sha256sum && ! need_cmd_optional shasum; then
    echo "warning: sha256sum/shasum not found, skipping checksum verification" >&2
    return
  fi

  expected="$(grep " ${file}$" "${dir}/${checksum_file}" | awk '{print $1}')"
  if [ -z "${expected}" ]; then
    echo "error: checksum for ${file} not found in ${checksum_file}" >&2
    exit 1
  fi

  if command -v sha256sum >/dev/null 2>&1; then
    actual="$(sha256sum "${dir}/${file}" | awk '{print $1}')"
  else
    actual="$(shasum -a 256 "${dir}/${file}" | awk '{print $1}')"
  fi

  if [ "${expected}" != "${actual}" ]; then
    echo "error: checksum mismatch for ${file}" >&2
    echo "  expected: ${expected}" >&2
    echo "  actual:   ${actual}" >&2
    exit 1
  fi
  echo "Checksum verified"
}

install_binary() {
  dir="$1"
  src="${dir}/${BINARY}"

  if [ ! -f "${src}" ]; then
    echo "error: binary not found in archive: ${src}" >&2
    exit 1
  fi

  if [ ! -d "${BINDIR}" ]; then
    mkdir -p "${BINDIR}"
  fi

  if [ -w "${BINDIR}" ]; then
    cp "${src}" "${BINDIR}/${BINARY}"
    chmod +x "${BINDIR}/${BINARY}"
  else
    echo "Requires sudo to install to ${BINDIR}"
    sudo cp "${src}" "${BINDIR}/${BINARY}"
    sudo chmod +x "${BINDIR}/${BINARY}"
  fi
}

http_download() {
  local_file="$1"
  url="$2"
  if ! curl -fsSL -o "${local_file}" "${url}"; then
    echo "error: download failed: ${url}" >&2
    exit 1
  fi
}

http_get() {
  url="$1"
  curl -fsSL "${url}"
}

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command not found: $1" >&2
    exit 1
  fi
}

need_cmd_optional() {
  command -v "$1" >/dev/null 2>&1
}

main "$@"
