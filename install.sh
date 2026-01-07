#!/bin/bash
# install.sh - Install Claude skills management tools
#
# This script configures mise tasks for Claude skills management.
# Run this script from your cloned skills repository.

set -euo pipefail

# Detect the directory where this script is located (= cloned repository)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SKILLS_DIR="$SCRIPT_DIR/.claude/skills"
MISE_CONFIG="$HOME/.config/mise/config.toml"
SKILLS_CONFIG="$HOME/.config/skills/config"
SKILLS_CONFIG_SAMPLE="$SCRIPT_DIR/config.sample"

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

# Verify skills directory exists
if [ ! -d "$SKILLS_DIR" ]; then
    print_error "Skills directory not found: $SKILLS_DIR"
    print_info "Make sure you're running this script from the skills repository root."
    exit 1
fi

print_info "Skills repository: $SCRIPT_DIR"
print_info "Skills directory:  $SKILLS_DIR"
print_info ""

# Make scripts executable
chmod +x "$SCRIPT_DIR/bin/"*.sh 2>/dev/null || true

# Copy config sample
print_info "Setting up skills config..."

mkdir -p "$(dirname "$SKILLS_CONFIG")"

if [ -f "$SKILLS_CONFIG" ]; then
    print_warning "  Config file already exists: $SKILLS_CONFIG"
    read -p "  Replace existing config? [y/N] " answer
    if [[ "$answer" =~ ^[Yy]$ ]]; then
        cp "$SKILLS_CONFIG_SAMPLE" "$SKILLS_CONFIG"
        print_success "  Config file replaced."
    else
        print_info "  Skipping config file."
    fi
else
    cp "$SKILLS_CONFIG_SAMPLE" "$SKILLS_CONFIG"
    print_success "  Config file created: $SKILLS_CONFIG"
fi

print_info ""

# Update mise config
print_info "Configuring mise..."

# Create config directory if needed
mkdir -p "$(dirname "$MISE_CONFIG")"

# Generate the new configuration content
# Note: We embed the actual path since it's determined at install time
generate_skills_config() {
    cat << EOF
# BEGIN claude-skills
# Claude skills management - https://github.com/takoeight0821/skills
# Skills repository: $SCRIPT_DIR
# All sync tasks run in dry-run mode by default.
# Use *-apply tasks or set apply=true in config file to actually sync.

[tasks.update-shared-skills]
description = "Update shared Claude skills repository"
run = """
SKILLS_DIR="$SCRIPT_DIR"
if [ -d "\$SKILLS_DIR" ]; then
  echo "Updating skills repository..."
  git -C "\$SKILLS_DIR" pull --quiet
  echo "Done."
else
  echo "Skills repository not found: \$SKILLS_DIR"
  exit 1
fi
"""

[tasks.sync-skills-global]
description = "Preview sync to ~/.claude/skills (dry-run)"
run = """
"$SCRIPT_DIR/bin/sync-skills.sh" --global "\$@"
"""

[tasks.sync-skills-global-apply]
description = "Sync shared skills to ~/.claude/skills"
run = """
"$SCRIPT_DIR/bin/sync-skills.sh" --global --apply "\$@"
"""

[tasks.sync-skills-project]
description = "Preview sync to .claude/skills (dry-run)"
run = """
"$SCRIPT_DIR/bin/sync-skills.sh" --project "\$@"
"""

[tasks.sync-skills-project-apply]
description = "Sync shared skills to .claude/skills"
run = """
"$SCRIPT_DIR/bin/sync-skills.sh" --project --apply "\$@"
"""
# END claude-skills
EOF
}

# Check if already configured
MARKER="# BEGIN claude-skills"
END_MARKER="# END claude-skills"

if grep -qF "$MARKER" "$MISE_CONFIG" 2>/dev/null; then
    print_warning "  mise config already contains skills configuration."
    read -p "  Replace existing configuration? [y/N] " answer
    if [[ "$answer" =~ ^[Yy]$ ]]; then
        # Update existing configuration in place
        TMP_FILE=$(mktemp)
        trap "rm -f '$TMP_FILE'" EXIT

        inside_section=false
        section_replaced=false

        while IFS= read -r line || [[ -n "$line" ]]; do
            if [[ "$line" == "$MARKER"* ]]; then
                # Start of claude-skills section - write new config
                inside_section=true
                section_replaced=true
                generate_skills_config >> "$TMP_FILE"
            elif [[ "$line" == "$END_MARKER"* ]]; then
                # End of claude-skills section - skip this line (already written by generate_skills_config)
                inside_section=false
            elif [[ "$inside_section" == false ]]; then
                # Outside section - copy line as-is
                printf '%s\n' "$line" >> "$TMP_FILE"
            fi
            # Inside section - skip lines (will be replaced)
        done < "$MISE_CONFIG"

        # Replace original file
        mv "$TMP_FILE" "$MISE_CONFIG"
        trap - EXIT

        print_success "  mise configuration updated (in place)."
    else
        print_info "  Skipping mise configuration update."
        print_success ""
        print_success "Installation complete!"
        exit 0
    fi
else
    # No existing configuration - append to file
    # Only add blank line if file is non-empty and doesn't end with empty line
    if [ -s "$MISE_CONFIG" ]; then
        tail -1 "$MISE_CONFIG" | grep -q . && echo "" >> "$MISE_CONFIG"
    fi
    generate_skills_config >> "$MISE_CONFIG"
    print_success "  mise configuration added."
fi

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
