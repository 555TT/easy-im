// 帧类型
export const FrameType = {
  FrameData: 0x0,
  FramePing: 0x1,
  FrameAck: 0x2,
  FrameNoAck: 0x3,
  FrameError: 0x9,
} as const
export type FrameType = (typeof FrameType)[keyof typeof FrameType]

// 聊天类型
export const ChatType = {
  SingleChatType: 1,
  GroupChatType: 2,
} as const
export type ChatType = (typeof ChatType)[keyof typeof ChatType]

// 消息类型
export const MType = {
  TextMType: 0,
} as const
export type MType = (typeof MType)[keyof typeof MType]

// 内容类型
export const ContentType = {
  ContentChatMsg: 1,
} as const
export type ContentType = (typeof ContentType)[keyof typeof ContentType]

// WebSocket 消息格式
export interface WSMessage {
  id?: string
  frameType: FrameType
  ackSeq?: number
  method?: string
  userId?: string
  formId?: string
  data?: any
}

// 消息体
export interface Msg {
  msgId: string
  mType: MType
  content: string
  readRecords: Record<string, string>
}

// 聊天消息
export interface Chat {
  conversationId: string
  chatType: ChatType
  sendId: string
  recvId: string
  sendTime: number
  msg: Msg
}

// 已读标记
export interface MarkRead {
  chatType: ChatType
  conversationId: string
  recvId: string
  msgIds: string[]
}

// 推送消息
export interface Push {
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

// 用户信息
export interface User {
  userId: string
  nickname?: string
  avatar?: string
  online?: boolean
}

// 会话信息
export interface Conversation {
  conversationId: string
  chatType: ChatType
  lastMsg?: Msg
  lastMsgTime?: number
  unreadCount?: number
  targetUser?: User
}
