#!/bin/bash
# =============================================================================
# Common Shell Library
# =============================================================================
# Shared functions for all shell scripts in this repository.
# Source this file at the beginning of each script.
#
# Usage:
#   SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
#   source "${SCRIPT_DIR}/lib/common.sh"
# =============================================================================

# Strict error handling (can be overridden after sourcing if needed)
set -euo pipefail

# =============================================================================
# Colors
# =============================================================================

export RED='\033[0;31m'
export GREEN='\033[0;32m'
export YELLOW='\033[0;33m'
export NC='\033[0m'  # No Color

# =============================================================================
# Logging Functions
# =============================================================================

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

log_success() {
    echo -e "${GREEN}$1${NC}"
}

# =============================================================================
# Utility Functions
# =============================================================================

# Check if a command exists
command_exists() {
    command -v "$1" &>/dev/null
}
