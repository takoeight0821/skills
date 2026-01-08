#!/bin/bash
# install.sh - Install Docker-based Claude coding agent
#
# This script configures mise tasks for the Docker-based coding agent.
# Run this script from your cloned skills repository.

# Detect the directory where this script is located (= cloned repository)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=docker/lib/common.sh
source "${SCRIPT_DIR}/docker/lib/common.sh"

MISE_CONFIG="$HOME/.config/mise/config.toml"

# Check for mise
if ! command -v mise &> /dev/null; then
    log_error "mise is not installed. Please install mise first."
    echo "See: https://mise.jdx.dev/getting-started.html"
    exit 1
fi

echo "Skills repository: $SCRIPT_DIR"
echo ""

# Update mise config
echo "Configuring mise..."

# Create config directory if needed
mkdir -p "$(dirname "$MISE_CONFIG")"

# Generate the new configuration content
generate_docker_config() {
    cat << EOF
# BEGIN claude-skills
# Docker-based coding agent - https://github.com/takoeight0821/skills
# Repository: $SCRIPT_DIR

[tasks.docker-build]
description = "Build the coding agent Docker image"
run = """
cd "$SCRIPT_DIR/docker" && docker compose build "\$@"
"""

[tasks.docker-up]
description = "Start the coding agent container"
run = """
"$SCRIPT_DIR/docker/run.sh" up "\$@"
"""

[tasks.docker-down]
description = "Stop the coding agent container"
run = """
"$SCRIPT_DIR/docker/run.sh" down "\$@"
"""

[tasks.docker-claude]
description = "Run Claude Code in the container"
run = """
"$SCRIPT_DIR/docker/run.sh" claude "\$@"
"""

[tasks.docker-gemini]
description = "Run Gemini CLI in the container"
run = """
"$SCRIPT_DIR/docker/run.sh" gemini "\$@"
"""

[tasks.docker-shell]
description = "Open a shell in the container"
run = """
"$SCRIPT_DIR/docker/run.sh" shell "\$@"
"""
# END claude-skills
EOF
}

# Check if already configured
MARKER="# BEGIN claude-skills"
END_MARKER="# END claude-skills"

if grep -qF "$MARKER" "$MISE_CONFIG" 2>/dev/null; then
    log_warn "  mise config already contains skills configuration."
    read -rp "  Replace existing configuration? [y/N] " answer
    if [[ "$answer" =~ ^[Yy]$ ]]; then
        # Update existing configuration in place
        TMP_FILE=$(mktemp)
        trap 'rm -f "$TMP_FILE"' EXIT

        inside_section=false

        while IFS= read -r line || [[ -n "$line" ]]; do
            if [[ "$line" == "$MARKER"* ]]; then
                # Start of claude-skills section - write new config
                inside_section=true
                generate_docker_config >> "$TMP_FILE"
            elif [[ "$line" == "$END_MARKER"* ]]; then
                # End of claude-skills section - skip this line (already written by generate_docker_config)
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

        log_success "  mise configuration updated."
    else
        echo "  Skipping mise configuration update."
        log_success ""
        log_success "Installation complete!"
        exit 0
    fi
else
    # No existing configuration - append to file
    # Only add blank line if file is non-empty and doesn't end with empty line
    if [ -s "$MISE_CONFIG" ]; then
        tail -1 "$MISE_CONFIG" | grep -q . && echo "" >> "$MISE_CONFIG"
    fi
    generate_docker_config >> "$MISE_CONFIG"
    log_success "  mise configuration added."
fi

# Done
log_success ""
log_success "Installation complete!"
echo ""
echo "Available commands:"
echo "  mise run docker-build   - Build the Docker image"
echo "  mise run docker-up      - Start the container"
echo "  mise run docker-down    - Stop the container"
echo "  mise run docker-claude  - Run Claude Code"
echo "  mise run docker-gemini  - Run Gemini CLI"
echo "  mise run docker-shell   - Open a shell"
echo ""
echo "Quick start:"
echo "  mise run docker-build && mise run docker-up && mise run docker-claude"
