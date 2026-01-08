#!/bin/bash
# =============================================================================
# Coding Agent Runner
# =============================================================================
# Cross-platform wrapper script for docker compose.
# Handles OS differences automatically.
#
# Usage:
#   ./run.sh up        # Start container
#   ./run.sh down      # Stop container
#   ./run.sh claude    # Run Claude Code
#   ./run.sh gemini    # Run Gemini CLI
#   ./run.sh shell     # Interactive shell
#   ./run.sh exec CMD  # Execute arbitrary command
# =============================================================================

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=lib/common.sh
source "${SCRIPT_DIR}/lib/common.sh"
cd "$SCRIPT_DIR" || exit 1

# =============================================================================
# Platform Detection
# =============================================================================

detect_platform() {
    case "$(uname -s)" in
        Linux*)
            echo "linux"
            ;;
        Darwin*)
            # On macOS, detect Docker Desktop vs Rancher Desktop
            if command -v rdctl &>/dev/null && rdctl list-settings &>/dev/null 2>&1; then
                echo "macos-rancher"
            else
                echo "macos-docker-desktop"
            fi
            ;;
        *)
            echo "linux"
            ;;
    esac
}

# =============================================================================
# Build docker compose command
# =============================================================================

get_compose_files() {
    local platform="$1"
    local files="-f docker-compose.yml"

    case "$platform" in
        macos-docker-desktop)
            files="$files -f docker-compose.macos.yml"
            ;;
        macos-rancher)
            files="$files -f docker-compose.rancher.yml"
            ;;
        linux)
            files="$files -f docker-compose.linux.yml"
            ;;
    esac

    echo "$files"
}

# =============================================================================
# Main
# =============================================================================

PLATFORM=$(detect_platform)
# shellcheck disable=SC2086 # Intentional word splitting for compose file flags
COMPOSE_FILES=$(get_compose_files "$PLATFORM")

show_help() {
    echo "Usage: $0 <command> [args...]"
    echo ""
    echo "Commands:"
    echo "  up       Start container (background)"
    echo "  down     Stop container"
    echo "  claude   Run Claude Code"
    echo "  gemini   Run Gemini CLI"
    echo "  shell    Interactive shell"
    echo "  exec     Execute arbitrary command"
    echo "  logs     Show logs"
    echo "  status   Show container status"
    echo "  build    Build image"
    echo ""
    echo "Detected platform: $PLATFORM"
}

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
    up)
        log_info "Platform: $PLATFORM"
        log_info "Starting container..."
        docker compose $COMPOSE_FILES up -d "$@"
        log_info "Container started. Run './run.sh claude' to start Claude Code."
        ;;
    down)
        log_info "Stopping container..."
        docker compose $COMPOSE_FILES down "$@"
        ;;
    claude)
        docker compose $COMPOSE_FILES exec coding-agent claude "$@"
        ;;
    gemini)
        docker compose $COMPOSE_FILES exec coding-agent gemini "$@"
        ;;
    shell|bash)
        docker compose $COMPOSE_FILES exec coding-agent bash "$@"
        ;;
    exec)
        docker compose $COMPOSE_FILES exec coding-agent "$@"
        ;;
    logs)
        docker compose $COMPOSE_FILES logs "$@"
        ;;
    status|ps)
        docker compose $COMPOSE_FILES ps "$@"
        ;;
    build)
        log_info "Building image..."
        docker compose $COMPOSE_FILES build "$@"
        ;;
    *)
        # Pass through other docker compose commands
        docker compose $COMPOSE_FILES "$COMMAND" "$@"
        ;;
esac
