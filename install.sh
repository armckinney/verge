#!/usr/bin/env bash
# Verge CLI Installer Script
# Works on Linux (amd64/arm64), macOS (amd64/arm64), and Windows (git bash/WSL)
# Pipe-chainable: curl -sSL https://raw.githubusercontent.com/armckinney/verge/main/install.sh | bash

set -euo pipefail

REPO="armckinney/verge"
PROJECT_NAME="verge"

# Diagnostic trace helper to stderr
log() {
  echo -e "\033[1;32m=>\033[0m $@" >&2
}

log_err() {
  echo -e "\033[1;31mError:\033[0m $@" >&2
}

# 1. Detect OS
OS_RAW="$(uname -s)"
case "${OS_RAW}" in
  Linux*)   OS="linux" ;;
  Darwin*)  OS="darwin" ;;
  MINGW*|MSYS*|CYGWIN*) OS="windows" ;;
  *)
    log_err "Unsupported operating system raw: ${OS_RAW}"
    exit 1
    ;;
esac

# 2. Detect Architecture
ARCH_RAW="$(uname -m)"
case "${ARCH_RAW}" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    log_err "Unsupported CPU architecture raw: ${ARCH_RAW}"
    exit 1
    ;;
esac

log "Detected Environment: OS=${OS}, Architecture=${ARCH}"

# 3. Retrieve Version
if [[ -n "${VERGE_VERSION:-}" ]]; then
  TAG="${VERGE_VERSION}"
  VERSION="${TAG#v}"
  log "Using specified version: v${VERSION}"
else
  log "Querying latest release version from GitHub..."
  REDIRECT_URL="$(curl -sIL -o /dev/null -w '%{url_effective}' "https://github.com/${REPO}/releases/latest")"
  TAG="${REDIRECT_URL##*/}"

  if [[ -z "${TAG}" || "${TAG}" == "latest" ]]; then
    log "Fallback to GitHub API for latest release..."
    TAG="$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"tag_name":\s*"v?([0-9.]+)"/\1/g' | tr -d 'v' || echo "")"
  fi

  # Strip potential 'v' prefix for release download formatting, but keep it for tags
  VERSION="${TAG#v}"
  if [[ -z "${VERSION}" ]]; then
    log_err "Failed to resolve latest version tag."
    exit 1
  fi

  log "Resolved Latest Version: v${VERSION}"
fi

# 4. Formulate Archive URL
ARCHIVE_NAME="${PROJECT_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${ARCHIVE_NAME}"

log "Downloading Archive: ${DOWNLOAD_URL}"

# 5. Download and Extract
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

if ! curl -sSL -o "${TMP_DIR}/${ARCHIVE_NAME}" "${DOWNLOAD_URL}"; then
  # Fallback to prefix-less download link if v-prefix fails
  DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"
  log "Retrying prefix-less download link: ${DOWNLOAD_URL}"
  curl -sSL -o "${TMP_DIR}/${ARCHIVE_NAME}" "${DOWNLOAD_URL}"
fi

tar -xzf "${TMP_DIR}/${ARCHIVE_NAME}" -C "${TMP_DIR}"

BINARY_NAME="${PROJECT_NAME}"
if [[ "${OS}" == "windows" ]]; then
  BINARY_NAME="${PROJECT_NAME}.exe"
fi

# 6. Determine Destination Path
if [[ "${OS}" == "windows" ]]; then
  # Windows git-bash / MSYS environment
  DEST_DIR="/usr/bin"
  if [[ ! -w "${DEST_DIR}" ]]; then
    DEST_DIR="/bin"
  fi
  if [[ ! -w "${DEST_DIR}" ]]; then
    DEST_DIR="."
  fi
else
  # Linux & macOS
  DEST_DIR="/usr/local/bin"
  if [[ ! -w "${DEST_DIR}" ]]; then
    # Try user local bin path if global path is not writable
    DEST_DIR="${HOME}/.local/bin"
    mkdir -p "${DEST_DIR}"
  fi
fi

DEST_PATH="${DEST_DIR}/${BINARY_NAME}"

log "Installing binary to ${DEST_PATH}..."
if [[ -w "${DEST_DIR}" ]]; then
  mv "${TMP_DIR}/${BINARY_NAME}" "${DEST_PATH}"
  chmod +x "${DEST_PATH}"
else
  # Elevate permissions with sudo if directory requires it
  log "Elevation required to write to ${DEST_DIR}. Requesting sudo..."
  sudo mv "${TMP_DIR}/${BINARY_NAME}" "${DEST_PATH}"
  sudo chmod +x "${DEST_PATH}"
fi

log "Verge CLI v${VERSION} has been successfully installed!"
"${DEST_PATH}" version || true
