#!/usr/bin/env sh
# Ohm installer — downloads the latest release binary from GitHub.
# Usage: curl -fsSL https://github.com/derKosi/Ohm/releases/latest/download/install.sh | sh
#
# Works on Linux and macOS. For Windows, download the .zip from the
# releases page and extract ohm.exe manually.
set -eu

OWNER="derKosi"
REPO="Ohm"
BINARY="ohm"

err() {
    echo "Error: $*" >&2
    exit 1
}

# --- Detect platform ---
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux*)  goos="linux" ;;
    Darwin*) goos="darwin" ;;
    *)       err "Unsupported OS: $OS. Download a binary from https://github.com/$OWNER/$REPO/releases" ;;
esac

case "$ARCH" in
    x86_64|amd64)  goarch="x86_64" ;;
    arm64|aarch64) goarch="arm64" ;;
    *)             err "Unsupported architecture: $ARCH" ;;
esac

# --- Find the latest release tag (follow redirect) ---
echo "Fetching latest release..."
if command -v curl >/dev/null 2>&1; then
    LATEST_TAG="$(curl -fsSL -o /dev/null -w '%{url_effective}' "https://github.com/$OWNER/$REPO/releases/latest" | sed 's|.*/tag/||')"
elif command -v wget >/dev/null 2>&1; then
    LATEST_TAG="$(wget -qO- --spider "https://github.com/$OWNER/$REPO/releases/latest" 2>&1 | grep -i location | tail -1 | sed 's|.*/tag/||;s|\r||')"
else
    err "Neither curl nor wget is installed."
fi

[ -z "$LATEST_TAG" ] && err "Could not determine latest release tag."

VERSION="${LATEST_TAG#v}"

# GoReleaser archive names use title-cased OS: Linux, Darwin
case "$goos" in
    linux)  asset_os="Linux"  ;;
    darwin) asset_os="Darwin" ;;
esac

ASSET="${BINARY}_${VERSION}_${asset_os}_${goarch}.tar.gz"
URL="https://github.com/$OWNER/$REPO/releases/download/${LATEST_TAG}/${ASSET}"

echo "Downloading Ohm ${LATEST_TAG} for ${asset_os}/${goarch}..."

# --- Download ---
TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

if command -v curl >/dev/null 2>&1; then
    curl -fsSL -o "$TMPDIR/$ASSET" "$URL"
else
    wget -qO "$TMPDIR/$ASSET" "$URL"
fi

# --- Extract ---
tar -xzf "$TMPDIR/$ASSET" -C "$TMPDIR"

# --- Install ---
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    echo "/usr/local/bin not writable — installing to $INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"
fi

mv "$TMPDIR/$BINARY" "$INSTALL_DIR/$BINARY"
chmod +x "$INSTALL_DIR/$BINARY"

echo ""
echo "✓ Installed Ohm ${LATEST_TAG} to $INSTALL_DIR/$BINARY"
echo ""

# --- PATH hint ---
case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *)
        echo "⚠️  $INSTALL_DIR is not in your PATH."
        echo "   Add it to your shell profile:"
        echo "     export PATH=\"$INSTALL_DIR:\$PATH\""
        echo ""
        ;;
esac

echo "Run 'ohm scan' to scan your system for AI software."
echo "https://github.com/$OWNER/$REPO"
