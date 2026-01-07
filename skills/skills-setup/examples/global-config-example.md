# Global Configuration Example

Example global configuration file for `~/.config/skills/config`.

## Recommended: Safe Defaults

For users who want to preview changes before applying:

```ini
# ~/.config/skills/config
# Global defaults for skills sync - safe mode

# Preview changes before applying
apply=false

# Ask before overwriting modified files
force=false

# Keep files even if removed from shared repo
prune=false

# Not applicable for global mode
exclude=false
```

## Alternative: Auto-sync Mode

For users who want automatic syncing without prompts:

```ini
# ~/.config/skills/config
# Global defaults for skills sync - auto mode

# Automatically apply changes
apply=true

# Overwrite without asking (shared repo is canonical)
force=true

# Keep removed files (safer)
prune=false

# Not applicable for global mode
exclude=false
```

## Alternative: Full Auto Mode

For advanced users who want complete automation:

```ini
# ~/.config/skills/config
# Global defaults for skills sync - full auto

# Automatically apply changes
apply=true

# Overwrite without asking
force=true

# Remove deleted skills automatically
prune=true

# Not applicable for global mode
exclude=false
```

## Setup Instructions

1. Create the config directory:
   ```bash
   mkdir -p ~/.config/skills
   ```

2. Create the config file:
   ```bash
   cat > ~/.config/skills/config << 'EOF'
   apply=false
   force=false
   prune=false
   exclude=false
   EOF
   ```

3. Verify:
   ```bash
   cat ~/.config/skills/config
   ```
