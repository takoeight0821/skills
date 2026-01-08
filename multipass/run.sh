#!/bin/bash
# =============================================================================
# Coding Agent Runner (Multipass)
# =============================================================================
# Wrapper script for managing Multipass VM with Claude Code and Gemini CLI.
# Supports SSH Agent Forwarding for secure git commit signing.
#
# Usage:
#   ./run.sh launch      # Create and start VM
#   ./run.sh ssh         # SSH with agent forwarding
#   ./run.sh claude      # Run Claude Code
#   ./run.sh gemini      # Run Gemini CLI
#   ./run.sh stop        # Stop VM
#   ./run.sh delete      # Delete VM
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

VM_NAME="${MULTIPASS_VM_NAME:-coding-agent}"
VM_CPUS="${MULTIPASS_VM_CPUS:-2}"
VM_MEMORY="${MULTIPASS_VM_MEMORY:-4G}"
VM_DISK="${MULTIPASS_VM_DISK:-20G}"
CLOUD_INIT="${SCRIPT_DIR}/cloud-init.yaml"

# =============================================================================
# Helper Functions
# =============================================================================

get_vm_ip() {
    multipass info "$VM_NAME" 2>/dev/null | awk '/IPv4/ {print $2}'
}

is_vm_running() {
    multipass info "$VM_NAME" 2>/dev/null | grep -q "State:.*Running"
}

is_vm_exists() {
    multipass info "$VM_NAME" &>/dev/null
}

# パスからユニークなマウントポイントを生成
get_mount_point() {
    local src_path="$1"
    # 先頭の / を除去し、残りの / を - に置換
    local sanitized="${src_path#/}"
    sanitized="${sanitized//\//-}"
    echo "/mnt/${sanitized}"
}

# 指定パスがマウント済みかチェック
is_path_mounted() {
    local mount_point="$1"
    multipass info "$VM_NAME" 2>/dev/null | grep -q "$mount_point"
}

# 自動マウント（カレントディレクトリ → /mnt/xxx）
auto_mount() {
    local src_path
    src_path="$(pwd)"
    local mount_point
    mount_point=$(get_mount_point "$src_path")

    if ! is_path_mounted "$mount_point"; then
        log_info "Auto-mounting $src_path to $mount_point..." >&2
        multipass exec "$VM_NAME" -- sudo mkdir -p "$mount_point" >&2
        multipass mount "$src_path" "${VM_NAME}:${mount_point}" >&2
    fi

    echo "$mount_point"
}

# 自動アンマウント
auto_umount() {
    local mount_point="$1"
    if is_path_mounted "$mount_point"; then
        log_info "Auto-unmounting $mount_point..."
        multipass umount "${VM_NAME}:${mount_point}" 2>/dev/null || true
    fi
}

show_cloud_init_logs() {
    local lines="${1:-100}"
    log_info "Cloud-init status:"
    multipass exec "$VM_NAME" -- cloud-init status --long 2>&1 || true
    echo ""
    log_info "Cloud-init output log (last $lines lines):"
    multipass exec "$VM_NAME" -- tail -"$lines" /var/log/cloud-init-output.log 2>/dev/null || echo "  (log not available)"
}

wait_for_cloud_init() {
    log_info "Waiting for cloud-init to complete..."
    log_info "You can check logs anytime with: ./run.sh logs"
    local max_wait=300
    local waited=0
    local status_output
    while [ $waited -lt $max_wait ]; do
        status_output=$(multipass exec "$VM_NAME" -- cloud-init status --long 2>&1 || true)

        if echo "$status_output" | grep -q "status: done"; then
            echo ""
            log_success "Cloud-init completed."
            return 0
        fi

        if echo "$status_output" | grep -q "status: error"; then
            echo ""
            log_error "Cloud-init failed!"
            show_cloud_init_logs 100
            return 1
        fi

        sleep 5
        waited=$((waited + 5))
        echo -n "."

        # Show progress every 30 seconds
        if [ $((waited % 30)) -eq 0 ]; then
            echo ""
            log_info "Still waiting... ($waited/${max_wait}s)"
            log_info "Current status: $(echo "$status_output" | grep -E '^status:' || echo 'unknown')"
        fi
    done
    echo ""
    log_warn "Cloud-init timed out after ${max_wait}s."
    show_cloud_init_logs 100
    return 1
}

configure_git_in_vm() {
    log_info "Configuring Git in VM..."

    # Set user name
    if [ -n "${GIT_USER_NAME:-}" ]; then
        multipass exec "$VM_NAME" -- git config --global user.name "$GIT_USER_NAME"
        log_info "Git user.name set to: $GIT_USER_NAME"
    else
        log_warn "GIT_USER_NAME not set. Set it in ~/.config/skills/config or .skills.conf"
    fi

    # Set user email
    if [ -n "${GIT_USER_EMAIL:-}" ]; then
        multipass exec "$VM_NAME" -- git config --global user.email "$GIT_USER_EMAIL"
        log_info "Git user.email set to: $GIT_USER_EMAIL"
    else
        log_warn "GIT_USER_EMAIL not set. Set it in ~/.config/skills/config or .skills.conf"
    fi

    # Configure SSH signing key
    local pubkey_path="${SSH_SIGNING_KEY:-$HOME/.ssh/id_ed25519.pub}"
    if [ -f "$pubkey_path" ]; then
        local pubkey_content
        pubkey_content=$(cat "$pubkey_path")

        # Write public key to VM and add to authorized_keys for SSH access
        multipass exec "$VM_NAME" -- bash -c "echo '$pubkey_content' > ~/.ssh/id_ed25519.pub && chmod 644 ~/.ssh/id_ed25519.pub"
        multipass exec "$VM_NAME" -- bash -c "grep -qxF '$pubkey_content' ~/.ssh/authorized_keys 2>/dev/null || echo '$pubkey_content' >> ~/.ssh/authorized_keys"
        log_info "SSH public key added to authorized_keys."

        # Configure Git to use this key for signing
        multipass exec "$VM_NAME" -- bash -c "git config --global user.signingkey ~/.ssh/id_ed25519.pub"
        log_info "Git signing key configured."

        # Set up allowed signers for verification
        if [ -n "${GIT_USER_EMAIL:-}" ]; then
            multipass exec "$VM_NAME" -- bash -c "echo '$GIT_USER_EMAIL namespaces=\"git\" $pubkey_content' > ~/.ssh/allowed_signers"
            multipass exec "$VM_NAME" -- bash -c "git config --global gpg.ssh.allowedSignersFile ~/.ssh/allowed_signers"
            log_info "Allowed signers file configured."
        fi
    else
        log_warn "SSH public key not found at $pubkey_path"
        log_warn "Git signing will not work until key is configured."
    fi
}

sync_skills_to_vm() {
    local skills_source="${SCRIPT_DIR}/../skills"
    if [ -d "$skills_source" ]; then
        log_info "Syncing skills to VM..."
        multipass transfer -r "$skills_source"/* "${VM_NAME}:/home/ubuntu/.claude/skills/"
        log_success "Skills synced."
    fi
}

# =============================================================================
# Commands
# =============================================================================

cmd_launch() {
    if is_vm_exists; then
        log_warn "VM '$VM_NAME' already exists."
        if ! is_vm_running; then
            log_info "Starting existing VM..."
            multipass start "$VM_NAME"
        fi
    else
        log_info "Creating VM '$VM_NAME'..."
        log_info "  CPUs: $VM_CPUS"
        log_info "  Memory: $VM_MEMORY"
        log_info "  Disk: $VM_DISK"

        multipass launch \
            --name "$VM_NAME" \
            --cpus "$VM_CPUS" \
            --memory "$VM_MEMORY" \
            --disk "$VM_DISK" \
            --cloud-init "$CLOUD_INIT"

        wait_for_cloud_init
    fi

    configure_git_in_vm
    sync_skills_to_vm

    local ip
    ip=$(get_vm_ip)
    log_success ""
    log_success "VM '$VM_NAME' is ready!"
    log_success "IP Address: $ip"
    echo ""
    echo "Connect with:"
    echo "  ssh -A ubuntu@$ip"
    echo ""
    echo "Or use:"
    echo "  ./run.sh ssh      # Interactive shell"
    echo "  ./run.sh claude   # Run Claude Code"
    echo "  ./run.sh gemini   # Run Gemini CLI"
}

cmd_start() {
    if ! is_vm_exists; then
        log_error "VM '$VM_NAME' does not exist. Run './run.sh launch' first."
        exit 1
    fi

    if is_vm_running; then
        log_info "VM '$VM_NAME' is already running."
    else
        log_info "Starting VM '$VM_NAME'..."
        multipass start "$VM_NAME"
        log_success "VM started."
    fi

    local ip
    ip=$(get_vm_ip)
    log_info "IP Address: $ip"
}

cmd_stop() {
    if ! is_vm_exists; then
        log_warn "VM '$VM_NAME' does not exist."
        return 0
    fi

    log_info "Stopping VM '$VM_NAME'..."
    multipass stop "$VM_NAME"
    log_success "VM stopped."
}

cmd_delete() {
    if ! is_vm_exists; then
        log_warn "VM '$VM_NAME' does not exist."
        return 0
    fi

    log_warn "This will permanently delete VM '$VM_NAME' and all its data."
    read -rp "Are you sure? [y/N] " answer
    if [[ "$answer" =~ ^[Yy]$ ]]; then
        log_info "Deleting VM '$VM_NAME'..."
        multipass delete "$VM_NAME"
        multipass purge
        log_success "VM deleted."
    else
        log_info "Cancelled."
    fi
}

cmd_ssh() {
    if ! is_vm_running; then
        log_error "VM '$VM_NAME' is not running. Run './run.sh launch' or './run.sh start' first."
        exit 1
    fi

    local mount_point
    mount_point=$(auto_mount)
    trap "auto_umount '$mount_point'" EXIT

    local ip
    ip=$(get_vm_ip)
    log_info "Connecting to VM with SSH Agent Forwarding..."
    log_info "IP: $ip"
    log_info "Working directory: $mount_point"

    if [ $# -eq 0 ]; then
        # Interactive shell
        ssh -A -t "ubuntu@$ip" "cd '$mount_point' && exec \$SHELL"
    else
        # Execute command
        ssh -A "ubuntu@$ip" "cd '$mount_point' && $*"
    fi
}

cmd_claude() {
    if ! is_vm_running; then
        log_error "VM '$VM_NAME' is not running."
        exit 1
    fi

    local mount_point
    mount_point=$(auto_mount)
    trap "auto_umount '$mount_point'" EXIT

    local ip
    ip=$(get_vm_ip)
    log_info "Working directory: $mount_point"
    ssh -A -t "ubuntu@$ip" "cd '$mount_point' && claude $*"
}

cmd_gemini() {
    if ! is_vm_running; then
        log_error "VM '$VM_NAME' is not running."
        exit 1
    fi

    local mount_point
    mount_point=$(auto_mount)
    trap "auto_umount '$mount_point'" EXIT

    local ip
    ip=$(get_vm_ip)
    log_info "Working directory: $mount_point"
    ssh -A -t "ubuntu@$ip" "cd '$mount_point' && gemini $*"
}

cmd_exec() {
    if ! is_vm_running; then
        log_error "VM '$VM_NAME' is not running."
        exit 1
    fi

    local ip
    ip=$(get_vm_ip)
    ssh -A "ubuntu@$ip" "$@"
}

cmd_mount() {
    local source_path="$1"
    local target_path="${2:-/workspace}"

    if [ -z "$source_path" ]; then
        log_error "Usage: ./run.sh mount <source_path> [target_path]"
        exit 1
    fi

    if ! is_vm_running; then
        log_error "VM '$VM_NAME' is not running."
        exit 1
    fi

    log_info "Mounting $source_path to $VM_NAME:$target_path..."
    multipass mount "$source_path" "${VM_NAME}:${target_path}"
    log_success "Mounted."
}

cmd_umount() {
    local target_path="${1:-/workspace}"

    if ! is_vm_running; then
        log_error "VM '$VM_NAME' is not running."
        exit 1
    fi

    log_info "Unmounting $VM_NAME:$target_path..."
    multipass umount "${VM_NAME}:${target_path}"
    log_success "Unmounted."
}

cmd_status() {
    if ! is_vm_exists; then
        log_warn "VM '$VM_NAME' does not exist."
        return 0
    fi

    multipass info "$VM_NAME"
}

cmd_configure_git() {
    if ! is_vm_running; then
        log_error "VM '$VM_NAME' is not running."
        exit 1
    fi

    configure_git_in_vm
}

cmd_logs() {
    local lines="${1:-100}"

    if ! is_vm_exists; then
        log_error "VM '$VM_NAME' does not exist."
        exit 1
    fi

    if ! is_vm_running; then
        log_error "VM '$VM_NAME' is not running."
        exit 1
    fi

    show_cloud_init_logs "$lines"
}

show_help() {
    echo "Usage: $0 <command> [args...]"
    echo ""
    echo "Commands:"
    echo "  launch          Create and start VM with cloud-init"
    echo "  start           Start existing VM"
    echo "  stop            Stop VM"
    echo "  delete          Delete VM permanently"
    echo "  ssh             SSH with agent forwarding"
    echo "  claude [args]   Run Claude Code"
    echo "  gemini [args]   Run Gemini CLI"
    echo "  exec <cmd>      Execute arbitrary command"
    echo "  mount <src>     Mount directory to /workspace"
    echo "  umount          Unmount /workspace"
    echo "  status          Show VM status"
    echo "  logs [lines]    Show cloud-init logs (default: 100 lines)"
    echo "  configure-git   Re-configure git settings"
    echo ""
    echo "Configuration files (loaded in order, later overrides earlier):"
    echo "  ~/.config/skills/config  Global configuration"
    echo "  .skills.conf             Project configuration (current directory)"
    echo ""
    echo "Environment variables:"
    echo "  MULTIPASS_VM_NAME    VM name (default: coding-agent)"
    echo "  MULTIPASS_VM_CPUS    CPU count (default: 2)"
    echo "  MULTIPASS_VM_MEMORY  Memory size (default: 4G)"
    echo "  MULTIPASS_VM_DISK    Disk size (default: 20G)"
    echo "  GIT_USER_NAME        Git user name"
    echo "  GIT_USER_EMAIL       Git user email"
    echo "  SSH_SIGNING_KEY      Path to SSH public key (default: ~/.ssh/id_ed25519.pub)"
}

# =============================================================================
# Main
# =============================================================================

# Check for multipass
if ! command_exists multipass; then
    log_error "multipass is not installed."
    echo "Install with: brew install --cask multipass"
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
    umount|unmount)
        cmd_umount "$@"
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
