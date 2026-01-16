# Implementation Plan: Skills CLI Go Replacement

This plan outlines the steps to create a unified Go CLI `skills` to replace the existing shell-script and `mise` task-based workflow.

## Phase 1: Project Initialization & Core Framework [checkpoint: ]
- [ ] Task: Initialize Go project `skills` in a new directory `cmd/skills`.
- [ ] Task: Set up `cobra` for command-line argument parsing.
- [ ] Task: Implement basic `version` command.
- [ ] Task: Implement configuration management using `viper` to read from `~/.config/skills/config`.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Project Initialization' (Protocol in workflow.md)

## Phase 2: Environment Setup & Configuration Wizard [checkpoint: ]
- [ ] Task: Implement `skills configure` command (interactive wizard).
- [ ] Task: Port logic from `install.sh` to Go (Git config, SSH key detection).
- [ ] Task: Implement configuration file generation/updating logic.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Environment Setup' (Protocol in workflow.md)

## Phase 3: Multipass Management [checkpoint: ]
- [ ] Task: Create `vm` subcommand group.
- [ ] Task: Implement `skills vm launch/start/stop/delete` (wrapping `multipass` commands).
- [ ] Task: Implement `skills vm ssh/claude/gemini` (SSH execution logic).
- [ ] Task: Implement `skills vm status/mount`.
- [ ] Task: Write tests for Multipass subcommand execution.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Multipass Management' (Protocol in workflow.md)

## Phase 4: Docker Management [checkpoint: ]
- [ ] Task: Create `docker` subcommand group.
- [ ] Task: Implement `skills docker launch/start/stop/delete` (wrapping `docker` commands).
- [ ] Task: Implement `skills docker ssh/claude/gemini/logs`.
- [ ] Task: Implement `skills docker configure-git`.
- [ ] Task: Write tests for Docker subcommand execution.
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Docker Management' (Protocol in workflow.md)

## Phase 5: Skill Synchronization [checkpoint: ]
- [ ] Task: Create `sync` command.
- [ ] Task: Implement skill synchronization logic (replacing shell sync scripts).
- [ ] Task: Add flags for `--dry-run` and `--force`.
- [ ] Task: Write tests for sync logic.
- [ ] Task: Conductor - User Manual Verification 'Phase 5: Skill Synchronization' (Protocol in workflow.md)

## Phase 6: Final Integration & Cleanup [checkpoint: ]
- [ ] Task: Update project documentation to point to the new CLI.
- [ ] Task: (Optional) Add `mise` tasks that wrap the new `skills` CLI for backward compatibility during transition.
- [ ] Task: Final end-to-end verification.
- [ ] Task: Conductor - User Manual Verification 'Phase 6: Final Integration' (Protocol in workflow.md)
