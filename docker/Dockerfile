# =============================================================================
# Coding Agent Base Image
# =============================================================================
# A minimal Docker image for AI coding agents (Claude Code, Gemini CLI)
# with SSH-based git commit signing and SSH agent forwarding support.
#
# Supported platforms:
# - Linux (native Docker)
# - macOS + Docker Desktop
# - macOS + Rancher Desktop
# =============================================================================

FROM ubuntu:24.04

# Prevent interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# =============================================================================
# Install System Dependencies
# =============================================================================

RUN apt-get update && apt-get install -y --no-install-recommends \
    # Core utilities
    ca-certificates \
    curl \
    gnupg \
    # Required tools
    git \
    jq \
    # SSH client for agent forwarding
    openssh-client \
    # sudo for optional elevated privileges
    sudo \
    && rm -rf /var/lib/apt/lists/*

# Verify git version >= 2.34 (required for SSH signing)
RUN git_version=$(git --version | grep -oE '[0-9]+\.[0-9]+' | head -1) && \
    echo "Git version: $git_version" && \
    major=$(echo "$git_version" | cut -d. -f1) && \
    minor=$(echo "$git_version" | cut -d. -f2) && \
    if [ "$major" -lt 2 ] || { [ "$major" -eq 2 ] && [ "$minor" -lt 34 ]; }; then \
        echo "ERROR: Git version must be >= 2.34 for SSH signing support" && exit 1; \
    fi

# =============================================================================
# Install Node.js via NodeSource
# =============================================================================

RUN mkdir -p /etc/apt/keyrings \
    && curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key \
       | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_22.x nodistro main" \
       > /etc/apt/sources.list.d/nodesource.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends nodejs \
    && rm -rf /var/lib/apt/lists/*

# Verify Node.js installation
RUN node --version && npm --version

# =============================================================================
# Create Non-Root User
# =============================================================================

ARG USERNAME=agent
ARG USER_UID=1000
ARG USER_GID=${USER_UID}

RUN groupadd --gid ${USER_GID} ${USERNAME} 2>/dev/null || groupmod -n ${USERNAME} $(getent group ${USER_GID} | cut -d: -f1) \
    && useradd --uid ${USER_UID} --gid ${USER_GID} -m ${USERNAME} 2>/dev/null || usermod -l ${USERNAME} -d /home/${USERNAME} -m $(getent passwd ${USER_UID} | cut -d: -f1) \
    && echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers.d/${USERNAME} \
    && chmod 0440 /etc/sudoers.d/${USERNAME}

# Create directories for user
RUN mkdir -p /home/${USERNAME}/.ssh \
    && mkdir -p /home/${USERNAME}/.config \
    && chown -R ${USERNAME}:${USERNAME} /home/${USERNAME}

# =============================================================================
# Install Coding Agents
# =============================================================================

# Install Claude Code and Gemini CLI globally
RUN npm install -g --unsafe-perm \
    @anthropic-ai/claude-code \
    @google/gemini-cli \
    && npm cache clean --force

# =============================================================================
# Configure Git for SSH Signing
# =============================================================================

# Set up global git configuration for SSH signing
RUN git config --system gpg.format ssh \
    && git config --system commit.gpgsign true \
    && git config --system tag.gpgsign true

# =============================================================================
# SSH Agent Forwarding Setup
# =============================================================================

# Note: SSH_AUTH_SOCK is set at runtime via docker-compose or docker run
# Default socket path: /tmp/ssh-agent.sock (mounted from host)

# Create directory for SSH agent socket mount point
RUN mkdir -p /tmp && chmod 1777 /tmp

# Configure SSH for common Git hosting services
RUN mkdir -p /etc/ssh/ssh_config.d && \
    echo 'Host github.com gitlab.com bitbucket.org' > /etc/ssh/ssh_config.d/agent-forwarding.conf && \
    echo '    ForwardAgent yes' >> /etc/ssh/ssh_config.d/agent-forwarding.conf && \
    echo '    StrictHostKeyChecking accept-new' >> /etc/ssh/ssh_config.d/agent-forwarding.conf && \
    echo '    IdentitiesOnly yes' >> /etc/ssh/ssh_config.d/agent-forwarding.conf

# =============================================================================
# Entrypoint and Final Configuration
# =============================================================================

# Copy entrypoint script
COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod 755 /usr/local/bin/entrypoint.sh

# =============================================================================
# Copy Skills
# =============================================================================

# Copy skills to user's .claude/skills directory
RUN mkdir -p /home/${USERNAME}/.claude/skills
COPY --chown=${USERNAME}:${USERNAME} skills/ /home/${USERNAME}/.claude/skills/

# Set working directory
WORKDIR /workspace

# Switch to non-root user
USER ${USERNAME}

# Default shell
SHELL ["/bin/bash", "-c"]

# Labels for image metadata
LABEL org.opencontainers.image.title="Coding Agent Base Image" \
      org.opencontainers.image.description="Base image for AI coding agents with SSH signing support" \
      org.opencontainers.image.version="1.0.0"

# Entrypoint handles SSH agent and git configuration
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Default command opens an interactive shell
CMD ["/bin/bash"]
