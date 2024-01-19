# Harbor - Docker Management Platform

[中文](README.md)

> Under development, not available at this time

Harbor is a lightweight, local Docker management platform that accelerates the delivery of single-machine services and simplifies subsequent operations.

## Function

* Engine: View engine information, listen to engine daemon events.
* Container: container start/stop control, log view, web terminal, resource statistics, port shortcut.
* Mirror: build information view, pull and delete mirrors
* Network: network creation and editing, container access
* Storage volume: create, edit and delete
* Monitoring: engine offline, container start/stop, burst load, resource exhaustion monitoring.
* Users: account password login, OAuth2 login.

## Features

* Small size: program body less than 30 MB
* Fast Deployment: Provides one-click scripts for quick deployment.
* Low overhead: ≈ 20 MB memory footprint

## Installation

Please copy the following commands to your server terminal, the script provides installation, uninstallation and startup settings for Harbor.

```shell
bash -c "$(curl -fsSL https://skye-z.github.io/harbor/install.sh)"
```

After the installation is complete, Harbor will run automatically. The first time you run Harbor, it will automatically generate an administrative account, and the password will be output in the log, as shown below:

```log
2024/01/01 00:00:00 [Tips] create default admin, password is ..........
```

## Control

```shell
# Start Harbor
systemctl start harbor
# Stop Harbor
systemctl stop harbor
# View Harbor status and logs
systemctl status harbor
```

## Build

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o harbor -ldflags '-s -w'
```
