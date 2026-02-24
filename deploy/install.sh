#!/bin/bash
#
# Harbor Docker 管理面板 - 一键安装脚本
# 支持系统: Ubuntu/Debian, RHEL/CentOS, Fedora
#

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 配置
APP_NAME="harbor"
APP_USER="harbor"
INSTALL_DIR="/opt/harbor"
HARBOR_FILE="$INSTALL_DIR/harbor"
SYSTEMD_DIR="/etc/systemd/system"
LOG_DIR="/var/log/harbor"
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

# ==================== systemd 单元文件 ====================

HARBOR_SERVICE='[Unit]
Description=Harbor Docker Management Panel
Documentation=https://github.com/skye-z/harbor
After=network.target docker.service
Wants=docker.service

[Service]
Type=simple
User=harbor
Group=harbor
WorkingDirectory=/opt/harbor
ExecStart=/opt/harbor/harbor
ExecReload=/bin/kill -HUP $MAINPID
KillMode=mixed
TimeoutStopSec=5
PrivateTmp=true
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/harbor /var/run/docker.sock
Restart=on-failure
RestartSec=5

StandardOutput=journal
StandardError=journal
SyslogIdentifier=harbor

LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
'

HARBOR_BACKUP_TIMER='[Unit]
Description=Harbor Database Backup Timer
Requires=harbor-backup.service

[Timer]
OnCalendar=daily
Persistent=true
RandomizedDelaySec=3600

[Install]
WantedBy=timers.target
'

HARBOR_BACKUP_SERVICE='[Unit]
Description=Harbor Database Backup Service
After=network.target

[Service]
Type=oneshot
User=harbor
Group=harbor
WorkingDirectory=/opt/harbor
ExecStart=/opt/harbor/scripts/backup.sh

StandardOutput=journal
StandardError=journal
SyslogIdentifier=harbor-backup

[Install]
WantedBy=multi-user.target
'

# ==================== 备份脚本 ====================

BACKUP_SCRIPT='#!/bin/bash
#
# Harbor 备份脚本
#

set -e

APP_DIR="/opt/harbor"
BACKUP_DIR="$APP_DIR/backups"
CONFIG_FILE="$APP_DIR/config.yaml"
LOG_DIR="/var/log/harbor"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/harbor_backup_$TIMESTAMP"

mkdir -p "$BACKUP_DIR"
chmod 700 "$BACKUP_DIR"

log() {
    echo "[$(date '"'"'+%Y-%m-%d %H:%M:%S'"'"')] $1" >> "$LOG_DIR/backup.log"
}

log_info() {
    log "INFO: $1"
    echo "INFO: $1"
}

log_error() {
    log "ERROR: $1"
    echo "ERROR: $1" >&2
}

backup_database() {
    log_info "开始数据库备份..."

    if [ -f "$APP_DIR/harbor.db" ]; then
        sqlite3 "$APP_DIR/harbor.db" ".backup" "$BACKUP_FILE.db"
        log_info "数据库备份完成: $BACKUP_FILE.db"
    else
        log_error "数据库文件未找到: $APP_DIR/harbor.db"
        return 1
    fi
}

backup_config() {
    log_info "开始配置备份..."

    tar -czf "$BACKUP_FILE.conf.tar.gz" -C "$(dirname "$CONFIG_FILE")" "$(basename "$CONFIG_FILE")"
    log_info "配置备份完成: $BACKUP_FILE.conf.tar.gz"
}

cleanup_old() {
    log_info "清理 $BACKUP_DIR 中 30 天前的备份..."
    find "$BACKUP_DIR" -name "harbor_backup_*" -type f -mtime +30 -delete
    log_info "旧备份清理完成"
}

verify_backup() {
    log_info "验证备份完整性..."

    if [ -f "$BACKUP_FILE.db" ]; then
        if ! sqlite3 "$BACKUP_FILE.db" "PRAGMA integrity_check;" | grep -q "ok"; then
            log_error "数据库完整性检查失败"
            rm -f "$BACKUP_FILE.db"
            return 1
        fi
        log_info "备份完整性检查通过"
    fi
}

main() {
    log_info "========================================"
    log_info "Harbor 备份开始"
    log_info "========================================"

    backup_database
    backup_config
    cleanup_old
    verify_backup

    log_info "========================================"
    log_info "Harbor 备份完成"
    log_info "========================================"

    echo ""
    echo "备份完成!"
    ls -lh "$BACKUP_FILE"*
}

main "$@"
'

# ==================== 日志轮转配置 ====================

LOGROTATE_CONFIG='
/var/log/harbor/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 harbor harbor
    postrotate
        systemctl reload harbor > /dev/null 2>&1 || true
    endscript
}
'

# ==================== 安装函数 ====================

detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
    else
        OS="unknown"
    fi
}

download_harbor() {
    local version
    version=$(curl -s "https://api.github.com/repos/skye-z/harbor/releases/latest" 2>/dev/null | grep '"'"'tag_name'"'"' | sed 's/.*"\([^"]*\)".*/\1/')

    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local current_binary="$script_dir/harbor"

    # 优先使用当前目录的 harbor 文件
    if [ -f "$current_binary" ]; then
        log_info "使用当前目录的 harbor 文件"
        chmod +x "$current_binary"
        return 0
    fi

    # 从 GitHub 下载
    if [ -n "$version" ]; then
        log_info "正在下载 Harbor v$version ..."
        local github_url="$GITHUB_REPO/$version/harbor-linux-amd64"

        if curl -L -o "$current_binary" "$github_url" 2>/dev/null; then
            chmod +x "$current_binary"
            log_info "下载成功"
            return 0
        fi

        # 失败则使用 gh-proxy
        log_warn "GitHub 下载失败，尝试代理..."
        local proxy_url="$GH_PROXY/$github_url"
        if curl -L -o "$current_binary" "$proxy_url" 2>/dev/null; then
            chmod +x "$current_binary"
            log_info "代理下载成功"
            return 0
        fi
    fi

    log_error "无法下载 Harbor，请检查网络或手动下载"
    return 1
}

create_user() {
    if id "$APP_USER" &>/dev/null; then
        log_info "用户 $APP_USER 已存在"
    else
        useradd -r -s /sbin/nologin -d "$INSTALL_DIR" -M "$APP_USER" 2>/dev/null || \
        useradd -r -s /usr/sbin/nologin -d "$INSTALL_DIR" -M "$APP_USER" 2>/dev/null || true
        log_info "用户 $APP_USER 创建成功"
    fi
    
    # 将 harbor 用户加入 docker 组（解决 socket 权限问题）
    if getent group docker >/dev/null 2>&1; then
        usermod -aG docker "$APP_USER" 2>/dev/null || true
        log_info "已添加 $APP_USER 到 docker 组"
    fi
}

create_directories() {
    mkdir -p "$INSTALL_DIR"
    mkdir -p "$LOG_DIR"
    mkdir -p "$INSTALL_DIR/scripts"
    chown -R "$APP_USER:$APP_USER" "$INSTALL_DIR"
    chmod 755 "$INSTALL_DIR"
}

install_systemd() {
    log_info "安装 systemd 服务..."

    echo "$HARBOR_SERVICE" > "$SYSTEMD_DIR/harbor.service"
    echo "$HARBOR_BACKUP_TIMER" > "$SYSTEMD_DIR/harbor-backup.timer"
    echo "$HARBOR_BACKUP_SERVICE" > "$SYSTEMD_DIR/harbor-backup.service"

    chmod 644 "$SYSTEMD_DIR/harbor.service"
    chmod 644 "$SYSTEMD_DIR/harbor-backup.timer"
    chmod 644 "$SYSTEMD_DIR/harbor-backup.service"

    systemctl daemon-reload
    log_info "systemd 服务安装完成"
}

install_backup_script() {
    log_info "安装备份脚本..."

    mkdir -p "$INSTALL_DIR/scripts"
    echo "$BACKUP_SCRIPT" > "$INSTALL_DIR/scripts/backup.sh"
    chmod +x "$INSTALL_DIR/scripts/backup.sh"
    chown -R "$APP_USER:$APP_USER" "$INSTALL_DIR/scripts"
}

install_logrotate() {
    log_info "配置日志轮转..."
    echo "$LOGROTATE_CONFIG" > /etc/logrotate.d/harbor
    chmod 644 /etc/logrotate.d/harbor
}

fix_docker_socket() {
    if [ -S /var/run/docker.sock ]; then
        # 确保 docker 组对 socket 有读写权限
        if getent group docker >/dev/null 2>&1; then
            chmod 660 /var/run/docker.sock 2>/dev/null || true
            chown root:docker /var/run/docker.sock 2>/dev/null || true
            log_info "已修复 Docker socket 权限"
        fi
    fi
}

install_files() {
    log_info "复制 Harbor 文件..."

    # 停止服务
    if systemctl is-active --quiet "$APP_NAME" 2>/dev/null; then
        log_info "停止 Harbor 服务..."
        systemctl stop "$APP_NAME" 2>/dev/null || true
    fi

    # 获取源码目录
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local current_binary="$script_dir/harbor"

    # 复制二进制文件
    cp "$current_binary" "$HARBOR_FILE"
    chmod +x "$HARBOR_FILE"

    # 复制前端文件
    if [ -d "$script_dir/dist" ]; then
        rm -rf "$INSTALL_DIR/dist" 2>/dev/null || true
        cp -r "$script_dir/dist" "$INSTALL_DIR/"
    fi

    chown -R "$APP_USER:$APP_USER" "$INSTALL_DIR"
    log_info "文件复制完成"
}

start_service() {
    log_info "启动 Harbor 服务..."
    
    # 修复 Docker socket 权限
    fix_docker_socket

    systemctl enable "$APP_NAME" 2>/dev/null || true
    systemctl enable "$APP_NAME-backup.timer" 2>/dev/null || true

    systemctl start "$APP_NAME"
    log_info "Harbor 服务已启动"
}

verify_service() {
    sleep 2

    if systemctl is-active --quiet "$APP_NAME"; then
        log_info "Harbor 服务运行正常"

        local port
        port=$(grep -A2 "^server:" "$INSTALL_DIR/config.yaml" 2>/dev/null | grep "port:" | awk '{print $2}' | tr -d ' ')
        port=${port:-8080}
        log_info "访问地址: http://localhost:$port"
    else
        log_warn "Harbor 服务可能未正常运行"
        systemctl status "$APP_NAME" --no-pager || true
    fi
}

# ==================== 主逻辑 ====================

upgrade_mode() {
    echo ""
    echo -e "${YELLOW}检测到 Harbor 已安装，进入升级模式${NC}"
    echo "============================================"
    echo ""

    download_harbor

    create_user
    create_directories

    # 备份旧版本
    local backup_dir="$INSTALL_DIR/backups/upgrade_$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$backup_dir"
    cp "$HARBOR_FILE" "$backup_dir/harbor.bak"
    log_info "已备份旧版本到: $backup_dir"

    install_files
    install_backup_script

    # systemd 已存在则跳过
    if [ -f "$SYSTEMD_DIR/$APP_NAME.service" ]; then
        log_info "Harbor 已注册 systemd，跳过注册"
    else
        install_systemd
    fi

    start_service
    verify_service

    echo ""
    echo "============================================"
    echo -e "${GREEN}升级完成${NC}"
    echo "============================================"
}

install_mode() {
    echo ""
    echo -e "${GREEN}未检测到 Harbor，进入安装模式${NC}"
    echo "============================================"
    echo ""

    download_harbor

    create_user
    create_directories
    install_files
    install_backup_script
    install_systemd
    install_logrotate
    start_service
    verify_service

    echo ""
    echo "============================================"
    echo -e "${GREEN}安装完成${NC}"
    echo "============================================"
    echo ""
    echo "默认账号:"
    echo "  用户名: admin"
    echo "  密码:   HarborAdmin2026!"
    echo ""
    echo "请在首次登录后修改默认密码!"
    echo ""
}

main() {
    echo "============================================"
    echo "Harbor Docker 管理面板"
    echo "============================================"
    echo ""

    detect_os

    # 检测安装状态
    if [ -d "$INSTALL_DIR" ] && [ -f "$HARBOR_FILE" ]; then
        upgrade_mode
    else
        install_mode
    fi
}

main "$@"
