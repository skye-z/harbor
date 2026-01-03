#!/bin/bash
#
# Harbor 更新脚本 - 零停机更新
#

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

APP_NAME="harbor"
INSTALL_DIR="/opt/harbor"
HARBOR_FILE="$INSTALL_DIR/harbor"
GITHUB_REPO="https://github.com/skye-z/harbor/releases/download"
GH_PROXY="https://gh-proxy.com"

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

get_latest_version() {
    local version
    version=$(curl -s "https://api.github.com/repos/skye-z/harbor/releases/latest" 2>/dev/null | grep '"'"'tag_name'"'"' | sed 's/.*"\([^"]*\)".*/\1/')
    echo "$version"
}

download_harbor() {
    local version
    version=$(get_latest_version)
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local current_binary="$script_dir/harbor"

    # 优先使用当前目录
    if [ -f "$current_binary" ]; then
        log_info "使用当前目录的 harbor 文件"
        cp "$current_binary" "./harbor"
        chmod +x "./harbor"
        return 0
    fi

    # 从 GitHub 下载
    if [ -n "$version" ]; then
        log_info "正在下载 Harbor v$version ..."
        local github_url="$GITHUB_REPO/$version/harbor-linux-amd64"

        if curl -L -o "harbor" "$github_url" 2>/dev/null; then
            chmod +x "./harbor"
            log_info "下载成功"
            return 0
        fi

        log_warn "GitHub 下载失败，尝试代理..."
        local proxy_url="$GH_PROXY/$github_url"
        if curl -L -o "harbor" "$proxy_url" 2>/dev/null; then
            chmod +x "./harbor"
            log_info "代理下载成功"
            return 0
        fi
    fi

    log_error "无法下载 Harbor"
    return 1
}

check_prerequisites() {
    if [ "$EUID" -ne 0 ]; then
        log_error "请使用 root 权限运行"
        exit 1
    fi

    if [ ! -f "$HARBOR_FILE" ]; then
        log_error "Harbor 未安装: $HARBOR_FILE"
        exit 1
    fi

    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    if [ ! -f "$script_dir/harbor" ]; then
        download_harbor
    fi
}

create_backup() {
    log_info "备份旧版本..."

    local backup_dir="$INSTALL_DIR/backups/update_$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$backup_dir"

    cp "$HARBOR_FILE" "$backup_dir/harbor.bak"
    cp -r "$INSTALL_DIR/dist" "$backup_dir/dist.bak" 2>/dev/null || true

    log_info "备份已保存到: $backup_dir"
}

stop_service() {
    log_info "停止 Harbor 服务..."

    systemctl stop "$APP_NAME.service"
    sleep 2

    if systemctl is-active --quiet "$APP_NAME.service"; then
        log_error "停止服务失败"
        exit 1
    fi

    log_info "服务已停止"
}

update_binary() {
    log_info "更新 Harbor 二进制文件..."

    cp "./harbor" "$HARBOR_FILE"
    chmod +x "$HARBOR_FILE"

    log_info "二进制文件已更新"
}

sync_files() {
    log_info "同步前端文件..."

    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    if [ -d "$script_dir/dist" ]; then
        rm -rf "$INSTALL_DIR/dist"
        cp -r "$script_dir/dist" "$INSTALL_DIR/"
    fi

    # 保留现有配置
    if [ -f "$INSTALL_DIR/config.yaml" ]; then
        log_info "保留现有配置文件"
    fi

    chown -R "$APP_USER:$APP_USER" "$INSTALL_DIR"
    log_info "文件同步完成"
}

start_service() {
    log_info "启动 Harbor 服务..."

    systemctl start "$APP_NAME.service"
    sleep 3

    if systemctl is-active --quiet "$APP_NAME.service"; then
        log_info "服务启动成功"
    else
        log_error "服务启动失败，尝试回滚..."
        rollback
        exit 1
    fi
}

get_port() {
    local port
    port=$(grep -A2 "^server:" "$INSTALL_DIR/config.yaml" 2>/dev/null | grep "port:" | awk '{print $2}' | tr -d ' ')
    echo "${port:-8080}"
}

verify_update() {
    log_info "验证更新..."

    local port
    port=$(get_port)
    if curl -s -o /dev/null -w "%{http_code}" "http://localhost:$port/api/docker/info" 2>/dev/null | grep -q "200\|401"; then
        log_info "更新验证成功 - 服务正常运行"
    else
        log_warn "服务可能需要更多时间启动"
    fi
}

rollback() {
    log_info "回滚中..."

    local backup_dir="$INSTALL_DIR/backups/update_"*
    if [ -d "$backup_dir" ]; then
        local latest
        latest=$(ls -td "$INSTALL_DIR/backups"/update_*/ 2>/dev/null | head -1)
        if [ -n "$latest" ] && [ -f "$latest/harbor.bak" ]; then
            cp "$latest/harbor.bak" "$HARBOR_FILE"
            cp -r "$latest/dist.bak" "$INSTALL_DIR/dist" 2>/dev/null || true
            chown -R "$APP_NAME:$APP_NAME" "$INSTALL_DIR"
            systemctl restart "$APP_NAME.service"
            log_info "回滚完成"
            return
        fi
    fi
    log_error "未找到备份，无法回滚"
}

main() {
    echo "============================================"
    echo "Harbor Docker 管理面板更新"
    echo "============================================"
    echo ""

    check_prerequisites
    create_backup
    stop_service
    update_binary
    sync_files
    start_service
    verify_update

    echo ""
    echo "============================================"
    echo -e "${GREEN}更新完成${NC}"
    echo "============================================"
    echo ""
    log_info "日志: journalctl -u harbor -f"
}

main "$@"
