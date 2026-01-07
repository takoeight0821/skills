# Configuration Reference

Complete reference for Claude Skills sync configuration options.

## Configuration Files

### Global Configuration
- **Location**: `~/.config/skills/config`
- **Scope**: Applies to all sync operations unless overridden
- **Create**: `mkdir -p ~/.config/skills && touch ~/.config/skills/config`

### Project Configuration
- **Location**: `.claude/.skills.conf` (in project root)
- **Scope**: Applies only to this project, overrides global config
- **Create**: `mkdir -p .claude && touch .claude/.skills.conf`

## Configuration Options

### apply

Controls whether sync operations actually modify files.

| Value | Behavior |
|-------|----------|
| `false` (default) | Dry-run mode: preview changes without applying |
| `true` | Apply changes: actually copy/delete files |

```ini
# Preview mode (safe, see what would happen)
apply=false

# Auto-apply mode (actually make changes)
apply=true
```

### force

Controls how file conflicts are handled.

| Value | Behavior |
|-------|----------|
| `false` (default) | Ask for confirmation when local file differs from shared |
| `true` | Overwrite without confirmation |

```ini
# Interactive mode (ask before overwriting)
force=false

# Force mode (always use shared version)
force=true
```

### prune

Controls cleanup of files removed from shared repository.

| Value | Behavior |
|-------|----------|
| `false` (default) | Keep local files even if deleted from shared repo |
| `true` | Delete local files that no longer exist in shared repo |

```ini
# Keep removed files (safer)
prune=false

# Remove deleted files (stay in sync)
prune=true
```

### exclude

Controls git tracking for synced files (project mode only).

| Value | Behavior |
|-------|----------|
| `false` (default) | Synced files are tracked by git normally |
| `true` | Add synced files to `.git/info/exclude` |

```ini
# Track synced files in git
exclude=false

# Exclude synced files from git
exclude=true
```

## Environment Variables

### SKILLS_SHARED_DIR

Location of the shared skills repository.

- **Default**: `~/.local/share/skills`
- **Usage**: Override to use a different location

```bash
export SKILLS_SHARED_DIR="$HOME/my-skills-repo"
```

### SKILLS_REPO

Git URL for the skills repository.

- **Default**: `https://github.com/takoeight0821/skills.git`
- **Usage**: Override to use a fork or private repository

```bash
export SKILLS_REPO="https://github.com/myuser/my-skills.git"
```

## Command Line Flags

Override configuration file settings from the command line:

| Flag | Effect |
|------|--------|
| `--apply` | Set `apply=true` |
| `--force` | Set `force=true` |
| `--prune` | Set `prune=true` |
| `--exclude` | Set `exclude=true` |
| `--global` | Sync to `~/.claude/skills` |
| `--project` | Sync to `.claude/skills` |

Example:
```bash
~/.local/share/skills/bin/sync-skills.sh --project --apply --force
```

## Configuration Precedence

Settings are resolved in this order (later overrides earlier):

1. **Built-in defaults** (lowest priority)
   - apply=false, force=false, prune=false, exclude=false

2. **Global config** (`~/.config/skills/config`)

3. **Project config** (`.claude/.skills.conf`)

4. **Command line flags** (highest priority)

## Manifest Files

The sync script maintains manifest files to track which files were synced:

- **Global**: `~/.claude/.skills-manifest`
- **Project**: `.claude/.skills-manifest`

These files are used to:
- Distinguish synced files from locally-created skills
- Enable proper pruning of deleted files
- Detect conflicts between local changes and shared updates

Do not edit manifest files manually.

## Troubleshooting

### "mise not found"
Install mise: https://mise.jdx.dev/getting-started.html

### "Permission denied" on sync
Check that scripts are executable:
```bash
chmod +x ~/.local/share/skills/bin/sync-skills.sh
```

### Config file not being read
Verify file locations and permissions:
```bash
ls -la ~/.config/skills/config
ls -la .claude/.skills.conf
```

### Unexpected overwriting
Check configuration precedence - project config or CLI flags may be overriding global settings.
