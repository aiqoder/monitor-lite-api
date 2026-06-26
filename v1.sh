#!/bin/bash
# 兼容旧版入口，转发到 scripts/install.sh
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec bash "${SCRIPT_DIR}/scripts/install.sh" "$@"
