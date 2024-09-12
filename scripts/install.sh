#!/bin/bash

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Define the installation directory
INSTALL_DIR="$HOME/.local/kthcloud/bin"
BINARY_NAME="kthcloud"
GITHUB_REPO="Phillezi/kthcloud-cli"

# Detect OS and architecture
OS=$(uname -s)
ARCH=$(uname -m)

# Determine OS and ARCH for the binary
case $OS in
  Linux)
    OS="linux"
    ;;
  Darwin)
    OS="darwin"
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

case $ARCH in
  x86_64)
    ARCH="amd64"
    ;;
  aarch64)
    ARCH="arm64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Construct the download URL for the binary
BINARY_URL="https://github.com/$GITHUB_REPO/releases/latest/download/${BINARY_NAME}_${ARCH}_${OS}"

# Function to show a loading spinner
spinner() {
  echo -e "${GREEN}"
  local pid=$1
  local delay=0.1
  local spinstr='|/-\'
  while [ "$(ps a | awk '{print $1}' | grep $pid)" ]; do
    local temp=${spinstr#?}
    printf " [%c]  " "$spinstr"
    local spinstr=$temp${spinstr%"$temp"}
    sleep $delay
    printf "\b\b\b\b\b\b"
  done
  printf "    \b\b\b\b"
  echo -e "${NC}"
}

# Create the install directory if it doesn't exist
mkdir -p "$INSTALL_DIR"

# Download the binary
echo -e "${GREEN}Downloading $BINARY_NAME for $OS $ARCH...${NC}"
curl -fSslL -o "$INSTALL_DIR/$BINARY_NAME" "$BINARY_URL" &  # Start the download in the background
CURL_PID=$!  # Capture the PID of the curl process

# Call the spinner function with the curl process ID
spinner $CURL_PID

# Wait for the curl process to finish and capture its exit status
wait $CURL_PID 2>/dev/null
CURL_STATUS=$?

if [ "$CURL_STATUS" -ne "0" ]; then
  echo -e "${RED}Failed to download the binary... :(${NC}"
  echo "Check if it exists:"
  echo "$BINARY_URL"
  exit 1
fi

# Make the binary executable
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Add binary path to .bashrc or .zshrc for future use
if [[ $SHELL == *"zsh"* ]]; then
  SHELL_CONFIG="$HOME/.zshrc"
else
  SHELL_CONFIG="$HOME/.bashrc"
fi

# Check if the path is already added
if ! grep -q "$INSTALL_DIR" "$SHELL_CONFIG"; then
  echo -e "\n# Created by kthcloud installation script on $(date) \nexport PATH=\$PATH:$INSTALL_DIR \n" >> "$SHELL_CONFIG"
  echo -e "${GREEN}Added $INSTALL_DIR to PATH in $SHELL_CONFIG${NC}"
  echo "Please run 'source $SHELL_CONFIG' or open a new terminal to apply the changes."
else
  echo "Path $INSTALL_DIR already added to $SHELL_CONFIG"
fi

echo ""
echo -e "${GREEN}$BINARY_NAME installed successfully!${NC}"
