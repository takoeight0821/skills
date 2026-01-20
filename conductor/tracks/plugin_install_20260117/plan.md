# Implementation Plan: Remote Plugin/Extension Installation & Command Conversion

## Phase 1: Infrastructure & Discovery
- [x] Task: Define Remote Source Handler interface and GitHub implementation 2073a4f
    - [ ] Create interface for fetching remote manifests (`plugin.json`, `marketplace.json`, `gemini.json`)
    - [ ] Implement GitHub-specific fetcher using `git` CLI (to leverage local credentials)
- [~] Task: Implement Auto-Detection logic for Claude Plugins and Gemini Extensions
    - [ ] Add logic to distinguish between Claude and Gemini based on manifest files
    - [ ] Create unit tests for various repo structures (Claude-only, Gemini-only, Mixed)
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Infrastructure & Discovery' (Protocol in workflow.md)

## Phase 2: Content Processing & Conversion
- [ ] Task: Implement "Command" detection and metadata extraction
    - [ ] Logic to identify standalone markdown files as "commands"
    - [ ] Logic to extract name and description from these files
- [ ] Task: Implement `SKILL.md` Generator for Commands
    - [ ] Create a template for auto-generated `SKILL.md`
    - [ ] Implement the conversion logic (TDD: Write tests for various markdown inputs)
- [ ] Task: Implement Namespaced Directory logic
    - [ ] Logic to generate `<plugin-name>-<skill-name>` paths
    - [ ] Logic to ensure entire skill subdirectories are included in the sync list
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Content Processing & Conversion' (Protocol in workflow.md)

## Phase 3: CLI Integration & Deployment
- [ ] Task: Implement `skills install <url>` command
    - [ ] Integrate Phase 1 and 2 logic into the Cobra command
    - [ ] Support `--dry-run`, `--force`, and `--target` (global/project) flags
- [ ] Task: Implement Claude Code Deployment
    - [ ] Logic to write processed skills to `.claude/skills/`
- [ ] Task: Implement Gemini CLI Deployment
    - [ ] Verify Gemini's extension path and format
    - [ ] Implement logic to deploy skills/commands to Gemini's target directory
- [ ] Task: Conductor - User Manual Verification 'Phase 3: CLI Integration & Deployment' (Protocol in workflow.md)

## Phase 4: Verification & Refinement
- [ ] Task: Integration testing with real-world repositories
    - [ ] Test installation from a public GitHub repo
    - [ ] Verify that installed skills/commands are recognized by Claude Code and Gemini CLI
- [ ] Task: Final code cleanup and documentation
    - [ ] Ensure all public Go functions are documented
    - [ ] Verify code coverage meets >80% requirement
- [ ] Task: Conductor - User Manual Verification 'Phase 4: Verification & Refinement' (Protocol in workflow.md)
