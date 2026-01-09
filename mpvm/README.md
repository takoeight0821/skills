# mpvm - Multipass VM Manager

A single-binary CLI tool for managing Multipass VMs configured for AI coding agents (Claude Code, Gemini CLI).

## Features

- Single binary with embedded cloud-init configuration
- XDG Base Directory specification compliant
- SSH Agent Forwarding for secure git commit signing
- Auto-mount working directories
- Auto-detects and copies SSH keys for authentication
- Shell completion support (bash, zsh, fish, powershell)

## Installation

### go install (Recommended)

```bash
go install github.com/takoeight0821/skills/mpvm/cmd/mpvm@latest
```

### From Source

```bash
git clone https://github.com/takoeight0821/skills.git
cd skills
go install -C mpvm ./cmd/mpvm
```

### Build Locally

```bash
cd mpvm
make build
./bin/mpvm --help
```

## Quick Start

```bash
# Initialize configuration
mpvm config init

# Create and start VM
mpvm launch

# SSH into VM
mpvm ssh

# Run Claude Code (First time will start authentication flow)
mpvm claude

# Run Gemini CLI
mpvm gemini

# Stop VM
mpvm stop

# Delete VM
mpvm delete
```

## Configuration

Configuration is stored at `~/.config/mpvm/config.toml` on all platforms.

### Example Configuration

```toml
[vm]
name = "coding-agent"
cpus = 2
memory = "4G"
disk = "20G"

[git]
user_name = "Your Name"
user_email = "you@example.com"

[ssh]
signing_key = "~/.ssh/id_ed25519.pub"

[claude]
marketplaces = [
    "anthropics/claude-plugins-official",
    "anthropics/skills"
]
```

## Commands

| Command | Description |
|---------|-------------|
| `mpvm launch` | Create and start VM with cloud-init |
| `mpvm start` | Start stopped VM |
| `mpvm stop` | Stop running VM |
| `mpvm delete` | Delete VM (with confirmation) |
| `mpvm ssh` | SSH into VM with agent forwarding |
| `mpvm claude` | Run Claude Code in VM |
| `mpvm gemini` | Run Gemini CLI in VM |
| `mpvm exec <cmd>` | Execute arbitrary command in VM |
| `mpvm mount [path]` | Mount directory to VM |
| `mpvm umount [path]` | Unmount directory from VM |
| `mpvm status` | Show VM status |
| `mpvm logs` | Show cloud-init logs |
| `mpvm configure-git` | Re-configure git in VM |
| `mpvm config init` | Initialize configuration file |
| `mpvm config show` | Show current configuration |
| `mpvm completion` | Generate shell completion |

## Shell Completion

```bash
# Bash
source <(mpvm completion bash)

# Zsh
mpvm completion zsh > "${fpath[1]}/_mpvm"

# Fish
mpvm completion fish | source
```

## Requirements

- [Multipass](https://multipass.run/) installed
- SSH key for authentication and git signing (default: `~/.ssh/id_ed25519.pub` or `~/.ssh/id_rsa.pub`)

## Development

### Prerequisites
- Go 1.22+
- Make

### Build and Test
```bash
# Run unit tests
make test

# Build binary
make build

# Install locally
make install
```

### Linting
```bash
make lint
```

