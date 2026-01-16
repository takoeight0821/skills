# Initial Concept
# Claude Skills

Claude Codeのagent skillを管理する共用リポジトリです。miseを使って各プロジェクトの`.claude/skills`やグローバルの`~/.claude/skills`に展開して利用します。

## 必要条件

# Product Definition

## Core Features
-   **Unified CLI (`skills`):** A custom Go-based CLI tool to manage the entire ecosystem. It replaces the previous `mpvm` tool and shell-based sync scripts. Features include:
    -   **Skill Synchronization:** Sync skills from a central repository to `~/.claude/skills` (global) or `.claude/skills` (project-specific). Supports dry-runs and force updates.
    -   **Multipass VM Management:** VM lifecycle management (launch, start, stop, delete), SSH connectivity with agent forwarding, and direct agent execution.
    -   **Docker Management:** Container lifecycle management, interactive shells, and agent execution.
-   **Docker Development Environment:** A containerized alternative for agent execution, supporting:
    -   Cross-platform compatibility (Linux, macOS Docker Desktop, Rancher Desktop).
    -   SSH agent forwarding for signed Git commits.
    -   Dev Container support for VS Code integration.
-   **Task Automation:** Optional `mise` configuration provided for convenience, wrapping the `skills` CLI commands.

## Target Audience
-   Developers using AI coding agents (Claude Code, Gemini CLI) who need a standardized way to manage custom skills.
-   Teams requiring shared, version-controlled agent skills.
-   Users needing isolated, reproducible environments for running AI agents to prevent host system contamination or to ensure consistent tooling.

## User Experience
-   **CLI-First:** All interactions are driven through the command line using \`mise\` tasks or the \`mpvm\` binary.
-   **Automated Setup:** \`install.sh\` and cloud-init scripts automate the provisioning of environments and configuration, minimizing manual setup.
-   **Transparency:** Sync operations provide detailed previews (dry-runs) before applying changes, ensuring users are always in control of file system modifications.
