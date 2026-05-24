1. 本项目仓库是我的名为《基于Go-zero和WebSocket的IM即时通讯系统》的毕业设计及论文。paper目录下是学校发的论文的模版以及一些要求。
2. 本项目的所有组件服务（MySQL、redis等）、前后端项目都在另一台 windows 个人电脑中安装的ubuntu24虚拟机上启动。我会在虚拟机上也把本仓库拉取下来，然后启动docker
   compose以及前后端服务。
3. ubuntu24虚拟机中已经用了clash代理，可以访问到外网
   4`go-im/infra.host.env` 中的 `VM_IP` 只用于 docker compose 里对外公布的地址（如 etcd advertise、kafka advertised
   listeners）以及 Windows 宿主机中的外部工具访问；`go-im/apps/**/etc/*.yaml` 中的后端运行时依赖地址统一使用 `127.0.0.1`
   ，因为所有 Go 服务和基础设施容器都运行在同一台 Ubuntu 虚拟机内。需要用 Navicat 等工具从 Windows 宿主机访问数据库时，连接
   `VM_IP:13306`（MySQL）或 `VM_IP:47017`（MongoDB）。

## Claude Code 权限配置

### 允许执行的命令

- **目录操作**: cd, ls, pwd, mkdir, rm, cp, mv, find
- **文件操作**: cat, head, tail, less, grep, wc, sort, uniq, diff
- **Go 相关**: go build, go run, go test, go mod, go vet, go fmt
- **网络工具**: curl, wget, ping, netstat, ss
- **Git 命令** (除 commit 和 push 外):
    - ✅ git status, git log, git diff, git branch, git checkout, git pull, git fetch, git merge, git stash, git tag, git
      remote, git show, git blame
    - ❌ git commit, git push (需要用户手动执行)
- **其他常用**: make, chmod, chown, tar, zip, unzip, env, echo, date

### 禁止的操作

- **Git 提交**: git commit, git push (由用户自行决定提交时机)
- **危险命令**: sudo, rm -rf /, mkfs, fdisk 等系统级破坏命令
- **包管理器安装**: apt install, yum install 等需要确认的操作（需用户授权）

# 功能列表

- 登录，注册、修改密码
- 修改个人信息
- 申请加别人好友、同意/拒绝别人加自己好友、查询好友列表（显示好友的昵称）、查看好友的在线状态
- 发送信息给好友、查询和好友的聊天记录最近的聊天记录
  表结构文件:./go-im/deploy/sql/init.sql