# =============================================================================
# Coding Agent Docker Image
# =============================================================================
# Docker image for Claude Code and Gemini CLI development environment.
# Based on Ubuntu 24.04 with Node.js 22.x, mise, and GitHub CLI.
#
# Build:
#   docker build --build-arg USER_ID=$(id -u) --build-arg GROUP_ID=$(id -g) -t coding-agent .
#
# Run:
#   ./run.sh launch
# =============================================================================

FROM ubuntu:24.04

# Build arguments for user creation (match host user for file permissions)
ARG USER_ID=1000
ARG GROUP_ID=1000
ARG USERNAME=agent

# Prevent interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# =============================================================================
# System Packages
# =============================================================================

RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    git \
    jq \
    openssh-client \
    build-essential \
    trash-cli \
    && rm -rf /var/lib/apt/lists/*

# =============================================================================
# Node.js 22.x
# =============================================================================

RUN mkdir -p /etc/apt/keyrings \
    && curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_22.x nodistro main" > /etc/apt/sources.list.d/nodesource.list \
    && apt-get update \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# =============================================================================
# GitHub CLI
# =============================================================================

RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg \
    && chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" > /etc/apt/sources.list.d/github-cli.list \
    && apt-get update \
    && apt-get install -y gh \
    && rm -rf /var/lib/apt/lists/*

# =============================================================================
# Claude Code and Gemini CLI
# =============================================================================

RUN npm install -g @anthropic-ai/claude-code @google/gemini-cli

# =============================================================================
# SSH Configuration for Git Hosts
# =============================================================================

RUN mkdir -p /etc/ssh/ssh_config.d \
    && echo "Host github.com gitlab.com bitbucket.org\n    StrictHostKeyChecking accept-new\n    IdentitiesOnly yes" > /etc/ssh/ssh_config.d/git-hosts.conf

# =============================================================================
# Git SSH Signing Configuration (system-wide)
# =============================================================================

RUN git config --system gpg.format ssh \
    && git config --system commit.gpgsign true \
    && git config --system tag.gpgsign true

# =============================================================================
# Create User with Matching UID/GID
# =============================================================================

RUN groupadd -g ${GROUP_ID} ${USERNAME} 2>/dev/null || groupmod -n ${USERNAME} $(getent group ${GROUP_ID} | cut -d: -f1) \
    && useradd -m -u ${USER_ID} -g ${USERNAME} -s /bin/bash ${USERNAME} 2>/dev/null || true

# =============================================================================
# User Environment Setup
# =============================================================================

USER ${USERNAME}
WORKDIR /home/${USERNAME}

# Install mise (runtime manager)
RUN curl https://mise.run | sh \
    && echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc \
    && echo 'eval "$(~/.local/bin/mise activate bash)"' >> ~/.bashrc \
    && echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.profile \
    && echo 'eval "$(~/.local/bin/mise activate bash)"' >> ~/.profile

# Create necessary directories
RUN mkdir -p ~/.claude ~/.ssh \
    && chmod 700 ~/.ssh

# Add github.com to known_hosts
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts 2>/dev/null

# =============================================================================
# Entrypoint
# =============================================================================

COPY --chown=${USERNAME}:${USERNAME} entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["bash"]
