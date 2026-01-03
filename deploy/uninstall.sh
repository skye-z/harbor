#!/bin/bash
#
# Harbor 卸载脚本
#

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

APP_NAME="harbor"
INSTALL_DIR="/opt/harbor"
SYSTEMD_DIR="/etc/systemd/system"

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

confirm_action() {
    read -p "$1 [y/N]: " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "操作已取消"
        exit 0
    fi
}

stop_services() {
    log_info "停止服务..."

    systemctl stop harbor.service 2>/dev/null || true
    systemctl stop harbor-backup.timer 2>/dev/null || true

    systemctl disable harbor.service 2>/dev/null || true
    systemctl disable harbor-backup.timer 2>/dev/null || true

    systemctl daemon-reload
    log_info "服务已停止并禁用"
}

remove_files() {
    log_info "移除应用文件..."
    rm -rf "$INSTALL_DIR"
    log_info "应用文件已移除"
}

remove_systemd_files() {
    log_info "移除 systemd 服务文件..."

    rm -f "$SYSTEMD_DIR/harbor.service"
    rm -f "$SYSTEMD_DIR/harbor-backup.timer"
    rm -f "$SYSTEMD_DIR/harbor-backup.service"

    systemctl daemon-reload
    log_info "systemd 服务文件已移除"
}

remove_logrotate() {
    log_info "移除日志轮转配置..."
    rm -f /etc/logrotate.d/harbor
}

remove_user() {
    if id "$APP_NAME" &>/dev/null; then
        log_info "移除系统用户: $APP_NAME"
        userdel "$APP_NAME" 2>/dev/null || true
        log_info "用户已移除"
    fi
}

main() {
    echo "============================================"
    echo -e "${RED}Harbor 卸载${NC}"
    echo "============================================"
    echo ""

    if [ "$EUID" -ne 0 ]; then
        log_error "请使用 root 权限运行"
        exit 1
    fi

    if [ ! -d "$INSTALL_DIR" ]; then
        log_info "Harbor 未安装或已移除"
        exit 0
    fi

    confirm_action "确认卸载 Harbor?"

    stop_services
    remove_systemd_files
    remove_logrotate
    remove_files
    remove_user

    echo ""
    echo "============================================"
    echo -e "${GREEN}Harbor 卸载完成${NC}"
    echo "============================================"
}

main "$@"
