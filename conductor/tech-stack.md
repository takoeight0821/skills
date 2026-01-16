# Technology Stack

## Core Technologies
- **Go:** Used for building the \`mpvm\` CLI tool, providing a performant and statically typed language for system-level operations.
- **Shell (Bash/POSIX):** Used for synchronization scripts, installation routines, and environment setup, ensuring high portability across Unix-like systems.
- **Node.js:** The primary runtime for the AI agents (Claude Code, Gemini CLI).

## Infrastructure & Environment
- **Multipass:** Utilized for creating and managing Ubuntu-based virtual machines for isolated agent execution.
- **Docker:** Provides container-based environments as an alternative to VMs, supporting cross-platform deployment.
- **mise:** Used as a task runner and environment manager to unify developer workflows.

## Development Tools
- **Git:** Essential for version control and for identifying project root/maturity during agent operations.
- **SSH:** Used for secure communication with Multipass VMs and Docker containers, with agent forwarding for Git authentication.
- **Cloud-init:** Used for automated provisioning of Multipass VMs.

## Key Dependencies (Go)
- \`github.com/spf13/cobra\`: For building the CLI interface.
- \`github.com/spf13/viper\`: For configuration management.
