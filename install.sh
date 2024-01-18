#!/bin/bash

# 设置服务名称和描述
SERVICE_NAME="harbor"
SERVICE_DESCRIPTION="Harbor"

# 设置可执行文件的下载 URL 和版本
HARBOR_VERSION=$(curl -sL https://api.github.com/repos/skye-z/harbor/releases/latest | grep '"tag_name":' | cut -d'"' -f4)
if [ -z "$HARBOR_VERSION" ]; then
    echo "Failed to retrieve the latest Harbor version."
    exit 1
fi
HARBOR_DOWNLOAD_URL="https://github.com/skye-z/harbor/releases/download/${HARBOR_VERSION}/harbor-linux-amd64"

# 设置工作目录
WORKING_DIRECTORY="/opt/harbor"

# 创建 Systemd 服务单元文件
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"
echo "[Unit]" > $SERVICE_FILE
echo "Description=${SERVICE_DESCRIPTION}" >> $SERVICE_FILE
echo "After=network.target" >> $SERVICE_FILE
echo "" >> $SERVICE_FILE
echo "[Service]" >> $SERVICE_FILE
echo "ExecStart=${WORKING_DIRECTORY}/harbor >> /var/log/harbor.log 2>&1" >> $SERVICE_FILE
echo "WorkingDirectory=${WORKING_DIRECTORY}" >> $SERVICE_FILE
echo "Restart=always" >> $SERVICE_FILE
echo "" >> $SERVICE_FILE
echo "[Install]" >> $SERVICE_FILE
echo "WantedBy=multi-user.target" >> $SERVICE_FILE

install_harbor() {
    # 创建工作目录
    sudo mkdir -p $WORKING_DIRECTORY

    # 下载 Harbor 可执行文件
    curl -LJ $HARBOR_DOWNLOAD_URL -o ${WORKING_DIRECTORY}/harbor

    # 赋予可执行权限
    chmod +x ${WORKING_DIRECTORY}/harbor

    # 重载 Systemd 配置
    sudo systemctl daemon-reload

    # 启动服务
    sudo systemctl start $SERVICE_NAME

    echo "Harbor service installed successfully!"
}

uninstall_harbor() {
    # 停止服务
    sudo systemctl stop $SERVICE_NAME

    # 禁用开机自启
    sudo systemctl disable $SERVICE_NAME

    # 删除工作目录
    sudo rm -rf $WORKING_DIRECTORY

    # 删除 Systemd 服务文件
    sudo rm -f $SERVICE_FILE

    # 重载 Systemd 配置
    sudo systemctl daemon-reload

    echo "Harbor service uninstalled successfully!"
}

enable_autostart() {
    # 启用开机自启
    sudo systemctl enable $SERVICE_NAME

    echo "Autostart enabled for Harbor service!"
}

disable_autostart() {
    # 禁用开机自启
    sudo systemctl disable $SERVICE_NAME

    echo "Autostart disabled for Harbor service!"
}

# 显示选项
echo "Select an option:"
echo "1. Install Harbor"
echo "2. Uninstall Harbor"
echo "3. Enable Autostart"
echo "4. Disable Autostart"
read -p "Enter option number: " option

# 根据选项执行相应操作
case $option in
    1) install_harbor;;
    2) uninstall_harbor;;
    3) enable_autostart;;
    4) disable_autostart;;
    *) echo "Invalid option";;
esac