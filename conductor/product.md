# Initial Concept
# Claude Skills

Claude Codeのagent skillを管理する共用リポジトリです。miseを使って各プロジェクトの`.claude/skills`やグローバルの`~/.claude/skills`に展開して利用します。

## 必要条件

# Product Definition

## Core Features
-   **Skill Synchronization:** A robust system to sync skills from a central repository to \`~/.claude/skills\` (global) or \`.claude/skills\` (project-specific) using \`mise\` tasks and shell scripts. Supports dry-runs, force updates, pruning, and git-exclusion.
-   **Multipass VM Manager (\`mpvm\`):** A custom Go-based CLI tool to manage Multipass virtual machines. Features include:
    -   VM lifecycle management (launch, start, stop, delete).
    -   Automated cloud-init configuration for setting up dependencies (Node.js, Claude Code, Gemini CLI).
    -   SSH connectivity with agent forwarding.
    -   Seamless file mounting.
    -   Direct execution of Claude Code and Gemini CLI within the VM.
-   **Docker Development Environment:** A containerized alternative for agent execution, supporting:
    -   Cross-platform compatibility (Linux, macOS Docker Desktop, Rancher Desktop).
    -   SSH agent forwarding for signed Git commits.
    -   Dev Container support for VS Code integration.
-   **Task Automation:** Integrated \`mise\` configuration for unified command execution across skill management, VM operations, and Docker workflows.

## Target Audience
-   Developers using AI coding agents (Claude Code, Gemini CLI) who need a standardized way to manage custom skills.
-   Teams requiring shared, version-controlled agent skills.
-   Users needing isolated, reproducible environments for running AI agents to prevent host system contamination or to ensure consistent tooling.

## User Experience
-   **CLI-First:** All interactions are driven through the command line using \`mise\` tasks or the \`mpvm\` binary.
-   **Automated Setup:** \`install.sh\` and cloud-init scripts automate the provisioning of environments and configuration, minimizing manual setup.
-   **Transparency:** Sync operations provide detailed previews (dry-runs) before applying changes, ensuring users are always in control of file system modifications.
