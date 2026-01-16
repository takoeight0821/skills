# Specification: Skills CLI (Go Replacement for mise tasks)

## 1. Overview
This track involves creating a new Go-based CLI tool named `skills`. The primary objective is to replace the current shell script and `mise` task-based workflow with a unified, cross-platform binary. This tool will centralize the management of Multipass VMs, Docker containers, skill synchronization, and environment setup.

## 2. Functional Requirements

### 2.1. CLI Structure
- **Name:** `skills`
- **Subcommands:** The CLI must support subcommands to organize functionality logically (e.g., `skills vm launch`, `skills docker start`).

### 2.2. Multipass Management (Replacing `multipass/run.sh`)
The CLI must verify if `multipass` is installed and manage the VM lifecycle:
- `launch`: Create and start the coding agent VM with cloud-init.
- `start`: Start an existing VM.
- `stop`: Stop the VM.
- `delete`: Delete the VM.
- `ssh`: SSH into the VM with agent forwarding.
- `claude`: Run Claude Code in the VM.
- `gemini`: Run Gemini CLI in the VM.
- `status`: Show VM status.
- `mount`: Mount directories to the VM.

### 2.3. Docker Management (Replacing `docker/run.sh`)
The CLI must verify if `docker` is installed and manage the container lifecycle:
- `launch`: Build and start the coding agent container.
- `start`: Start the container.
- `stop`: Stop the container.
- `delete`: Delete the container.
- `ssh`: Open an interactive shell in the container.
- `claude`: Run Claude Code in the container.
- `gemini`: Run Gemini CLI in the container.
- `status`: Show container status.
- `logs`: Show container logs.
- `configure-git`: Re-configure git in the container.

### 2.4. Skill Synchronization (Replacing Sync Scripts)
- Implement logic to synchronize skills from the central repository to local directories (`~/.claude/skills` or `.claude/skills`).
- Support "dry-run" capability.
- Support "force" update.

### 2.5. Environment Setup (Replacing `install.sh`)
- Interactive configuration wizard.
- Setup global configuration (Git user/email, SSH keys).
- Generate/update configuration files.

## 3. Non-Functional Requirements
- **Language:** Go (Golang).
- **Architecture:** Modular design to allow easy addition of new features.
- **Error Handling:** robust error handling and user-friendly error messages.
- **Cross-Platform:** Must work on macOS and Linux.

## 4. Acceptance Criteria
- [ ] The `skills` binary can be built successfully.
- [ ] All `mise` tasks currently defined in `mise.toml` have a corresponding `skills` command.
- [ ] The `install.sh` functionality is fully replicated in the CLI.
- [ ] Users can successfully launch, manage, and connect to both Multipass VMs and Docker containers using the new CLI.
