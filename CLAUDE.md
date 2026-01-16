# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

> **Project Context:** Please refer to [conductor/index.md](conductor/index.md) for the most up-to-date documentation, workflow, and architecture details.

This repository manages Claude Code agent skills and provides a Multipass-based development environment for AI coding agents (Claude Code, Gemini CLI). Skills are synchronized to `~/.claude/skills` (global) or `.claude/skills` (project-level) using mise tasks.

## Prerequisites

- [mise](https://mise.jdx.dev/) installed
- Git
- [Multipass](https://multipass.run/) (for VM-based development)

## Common Commands

### Installation and Setup
```bash
./install.sh                    # Install mise tasks to ~/.config/mise/config.toml
```

### Multipass Development Environment
```bash
mise run tk:vm-launch              # Create and start VM with cloud-init
mise run tk:vm-start               # Start existing VM
mise run tk:vm-stop                # Stop VM
mise run tk:vm-delete              # Delete VM
mise run tk:vm-ssh                 # SSH with agent forwarding
mise run tk:vm-claude              # Run Claude Code in VM
mise run tk:vm-gemini              # Run Gemini CLI in VM
mise run tk:vm-status              # Show VM status
mise run tk:vm-mount <path>        # Mount directory to VM

# Or use run.sh directly:
./multipass/run.sh launch       # Create and start VM
./multipass/run.sh ssh          # SSH with agent forwarding
./multipass/run.sh claude       # Run Claude Code
./multipass/run.sh stop         # Stop VM
```

### Skills Sync (when using the sync feature)
```bash
mise run tk:update-shared-skills         # Pull latest from remote
mise run tk:sync-skills-global           # Preview global sync (dry-run)
mise run tk:sync-skills-global-apply     # Apply global sync
mise run tk:sync-skills-project          # Preview project sync (dry-run)
mise run tk:sync-skills-project-apply    # Apply project sync
```

## Architecture

### Directory Structure
- `skills/` - Shareable Claude Code skills (SKILL.md files)
- `multipass/` - Multipass VM environment for coding agents
  - `cloud-init.yaml` - VM initialization (Ubuntu 24.04, Node.js, Claude Code, Gemini CLI)
  - `run.sh` - Wrapper script for VM management
  - `lib/common.sh` - Shared shell functions
  - `config.example` - Example configuration file
- `.claude-plugin/` - Claude Code plugin configuration
- `install.sh` - Adds mise tasks to global config

### Sync System
Skills are synchronized via `bin/sync-skills.sh` (referenced in mise.sample.toml):
- Default mode is dry-run; use `--apply` or `*-apply` tasks to execute
- Tracking via `.skills-manifest` files
- Options: `--force` (overwrite), `--prune` (remove deleted), `--exclude` (git exclude)
- Config files: `~/.config/skills/config` (global), `.claude/.skills.conf` (project)

### Multipass VM Features
- **SSH Agent Forwarding**: Secure git commit signing without copying private keys
- **cloud-init**: Automatic VM configuration at launch
- **Git SSH Signing**: Pre-configured for GitHub Verified commits
- **Skills Sync**: Skills are copied to VM at launch

## Adding New Skills

Create a directory under `skills/` with a `SKILL.md` file:
```bash
mkdir -p skills/my-skill
# Create skills/my-skill/SKILL.md with YAML frontmatter (name, description) and skill content
```

Skill files use YAML frontmatter for metadata (name, description) followed by markdown content describing the skill workflow.

## VM Configuration

Configuration is loaded from these files (later overrides earlier):

| Priority | File | Description |
|----------|------|-------------|
| 1 (low) | `~/.config/skills/config` | Global configuration |
| 2 (high) | `.skills.conf` | Project configuration (current directory) |

### Setup

```bash
# Copy example to global config
mkdir -p ~/.config/skills
cp multipass/config.example ~/.config/skills/config
# Edit and set GIT_USER_NAME, GIT_USER_EMAIL
```

### Available Variables

```bash
# VM resources
MULTIPASS_VM_NAME=coding-agent    # VM name (default: coding-agent)
MULTIPASS_VM_CPUS=2               # CPU count (default: 2)
MULTIPASS_VM_MEMORY=4G            # Memory size (default: 4G)
MULTIPASS_VM_DISK=20G             # Disk size (default: 20G)

# Git configuration (required for signed commits)
GIT_USER_NAME="Your Name"
GIT_USER_EMAIL="you@example.com"

# SSH key for signing (default: ~/.ssh/id_ed25519.pub)
SSH_SIGNING_KEY="$HOME/.ssh/id_ed25519.pub"
```

### Project-specific VM

To use a different VM per project, create `.skills.conf` in the project directory:

```bash
# In your project directory
echo 'MULTIPASS_VM_NAME=my-project-vm' > .skills.conf
```
