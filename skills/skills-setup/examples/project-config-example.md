# Project Configuration Example

Example project configuration file for `.claude/.skills.conf`.

## Recommended: Development Project

For active development where you want shared skills but also local modifications:

```ini
# .claude/.skills.conf
# Project-specific skills sync settings

# Preview changes first
apply=false

# Ask before overwriting local changes
force=false

# Keep removed files (may have local modifications)
prune=false

# Track synced skills in git
exclude=false
```

## Alternative: Shared Skills Only

For projects where shared skills are canonical and shouldn't be modified locally:

```ini
# .claude/.skills.conf
# Project-specific skills sync settings

# Automatically apply updates
apply=true

# Always use shared version
force=true

# Remove skills deleted from shared repo
prune=true

# Exclude from git (shared repo is source of truth)
exclude=true
```

## Alternative: Mixed Mode

For projects with both shared and local skills:

```ini
# .claude/.skills.conf
# Project-specific skills sync settings

# Preview changes first
apply=false

# Ask before overwriting (may have local modifications)
force=false

# Keep removed files (may still be useful)
prune=false

# Exclude synced files from git (avoids duplication)
exclude=true
```

## Setup Instructions

1. Create the .claude directory:
   ```bash
   mkdir -p .claude
   ```

2. Create the config file:
   ```bash
   cat > .claude/.skills.conf << 'EOF'
   apply=false
   force=false
   prune=false
   exclude=false
   EOF
   ```

3. Verify:
   ```bash
   cat .claude/.skills.conf
   ```

## Git Considerations

If using `exclude=true`:
- Synced skills are added to `.git/info/exclude`
- They won't appear in `git status`
- They won't be committed to the project repo
- The shared skills repo remains the source of truth

If using `exclude=false`:
- Synced skills are tracked by git
- Changes appear in `git status`
- You can commit modified skills to the project repo
- Project repo has its own copy of skills
