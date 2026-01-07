#!/bin/bash
# sync-skills.sh - Sync shared Claude skills to local directory
#
# Usage:
#   sync-skills.sh --global [--apply] [--force] [--prune] [--exclude]
#   sync-skills.sh --project [--apply] [--force] [--prune] [--exclude]
#
# Options:
#   --global    Sync to ~/.claude/skills
#   --project   Sync to .claude/skills
#   --apply     Actually perform the sync (default is dry-run)
#   --force     Overwrite existing files without confirmation
#   --prune     Remove files that no longer exist in shared repository
#   --exclude   Add synced files to .git/exclude (project mode only)
#
# Config files (loaded in order, later overrides earlier):
#   ~/.config/skills/config    Global config
#   .claude/.skills.conf       Project config
#
# Config file format:
#   apply=true
#   force=true
#   prune=false
#   exclude=true

set -euo pipefail

SHARED_DIR="${SKILLS_SHARED_DIR:-$HOME/.local/share/skills/skills}"
GLOBAL_CONFIG="$HOME/.config/skills/config"
PROJECT_CONFIG=".claude/.skills.conf"

# Default values (dry-run by default)
DRY_RUN=true
FORCE=false
PRUNE=false
EXCLUDE=false
TARGET_DIR=""
MANIFEST_FILE=""
MODE=""

# Track if flags were explicitly set on command line
FLAG_APPLY_SET=false
FLAG_FORCE_SET=false
FLAG_PRUNE_SET=false
FLAG_EXCLUDE_SET=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_error() { echo -e "${RED}Error: $1${NC}" >&2; }
print_success() { echo -e "${GREEN}$1${NC}"; }
print_warning() { echo -e "${YELLOW}Warning: $1${NC}"; }
print_dry_run() { echo -e "${BLUE}[dry-run] $1${NC}"; }

# Load config file if exists
load_config() {
    local config_file="$1"
    if [ -f "$config_file" ]; then
        while IFS='=' read -r key value || [ -n "$key" ]; do
            # Skip comments and empty lines
            [[ "$key" =~ ^[[:space:]]*# ]] && continue
            [[ -z "$key" ]] && continue

            # Trim whitespace
            key=$(echo "$key" | xargs)
            value=$(echo "$value" | xargs)

            case "$key" in
                apply)
                    if [ "$FLAG_APPLY_SET" = false ]; then
                        [ "$value" = "true" ] && DRY_RUN=false
                    fi
                    ;;
                force)
                    [ "$FLAG_FORCE_SET" = false ] && FORCE="$value"
                    ;;
                prune)
                    [ "$FLAG_PRUNE_SET" = false ] && PRUNE="$value"
                    ;;
                exclude)
                    [ "$FLAG_EXCLUDE_SET" = false ] && EXCLUDE="$value"
                    ;;
            esac
        done < "$config_file"
    fi
}

usage() {
    cat <<EOF
Usage: $(basename "$0") --global|--project [--apply] [--force] [--prune] [--exclude]

Sync shared Claude skills to local directory.

By default, runs in dry-run mode (no changes made). Use --apply to actually sync.

Options:
  --global    Sync to ~/.claude/skills
  --project   Sync to .claude/skills
  --apply     Actually perform the sync (default is dry-run)
  --force     Overwrite existing files without confirmation
  --prune     Remove files that no longer exist in shared repository
  --exclude   Add synced files to .git/exclude (project mode only)
  --help      Show this help message

Config files (loaded in order, flags override config):
  ~/.config/skills/config    Global config
  .claude/.skills.conf       Project config

Config file format:
  apply=true
  force=true
  prune=false
  exclude=true

Environment:
  SKILLS_SHARED_DIR   Override shared skills directory
                      (default: ~/.local/share/skills/skills)
EOF
    exit 0
}

# Parse arguments (first pass to get flags)
ARGS=("$@")
for arg in "${ARGS[@]}"; do
    case $arg in
        --apply) FLAG_APPLY_SET=true; DRY_RUN=false ;;
        --force) FLAG_FORCE_SET=true; FORCE=true ;;
        --prune) FLAG_PRUNE_SET=true; PRUNE=true ;;
        --exclude) FLAG_EXCLUDE_SET=true; EXCLUDE=true ;;
    esac
done

# Load config files (global first, then project)
load_config "$GLOBAL_CONFIG"
load_config "$PROJECT_CONFIG"

# Parse arguments (second pass for mode and validation)
while [[ $# -gt 0 ]]; do
    case $1 in
        --apply) shift ;;
        --force) shift ;;
        --prune) shift ;;
        --exclude) shift ;;
        --global)
            MODE="global"
            TARGET_DIR="$HOME/.claude/skills"
            MANIFEST_FILE="$HOME/.claude/.skills-manifest"
            shift ;;
        --project)
            MODE="project"
            TARGET_DIR=".claude/skills"
            MANIFEST_FILE=".claude/.skills-manifest"
            shift ;;
        --help) usage ;;
        *)
            print_error "Unknown option: $1"
            usage
            ;;
    esac
done

# Validate arguments
if [ -z "$MODE" ]; then
    print_error "Either --global or --project is required"
    usage
fi

if [ ! -d "$SHARED_DIR" ]; then
    print_error "Shared skills directory not found: $SHARED_DIR"
    echo "Run 'mise run update-shared-skills' first."
    exit 1
fi

# Check if exclude option in non-git repo
if [ "$EXCLUDE" = true ] && [ ! -d .git ]; then
    print_warning "Not a git repository. --exclude option will be ignored."
    EXCLUDE=false
fi

# Show mode
if [ "$DRY_RUN" = true ]; then
    echo -e "${BLUE}=== DRY-RUN MODE (use --apply to actually sync) ===${NC}"
    echo ""
fi

# Create target directory (only if not dry-run)
if [ "$DRY_RUN" = false ]; then
    mkdir -p "$TARGET_DIR"
    mkdir -p "$(dirname "$MANIFEST_FILE")"
    touch "$MANIFEST_FILE"
fi

# Track changes
ADDED=0
UPDATED=0
SKIPPED=0
REMOVED=0

# Sync files
shopt -s nullglob
for f in "$SHARED_DIR"/*.md; do
    name=$(basename "$f")
    target="$TARGET_DIR/$name"

    if [ -f "$target" ]; then
        # File exists - check if different
        if ! diff -q "$f" "$target" > /dev/null 2>&1; then
            if [ "$DRY_RUN" = true ]; then
                print_dry_run "Would update: $name"
                ((UPDATED++))
            elif [ "$FORCE" = true ]; then
                cp "$f" "$target"
                print_success "Updated: $name"
                ((UPDATED++))
            else
                print_warning "$name differs from shared version"
                echo "  Local:  $(head -c 50 "$target" | tr '\n' ' ')..."
                echo "  Shared: $(head -c 50 "$f" | tr '\n' ' ')..."
                read -p "Overwrite? [y/N] " answer
                if [[ "$answer" =~ ^[Yy]$ ]]; then
                    cp "$f" "$target"
                    print_success "Updated: $name"
                    ((UPDATED++))
                else
                    echo "Skipped: $name"
                    ((SKIPPED++))
                fi
            fi
        fi
        # Ensure it's in manifest (only if not dry-run)
        if [ "$DRY_RUN" = false ]; then
            grep -qxF "$name" "$MANIFEST_FILE" || echo "$name" >> "$MANIFEST_FILE"
        fi
    else
        # New file
        if [ "$DRY_RUN" = true ]; then
            print_dry_run "Would add: $name"
            if [ "$EXCLUDE" = true ]; then
                print_dry_run "Would add to .git/exclude: .claude/skills/$name"
            fi
        else
            cp "$f" "$target"
            print_success "Added: $name"

            # Add to manifest
            grep -qxF "$name" "$MANIFEST_FILE" || echo "$name" >> "$MANIFEST_FILE"

            # Add to .git/exclude if --exclude option is set
            if [ "$EXCLUDE" = true ]; then
                exclude_entry=".claude/skills/$name"
                mkdir -p .git/info
                touch .git/info/exclude
                if ! grep -qxF "$exclude_entry" .git/info/exclude 2>/dev/null; then
                    echo "$exclude_entry" >> .git/info/exclude
                fi
            fi
        fi
        ((ADDED++))
    fi
done
shopt -u nullglob

# Prune removed files
if [ "$PRUNE" = true ]; then
    if [ "$DRY_RUN" = true ]; then
        # Dry-run: just show what would be removed
        if [ -f "$MANIFEST_FILE" ]; then
            while IFS= read -r name || [ -n "$name" ]; do
                [ -z "$name" ] && continue
                if [ ! -f "$SHARED_DIR/$name" ] && [ -f "$TARGET_DIR/$name" ]; then
                    print_dry_run "Would remove: $name"
                    if [ "$EXCLUDE" = true ]; then
                        print_dry_run "Would remove from .git/exclude: .claude/skills/$name"
                    fi
                    ((REMOVED++))
                fi
            done < "$MANIFEST_FILE"
        fi
    else
        # Actually prune
        NEW_MANIFEST=$(mktemp)

        while IFS= read -r name || [ -n "$name" ]; do
            [ -z "$name" ] && continue

            if [ ! -f "$SHARED_DIR/$name" ]; then
                # File no longer in shared repo
                if [ -f "$TARGET_DIR/$name" ]; then
                    rm -f "$TARGET_DIR/$name"
                    print_warning "Removed: $name"
                    ((REMOVED++))

                    # Remove from .git/exclude if --exclude option is set
                    if [ "$EXCLUDE" = true ] && [ -f .git/info/exclude ]; then
                        # macOS and Linux compatible sed
                        if [[ "$OSTYPE" == "darwin"* ]]; then
                            sed -i '' "/^\.claude\/skills\/${name//\//\\/}$/d" .git/info/exclude 2>/dev/null || true
                        else
                            sed -i "/^\.claude\/skills\/${name//\//\\/}$/d" .git/info/exclude 2>/dev/null || true
                        fi
                    fi
                fi
            else
                # Keep in manifest
                echo "$name" >> "$NEW_MANIFEST"
            fi
        done < "$MANIFEST_FILE"

        mv "$NEW_MANIFEST" "$MANIFEST_FILE"
    fi
fi

# Summary
echo ""
echo "Summary:"
if [ "$DRY_RUN" = true ]; then
    echo "  Would add:    $ADDED"
    echo "  Would update: $UPDATED"
    [ "$PRUNE" = true ] && echo "  Would remove: $REMOVED"
    echo ""
    echo -e "${BLUE}Run with --apply to actually perform these changes.${NC}"
else
    echo "  Added:   $ADDED"
    echo "  Updated: $UPDATED"
    echo "  Skipped: $SKIPPED"
    [ "$PRUNE" = true ] && echo "  Removed: $REMOVED"
fi
