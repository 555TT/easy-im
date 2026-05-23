import type { ChatTypeValue, MTypeValue } from '@/types/domain'

export const FrameType = {
  Data: 0x0,
  Ping: 0x1,
  Ack: 0x2,
  NoAck: 0x3,
  Error: 0x9,
} as const
export type FrameTypeValue = typeof FrameType[keyof typeof FrameType]

export interface WSMessage<T = unknown> {
  id?: string
  frameType: FrameTypeValue
  ackSeq?: number
  method?: string
  userId?: string
  fromId?: string
  data?: T
}

export interface WSChatPayload {
  conversationId: string
  chatType: ChatTypeValue
  sendId: string
  recvId: string
  sendTime: number
  msg: {
    msgId: string
    mType: MTypeValue
    content: string
    readRecords: Record<string, string>
  }
}

export interface WSMarkReadPayload {
  chatType: ChatTypeValue
  conversationId: string
  recvId: string
  msgIds: string[]
}

export interface WSPushPayload {
  conversationId: string
  chatType: ChatTypeValue
  sendId: string
  recvId: string
  recvIds: string[]
  sendTime: number
  msgId: string
  readRecords: Record<string, string>
  contentType: number
  mType: MTypeValue
  content: string
}
