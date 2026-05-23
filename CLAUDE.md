1. 本项目仓库是我的名为《基于Go-zero和WebSocket的IM即时通讯系统》的毕业设计及论文。paper目录下是学校发的论文的模版以及一些要求。
2. 本项目的所有组件服务（MySQL、redis等）、前后端项目都在另一台windows个人电脑中安装的ubuntu24虚拟机上启动。我会在虚拟机上也把本仓库拉取下来，然后启动docker compose以及前后端服务。
3. `go-im/infra.host.env` 中的 `VM_IP` 只用于 docker compose 里对外公布的地址（如 etcd advertise、kafka advertised listeners）以及 Windows 宿主机中的外部工具访问；`go-im/apps/**/etc/*.yaml` 中的后端运行时依赖地址统一使用 `127.0.0.1`，因为所有 Go 服务和基础设施容器都运行在同一台 Ubuntu 虚拟机内。需要用 Navicat 等工具从 Windows 宿主机访问数据库时，连接 `VM_IP:13306`（MySQL）或 `VM_IP:47017`（MongoDB）。
