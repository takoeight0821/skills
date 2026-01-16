# Specification: Remote Plugin/Extension Installation & Command Conversion

## Overview
Implement functionality in the `skills` CLI to interpret Claude plugins and Gemini extensions from remote repositories (URLs) and deploy them to `.claude/skills` (for Claude Code) and the appropriate location for Gemini CLI. This includes converting "command" markdown files (prompt files) into the standard `SKILL.md` format and supporting namespaced installation to avoid collisions.

## Functional Requirements
1.  **Remote Source Support**:
    -   Accept a URL (e.g., GitHub repository) as input.
    -   Support fetching from private repositories (using local git credentials).

2.  **Auto-Detection**:
    -   Identify **Claude Plugins** via `plugin.json` or `marketplace.json`.
    -   Identify **Gemini Extensions** via `manifest.json` or `gemini.json` (TBD).

3.  **Content Processing & Conversion**:
    -   **Standard Skills**: Copy existing `SKILL.md` based directories as-is.
    -   **Commands**: Identify "command" markdown files (files containing prompts but lacking `SKILL.md` metadata).
    -   **Conversion Logic**: Automatically generate a `SKILL.md` wrapper for these commands.
        -   *Metadata*: distinct `name` (derived from filename) and `description` (extracted from file content or default).
        -   *Content*: Embed or reference the original command content.

4.  **Deployment (Claude Code)**:
    -   Target: `.claude/skills/`
    -   Structure: Namespaced directories `<plugin-name>-<skill-name>/`.
    -   Copy entire skill subdirectories (including auxiliary files).

5.  **Deployment (Gemini CLI)**:
    -   Target: User-specified or default Gemini extension path (e.g., `~/.gemini/extensions/` or `.gemini/extensions/`).
    -   Structure: Ensure compatibility with Gemini CLI's expected format (converting `SKILL.md` back to Gemini prompts if necessary, or deploying as-is if Gemini supports it).

6.  **CLI Interface**:
    -   `skills install <url>`: Main entry point.
    -   `--dry-run`: Preview changes.
    -   `--force`: Overwrite existing files.

## Non-Functional Requirements
-   **Robostness**: Handle mixed repositories (containing both skills and commands).
-   **Clarity**: Generated `SKILL.md` files should clearly indicate they were auto-generated from a command.

## Acceptance Criteria
-   [ ] `skills install <url>` fetches a remote repo.
-   [ ] Valid `SKILL.md` directories are deployed to `.claude/skills/<plugin>-<skill>`.
-   [ ] Standalone "command" markdown files are converted to `SKILL.md` format and deployed.
-   [ ] Equivalent deployment occurs for Gemini CLI (location/format verified).
-   [ ] Namespacing prevents collisions between plugins.

## Out of Scope
-   Two-way sync (pushing changes back to the remote).
