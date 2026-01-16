# Documentation Audit Report

**Date:** 2026-01-16
**Status:** In Progress

## 1. Documentation Inventory

| File Path | Purpose | Status | Action Required |
|-----------|---------|--------|-----------------|
| `README.md` | Main project entry point. Covers installation, usage, and architecture of the skills repo and `jig`. | Active | Update to reflect new `conductor` structure and any recent changes. |
| `CLAUDE.md` | Context file for Claude Code. Contains project-specific instructions and skills. | Active | Update to reference `conductor/index.md` as the source of truth for project context. |
| `GEMINI.md` | Context file for Gemini CLI. Contains project-specific instructions. | Active | Update to reference `conductor/index.md` as the source of truth for project context. |
| `conductor/index.md` | New project context index. | New | Populate with links to consolidated documentation. |
| `conductor/product.md` | Product definition and core features. | New | Maintain as source of truth for product scope. |
| `conductor/tech-stack.md` | Technology stack definition. | New | Maintain as source of truth for tech stack. |
| `conductor/workflow.md` | Development workflow and processes. | New | Maintain as source of truth for workflow. |
| `docs/research/docker-ssh-agent-forwarding-macos.md` | Detailed guide on Docker + SSH Agent Forwarding on macOS. | Reference | Keep as a deep-dive reference. Link from `conductor/index.md` or a new "Guides" section. |
| `docs/research/macos-multipass-ssh-agent-forwarding.md` | Detailed guide on Multipass + SSH Agent Forwarding on macOS. | Reference | Keep as a deep-dive reference. Link from `conductor/index.md` or a new "Guides" section. |
| `skills/research/SKILL.md` | Likely a template or placeholder. | TBD | Verify if this is needed. If it's a template, move to `templates/` or document its usage. |
| `skills/review-plan/SKILL.md` | Review plan skill definition. | Active | Part of the skills library. |
| `skills/review-plan/references/categories.md` | Reference for the review plan skill. | Active | Part of the skills library. |

## 2. Redundancy & Outdated Information Analysis

- **`README.md` vs `conductor/product.md`**: `README.md` contains a good overview but `conductor/product.md` now holds the formal definition. `README.md` should focus on usage and installation, while linking to `conductor/` for architectural details if necessary, or `product.md` should summarize `README.md`. Currently `product.md` was seeded from `README.md`, so they are consistent.
- **Research Docs**: The files in `docs/research/` are valuable technical guides. They are not redundant but are currently "hidden" in a research folder. They should be elevated to a "Guides" or "Documentation" section in the index.
- **`CLAUDE.md` & `GEMINI.md`**: These contain redundant context that is now better structured in `conductor/`. They should be streamlined to point to the `conductor` files to avoid drift.

## 3. Recommendations

1.  **Centralize Index**: Use `conductor/index.md` as the main map.
2.  **Streamline Context Files**: Reduce `CLAUDE.md` and `GEMINI.md` to essential agent instructions and point them to `conductor/index.md` for project knowledge.
3.  **Promote Research**: Move `docs/research` to `docs/guides` or simply `docs/` to indicate they are permanent documentation, not just temporary research. Or keep `research` if they are static snapshots. Given their high quality, `docs/guides` seems appropriate.
4.  **Clarify Skills Structure**: Ensure the `skills/` directory structure is documented in `README.md` or a dedicated `docs/skills.md`.
