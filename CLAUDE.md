# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

> **Project Context:** Please refer to [conductor/index.md](conductor/index.md) for the most up-to-date documentation, workflow, and architecture details.

This repository is a Claude Code Plugin that provides skills for research and development. Skills include repository analysis, code review support, and comment optimization.

## Architecture

### Directory Structure

- `skills/` - Claude Code skills (SKILL.md files)
  - `clean-comments/` - Code comment optimization skill
  - `research/` - Repository and paper research skill
  - `review-plan/` - Code review response planning skill
- `.claude-plugin/` - Claude Code plugin configuration
  - `plugin.json` - Plugin metadata
- `conductor/` - Project management documentation

## Adding New Skills

Create a directory under `skills/` with a `SKILL.md` file:

```bash
mkdir -p skills/my-skill
# Create skills/my-skill/SKILL.md with YAML frontmatter (name, description) and skill content
```

Skill files use YAML frontmatter for metadata (name, description) followed by markdown content describing the skill workflow.

## Hints

- Use `Write` or `WriteFile` tools to create or modify files. Do not edit files by `echo` or `cat`.
