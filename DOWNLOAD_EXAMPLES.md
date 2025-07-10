# One-liner commands to download the latest gokcat release

## Download and install binary (recommended)
```bash
curl -sSL https://raw.githubusercontent.com/philipparndt/gokcat/main/install.sh | bash
```

## Manual download examples

### Get latest release info
```bash
curl -s https://api.github.com/repos/philipparndt/gokcat/releases/latest | grep browser_download_url
```

### Download specific architecture (AMD64)
```bash
# Binary archive
curl -L -o gokcat.tar.gz "$(curl -s https://api.github.com/repos/philipparndt/gokcat/releases/latest | grep browser_download_url | grep linux_x86_64 | grep tar.gz | cut -d '"' -f 4)"

# Debian package
curl -L -o gokcat.deb "$(curl -s https://api.github.com/repos/philipparndt/gokcat/releases/latest | grep browser_download_url | grep linux_amd64 | grep deb | cut -d '"' -f 4)"
```

### Download specific architecture (ARM64)
```bash
# Binary archive
curl -L -o gokcat.tar.gz "$(curl -s https://api.github.com/repos/philipparndt/gokcat/releases/latest | grep browser_download_url | grep linux_arm64 | grep tar.gz | cut -d '"' -f 4)"

### Simple wget examples
```bash
# Get the latest version first
VERSION=$(curl -s https://api.github.com/repos/philipparndt/gokcat/releases/latest | grep tag_name | cut -d '"' -f 4)

# Then download for your architecture
wget "https://github.com/philipparndt/gokcat/releases/download/$VERSION/gokcat_linux_x86_64.tar.gz"
```

### Extract and use
```bash
# Extract binary
tar -xzf gokcat_linux_x86_64.tar.gz
chmod +x gokcat
sudo mv gokcat /usr/local/bin/
```
