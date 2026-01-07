#!/bin/bash
# =============================================================================
# Coding Agent Runner
# =============================================================================
# OS差分を吸収してdocker composeを実行するラッパースクリプト
#
# Usage:
#   ./run.sh up        # コンテナを起動
#   ./run.sh down      # コンテナを停止
#   ./run.sh claude    # Claude Codeを実行
#   ./run.sh gemini    # Gemini CLIを実行
#   ./run.sh shell     # インタラクティブシェル
#   ./run.sh exec CMD  # 任意のコマンドを実行
# =============================================================================

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# =============================================================================
# プラットフォーム検出
# =============================================================================

detect_platform() {
    case "$(uname -s)" in
        Linux*)
            echo "linux"
            ;;
        Darwin*)
            # macOSの場合、Docker DesktopかRancher Desktopかを判定
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
# docker compose コマンド構築
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
# メイン
# =============================================================================

PLATFORM=$(detect_platform)
COMPOSE_FILES=$(get_compose_files "$PLATFORM")

# カラー出力
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }

show_help() {
    echo "Usage: $0 <command> [args...]"
    echo ""
    echo "Commands:"
    echo "  up       コンテナを起動（バックグラウンド）"
    echo "  down     コンテナを停止"
    echo "  claude   Claude Codeを実行"
    echo "  gemini   Gemini CLIを実行"
    echo "  shell    インタラクティブシェル"
    echo "  exec     任意のコマンドを実行"
    echo "  logs     ログを表示"
    echo "  status   コンテナの状態を表示"
    echo "  build    イメージをビルド"
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
        # その他のdocker composeコマンドをそのまま実行
        docker compose $COMPOSE_FILES "$COMMAND" "$@"
        ;;
esac
