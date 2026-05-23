# go-im

`easy-im` 项目的后端，基于 [go-zero](https://github.com/zeromicro/go-zero) 构建的 IM 微服务集群。提供用户、社交关系、消息收发、长连接网关与异步任务消费等能力。

## 技术栈

- **语言**：Go 1.23
- **微服务框架**：go-zero（HTTP API + gRPC）
- **长连接**：`gorilla/websocket`
- **存储**：MySQL（关系数据）、MongoDB（聊天消息）、Redis（缓存 / 在线状态）
- **消息队列**：Kafka（task-mq 异步消费）
- **服务发现**：etcd
- **配置中心**：sail（`HYY-yu/sail-client`）
- **网关**：APISIX（可选，组件在 `components/apisix/`）
- **鉴权**：JWT（`golang-jwt/jwt/v4`）
- **ID 生成**：Snowflake

## 模块划分

```
go-im/
├── apps/
│   ├── user/             # 用户中心：注册、登录、用户信息
│   │   ├── api/          # 对外 HTTP 接口（user.api）
│   │   └── rpc/          # 内部 gRPC 服务（user.proto）
│   ├── social/           # 社交关系：好友、群组
│   │   ├── api/
│   │   └── rpc/
│   ├── im/               # 即时通讯核心
│   │   ├── api/          # 会话、消息查询等 HTTP 接口
│   │   ├── rpc/          # 消息存储/查询 gRPC
│   │   └── ws/           # WebSocket 长连接网关（路由：user.online / conversation.chat / conversation.markChat / push）
│   └── task/
│       └── mq/           # Kafka 消费者，异步持久化与消息分发
├── pkg/                  # 通用基础库
│   ├── bitmap/           # 位图（已读状态等）
│   ├── configserver/     # 配置中心客户端封装
│   ├── ctxdata/          # context 透传工具
│   ├── encrypy/          # 加密相关
│   ├── interceptor/      # gRPC 拦截器
│   ├── middleware/       # HTTP 中间件
│   ├── resultx/          # 统一响应封装
│   ├── retry/            # 重试
│   ├── status/           # 错误码定义
│   ├── suid/             # ID 工具
│   ├── utils/            # 通用工具
│   └── xerr/             # 错误处理
├── components/           # 中间件相关配置/资源
│   ├── apisix/           # APISIX 网关
│   ├── apisix-dashboard/
│   ├── kibana/           # 日志可视化
│   └── sail/             # 配置中心
├── deploy/               # 部署相关
│   ├── dockerfile/       # 各服务 Dockerfile
│   ├── mk/               # 每个服务对应的 Makefile 片段
│   └── script/           # 一键发布脚本
├── test/                 # 测试代码
├── docker-compose.yaml   # 本地中间件编排
├── Makefile              # 全局构建入口
├── go.mod / go.sum
└── README.md
```

### 服务清单

| 服务 | 形态 | 说明 |
|------|------|------|
| `user-api` | HTTP | 用户对外接口（注册/登录/资料） |
| `user-rpc` | gRPC | 用户领域内部服务 |
| `social-api` | HTTP | 好友/群组对外接口 |
| `social-rpc` | gRPC | 社交领域内部服务 |
| `im-api` | HTTP | 会话与消息查询接口 |
| `im-rpc` | gRPC | 消息存储/读取等内部服务 |
| `im-ws` | WebSocket | 客户端长连接网关，承接四个 IM 路由 |
| `task-mq` | Consumer | Kafka 消费者，异步落库与分发 |

## WebSocket 协议

连接地址：`ws://{im-ws-host}/ws?userId={userId}`

帧类型与路由定义见 `apps/im/ws/internal/handler/router.go`：

| Method | 处理器 | 作用 |
|--------|--------|------|
| `user.online` | `user.OnLine` | 拉取在线用户列表 |
| `conversation.chat` | `conversation.Chat` | 发送聊天消息（单聊/群聊） |
| `conversation.markChat` | `conversation.MarkRead` | 标记消息已读 |
| `push` | `push.Push` | 服务端主动下发 |

详细 DTO 与帧定义见仓库根目录 [`README.md`](../README.md) 与 [`easy-im-web/CLAUDE.md`](../easy-im-web/CLAUDE.md)。

## 环境依赖

通过 `docker-compose.yaml` 一键启动本地中间件：

- etcd（端口 3379）：服务注册发现
- MySQL 8.0（端口 13306，密码 `easy-chat`）
- Redis（端口 16379，密码 `easy-chat`）
- MongoDB 4.0
- 其余 Kafka / APISIX 等见 `docker-compose.yaml`

```bash
docker compose up -d
```

## 构建与运行

### 一键构建发布（开发环境）

```bash
make release-test
```

该目标会依次执行 `deploy/mk/` 下各服务的发布脚本（`user-rpc / user-api / social-rpc / social-api / im-rpc / im-api / im-ws / task-mq`）。

### 单独构建某个服务

```bash
make user-rpc-dev
make user-api-dev
make social-rpc-dev
make social-api-dev
make im-rpc-dev
make im-api-dev
make im-ws-dev
make task-mq-dev
```

### 一键安装/部署到测试环境

```bash
make install-server
```

等价于执行 `deploy/script/release-test.sh`。

### 直接运行（开发调试）

每个服务目录都提供了 `exec.sh` 与对应的 `*.go` 入口，例如：

```bash
cd apps/user/rpc && go run user.go -f etc/user.yaml
cd apps/user/api && go run user.go -f etc/user.yaml
cd apps/im/ws    && go run im.go   -f etc/im.yaml
cd apps/task/mq  && go run task.go -f etc/task.yaml
```

> 各服务配置文件位于对应模块的 `etc/` 目录下，配置项可通过 sail 配置中心下发。

## 开发约定

- 接口使用 go-zero `*.api` DSL 定义，`*.proto` 定义 gRPC；修改后通过 `goctl` 重新生成代码。
- 所有 HTTP 响应统一通过 `pkg/resultx` 封装。
- 业务错误码定义在 `pkg/xerr` 与 `pkg/status`。
- 上下文中的用户身份等数据透传走 `pkg/ctxdata`。

## 相关文档

- 项目总览：[`../README.md`](../README.md)
- 前端项目：[`../easy-im-web/README.md`](../easy-im-web/README.md)
