# easy-im-web

`easy-im` 项目的 Web 前端，基于 Vue 3 + TypeScript + Vite 构建，与 [`go-im`](../go-im) 后端通过 HTTP 与 WebSocket 协作，提供登录、会话、联系人、单聊与群聊等界面。

## 技术栈

- **框架**：Vue 3（Composition API + `<script setup>`）
- **语言**：TypeScript
- **构建工具**：Vite
- **路由**：Vue Router 4
- **状态管理**：Pinia
- **UI 组件**：Element Plus（含 `@element-plus/icons-vue`）
- **样式**：UnoCSS（含 `@unocss/preset-uno` / `@unocss/preset-icons`）
- **HTTP**：Axios
- **长连接**：原生 `WebSocket`，封装在 `composables/useWebSocket.ts`

## 目录结构

```
easy-im-web/
├── index.html
├── package.json
├── vite.config.ts
├── uno.config.ts
├── tsconfig*.json
├── public/
└── src/
    ├── main.ts                  # 入口（注册 Pinia / Router / Element Plus / UnoCSS）
    ├── App.vue
    ├── style.css
    ├── api/                     # HTTP 请求封装
    ├── assets/                  # 静态资源
    ├── components/              # 业务组件
    │   ├── Chat/
    │   ├── Contact/
    │   └── Common/
    ├── composables/
    │   └── useWebSocket.ts      # WebSocket 连接、心跳与消息分发
    ├── layouts/                 # 布局组件
    ├── router/
    │   └── index.ts             # 路由：/login /chat /contact
    ├── stores/                  # Pinia
    │   ├── user.ts
    │   ├── chat.ts
    │   └── contact.ts
    ├── types/
    │   └── im.ts                # 与后端共享的 IM 协议类型
    ├── utils/                   # 工具函数
    └── views/
        ├── Login.vue
        ├── Chat/
        │   ├── index.vue
        │   └── Single.vue
        └── Contact/
            └── index.vue
```

## 路由设计

| 路径 | 视图 | 说明 |
|------|------|------|
| `/login` | `views/Login.vue` | 登录页 |
| `/` | 重定向到 `/chat` | 默认首页 |
| `/chat` | `views/Chat/index.vue` | 聊天主界面 |
| `/contact` | `views/Contact/index.vue` | 联系人 / 在线用户 |

## 与后端的协议约定

连接地址：`ws://{host}/ws?userId={userId}`

帧类型：

```ts
enum FrameType {
  FrameData  = 0x0,
  FramePing  = 0x1,
  FrameAck   = 0x2,
  FrameNoAck = 0x3,
  FrameError = 0x9,
}
```

支持的方法：

| Method | 方向 | 作用 |
|--------|------|------|
| `user.online` | C → S | 上线/获取在线用户列表 |
| `conversation.chat` | C → S | 发送聊天消息 |
| `conversation.markChat` | C → S | 已读标记 |
| `push` | S → C | 服务端推送消息 |

完整的 `WSMessage / Chat / MarkRead / Push / Msg / ChatType / MType` 等类型定义见 [`CLAUDE.md`](./CLAUDE.md) 与 `src/types/im.ts`。

约定：

- 前端通过 query 参数 `userId` 标识用户；
- 消息 `id` 用于 ACK 确认；
- 服务端采用 `OnlyAck` 模式：服务端回 ACK 即视为送达；
- 时间字段均为毫秒时间戳；
- 单聊 `conversationId` = `sendId + recvId`（按排序后拼接）。

## 开发与运行

### 安装依赖

```bash
npm install
```

### 本地开发

```bash
npm run dev
```

默认在 `http://localhost:5173` 启动 Vite 开发服务器，需要确保 `go-im` 的 `im-ws` 与各 `*-api` 服务已启动。

### 生产构建

```bash
npm run build
```

输出到 `dist/`，会先执行 `vue-tsc -b` 做类型检查，再用 Vite 打包。

### 本地预览生产包

```bash
npm run preview
```

## 开发约定

- 视图统一使用 `<script setup lang="ts">`；
- 全局状态用 Pinia store，按领域拆分（`user / chat / contact`）；
- HTTP 请求集中在 `src/api/` 下封装 Axios 实例；
- WebSocket 操作通过 `composables/useWebSocket.ts` 暴露的方法调用，避免视图层直接持有连接；
- IM 相关 DTO 全部从 `src/types/im.ts` 引入，类型以后端定义为准。

## 相关文档

- 项目总览：[`../README.md`](../README.md)
- 后端项目：[`../go-im/README.md`](../go-im/README.md)
- 前端协议与设计细节：[`./CLAUDE.md`](./CLAUDE.md)
