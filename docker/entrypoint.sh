#!/bin/bash
# =============================================================================
# Container Entrypoint Script
# =============================================================================
# Configures Git and SSH settings from environment variables at container start.
# This allows per-session configuration without rebuilding the image.
# =============================================================================

set -e

# =============================================================================
# Git User Configuration
# =============================================================================

if [ -n "${GIT_USER_NAME:-}" ]; then
    git config --global user.name "$GIT_USER_NAME"
fi

if [ -n "${GIT_USER_EMAIL:-}" ]; then
    git config --global user.email "$GIT_USER_EMAIL"
fi

# =============================================================================
# SSH Signing Key Configuration
# =============================================================================

if [ -n "${SSH_SIGNING_KEY_CONTENT:-}" ]; then
    # Write public key to ~/.ssh/id_ed25519.pub
    echo "$SSH_SIGNING_KEY_CONTENT" > ~/.ssh/id_ed25519.pub
    chmod 644 ~/.ssh/id_ed25519.pub

    # Configure Git to use this key for signing
    git config --global user.signingkey ~/.ssh/id_ed25519.pub

    # Set up allowed signers for verification
    if [ -n "${GIT_USER_EMAIL:-}" ]; then
        echo "${GIT_USER_EMAIL} namespaces=\"git\" ${SSH_SIGNING_KEY_CONTENT}" > ~/.ssh/allowed_signers
        git config --global gpg.ssh.allowedSignersFile ~/.ssh/allowed_signers
    fi
fi

# =============================================================================
# Execute Main Command
# =============================================================================

exec "$@"
