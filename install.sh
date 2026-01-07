#!/bin/bash
# install.sh - Install Claude skills management tools
#
# This script:
# 1. Clones/updates the skills repository to ~/.local/share/skills
# 2. Appends mise tasks to ~/.config/mise/config.toml

set -euo pipefail

SKILLS_REPO="${SKILLS_REPO:-https://github.com/takoeight0821/skills.git}"
SKILLS_DIR="${SKILLS_SHARED_DIR:-$HOME/.local/share/skills}"
MISE_CONFIG="$HOME/.config/mise/config.toml"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

print_error() { echo -e "${RED}Error: $1${NC}" >&2; }
print_success() { echo -e "${GREEN}$1${NC}"; }
print_warning() { echo -e "${YELLOW}$1${NC}"; }
print_info() { echo -e "$1"; }

# Check for mise
if ! command -v mise &> /dev/null; then
    print_error "mise is not installed. Please install mise first."
    echo "See: https://mise.jdx.dev/getting-started.html"
    exit 1
fi

# Step 1: Clone or update skills repository
print_info "Step 1: Setting up skills repository..."

if [ -d "$SKILLS_DIR" ]; then
    print_info "  Updating existing repository at $SKILLS_DIR"
    git -C "$SKILLS_DIR" pull --quiet
else
    print_info "  Cloning repository to $SKILLS_DIR"
    git clone "$SKILLS_REPO" "$SKILLS_DIR"
fi

# Make scripts executable
chmod +x "$SKILLS_DIR/bin/"*.sh 2>/dev/null || true

print_success "  Skills repository ready."

# Step 2: Update mise config
print_info ""
print_info "Step 2: Configuring mise..."

# Create config directory if needed
mkdir -p "$(dirname "$MISE_CONFIG")"

# Check if already configured
MARKER="# BEGIN claude-skills"
if grep -qF "$MARKER" "$MISE_CONFIG" 2>/dev/null; then
    print_warning "  mise config already contains skills configuration."
    read -p "  Replace existing configuration? [y/N] " answer
    if [[ "$answer" =~ ^[Yy]$ ]]; then
        # Remove existing configuration
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' "/$MARKER/,/# END claude-skills/d" "$MISE_CONFIG"
        else
            sed -i "/$MARKER/,/# END claude-skills/d" "$MISE_CONFIG"
        fi
    else
        print_info "  Skipping mise configuration update."
        print_success ""
        print_success "Installation complete!"
        exit 0
    fi
fi

# Append configuration
cat >> "$MISE_CONFIG" << 'EOF'

# BEGIN claude-skills
# Claude skills management - https://github.com/takoeight0821/skills
# All sync tasks run in dry-run mode by default.
# Use *-apply tasks or set apply=true in config file to actually sync.

[tasks.update-shared-skills]
description = "Update shared Claude skills repository"
run = """
SKILLS_DIR="${SKILLS_SHARED_DIR:-$HOME/.local/share/skills}"
if [ -d "$SKILLS_DIR" ]; then
  echo "Updating skills repository..."
  git -C "$SKILLS_DIR" pull --quiet
else
  echo "Skills not installed. Run install.sh first."
  exit 1
fi
echo "Done."
"""

[tasks.sync-skills-global]
description = "Preview sync to ~/.claude/skills (dry-run)"
run = """
"${SKILLS_SHARED_DIR:-$HOME/.local/share/skills}/bin/sync-skills.sh" --global "$@"
"""

[tasks.sync-skills-global-apply]
description = "Sync shared skills to ~/.claude/skills"
run = """
"${SKILLS_SHARED_DIR:-$HOME/.local/share/skills}/bin/sync-skills.sh" --global --apply "$@"
"""

[tasks.sync-skills-project]
description = "Preview sync to .claude/skills (dry-run)"
run = """
"${SKILLS_SHARED_DIR:-$HOME/.local/share/skills}/bin/sync-skills.sh" --project "$@"
"""

[tasks.sync-skills-project-apply]
description = "Sync shared skills to .claude/skills"
run = """
"${SKILLS_SHARED_DIR:-$HOME/.local/share/skills}/bin/sync-skills.sh" --project --apply "$@"
"""
# END claude-skills
EOF

print_success "  mise configuration updated."

# Done
print_success ""
print_success "Installation complete!"
print_info ""
print_info "Available commands:"
print_info "  mise run update-shared-skills      - Update skills repository"
print_info "  mise run sync-skills-global        - Preview sync to ~/.claude/skills (dry-run)"
print_info "  mise run sync-skills-global-apply  - Sync to ~/.claude/skills"
print_info "  mise run sync-skills-project       - Preview sync to .claude/skills (dry-run)"
print_info "  mise run sync-skills-project-apply - Sync to .claude/skills"
print_info ""
print_info "Additional flags (pass via command line):"
print_info "  --force   - Overwrite existing files without confirmation"
print_info "  --prune   - Remove files deleted from shared repository"
print_info "  --exclude - Add synced files to .git/exclude"
print_info ""
print_info "Or configure defaults in ~/.config/skills/config"
