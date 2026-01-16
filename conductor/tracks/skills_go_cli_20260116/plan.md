# Implementation Plan: Skills CLI Go Replacement

This plan outlines the steps to create a unified Go CLI `skills` to replace the existing shell-script and `mise` task-based workflow.

## Phase 1: Project Initialization & Core Framework [checkpoint: c950212]
- [x] Task: Initialize Go project `skills` in a new directory `cmd/skills`. [commit: c950212]
- [x] Task: Set up `cobra` for command-line argument parsing. [commit: c950212]
- [x] Task: Implement basic `version` command. [commit: c950212]
- [x] Task: Implement configuration management using `viper` to read from `~/.config/skills/config`. [commit: c950212]
- [x] Task: Conductor - User Manual Verification 'Phase 1: Project Initialization' (Protocol in workflow.md) [commit: c950212]

## Phase 2: Environment Setup & Configuration Wizard [checkpoint: c950212]
- [x] Task: Implement `skills configure` command (implemented as `skills config init`). [commit: c950212]
- [x] Task: Port logic from `install.sh` to Go (Git config, SSH key detection). [commit: c950212]
- [x] Task: Implement configuration file generation/updating logic. [commit: c950212]
- [x] Task: Conductor - User Manual Verification 'Phase 2: Environment Setup' (Protocol in workflow.md) [commit: c950212]

## Phase 3: Multipass Management [checkpoint: 981f267]
- [x] Task: Create `vm` subcommand group. [commit: c950212]
- [x] Task: Implement `skills vm launch/start/stop/delete` (wrapping `multipass` commands). [commit: c950212]
- [x] Task: Implement `skills vm ssh/claude/gemini` (SSH execution logic). [commit: 981f267]
- [x] Task: Implement `skills vm status/mount`. [commit: c950212]
- [x] Task: Write tests for Multipass subcommand execution. [commit: ba74811]
- [x] Task: Conductor - User Manual Verification 'Phase 3: Multipass Management' (Protocol in workflow.md) [commit: 981f267]

## Phase 4: Docker Management [checkpoint: 981f267]
- [x] Task: Create `docker` subcommand group. [commit: c950212]
- [x] Task: Implement `skills docker launch/start/stop/delete` (wrapping `docker` commands). [commit: c950212]
- [x] Task: Implement `skills docker ssh/claude/gemini/logs`. [commit: 981f267]
- [x] Task: Implement `skills docker configure-git`. [commit: c950212]
- [x] Task: Write tests for Docker subcommand execution. [commit: ba74811]
- [x] Task: Conductor - User Manual Verification 'Phase 4: Docker Management' (Protocol in workflow.md) [commit: 981f267]

## Phase 5: Skill Synchronization [checkpoint: ]
- [x] Task: Create `sync` command. [commit: 7da74cd]
- [ ] Task: Implement skill synchronization logic (replacing shell sync scripts).
- [ ] Task: Add flags for `--dry-run` and `--force`.
- [ ] Task: Write tests for sync logic.
- [ ] Task: Conductor - User Manual Verification 'Phase 5: Skill Synchronization' (Protocol in workflow.md)

## Phase 6: Final Integration & Cleanup [checkpoint: ]
- [ ] Task: Update project documentation to point to the new CLI.
- [ ] Task: (Optional) Add `mise` tasks that wrap the new `skills` CLI for backward compatibility during transition.
- [ ] Task: Final end-to-end verification.
- [ ] Task: Conductor - User Manual Verification 'Phase 6: Final Integration' (Protocol in workflow.md)