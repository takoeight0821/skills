# Specification: Project Foundation

## Goal
To organize the existing project structure, consolidate documentation, and verify the functionality of current `mise` tasks to establish a solid foundation for future development.

## Scope
- **Documentation:**
    - Audit existing `README.md`, `CLAUDE.md`, and `GEMINI.md`.
    - Move research notes from `docs/research/` and `skills/research/` to a centralized location if needed, or index them properly.
    - Ensure `install.sh` usage and `mise` task descriptions are clear and up-to-date in the main documentation.
- **Verification:**
    - Verify all `mise` tasks listed in `mise.toml` and `mise.sample.toml` function as expected (dry-run where applicable).
    - Ensure `mpvm` build process works.
    - Validate Docker build commands.

## Out of Scope
- Adding new features to `mpvm` or the sync logic.
- Major refactoring of the Go codebase (unless critical bugs are found during verification).

## Success Criteria
- A unified documentation structure referenced in `conductor/index.md`.
- All `mise` tasks are documented with usage examples.
- Confirmation that the build and sync pipelines are operational.
