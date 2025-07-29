#!/bin/bash

# Simple download script for gokcat binary
set -e

REPO="philipparndt/gokcat"
ARCH=$(uname -m)
OS="linux"

# Convert architecture names to match GoReleaser output
case $ARCH in
    x86_64)
        ARCH="x86_64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    armv7l)
        ARCH="armv7"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Downloading gokcat for $OS/$ARCH..."

# Get the latest release URL for tar.gz archive
LATEST_URL=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | \
    grep "browser_download_url.*${OS}_${ARCH}.*\.tar\.gz" | \
    head -n 1 | \
    cut -d '"' -f 4)

if [ -z "$LATEST_URL" ]; then
    echo "Warning: Could not get latest release from GitHub API (possibly rate limited)"
    echo "Falling back to known stable release v0.7.6..."
    LATEST_URL="https://github.com/$REPO/releases/download/v0.7.6/gokcat_${OS}_${ARCH}.tar.gz"
    echo "Using fallback URL: $LATEST_URL"
fi

# Create temporary directory
TMP_DIR=$(mktemp -d)
TMP_FILE="$TMP_DIR/gokcat.tar.gz"

echo "Downloading from: $LATEST_URL"
curl -L -o "$TMP_FILE" "$LATEST_URL"

# Extract and install
echo "Extracting..."
tar -xzf "$TMP_FILE" -C "$TMP_DIR"

# Install to /usr/local/bin (or ~/bin if no sudo)
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    echo "Cannot write to $INSTALL_DIR, trying to use sudo..."
    sudo mv "$TMP_DIR/gokcat" "$INSTALL_DIR/gokcat"
    sudo chmod +x "$INSTALL_DIR/gokcat"
else
    mv "$TMP_DIR/gokcat" "$INSTALL_DIR/gokcat"
    chmod +x "$INSTALL_DIR/gokcat"
fi

echo "✅ gokcat installed successfully to $INSTALL_DIR/gokcat"
echo "Run 'gokcat --help' to get started."

# Clean up
rm -rf "$TMP_DIR"
