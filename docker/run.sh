#!/bin/bash
# =============================================================================
# Coding Agent Runner (Docker)
# =============================================================================
# Wrapper script for managing Docker container with Claude Code and Gemini CLI.
# Supports SSH Agent Forwarding for secure git commit signing.
#
# Usage:
#   ./run.sh launch      # Build image and start container
#   ./run.sh ssh         # Interactive shell in container
#   ./run.sh claude      # Run Claude Code
#   ./run.sh gemini      # Run Gemini CLI
#   ./run.sh stop        # Stop container
#   ./run.sh delete      # Delete container and image
# =============================================================================

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=lib/common.sh
source "${SCRIPT_DIR}/lib/common.sh"

# =============================================================================
# Configuration Loading
# =============================================================================
# Priority (later overrides earlier):
#   1. Global config:  ~/.config/skills/config
#   2. Project config: $PWD/.skills.conf
# =============================================================================

# Load global config
if [ -f "$HOME/.config/skills/config" ]; then
    # shellcheck source=/dev/null
    source "$HOME/.config/skills/config"
fi

# Load project config (current working directory)
if [ -f "$PWD/.skills.conf" ]; then
    # shellcheck source=/dev/null
    source "$PWD/.skills.conf"
fi

# =============================================================================
# Configuration Defaults
# =============================================================================

CONTAINER_NAME="${DOCKER_CONTAINER_NAME:-coding-agent-docker}"
IMAGE_NAME="${DOCKER_IMAGE_NAME:-coding-agent:latest}"
CONTAINER_CPUS="${DOCKER_CPUS:-2}"
CONTAINER_MEMORY="${DOCKER_MEMORY:-4g}"

# Named volume for persistent data
CLAUDE_VOLUME="${CONTAINER_NAME}-claude-data"

# =============================================================================
# Helper Functions
# =============================================================================

# Convert TERM to a value that works in the container
get_container_term() {
    local term="${TERM:-xterm-256color}"
    case "$term" in
        xterm-ghostty|ghostty)
            echo "xterm-256color"
            ;;
        *)
            echo "$term"
            ;;
    esac
}

is_container_exists() {
    docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"
}

is_container_running() {
    docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"
}

is_image_exists() {
    docker image inspect "$IMAGE_NAME" &>/dev/null
}

# Setup SSH agent socket mount based on OS
setup_ssh_agent_mount() {
    case "$(uname -s)" in
        Darwin)
            # macOS: Use Docker Desktop's built-in SSH agent forwarding
            if [ -S /run/host-services/ssh-auth.sock ]; then
                echo "-v /run/host-services/ssh-auth.sock:/ssh-agent.sock -e SSH_AUTH_SOCK=/ssh-agent.sock"
            else
                log_warn "SSH agent socket not available on macOS"
                echo ""
            fi
            ;;
        Linux)
            # Linux: Mount the SSH agent socket directly
            if [ -n "${SSH_AUTH_SOCK:-}" ] && [ -S "$SSH_AUTH_SOCK" ]; then
                echo "-v $SSH_AUTH_SOCK:/ssh-agent.sock -e SSH_AUTH_SOCK=/ssh-agent.sock"
            else
                log_warn "SSH_AUTH_SOCK not set or invalid"
                echo ""
            fi
            ;;
        *)
            log_warn "Unknown OS, SSH agent forwarding may not work"
            echo ""
            ;;
    esac
}

# Get SSH public key content
get_ssh_pubkey_content() {
    local pubkey_path="${SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519.pub}"
    if [ -f "$pubkey_path" ]; then
        cat "$pubkey_path"
    else
        echo ""
    fi
}

# Build common environment variables for docker run/exec
get_env_args() {
    local args=""
    args="$args -e TERM=$(get_container_term)"
    args="$args -e COLORTERM=${COLORTERM:-truecolor}"

    if [ -n "${GIT_USER_NAME:-}" ]; then
        args="$args -e GIT_USER_NAME=$GIT_USER_NAME"
    fi

    if [ -n "${GIT_USER_EMAIL:-}" ]; then
        args="$args -e GIT_USER_EMAIL=$GIT_USER_EMAIL"
    fi

    local pubkey_content
    pubkey_content=$(get_ssh_pubkey_content)
    if [ -n "$pubkey_content" ]; then
        args="$args -e SSH_SIGNING_KEY_CONTENT=$pubkey_content"
    fi

    echo "$args"
}

# Setup Claude Code config in container
setup_claude_config() {
    local vm_settings_file="${SCRIPT_DIR}/vm-settings.json"
    if [ -f "$vm_settings_file" ]; then
        log_info "Copying VM-specific settings.json..."
        docker cp "$vm_settings_file" "${CONTAINER_NAME}:/home/agent/.claude/settings.json"
        log_success "VM settings.json copied."
    else
        log_warn "vm-settings.json not found at $vm_settings_file"
    fi
}

# Install Claude Code marketplaces in container
install_claude_marketplaces() {
    local marketplaces_file="${SCRIPT_DIR}/vm-marketplaces.txt"
    if [ ! -f "$marketplaces_file" ]; then
        log_warn "vm-marketplaces.txt not found, skipping marketplace installation."
        return 0
    fi

    log_info "Checking Claude Code marketplaces in container..."

    # Get list of installed marketplaces from container
    local installed
    installed=$(docker exec "$CONTAINER_NAME" cat /home/agent/.claude/plugins/known_marketplaces.json 2>/dev/null || echo "{}")

    local ssh_agent_opts
    ssh_agent_opts=$(setup_ssh_agent_mount)
    local env_args
    env_args=$(get_env_args)

    while IFS= read -r line || [ -n "$line" ]; do
        # Skip empty lines and comments
        [[ -z "$line" || "$line" =~ ^# ]] && continue

        local repo="$line"
        # Check if marketplace is already installed
        if echo "$installed" | grep -q "\"repo\": \"$repo\""; then
            log_info "  Already installed: $repo"
        else
            log_info "  Installing marketplace: $repo"
            # shellcheck disable=SC2086
            docker exec -it $ssh_agent_opts $env_args "$CONTAINER_NAME" \
                bash -ic "claude marketplace add github:$repo" 2>/dev/null || true
        fi
    done < "$marketplaces_file"

    log_success "Marketplaces check complete."
}

# Configure Git in container
configure_git_in_container() {
    log_info "Configuring Git in container..."

    local env_args
    env_args=$(get_env_args)

    # Run entrypoint to configure git (it reads env vars)
    # shellcheck disable=SC2086
    docker exec $env_args "$CONTAINER_NAME" /entrypoint.sh true

    if [ -n "${GIT_USER_NAME:-}" ]; then
        log_info "Git user.name set to: $GIT_USER_NAME"
    else
        log_warn "GIT_USER_NAME not set. Set it in ~/.config/skills/config or .skills.conf"
    fi

    if [ -n "${GIT_USER_EMAIL:-}" ]; then
        log_info "Git user.email set to: $GIT_USER_EMAIL"
    else
        log_warn "GIT_USER_EMAIL not set. Set it in ~/.config/skills/config or .skills.conf"
    fi

    local pubkey_content
    pubkey_content=$(get_ssh_pubkey_content)
    if [ -n "$pubkey_content" ]; then
        log_info "SSH signing key configured."
    else
        log_warn "SSH public key not found. Git signing will not work."
    fi
}

# Mount config directories
get_config_mount_args() {
    local args=""

    # Mount ~/.gemini if exists
    if [ -d "$HOME/.gemini" ]; then
        args="$args -v $HOME/.gemini:/home/agent/.gemini:ro"
    fi

    # Mount ~/.aws if exists
    if [ -d "$HOME/.aws" ]; then
        args="$args -v $HOME/.aws:/home/agent/.aws:ro"
    fi

    echo "$args"
}

# =============================================================================
# Commands
# =============================================================================

cmd_launch() {
    # Build image if needed
    if ! is_image_exists; then
        log_info "Building Docker image..."
        docker build \
            --build-arg USER_ID="$(id -u)" \
            --build-arg GROUP_ID="$(id -g)" \
            -t "$IMAGE_NAME" \
            "$SCRIPT_DIR"
    else
        log_info "Docker image '$IMAGE_NAME' already exists."
    fi

    # Create volume for persistent data
    docker volume create "$CLAUDE_VOLUME" 2>/dev/null || true

    if is_container_exists; then
        log_warn "Container '$CONTAINER_NAME' already exists."
        if is_container_running; then
            log_info "Container is already running."
        else
            log_info "Starting existing container..."
            docker start "$CONTAINER_NAME"
        fi
    else
        log_info "Creating container '$CONTAINER_NAME'..."
        log_info "  CPUs: $CONTAINER_CPUS"
        log_info "  Memory: $CONTAINER_MEMORY"

        local config_mounts
        config_mounts=$(get_config_mount_args)
        local ssh_agent_opts
        ssh_agent_opts=$(setup_ssh_agent_mount)
        local env_args
        env_args=$(get_env_args)

        # shellcheck disable=SC2086
        docker create \
            --name "$CONTAINER_NAME" \
            --hostname coding-agent \
            --cpus "$CONTAINER_CPUS" \
            --memory "$CONTAINER_MEMORY" \
            -v "$CLAUDE_VOLUME:/home/agent/.claude" \
            $config_mounts \
            $ssh_agent_opts \
            $env_args \
            -it \
            "$IMAGE_NAME" \
            tail -f /dev/null

        docker start "$CONTAINER_NAME"
    fi

    configure_git_in_container
    setup_claude_config
    install_claude_marketplaces

    log_success ""
    log_success "Container '$CONTAINER_NAME' is ready!"
    echo ""

    # Run claude /login for initial authentication
    log_info "Running Claude Code login..."
    local ssh_agent_opts
    ssh_agent_opts=$(setup_ssh_agent_mount)
    local env_args
    env_args=$(get_env_args)

    # shellcheck disable=SC2086
    docker exec -it $ssh_agent_opts $env_args "$CONTAINER_NAME" \
        bash -ic "claude /login"

    echo ""
    echo "Use:"
    echo "  ./run.sh ssh      # Interactive shell"
    echo "  ./run.sh claude   # Run Claude Code"
    echo "  ./run.sh gemini   # Run Gemini CLI"
}

cmd_start() {
    if ! is_container_exists; then
        log_error "Container '$CONTAINER_NAME' does not exist. Run './run.sh launch' first."
        exit 1
    fi

    if is_container_running; then
        log_info "Container '$CONTAINER_NAME' is already running."
    else
        log_info "Starting container '$CONTAINER_NAME'..."
        docker start "$CONTAINER_NAME"
        log_success "Container started."
    fi

    configure_git_in_container
    setup_claude_config
}

cmd_stop() {
    if ! is_container_exists; then
        log_warn "Container '$CONTAINER_NAME' does not exist."
        return 0
    fi

    log_info "Stopping container '$CONTAINER_NAME'..."
    docker stop "$CONTAINER_NAME"
    log_success "Container stopped."
}

cmd_delete() {
    if ! is_container_exists && ! is_image_exists; then
        log_warn "Container '$CONTAINER_NAME' and image '$IMAGE_NAME' do not exist."
        return 0
    fi

    log_warn "This will delete:"
    is_container_exists && echo "  - Container: $CONTAINER_NAME"
    is_image_exists && echo "  - Image: $IMAGE_NAME"
    echo "  - Volume: $CLAUDE_VOLUME (Claude credentials)"

    read -rp "Are you sure? [y/N] " answer
    if [[ "$answer" =~ ^[Yy]$ ]]; then
        if is_container_running; then
            log_info "Stopping container..."
            docker stop "$CONTAINER_NAME"
        fi
        if is_container_exists; then
            log_info "Removing container..."
            docker rm "$CONTAINER_NAME"
        fi
        if is_image_exists; then
            log_info "Removing image..."
            docker rmi "$IMAGE_NAME"
        fi
        log_info "Removing volume..."
        docker volume rm "$CLAUDE_VOLUME" 2>/dev/null || true
        log_success "Deleted."
    else
        log_info "Cancelled."
    fi
}

cmd_ssh() {
    if ! is_image_exists; then
        log_error "Image '$IMAGE_NAME' does not exist. Run './run.sh launch' first."
        exit 1
    fi

    local src_path
    src_path="$(pwd)"
    local mount_point="/workspace"

    local ssh_agent_opts
    ssh_agent_opts=$(setup_ssh_agent_mount)
    local env_args
    env_args=$(get_env_args)
    local config_mounts
    config_mounts=$(get_config_mount_args)

    log_info "Working directory: $mount_point (mounted from $src_path)"

    if [ $# -eq 0 ]; then
        # Interactive shell with current directory mounted
        # shellcheck disable=SC2086
        docker run --rm -it \
            --hostname coding-agent \
            -v "$CLAUDE_VOLUME:/home/agent/.claude" \
            -v "$src_path:$mount_point" \
            -w "$mount_point" \
            $config_mounts \
            $ssh_agent_opts \
            $env_args \
            "$IMAGE_NAME" \
            bash
    else
        # Execute command
        # shellcheck disable=SC2086
        docker run --rm -it \
            --hostname coding-agent \
            -v "$CLAUDE_VOLUME:/home/agent/.claude" \
            -v "$src_path:$mount_point" \
            -w "$mount_point" \
            $config_mounts \
            $ssh_agent_opts \
            $env_args \
            "$IMAGE_NAME" \
            bash -c "$*"
    fi
}

cmd_claude() {
    if ! is_image_exists; then
        log_error "Image '$IMAGE_NAME' does not exist. Run './run.sh launch' first."
        exit 1
    fi

    # Ensure container is running for marketplace installation
    if is_container_running; then
        setup_claude_config
        install_claude_marketplaces
    fi

    local src_path
    src_path="$(pwd)"
    local mount_point="/workspace"

    local ssh_agent_opts
    ssh_agent_opts=$(setup_ssh_agent_mount)
    local env_args
    env_args=$(get_env_args)
    local config_mounts
    config_mounts=$(get_config_mount_args)

    log_info "Working directory: $mount_point (mounted from $src_path)"

    # shellcheck disable=SC2086
    docker run --rm -it \
        --hostname coding-agent \
        -v "$CLAUDE_VOLUME:/home/agent/.claude" \
        -v "$src_path:$mount_point" \
        -w "$mount_point" \
        $config_mounts \
        $ssh_agent_opts \
        $env_args \
        "$IMAGE_NAME" \
        bash -ic "claude $*"
}

cmd_gemini() {
    if ! is_image_exists; then
        log_error "Image '$IMAGE_NAME' does not exist. Run './run.sh launch' first."
        exit 1
    fi

    local src_path
    src_path="$(pwd)"
    local mount_point="/workspace"

    local ssh_agent_opts
    ssh_agent_opts=$(setup_ssh_agent_mount)
    local env_args
    env_args=$(get_env_args)
    local config_mounts
    config_mounts=$(get_config_mount_args)

    log_info "Working directory: $mount_point (mounted from $src_path)"

    # shellcheck disable=SC2086
    docker run --rm -it \
        --hostname coding-agent \
        -v "$CLAUDE_VOLUME:/home/agent/.claude" \
        -v "$src_path:$mount_point" \
        -w "$mount_point" \
        $config_mounts \
        $ssh_agent_opts \
        $env_args \
        "$IMAGE_NAME" \
        bash -ic "gemini $*"
}

cmd_exec() {
    if ! is_container_running; then
        log_error "Container '$CONTAINER_NAME' is not running."
        exit 1
    fi

    local ssh_agent_opts
    ssh_agent_opts=$(setup_ssh_agent_mount)
    local env_args
    env_args=$(get_env_args)

    # shellcheck disable=SC2086
    docker exec -it $ssh_agent_opts $env_args "$CONTAINER_NAME" "$@"
}

cmd_mount() {
    local source_path="$1"
    local target_path="${2:-/workspace}"

    if [ -z "$source_path" ]; then
        log_error "Usage: ./run.sh mount <source_path> [target_path]"
        exit 1
    fi

    if ! is_container_running; then
        log_error "Container '$CONTAINER_NAME' is not running."
        exit 1
    fi

    log_warn "Docker containers cannot add mounts dynamically after creation."
    log_info "Instead, use docker cp to copy files:"
    log_info "  docker cp $source_path ${CONTAINER_NAME}:${target_path}"
    log_info ""
    log_info "Or use ./run.sh claude or ./run.sh ssh which automatically mount PWD."
}

cmd_status() {
    if ! is_container_exists; then
        log_warn "Container '$CONTAINER_NAME' does not exist."
        return 0
    fi

    docker inspect --format='Container: {{.Name}}
State: {{.State.Status}}
Created: {{.Created}}
Image: {{.Config.Image}}
' "$CONTAINER_NAME"

    # Show mounts
    echo "Mounts:"
    docker inspect --format='{{range .Mounts}}  {{.Type}}: {{.Source}} -> {{.Destination}}
{{end}}' "$CONTAINER_NAME"
}

cmd_logs() {
    local lines="${1:-100}"

    if ! is_container_exists; then
        log_error "Container '$CONTAINER_NAME' does not exist."
        exit 1
    fi

    docker logs --tail "$lines" "$CONTAINER_NAME"
}

cmd_configure_git() {
    if ! is_container_running; then
        log_error "Container '$CONTAINER_NAME' is not running."
        exit 1
    fi

    configure_git_in_container
}

show_help() {
    echo "Usage: $0 <command> [args...]"
    echo ""
    echo "Commands:"
    echo "  launch          Build image and start container"
    echo "  start           Start existing container"
    echo "  stop            Stop container"
    echo "  delete          Delete container, image, and volume"
    echo "  ssh             Interactive shell in container"
    echo "  claude [args]   Run Claude Code"
    echo "  gemini [args]   Run Gemini CLI"
    echo "  exec <cmd>      Execute arbitrary command"
    echo "  mount <src>     Show mount instructions (Docker limitation)"
    echo "  status          Show container status"
    echo "  logs [lines]    Show container logs (default: 100 lines)"
    echo "  configure-git   Re-configure git settings"
    echo ""
    echo "Configuration files (loaded in order, later overrides earlier):"
    echo "  ~/.config/skills/config  Global configuration"
    echo "  .skills.conf             Project configuration (current directory)"
    echo ""
    echo "Environment variables:"
    echo "  DOCKER_CONTAINER_NAME  Container name (default: coding-agent-docker)"
    echo "  DOCKER_IMAGE_NAME      Image name (default: coding-agent:latest)"
    echo "  DOCKER_CPUS            CPU limit (default: 2)"
    echo "  DOCKER_MEMORY          Memory limit (default: 4g)"
    echo "  GIT_USER_NAME          Git user name"
    echo "  GIT_USER_EMAIL         Git user email"
    echo "  SSH_SIGNING_KEY        Path to SSH public key (default: ~/.ssh/id_ed25519.pub)"
}

# =============================================================================
# Main
# =============================================================================

# Check for docker
if ! command_exists docker; then
    log_error "docker is not installed."
    echo "Install Docker Desktop from https://www.docker.com/products/docker-desktop/"
    exit 1
fi

if [ $# -eq 0 ]; then
    show_help
    exit 1
fi

COMMAND="$1"
shift

case "$COMMAND" in
    -h|--help|help)
        show_help
        exit 0
        ;;
    launch|create)
        cmd_launch "$@"
        ;;
    start)
        cmd_start "$@"
        ;;
    stop)
        cmd_stop "$@"
        ;;
    delete|destroy)
        cmd_delete "$@"
        ;;
    ssh|shell)
        cmd_ssh "$@"
        ;;
    claude)
        cmd_claude "$@"
        ;;
    gemini)
        cmd_gemini "$@"
        ;;
    exec)
        cmd_exec "$@"
        ;;
    mount)
        cmd_mount "$@"
        ;;
    status|info)
        cmd_status "$@"
        ;;
    logs)
        cmd_logs "$@"
        ;;
    configure-git)
        cmd_configure_git "$@"
        ;;
    *)
        log_error "Unknown command: $COMMAND"
        show_help
        exit 1
        ;;
esac
