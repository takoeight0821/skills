#!/bin/bash
# =============================================================================
# Container Entrypoint Script
# =============================================================================
# Handles:
# - SSH agent socket verification
# - Git user configuration (if provided via environment)
# - SSH signing key configuration
# =============================================================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# =============================================================================
# SSH Agent Verification
# =============================================================================

verify_ssh_agent() {
    if [ -n "$SSH_AUTH_SOCK" ]; then
        if [ -S "$SSH_AUTH_SOCK" ]; then
            log_info "SSH agent socket found at $SSH_AUTH_SOCK"
            # Test if we can communicate with the agent
            if ssh-add -l &>/dev/null || [ $? -eq 1 ]; then
                log_info "SSH agent is accessible"
                ssh-add -l 2>/dev/null || log_warn "No SSH keys loaded in agent"
            else
                log_warn "SSH agent socket exists but agent is not responding"
            fi
        else
            log_warn "SSH_AUTH_SOCK is set but socket does not exist: $SSH_AUTH_SOCK"
            log_warn "SSH agent forwarding may not be configured correctly"
        fi
    else
        log_warn "SSH_AUTH_SOCK is not set - SSH agent forwarding is disabled"
    fi
}

# =============================================================================
# Git Configuration
# =============================================================================

configure_git() {
    # Set git user name if provided
    if [ -n "$GIT_USER_NAME" ]; then
        git config --global user.name "$GIT_USER_NAME"
        log_info "Git user.name set to: $GIT_USER_NAME"
    fi

    # Set git user email if provided
    if [ -n "$GIT_USER_EMAIL" ]; then
        git config --global user.email "$GIT_USER_EMAIL"
        log_info "Git user.email set to: $GIT_USER_EMAIL"
    fi

    # Configure SSH signing key
    # Priority: GIT_SIGNING_KEY env var > ~/.ssh/id_ed25519.pub > ~/.ssh/id_rsa.pub
    if [ -n "$GIT_SIGNING_KEY" ]; then
        git config --global user.signingkey "$GIT_SIGNING_KEY"
        log_info "Git signing key set to: $GIT_SIGNING_KEY"
    elif [ -f "$HOME/.ssh/id_ed25519.pub" ]; then
        git config --global user.signingkey "$HOME/.ssh/id_ed25519.pub"
        log_info "Git signing key auto-detected: ~/.ssh/id_ed25519.pub"
    elif [ -f "$HOME/.ssh/id_rsa.pub" ]; then
        git config --global user.signingkey "$HOME/.ssh/id_rsa.pub"
        log_info "Git signing key auto-detected: ~/.ssh/id_rsa.pub"
    else
        log_warn "No SSH signing key configured"
        log_warn "Set GIT_SIGNING_KEY or mount your public key to ~/.ssh/"
    fi

    # Set up allowed signers file for verification (if email is configured)
    if [ -n "$GIT_USER_EMAIL" ]; then
        SIGNING_KEY=$(git config --global user.signingkey 2>/dev/null || true)
        if [ -n "$SIGNING_KEY" ] && [ -f "$SIGNING_KEY" ]; then
            mkdir -p "$HOME/.ssh"
            ALLOWED_SIGNERS="$HOME/.ssh/allowed_signers"
            # Add current user's key to allowed signers
            key_content=$(cat "$SIGNING_KEY")
            if [ -n "$key_content" ]; then
                printf '%s namespaces="git" %s\n' "$GIT_USER_EMAIL" "$key_content" > "$ALLOWED_SIGNERS"
            else
                log_warn "Signing key file is empty: $SIGNING_KEY"
            fi
            if [ -f "$ALLOWED_SIGNERS" ]; then
                git config --global gpg.ssh.allowedSignersFile "$ALLOWED_SIGNERS"
                log_info "Allowed signers file configured at: $ALLOWED_SIGNERS"
            fi
        fi
    fi
}

# =============================================================================
# Main
# =============================================================================

log_info "Initializing coding agent container..."

verify_ssh_agent
configure_git

log_info "Container initialization complete"
echo ""

# Execute the command passed to the container
exec "$@"
