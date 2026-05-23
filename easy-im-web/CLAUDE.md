# easy-im Frontend (Vue3)

## 项目概述

基于 Vue3 的即时通讯（IM）前端项目，与 `go-im` 后端配合使用。

## 技术栈

- **框架**: Vue 3 (Composition API + `<script setup>`)
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **WebSocket**: 原生 WebSocket 或 `socket.io-client`
- **HTTP 客户端**: Axios
- **构建工具**: Vite
- **UI 组件库**: Element Plus 或 Ant Design Vue
- **CSS**: UnoCSS / Tailwind CSS

## 后端 API 分析

### WebSocket 消息协议

**连接地址**: `ws://{host}/ws?userId={userId}`

**消息格式**:
```typescript
interface WSMessage {
  id?: string          // 消息唯一ID
  frameType: FrameType // 帧类型
  ackSeq?: number      // ACK 序号
  method?: string      // 路由方法
  userId?: string      // 用户ID
  formId?: string      // 表单ID
  data?: any           // 消息数据
}

enum FrameType {
  FrameData  = 0x0,   // 数据帧
  FramePing  = 0x1,   // 心跳
  FrameAck   = 0x2,   // 确认帧
  FrameNoAck = 0x3,   // 无需确认
  FrameError = 0x9,   // 错误帧
}
```

### 四个路由

| Method | Handler | 作用 | Data 类型 |
|--------|---------|------|-----------|
| `user.online` | 用户上线 | 建立连接时获取所有在线用户列表 | - (无) |
| `conversation.chat` | 聊天消息 | 发送聊天消息 | `Chat` |
| `conversation.markChat` | 已读标记 | 标记消息已读 | `MarkRead` |
| `push` | 消息推送 | 服务端主动推送 | `Push` |

### 消息数据结构

```typescript
// 消息体
interface Msg {
  msgId: string
  mType: MType       // 消息类型: 0=文本
  content: string
  readRecords: Record<string, string>
}

// 聊天消息
interface Chat {
  conversationId: string  // 会话ID
  chatType: ChatType      // 1=单聊, 2=群聊
  sendId: string          // 发送者ID
  recvId: string          // 接收者ID
  sendTime: number        // 发送时间戳
  msg: Msg
}

// 已读标记
interface MarkRead {
  chatType: ChatType
  conversationId: string
  recvId: string
  msgIds: string[]
}

// 推送消息
interface Push {
  conversationId: string
  chatType: ChatType
  sendId: string
  recvId: string
  recvIds: string[]
  sendTime: number
  msgId: string
  readRecords: Record<string, string>
  contentType: ContentType
  mType: MType
  content: string
}

// 枚举
enum ChatType {
  SingleChatType = 1,  // 单聊
  GroupChatType  = 2,  // 群聊
}

enum MType {
  TextMType = 0,  // 文本消息
}
```

## 目录结构

```
src/
├── api/                  # API 请求
├── assets/              # 静态资源
├── components/          # 公共组件
│   ├── Chat/            # 聊天相关组件
│   ├── Contact/          # 联系人组件
│   └── Common/           # 通用组件
├── composables/          # 组合式函数
│   └── useWebSocket.ts   # WebSocket 连接管理
├── layouts/              # 布局组件
├── router/               # 路由配置
├── stores/               # Pinia 状态管理
│   ├── user.ts           # 用户状态
│   ├── chat.ts           # 聊天状态
│   └── contact.ts        # 联系人状态
├── types/                # TypeScript 类型定义
│   └── im.ts             # IM 相关类型
├── utils/                # 工具函数
├── views/                # 页面
│   ├── Login.vue
│   ├── register.vue
│   ├── Chat/
│   │   ├── index.vue     # 聊天主页面
│   │   ├── Single.vue    # 单聊窗口
│   │   └── Group.vue     # 群聊窗口
│   └── Contact/
│       └── index.vue     # 联系人列表
├── App.vue
└── main.ts
```
目录结构可根据具体情况再做修改
## 核心功能

### 1. 用户登录
- 输入 userId 连接到 WebSocket 服务
- 建立连接后自动调用 `user.online` 获取在线用户

### 2. 单聊 (SingleChatType = 1)
- 选择联系人发起聊天
- 发送消息到 `conversation.chat`
- 接收消息显示在聊天窗口
- 支持消息已读回执

### 3. 群聊 (GroupChatType = 2)
- 选择群组发起群聊
- 通过 `RecvIds` 群发消息

### 4. 消息列表
- 展示所有会话
- 显示最新消息预览
- 未读消息计数

### 5. 联系人管理
- 展示所有在线用户
- 用户上线/下线状态更新

## WebSocket 连接管理

```typescript
// composables/useWebSocket.ts
class WebSocketManager {
  // 连接
  connect(userId: string): void
  // 断开
  disconnect(): void
  // 发送消息
  send(method: string, data: any): void
  // 发送聊天消息
  sendChat(chat: Chat): void
  // 发送已读标记
  sendMarkRead(markRead: MarkRead): void
  // 监听消息
  onMessage(callback: (msg: WSMessage) => void): void
}
```

## 开发阶段

### Phase 1: 基础搭建
- [ ] 初始化 Vue3 + Vite 项目
- [ ] 配置 TypeScript、Pinia、Vue Router
- [ ] 配置 Element Plus / Ant Design Vue
- [ ] 创建项目目录结构
- [ ] 定义 TypeScript 类型

### Phase 2: WebSocket 通信
- [ ] 实现 WebSocket 连接管理
- [ ] 实现消息发送/接收
- [ ] 实现心跳机制
- [ ] 实现 ACK 确认机制

### Phase 3: 登录与用户
- [ ] 登录页面
- [ ] 注册页面
- [ ] 用户在线状态
- [ ] 在线用户列表

### Phase 4: 聊天功能
- [ ] 会话列表
- [ ] 单聊界面
- [ ] 群聊界面
- [ ] 消息发送与接收
- [ ] 消息时间显示
- [ ] 已读未读状态

### Phase 5: 社交功能
- [ ] 好友申请、好友列表、好友在线情况
- [ ] 创建群聊、申请进入群聊、群成员列表
### Phase 6: 优化与完善
- [ ] 消息本地缓存
- [ ] 断线重连
- [ ] 消息加载更多
- [ ] 表情包支持
- [ ] 图片/文件发送

以上所有功能/接口以go-im项目中路由中的定义为准
## 注意事项

- 前端通过 query 参数 `userId` 标识用户
- 消息 ID (`msg.id`) 用于 ACK 确认
- 服务端使用 `OnlyAck` 确认模式：服务端收到消息后回复 ACK，即算送达
- 所有时间使用毫秒时间戳
- 会话 ID (`conversationId`) = 单聊时为 `sendId + recvId`（已排序）
## 命令执行权限

Claude Code 在以下场景中执行命令时**无需用户确认**：
- 读取文件内容（如 `cat`, `head`, `tail`）
- 列出目录结构（如 `ls`, `find`）
- 搜索文件内容（如 `grep`, `rg`）
- 一些常用命令（如 `cd`, `mkdir`）
