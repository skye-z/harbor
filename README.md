# Harbor - 轻量级本地Docker管理平台

[English](README_en.md)

[![CodeQL](https://github.com/skye-z/harbor/workflows/CodeQL/badge.svg)](https://github.com/skye-z/harbor/security/code-scanning)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=skye-z_harbor&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=skye-z_harbor)

[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=skye-z_harbor&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=skye-z_harbor)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=skye-z_harbor&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=skye-z_harbor)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=skye-z_harbor&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=skye-z_harbor)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=skye-z_harbor&metric=bugs)](https://sonarcloud.io/summary/new_code?id=skye-z_harbor)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=skye-z_harbor&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=skye-z_harbor)

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

* 体积小: 程序本体小于 30 MB
* 部署快: 提供一键脚本快速部署
* 开销低: 内存占用仅 ≈ 30 MB

## 安装

请复制下方命令到服务器终端中执行, 脚本提供了 Harbor 的安装、卸载于开启自启设置服务

```shell
bash -c "$(curl -fsSL https://skye-z.github.io/harbor/install.sh)"
```

安装完成后将会自动运行, 首次运行 Harbor 会自动生成管理账户, 密码将输出在日志中, 如下所示:

```log
2024/01/01 00:00:00 [Tips] create default admin, password is ..........
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

## 构建

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o harbor -ldflags '-s -w'
```
