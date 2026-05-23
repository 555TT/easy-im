# easy-im

一个面向毕业设计的轻量级即时通讯（IM）系统，包含基于 go-zero 的微服务后端与基于 Vue 3 的 Web 前端。

> 本仓库为本科毕业设计项目源码，主要演示从协议设计、微服务拆分、长连接网关、消息存储到前端交互的完整 IM 实现思路。

## 仓库结构

```
easy-im/
├── go-im/          # 后端：基于 go-zero 的微服务集群（user / social / im / task）
├── easy-im-web/    # 前端：Vue 3 + TypeScript + Vite Web 客户端
└── README.md
```

| 模块 | 技术栈 | 说明 |
|------|--------|------|
| [`go-im/`](./go-im) | Go 1.23 / go-zero / gRPC / WebSocket / MySQL / Redis / MongoDB / Kafka / etcd | 后端微服务集群，提供 HTTP API、RPC、WebSocket 长连接网关与异步任务消费者 |
| [`easy-im-web/`](./easy-im-web) | Vue 3 / TypeScript / Vite / Pinia / Element Plus / UnoCSS | 前端 SPA，承担登录、会话、联系人、单聊/群聊等界面 |

## 系统架构

```
                                ┌──────────────────────┐
                                │     easy-im-web      │
                                │ (Vue3 SPA, Browser)  │
                                └──────────┬───────────┘
                       HTTP / WebSocket    │
                                           ▼
                          ┌────────────────────────────────┐
                          │           API 网关             │
                          │     (APISIX / 直连 api 服务)    │
                          └──┬─────────┬──────────┬────────┘
                             │         │          │
                  ┌──────────▼──┐ ┌────▼─────┐ ┌──▼────────┐
                  │ user-api    │ │social-api│ │  im-api   │
                  └──────┬──────┘ └────┬─────┘ └─────┬─────┘
                         │ gRPC        │ gRPC        │ gRPC
                  ┌──────▼──────┐ ┌────▼─────┐ ┌─────▼─────┐     ┌─────────────┐
                  │ user-rpc    │ │social-rpc│ │  im-rpc   │◀───▶│  im-ws (长  │
                  └──────┬──────┘ └────┬─────┘ └─────┬─────┘     │  连接网关)  │
                         │             │             │           └─────┬───────┘
                         └─────────────┼─────────────┘                 │
                                       ▼                               │
                       ┌───────────────────────────────┐               │
                       │ MySQL / MongoDB / Redis       │◀──────────────┘
                       │ Kafka  (task-mq 消费者写入)   │
                       │ etcd  (服务注册发现)          │
                       └───────────────────────────────┘
```

主要数据流：

1. 浏览器登录后通过 `ws://{host}/ws?userId={uid}` 接入 `im-ws` 网关；
2. 用户在前端发送的聊天消息通过 WebSocket 投递到 `im-ws`，再经 `im-rpc` 处理；
3. 持久化与异步分发由 `task-mq` 消费 Kafka 消息后写入 MySQL/Mongo；
4. `im-ws` 通过 push 路由把消息主动下发给目标在线用户。

## WebSocket 协议（前后端约定）

连接地址：`ws://{host}/ws?userId={userId}`

消息帧：

```ts
interface WSMessage {
  id?: string          // 消息唯一 ID（用于 ACK）
  frameType: FrameType // 帧类型
  ackSeq?: number      // ACK 序号
  method?: string      // 路由方法
  userId?: string      // 用户 ID
  formId?: string      // 表单 ID
  data?: any           // 消息数据
}

enum FrameType {
  FrameData  = 0x0, // 数据帧
  FramePing  = 0x1, // 心跳
  FrameAck   = 0x2, // 确认帧
  FrameNoAck = 0x3, // 无需确认
  FrameError = 0x9, // 错误帧
}
```

路由列表：

| Method | 方向 | 作用 |
|--------|------|------|
| `user.online` | C → S | 上线/拉取在线用户 |
| `conversation.chat` | C → S | 发送聊天消息 |
| `conversation.markChat` | C → S | 标记已读 |
| `push` | S → C | 服务端主动推送 |

详细 DTO 定义见 [`easy-im-web/CLAUDE.md`](./easy-im-web/CLAUDE.md) 与 `go-im/apps/im/ws/`。

## 快速开始

### 1. 启动依赖中间件

```bash
cd go-im
docker compose up -d        # 拉起 etcd / mysql / redis / mongo / kafka 等
```

### 2. 启动后端服务

参考 [`go-im/README.md`](./go-im/README.md)，依次启动：

```bash
cd go-im
make release-test           # 一键编译/启动 user / social / im / task 全部服务
```

或按模块分别启动 `user-api / user-rpc / social-api / social-rpc / im-api / im-rpc / im-ws / task-mq`。

### 3. 启动前端

```bash
cd easy-im-web
npm install
npm run dev                 # 默认 http://localhost:5173
```

## 模块文档

- 后端开发说明：[`go-im/README.md`](./go-im/README.md)
- 前端开发说明：[`easy-im-web/README.md`](./easy-im-web/README.md)
- 前端协议细节：[`easy-im-web/CLAUDE.md`](./easy-im-web/CLAUDE.md)

## 许可证

本项目为毕业设计学习用途，仅供个人学习与课程展示使用。
