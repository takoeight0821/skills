#!/bin/bash
# install.sh - Install Multipass-based Claude coding agent
#
# This script configures mise tasks for the Multipass-based coding agent.
# Run this script from your cloned skills repository.

# Detect the directory where this script is located (= cloned repository)
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=multipass/lib/common.sh
source "${SCRIPT_DIR}/multipass/lib/common.sh"

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
generate_multipass_config() {
    cat << EOF
# BEGIN claude-skills
# Multipass-based coding agent - https://github.com/takoeight0821/skills
# Repository: $SCRIPT_DIR

[tasks."tk:vm-launch"]
description = "Create and start the coding agent VM"
run = """
"$SCRIPT_DIR/multipass/run.sh" launch "\$@"
"""

[tasks."tk:vm-start"]
description = "Start the coding agent VM"
run = """
"$SCRIPT_DIR/multipass/run.sh" start "\$@"
"""

[tasks."tk:vm-stop"]
description = "Stop the coding agent VM"
run = """
"$SCRIPT_DIR/multipass/run.sh" stop "\$@"
"""

[tasks."tk:vm-delete"]
description = "Delete the coding agent VM"
run = """
"$SCRIPT_DIR/multipass/run.sh" delete "\$@"
"""

[tasks."tk:vm-ssh"]
description = "SSH into the VM with agent forwarding"
run = """
"$SCRIPT_DIR/multipass/run.sh" ssh "\$@"
"""

[tasks."tk:vm-claude"]
description = "Run Claude Code in the VM"
run = """
"$SCRIPT_DIR/multipass/run.sh" claude "\$@"
"""

[tasks."tk:vm-gemini"]
description = "Run Gemini CLI in the VM"
run = """
"$SCRIPT_DIR/multipass/run.sh" gemini "\$@"
"""

[tasks."tk:vm-status"]
description = "Show VM status"
run = """
"$SCRIPT_DIR/multipass/run.sh" status "\$@"
"""

[tasks."tk:vm-mount"]
description = "Mount directory to VM"
run = """
"$SCRIPT_DIR/multipass/run.sh" mount "\$@"
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
                generate_multipass_config >> "$TMP_FILE"
            elif [[ "$line" == "$END_MARKER"* ]]; then
                # End of claude-skills section - skip this line (already written by generate_multipass_config)
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
    generate_multipass_config >> "$MISE_CONFIG"
    log_success "  mise configuration added."
fi

# =============================================================================
# Configuration Wizard
# =============================================================================

SKILLS_CONFIG="$HOME/.config/skills/config"

setup_config() {
    echo ""
    echo "Setting up global configuration..."
    echo ""

    # Get defaults from git config if available
    local default_name default_email
    default_name=$(git config --global user.name 2>/dev/null || echo "")
    default_email=$(git config --global user.email 2>/dev/null || echo "")

    # Prompt for Git user name
    if [ -n "$default_name" ]; then
        read -rp "Git user name [$default_name]: " git_name
        git_name="${git_name:-$default_name}"
    else
        read -rp "Git user name: " git_name
        while [ -z "$git_name" ]; do
            log_warn "Git user name is required for signed commits."
            read -rp "Git user name: " git_name
        done
    fi

    # Prompt for Git user email
    if [ -n "$default_email" ]; then
        read -rp "Git user email [$default_email]: " git_email
        git_email="${git_email:-$default_email}"
    else
        read -rp "Git user email: " git_email
        while [ -z "$git_email" ]; do
            log_warn "Git user email is required for signed commits."
            read -rp "Git user email: " git_email
        done
    fi

    # Prompt for SSH signing key
    local default_key="$HOME/.ssh/id_ed25519.pub"
    if [ -f "$default_key" ]; then
        read -rp "SSH signing key [$default_key]: " ssh_key
        ssh_key="${ssh_key:-$default_key}"
    else
        # Try other common key types
        for key in "$HOME/.ssh/id_rsa.pub" "$HOME/.ssh/id_ecdsa.pub"; do
            if [ -f "$key" ]; then
                default_key="$key"
                break
            fi
        done
        if [ -f "$default_key" ]; then
            read -rp "SSH signing key [$default_key]: " ssh_key
            ssh_key="${ssh_key:-$default_key}"
        else
            log_warn "No SSH public key found. Git signing will not work."
            read -rp "SSH signing key path (or leave empty): " ssh_key
        fi
    fi

    # Create config directory
    mkdir -p "$(dirname "$SKILLS_CONFIG")"

    # Write config file
    cat > "$SKILLS_CONFIG" << EOF
# Skills VM Configuration
# Generated by install.sh on $(date +%Y-%m-%d)

# Git configuration (required for signed commits)
GIT_USER_NAME="$git_name"
GIT_USER_EMAIL="$git_email"
EOF

    if [ -n "$ssh_key" ]; then
        cat >> "$SKILLS_CONFIG" << EOF

# SSH signing key
SSH_SIGNING_KEY="$ssh_key"
EOF
    fi

    cat >> "$SKILLS_CONFIG" << EOF

# VM resources (uncomment to customize)
# MULTIPASS_VM_NAME=coding-agent
# MULTIPASS_VM_CPUS=2
# MULTIPASS_VM_MEMORY=4G
# MULTIPASS_VM_DISK=20G
EOF

    log_success "  Configuration saved to $SKILLS_CONFIG"
}

# Check if config exists
if [ -f "$SKILLS_CONFIG" ]; then
    echo ""
    log_info "Global configuration already exists at $SKILLS_CONFIG"
    read -rp "Reconfigure? [y/N] " answer
    if [[ "$answer" =~ ^[Yy]$ ]]; then
        setup_config
    fi
else
    echo ""
    read -rp "Set up global configuration for VM? [Y/n] " answer
    if [[ ! "$answer" =~ ^[Nn]$ ]]; then
        setup_config
    else
        log_warn "Skipping configuration. You can set it up later by creating:"
        echo "  $SKILLS_CONFIG"
        echo ""
        echo "Or copy the example:"
        echo "  mkdir -p ~/.config/skills"
        echo "  cp $SCRIPT_DIR/multipass/config.example ~/.config/skills/config"
    fi
fi

# =============================================================================
# Done
# =============================================================================

log_success ""
log_success "Installation complete!"
echo ""
echo "Available commands:"
echo "  mise run tk:vm-launch   - Create and start VM"
echo "  mise run tk:vm-start    - Start VM"
echo "  mise run tk:vm-stop     - Stop VM"
echo "  mise run tk:vm-delete   - Delete VM"
echo "  mise run tk:vm-ssh      - SSH with agent forwarding"
echo "  mise run tk:vm-claude   - Run Claude Code"
echo "  mise run tk:vm-gemini   - Run Gemini CLI"
echo "  mise run tk:vm-status   - Show VM status"
echo "  mise run tk:vm-mount    - Mount directory"
echo ""
echo "Quick start:"
echo "  mise run tk:vm-launch && mise run tk:vm-claude"
