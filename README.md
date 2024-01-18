# Harbor - 轻量级本地Docker管理平台

[English](README_en.md)

> 正在开发中，暂时无法使用

Harbor 是一个轻量级的本地 Docker 管理平台, 加快了单机服务的交付速度, 简化后续的运维工作.

## 功能

* 引擎: 查看引擎信息、监听引擎守护线程事件
* 容器: 容器启停控制、日志查看、Web终端、资源统计、端口快捷跳转
* 镜像: 构建信息查看、镜像拉取与删除
* 网络: 网络创建编辑、容器接入
* 存储卷: 创建编辑与删除
* 监控: 引擎下线、容器启停、突发负载、资源耗尽等监控
* 用户: 账号密码登陆、OAuth2登陆

## 特性

* 体积小: 程序本体小于30 MB
* 部署快: 提供一键脚本快速部署
* 开销低: 内存占用仅 ≈ 20 MB

## 安装

请复制下方命令到服务器终端中执行, 脚本提供了 Harbor 的安装、卸载于开启自启设置服务

```shell
bash -c "$(curl -fsSL https://raw.githubusercontent.com/skye-z/harbor/main/install.sh)"
```

## 控制

```shell
# 启动 Harbor
systemctl start harbor
# 停止 Harbor
systemctl stop harbor
# 查看 Harbor 状态与日志
systemctl status harbor
```