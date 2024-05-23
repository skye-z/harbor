#!/bin/bash

# 检查是否以root用户执行
if [ "$EUID" -ne 0 ]; then
    echo "Error: This script must be run as root."
    exit 1
fi

# 检查是否安装了 Docker
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker and try again."
    exit 1
fi

# 设置服务名称和描述
SERVICE_NAME="harbor"
SERVICE_DESCRIPTION="Harbor"

# 获取系统架构
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    HARBOR_FILE="harbor_linux_amd64"
elif [ "$ARCH" = "aarch64" ]; then
    HARBOR_FILE="harbor_linux_arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

# 设置可执行文件的下载 URL 和版本
HARBOR_VERSION=$(curl -sL https://api.github.com/repos/skye-z/harbor/releases/latest | grep '"tag_name":' | cut -d'"' -f4)
if [ -z "$HARBOR_VERSION" ]; then
    echo "Failed to retrieve the latest Harbor version."
    exit 1
fi
HARBOR_DOWNLOAD_URL="https://github.com/skye-z/harbor/releases/download/${HARBOR_VERSION}/${HARBOR_FILE}"

# 设置工作目录
WORKING_DIRECTORY="/opt/harbor"

# 创建 Systemd 服务单元文件
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
cat <<EOF > $SERVICE_FILE
[Unit]
Description=${SERVICE_DESCRIPTION}
After=docker.service

[Service]
ExecStart=${WORKING_DIRECTORY}/harbor >> /var/log/harbor.log 2>&1
WorkingDirectory=${WORKING_DIRECTORY}
Restart=always

[Install]
WantedBy=multi-user.target
EOF

install_harbor_online() {
    # 创建工作目录
    mkdir -p $WORKING_DIRECTORY

    # 下载 Harbor 可执行文件
    curl -LJ $HARBOR_DOWNLOAD_URL -o ${WORKING_DIRECTORY}/harbor

    # 赋予可执行权限
    chmod +x ${WORKING_DIRECTORY}/harbor

    # 重载 Systemd 配置
    systemctl daemon-reload

    # 启动服务
    systemctl start $SERVICE_NAME

    echo "Harbor service installed successfully!"
}

install_harbor_offline() {
    # 检查是否存在离线安装文件
    if [ -f $HARBOR_FILE ]; then
        # 创建工作目录
        mkdir -p $WORKING_DIRECTORY

        # 复制离线安装文件到工作目录
        cp ${HARBOR_FILE} ${WORKING_DIRECTORY}/harbor

        # 赋予可执行权限
        chmod +x ${WORKING_DIRECTORY}/harbor

        # 重载 Systemd 配置
        systemctl daemon-reload

        # 启动服务
        systemctl start $SERVICE_NAME

        echo "Harbor service installed successfully!"
    else
        echo "Error: Offline installation file '$HARBOR_FILE' not found. Please download it manually to the current directory."
        exit 1
    fi
}

uninstall_harbor() {
    # 停止服务
    systemctl stop $SERVICE_NAME

    # 禁用开机自启
    systemctl disable $SERVICE_NAME

    # 删除工作目录
    rm -rf $WORKING_DIRECTORY

    # 删除 Systemd 服务文件
    rm -f $SERVICE_FILE

    # 重载 Systemd 配置
    systemctl daemon-reload

    echo "Harbor service uninstalled successfully!"
}

enable_autostart() {
    # 启用开机自启
    systemctl enable $SERVICE_NAME

    echo "Autostart enabled for Harbor service!"
}

disable_autostart() {
    # 禁用开机自启
    systemctl disable $SERVICE_NAME

    echo "Autostart disabled for Harbor service!"
}

# 显示选项
echo "Select an option:"
echo "1. Install Harbor (Online)"
echo "2. Install Harbor (Offline)"
echo "3. Uninstall Harbor"
echo "4. Enable Autostart"
echo "5. Disable Autostart"
read -p "Enter option number: " option

# 根据选项执行相应操作
case $option in
    1) install_harbor_online ;;
    2) install_harbor_offline ;;
    3) uninstall_harbor ;;
    4) enable_autostart ;;
    5) disable_autostart ;;
    *) echo "Invalid option" ;;
esac
