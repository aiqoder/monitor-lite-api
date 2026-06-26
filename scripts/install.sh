#!/bin/bash
set -euo pipefail

# wtv (monitor-lite-api) Linux 一键安装脚本
# 用法:
#   curl -fsSL "https://github.com/<owner>/<repo>/releases/latest/download/install.sh" | sudo bash -s install
#   curl -fsSL "https://github.com/<owner>/<repo>/releases/latest/download/install.sh" | sudo bash -s update
#   curl -fsSL "https://github.com/<owner>/<repo>/releases/latest/download/install.sh" | sudo bash -s uninstall
#   curl -fsSL "..." | sudo bash -s install /opt/custom

GITHUB_REPO="${GITHUB_REPO:-aiqoder/monitor-lite-api}"
GH_PROXY="${GH_PROXY:-}"
BINARY_PREFIX="${BINARY_PREFIX:-monitor-lite-api}"
SERVICE_NAME="${SERVICE_NAME:-wtv}"
INSTALL_PATH="${INSTALL_PATH:-/opt/wtv}"
VERSION="${VERSION:-latest}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { echo -e "${GREEN}$*${NC}"; }
warn() { echo -e "${YELLOW}$*${NC}"; }
error() { echo -e "${RED}$*${NC}" >&2; }

resolve_install_path() {
  if [ -n "${2:-}" ]; then
    local custom_path="$2"
    if [[ "$custom_path" == */ ]]; then
      custom_path="${custom_path%/}"
    fi
    if [[ "$custom_path" != */wtv ]]; then
      custom_path="$custom_path/wtv"
    fi
    INSTALL_PATH="$custom_path"
  fi
}

detect_arch() {
  local platform
  if command -v arch >/dev/null 2>&1; then
    platform="$(arch)"
  else
    platform="$(uname -m)"
  fi

  case "$platform" in
    x86_64|amd64) ARCH=amd64 ;;
    aarch64|arm64) ARCH=arm64 ;;
    *)
      error "不支持的架构: $platform（仅支持 amd64 / arm64）"
      exit 1
      ;;
  esac
}

require_root() {
  if [ "$(id -u)" -ne 0 ]; then
    error "请使用 root 权限运行（sudo）"
    exit 1
  fi
}

require_systemd() {
  if ! command -v systemctl >/dev/null 2>&1; then
    error "未检测到 systemd，当前脚本仅支持 systemd 发行版"
    exit 1
  fi
}

ensure_ffmpeg() {
  if command -v ffmpeg >/dev/null 2>&1; then
    return
  fi

  warn "未检测到 ffmpeg，正在尝试安装..."
  if command -v apt-get >/dev/null 2>&1; then
    apt-get update -qq
    apt-get install -y ffmpeg
  elif command -v yum >/dev/null 2>&1; then
    yum install -y ffmpeg || yum install -y epel-release && yum install -y ffmpeg
  elif command -v dnf >/dev/null 2>&1; then
    dnf install -y ffmpeg
  else
    error "无法自动安装 ffmpeg，请先手动安装后重试"
    exit 1
  fi
}

download_url() {
  local asset="${BINARY_PREFIX}-linux-${ARCH}"
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
  info "下载 ${BINARY_PREFIX}-linux-${ARCH} ..."
  if ! curl -fsSL "$url" -o "$dest"; then
    error "下载失败: $url"
    return 1
  fi
  chmod 755 "$dest"
}

install_service() {
  cat >"/etc/systemd/system/${SERVICE_NAME}.service" <<EOF
[Unit]
Description=wtv IPTV monitor service
Wants=network-online.target
After=network-online.target

[Service]
Type=simple
WorkingDirectory=${INSTALL_PATH}
ExecStart=${INSTALL_PATH}/wtv
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

  systemctl daemon-reload
  systemctl enable "${SERVICE_NAME}" >/dev/null 2>&1
}

do_install() {
  if [ -f "${INSTALL_PATH}/wtv" ]; then
    warn "检测到已安装在 ${INSTALL_PATH}，请使用 update 命令更新"
    exit 0
  fi

  mkdir -p "${INSTALL_PATH}"
  download_binary "${INSTALL_PATH}/wtv"
  install_service
  systemctl restart "${SERVICE_NAME}"

  info "安装完成"
  echo
  echo "管理地址: http://YOUR_IP:9876/admin"
  echo "配置文件: ${INSTALL_PATH}/etc/tv.yaml"
  echo "默认账号: admin / admin123（请尽快修改）"
  echo
  echo "常用命令:"
  echo "  systemctl status ${SERVICE_NAME}"
  echo "  systemctl restart ${SERVICE_NAME}"
}

do_update() {
  if [ ! -f "${INSTALL_PATH}/wtv" ]; then
    error "未检测到已安装的 wtv，请先执行 install"
    exit 1
  fi

  systemctl stop "${SERVICE_NAME}" || true
  cp "${INSTALL_PATH}/wtv" /tmp/wtv.bak

  if ! download_binary /tmp/wtv.new; then
    mv /tmp/wtv.bak "${INSTALL_PATH}/wtv"
    systemctl start "${SERVICE_NAME}" || true
    exit 1
  fi

  mv /tmp/wtv.new "${INSTALL_PATH}/wtv"
  rm -f /tmp/wtv.bak
  systemctl start "${SERVICE_NAME}"
  info "已更新到最新版本"
}

do_uninstall() {
  systemctl disable "${SERVICE_NAME}" >/dev/null 2>&1 || true
  systemctl stop "${SERVICE_NAME}" >/dev/null 2>&1 || true
  rm -f "/etc/systemd/system/${SERVICE_NAME}.service" \
        "/lib/systemd/system/${SERVICE_NAME}.service"
  rm -rf "${INSTALL_PATH}"
  systemctl daemon-reload
  info "wtv 已卸载"
}

main() {
  local action="${1:-}"
  resolve_install_path "$@"
  require_root
  require_systemd
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
