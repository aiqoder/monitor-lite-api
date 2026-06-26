#!/bin/bash
set -euo pipefail

# wtv (monitor-lite-api) macOS 一键安装脚本
# 用法:
#   curl -fsSL "https://github.com/<owner>/<repo>/releases/latest/download/install-macos.sh" | bash -s install
#   curl -fsSL "..." | bash -s update
#   curl -fsSL "..." | bash -s uninstall

GITHUB_REPO="${GITHUB_REPO:-aiqoder/monitor-lite-api}"
GH_PROXY="${GH_PROXY:-}"
BINARY_PREFIX="${BINARY_PREFIX:-monitor-lite-api}"
SERVICE_NAME="${SERVICE_NAME:-com.wtv.server}"
INSTALL_PATH="${INSTALL_PATH:-$HOME/.wtv}"
VERSION="${VERSION:-latest}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}$*${NC}"; }
warn() { echo -e "${YELLOW}$*${NC}"; }
error() { echo -e "${RED}$*${NC}" >&2; }

detect_arch() {
  local platform
  platform="$(uname -m)"
  case "$platform" in
    x86_64) ARCH=amd64 ;;
    arm64) ARCH=arm64 ;;
    *)
      error "不支持的架构: $platform"
      exit 1
      ;;
  esac
}

ensure_ffmpeg() {
  if command -v ffmpeg >/dev/null 2>&1; then
    return
  fi

  warn "未检测到 ffmpeg"
  if command -v brew >/dev/null 2>&1; then
    info "正在通过 Homebrew 安装 ffmpeg..."
    brew install ffmpeg
  else
    error "请先安装 Homebrew 并执行: brew install ffmpeg"
    exit 1
  fi
}

download_url() {
  local asset="${BINARY_PREFIX}-darwin-${ARCH}"
  if [ "$VERSION" = "latest" ]; then
    echo "${GH_PROXY}https://github.com/${GITHUB_REPO}/releases/latest/download/${asset}"
  else
    echo "${GH_PROXY}https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${asset}"
  fi
}

download_binary() {
  local dest="$1"
  local url
  url="$(download_url)"
  info "下载 ${BINARY_PREFIX}-darwin-${ARCH} ..."
  curl -fsSL "$url" -o "$dest"
  chmod 755 "$dest"
}

install_launchd() {
  local plist_path="$HOME/Library/LaunchAgents/${SERVICE_NAME}.plist"
  mkdir -p "$HOME/Library/LaunchAgents"

  cat >"$plist_path" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>${SERVICE_NAME}</string>
  <key>ProgramArguments</key>
  <array>
    <string>${INSTALL_PATH}/wtv</string>
  </array>
  <key>WorkingDirectory</key>
  <string>${INSTALL_PATH}</string>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <true/>
  <key>StandardOutPath</key>
  <string>${INSTALL_PATH}/wtv.log</string>
  <key>StandardErrorPath</key>
  <string>${INSTALL_PATH}/wtv.err.log</string>
</dict>
</plist>
EOF

  launchctl bootout "gui/$(id -u)/${SERVICE_NAME}" >/dev/null 2>&1 || true
  launchctl bootstrap "gui/$(id -u)" "$plist_path"
  launchctl enable "gui/$(id -u)/${SERVICE_NAME}"
  launchctl kickstart -k "gui/$(id -u)/${SERVICE_NAME}"
}

do_install() {
  if [ -f "${INSTALL_PATH}/wtv" ]; then
    warn "检测到已安装在 ${INSTALL_PATH}，请使用 update 命令更新"
    exit 0
  fi

  mkdir -p "${INSTALL_PATH}"
  download_binary "${INSTALL_PATH}/wtv"
  install_launchd

  info "安装完成"
  echo
  echo "管理地址: http://localhost:9876/admin"
  echo "配置文件: ${INSTALL_PATH}/etc/tv.yaml"
  echo "默认账号: admin / admin123（请尽快修改）"
  echo "日志文件: ${INSTALL_PATH}/wtv.log"
}

do_update() {
  if [ ! -f "${INSTALL_PATH}/wtv" ]; then
    error "未检测到已安装的 wtv，请先执行 install"
    exit 1
  fi

  launchctl bootout "gui/$(id -u)/${SERVICE_NAME}" >/dev/null 2>&1 || true
  cp "${INSTALL_PATH}/wtv" /tmp/wtv.bak

  if ! download_binary /tmp/wtv.new; then
    mv /tmp/wtv.bak "${INSTALL_PATH}/wtv"
    install_launchd
    exit 1
  fi

  mv /tmp/wtv.new "${INSTALL_PATH}/wtv"
  rm -f /tmp/wtv.bak
  install_launchd
  info "已更新到最新版本"
}

do_uninstall() {
  launchctl bootout "gui/$(id -u)/${SERVICE_NAME}" >/dev/null 2>&1 || true
  rm -f "$HOME/Library/LaunchAgents/${SERVICE_NAME}.plist"
  rm -rf "${INSTALL_PATH}"
  info "wtv 已卸载"
}

main() {
  local action="${1:-}"
  if [ -n "${2:-}" ]; then
    INSTALL_PATH="$2"
  fi

  detect_arch
  ensure_ffmpeg

  case "$action" in
    install) do_install ;;
    update) do_update ;;
    uninstall) do_uninstall ;;
    *)
      error "用法: $0 {install|update|uninstall} [安装目录]"
      exit 1
      ;;
  esac
}

main "$@"
