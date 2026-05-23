// View models — what components consume.

export const ChatType = { Single: 1, Group: 2 } as const
export type ChatTypeValue = typeof ChatType[keyof typeof ChatType]

export const MType = { Text: 0 } as const
export type MTypeValue = typeof MType[keyof typeof MType]

export interface Friend {
  id: string
  userId: string
  nickname: string
  avatar: string
  remark: string
  online: boolean
}

export interface FriendRequest {
  id: string
  fromUserId: string
  message: string
  reqTime: number
  handleResult: number
}

export interface Conversation {
  conversationId: string
  chatType: ChatTypeValue
  peerUserId: string
  peerNickname: string
  peerAvatar: string
  lastContent: string
  lastTime: number
  unread: number
}

export type MessageStatus = 'sending' | 'sent' | 'failed'

export interface Message {
  msgId: string
  clientMsgId?: string
  conversationId: string
  chatType: ChatTypeValue
  sendId: string
  recvId: string
  content: string
  mType: MTypeValue
  sendTime: number
  status: MessageStatus
}
