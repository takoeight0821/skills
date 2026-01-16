# skills - Unified CLI for Coding Agent Environments

A single-binary CLI tool for managing Multipass VMs and Docker containers configured for AI coding agents (Claude Code, Gemini CLI).

## Features

- **Unified CLI** - One tool for both VMs and containers
- Single binary with embedded cloud-init configuration
- XDG Base Directory specification compliant
- SSH Agent Forwarding for secure git commit signing
- Auto-mount working directories
- Auto-detects and copies SSH keys for authentication
- Shell completion support (bash, zsh, fish, powershell)

## Installation

### go install (Recommended)

```bash
go install github.com/takoeight0821/skills/jig/cmd/skills@latest
```

### From Source

```bash
git clone https://github.com/takoeight0821/skills.git
cd skills/jig
make install
```

### Build Locally

```bash
cd jig
make build
./bin/skills --help
```

## Quick Start

### VM (Multipass)

```bash
# Initialize configuration
skills config init

# Create and start VM
skills vm launch

# SSH into VM
skills vm ssh

# Run Claude Code
skills vm claude

# Run Gemini CLI
skills vm gemini

# Stop VM
skills vm stop

# Delete VM
skills vm delete
```

### Docker

```bash
# Build image and start container
skills docker launch

# Interactive shell
skills docker ssh

# Run Claude Code
skills docker claude

# Run Gemini CLI
skills docker gemini

# Stop container
skills docker stop

# Delete container and image
skills docker delete
```

## Configuration

Configuration is stored at `~/.config/skills/config.toml` on all platforms.

### Example Configuration

```toml
[vm]
name = "coding-agent"
cpus = 2
memory = "4G"
disk = "20G"

[docker]
container_name = "coding-agent-docker"
image_name = "coding-agent:latest"
cpus = "2"
memory = "4g"

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

### VM Commands (Multipass)

| Command | Description |
|---------|-------------|
| `skills vm launch` | Create and start VM with cloud-init |
| `skills vm start` | Start stopped VM |
| `skills vm stop` | Stop running VM |
| `skills vm delete` | Delete VM (with confirmation) |
| `skills vm ssh` | SSH into VM with agent forwarding |
| `skills vm claude` | Run Claude Code in VM |
| `skills vm gemini` | Run Gemini CLI in VM |
| `skills vm exec <cmd>` | Execute arbitrary command in VM |
| `skills vm mount [path]` | Mount directory to VM |
| `skills vm umount [path]` | Unmount directory from VM |
| `skills vm status` | Show VM status |
| `skills vm logs` | Show cloud-init logs |
| `skills vm configure-git` | Re-configure git in VM |

### Docker Commands

| Command | Description |
|---------|-------------|
| `skills docker launch` | Build image and start container |
| `skills docker start` | Start existing container |
| `skills docker stop` | Stop container |
| `skills docker delete` | Delete container, image, and volume |
| `skills docker ssh` | Interactive shell in container |
| `skills docker claude` | Run Claude Code in container |
| `skills docker gemini` | Run Gemini CLI in container |
| `skills docker status` | Show container status |
| `skills docker logs` | Show container logs |
| `skills docker configure-git` | Re-configure git in container |

### Sync Commands

| Command | Description |
|---------|-------------|
| `skills sync --global` | Sync skills to ~/.claude/skills |
| `skills sync --project` | Sync skills to .claude/skills |
| `skills sync --dry-run` | Preview changes without applying |
| `skills sync --force` | Overwrite existing files |
| `skills sync --source <path>` | Specify custom source directory |

### Config Commands

| Command | Description |
|---------|-------------|
| `skills config init` | Initialize configuration file |
| `skills config show` | Show current configuration |
| `skills config path` | Show configuration file path |
| `skills completion` | Generate shell completion |

## Shell Completion

```bash
# Bash
source <(skills completion bash)

# Zsh
skills completion zsh > "${fpath[1]}/_skills"

# Fish
skills completion fish | source
```

## Requirements

- [Multipass](https://multipass.run/) for VM commands
- [Docker](https://www.docker.com/products/docker-desktop/) for container commands
- SSH key for authentication and git signing

## Development

### Prerequisites
- Go 1.22+
- Make

### Build and Test
```bash
make test    # Run unit tests
make build   # Build binary
make install # Install locally
```

### Linting
```bash
make lint
```
